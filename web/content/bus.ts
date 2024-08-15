import { proto3 } from "@bufbuild/protobuf";
import * as buspb from "/pb/bus/bus_pb.js";

export type handler = (msg: buspb.BusMessage) => void;

class BusClient {
    private socket: WebSocket;
    private handlers: { [key: string]: handler };
    private pendingReplies: { [key: string]: handler };
    private wsAddr: URL;

    reconnect: boolean;

    constructor() {
        this.reconnect = true;

        this.wsAddr = new URL(document.location.toString());
        this.wsAddr.protocol = "ws:";
        this.wsAddr.pathname = "/ws";
        this.wsAddr.hash = "";

        this.handlers = {};
        this.pendingReplies = {};
    }

    connect() {
        if (!this.reconnect) {
            return;
        }
        this.socket = new WebSocket(this.wsAddr.toString());
        this.socket.addEventListener("open", this.socketOpened);
        this.socket.addEventListener("close", this.socketClosed);
        this.socket.addEventListener("error", this.socketError);
        this.socket.addEventListener("message", (event) => this.socketMessage(event));

        for (let topic in this.handlers) {
            this.sendsubscribe(topic);
        }
    }

    socketOpened(event: Event) {
        console.log("websocket opened");
    }
    socketClosed(event: Event) {
        console.log("websocket closed");
        // wait a second then try to reconnect
        window.setTimeout(this.connect, 1000);
    }
    socketError(event: Event) {
        console.log("websocket error: ", event);
    }
    async socketMessage(event: MessageEvent) {
        let buffer = await event.data.arrayBuffer();
        let uintBuf = new Uint8Array(buffer);
        let bm = buspb.BusMessage.fromBinary(uintBuf);
        if (bm.replyTo) {
            let handlerFn = this.pendingReplies[bm.replyTo.toString()];
            if (handlerFn) {
                handlerFn(bm);
                delete this.pendingReplies[bm.replyTo.toString()];
            } else {
                console.log(`no reply handler for ${bm.replyTo}`);
            }
            return;
        }
        let handlerFn = this.handlers[bm.topic];
        if (!handlerFn) {
            console.log("no handler for topic ", bm.topic, bm);
            return;
        }
        handlerFn(bm);
    }

    private sendsubscribe(topic: string) {
        if (this.socket.readyState != WebSocket.OPEN) {
            // not ready yet, try again in a moment
            setTimeout(() => {this.sendsubscribe(topic)}, 250);
            return;
        }
        let sub = new buspb.SubscribeRequest();
        sub.topic = topic;
        let bm = new buspb.BusMessage();
        bm.type = buspb.ExternalMessageType.SUBSCRIBE;
        bm.message = sub.toBinary();
        this.socket.send(bm.toBinary());
    }
    subscribe(topic: string, handler: handler) {
        this.sendsubscribe(topic);
        this.handlers[topic] = handler;
    }
    unsubscribe(topic: string) {
        let unsub = new buspb.UnsubscribeRequest();
        unsub.topic = topic;
        let bm = new buspb.BusMessage();
        bm.type = buspb.ExternalMessageType.UNSUBSCRIBE;
        bm.message = unsub.toBinary();
        this.socket.send(bm.toBinary());
        delete this.handlers[topic];
    }
    send(msg: buspb.BusMessage) {
        this.socket.send(msg.toBinary());
    }
    sendWithReply(msg: buspb.BusMessage, cb: handler) {
        msg.replyTo = BigInt(Math.floor(Math.random() * 0xFFFFFFFF));
        this.pendingReplies[msg.replyTo.toString()] = cb;
        this.send(msg);
    }
    waitForTopic(topic: string, timeout: number): Promise<string> {
        let expiration = new Date(new Date().getTime() + timeout).getTime();
        let htr = new buspb.HasTopicRequest();
        htr.topic = topic;
        htr.timeoutMs = 50;
        let interval = 50;
        let b = htr.toBinary();
        return new Promise<string>((resolve, reject) => {
            let checkTopic = () => {
                let now = new Date().getTime();
                if (expiration < now) {
                    reject('expired'); 
                    return;
                }
                if (this.socket.readyState != WebSocket.OPEN) {
                    setTimeout(() => checkTopic(), interval);
                    return;
                }
                let msg = new buspb.BusMessage();
                msg.type = buspb.ExternalMessageType.HAS_TOPIC;
                msg.message = b;
                this.sendWithReply(msg, (reply: buspb.BusMessage) => {
                    if (reply.error) {
                        setTimeout(() => checkTopic(), interval);
                        return;
                    }
                    let htResp = buspb.HasTopicResponse.fromBinary(reply.message);
                    if (htResp.hasTopic) {
                        resolve(htResp.topic);
                        return;
                    }
                    setTimeout(() => checkTopic(), interval);
                });
            };
            checkTopic();
        });
    }
}

interface EnumObject {
    [key: number]: string;
    [k: string]: number | string;
}

function enumName(enumT: EnumObject, value: number ) {
    return proto3.getEnumType(enumT).values[value].localName;
}

let bus = new BusClient();
bus.connect();
export { bus, enumName };