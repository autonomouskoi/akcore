import { bus, enumName } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as intcfgpb from "/pb/internal/config_pb.js";
import { ValueUpdater } from "./vu.js";

const TOPIC_REQUEST = enumName(intcfgpb.BusTopic, intcfgpb.BusTopic.INTERNAL_REQUEST);
const TOPIC_COMMAND = enumName(intcfgpb.BusTopic, intcfgpb.BusTopic.INTERNAL_COMMAND);

class Cfg extends ValueUpdater<intcfgpb.Config> {
    constructor() {
        super(new intcfgpb.Config());
    }

    refresh() {
        bus.sendAnd(new buspb.BusMessage({
            topic: TOPIC_REQUEST,
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
        msg.topic = TOPIC_COMMAND;
        msg.type = intcfgpb.MessageTypeCommand.CONFIG_SET_REQ;
        msg.message = csr.toBinary();
        return bus.sendAnd(msg)
            .then((reply) => {
                let csResp = intcfgpb.ConfigSetResponse.fromBinary(reply.message);
                this.update(csResp.config);
            });
    }
}

class Listen extends ValueUpdater<boolean> {
    private _cfg: Cfg;

    constructor(cfg: Cfg) {
        super(false);
        this._cfg = cfg;
        this._cfg.subscribe((cfg) => {
            this.update(cfg.listenAddress === '0.0.0.0:8011');
        })
    }

    save(v: boolean): Promise<void> {
        let cfg = this._cfg.last.clone();
        cfg.listenAddress = v ? '0.0.0.0:8011' : '';
        return this._cfg.save(cfg);
    }
}

class Controller {
    private _cfg: Cfg;
    private _ready: Promise<void>;

    listenAddress: Listen;

    constructor() {
        this._cfg = new Cfg();

        this.listenAddress = new Listen(this._cfg);
        this._ready = new Promise<void>((resolve) => {
            bus.waitForTopic(TOPIC_REQUEST, 5000)
                .then(() => {
                    this._cfg.refresh();
                    let unsub = this._cfg.subscribe(() => {
                        resolve();
                        unsub();
                    })
                });
        })
    }

    ready(): Promise<void> {
        return this._ready;
    }
}

export { Controller };