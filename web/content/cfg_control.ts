import { bus, enumName } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as intcfgpb from "/pb/svc/pb/config_pb.js";
import * as controlpb from "/pb/modules/control_pb.js";
import { ValueUpdater } from "./vu.js";

const TOPIC_INT_REQUEST = enumName(intcfgpb.BusTopic, intcfgpb.BusTopic.INTERNAL_REQUEST);
const TOPIC_INT_COMMAND = enumName(intcfgpb.BusTopic, intcfgpb.BusTopic.INTERNAL_COMMAND);
const TOPIC_CTRL_COMMAND = enumName(controlpb.BusTopics, controlpb.BusTopics.MODULE_COMMAND);
const TOPIC_CTRL_EVENT = enumName(controlpb.BusTopics, controlpb.BusTopics.MODULE_EVENT);
const TOPIC_CTRL_REQUEST = enumName(controlpb.BusTopics, controlpb.BusTopics.MODULE_REQUEST);

class InternalConfig extends ValueUpdater<intcfgpb.Config> {
    constructor() {
        super(new intcfgpb.Config());
    }

    refresh() {
        bus.sendAnd(new buspb.BusMessage({
            topic: TOPIC_INT_REQUEST,
            type: intcfgpb.MessageTypeRequest.CONFIG_GET_REQ,
            message: new intcfgpb.ConfigGetRequest().toBinary(),
        })).then((reply) => {
            let cgResp = intcfgpb.ConfigGetResponse.fromBinary(reply.message);
            this.update(cgResp.config);
        });
    }

    save(cfg: intcfgpb.Config): Promise<void> {
        let csr = new intcfgpb.ConfigSetRequest();
        csr.config = cfg;
        let msg = new buspb.BusMessage();
        msg.topic = TOPIC_INT_COMMAND;
        msg.type = intcfgpb.MessageTypeCommand.CONFIG_SET_REQ;
        msg.message = csr.toBinary();
        return bus.sendAnd(msg)
            .then((reply) => {
                let csResp = intcfgpb.ConfigSetResponse.fromBinary(reply.message);
                this.update(csResp.config);
            });
    }
}

type ModuleCurrentStateReceiver = (state: controlpb.ModuleCurrentStateEvent) => void;

class Controller {
    private _cfg: InternalConfig;
    private _ready: Promise<void>;

    private _moduleStateListeners: ModuleCurrentStateReceiver[] = [];

    onManifestSelect = (entry: controlpb.ModuleListEntry) => { };

    constructor() {
        this._cfg = new InternalConfig();

        this._ready = new Promise<void>((resolve) => {
            bus.waitForTopic(TOPIC_INT_REQUEST, 5000)
                .then(() => {
                    this._cfg.refresh();
                    let unsub = this._cfg.subscribe(() => {
                        resolve();
                        unsub();
                    })
                });
            bus.subscribe(TOPIC_CTRL_EVENT, (msg) => this._recvCtrlEvent(msg));
        })
    }

    cfg(): InternalConfig {
        return this._cfg;
    }

    ready(): Promise<void> {
        return this._ready;
    }

    list_modules(): Promise<controlpb.ModuleListEntry[]> {
        return this.ready()
            .then(() => bus.sendAnd(new buspb.BusMessage({
                topic: TOPIC_CTRL_REQUEST,
                type: controlpb.MessageTypeRequest.MODULES_LIST_REQ,
                message: new controlpb.ModulesListRequest().toBinary(),
            })).then((reply) => {
                let resp = controlpb.ModulesListResponse.fromBinary(reply.message);
                return resp.entries;
            }));
    }

    select_module(id: string) {
        this.list_modules()
            .then((entries) => {
                let entry = entries.find((entry) => entry.manifest.id === id);
                if (entry) {
                    this.onManifestSelect(entry);
                }
            })
    }

    private _recvCtrlEvent(msg: buspb.BusMessage) {
        switch (msg.type) {
            case controlpb.MessageTypeEvent.MODULE_CURRENT_STATE:
                this._recvModuleState(msg);
                break;
        }
    }

    subscribeModuleState(fn: ModuleCurrentStateReceiver) {
        this._moduleStateListeners.push(fn)
    }

    unsubscribeModuleState(fn: ModuleCurrentStateReceiver) {
        this._moduleStateListeners = this._moduleStateListeners.filter((theFn) => theFn !== fn)
    }

    private _recvModuleState(msg: buspb.BusMessage) {
        let state = controlpb.ModuleCurrentStateEvent.fromBinary(msg.message);
        this._moduleStateListeners.forEach((fn) => fn(state));
    }

    changeState(moduleID: string, newState: controlpb.ModuleState) {
        bus.send(new buspb.BusMessage({
            topic: TOPIC_CTRL_COMMAND,
            type: controlpb.MessageTypeCommand.MODULE_STATE_SET_REQ,
            message: new controlpb.ModuleStateSetRequest({
                moduleId: moduleID,
                state: newState,
            }).toBinary(),
        }));
    }

    setAutostart(moduleID: string, autostart: boolean) {
        bus.send(new buspb.BusMessage({
            topic: TOPIC_CTRL_COMMAND,
            type: controlpb.MessageTypeCommand.MODULE_AUTOSTART_SET_REQ,
            message: new controlpb.ModuleAutostartSetRequest({
                moduleId: moduleID,
                autostart: autostart,
            }).toBinary(),
        }))
    }
}

export { Controller, InternalConfig };