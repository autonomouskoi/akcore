import { bus } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as svcpb from "/pb/svc/pb/svc_pb.js";
import * as oscpb from "/pb/svc/pb/osc_pb.js";
import * as intcfgpb from "/pb/svc/pb/svc_config_pb.js";
import { SectionHelp } from './help.js';
import { ValueSubscriber } from './vu.js';
import { ctrl } from './cfg_control.js';

interface ControlPanelParams {
    title: string;
    help: HTMLElement;
}

class ControlPanel extends HTMLFieldSetElement {
    private _title: string;
    private _help: HTMLElement;

    constructor(params: ControlPanelParams) {
        super();

        this._title = params.title;
        this._help = params.help;

        this.innerHTML = '';
    }

    set innerHTML(html: string) {
        super.innerHTML = html;
        let legend = document.createElement('legend');
        legend.innerHTML = `${this._title} &#9432;`;
        legend.style.fontWeight = 'bold';
        legend.style.fontSize = '1.3rem';

        let help = SectionHelp(legend, this._help);
        this.prepend(legend, help);
    }
}
customElements.define('ak-ctrl-panel', ControlPanel, { extends: 'fieldset' });

interface CfgUpdater<T> {
    subscribe: (f: ValueSubscriber<T>) => void;
    save: (cfg: T) => Promise<void>;
    last: T;
}

interface UpdatingControlPanelParams<T> extends ControlPanelParams {
    data: CfgUpdater<T>;
}

class UpdatingControlPanel<T> extends ControlPanel {
    private _data: CfgUpdater<T>;

    constructor(params: UpdatingControlPanelParams<T>) {
        super(params);

        this._data = params.data;
    }

    connectedCallback() {
        this._data.subscribe((data) => this.update(data));
        this.update(this.last);
    }

    update(v: T) { }
    get last(): T { return this._data.last }
    save(v: T) { this._data.save(v) }
    get updater(): CfgUpdater<T> { return this._data };
}
customElements.define('ak-updating-ctrl-panel', UpdatingControlPanel, { extends: 'fieldset' });

class VoidPromise {
    wait: Promise<void>;
    resolve: (value: void | PromiseLike<void>) => void;
    reject: (reason?: any) => void;

    constructor() {
        this.wait = new Promise<void>((resolve, reject) => {
            this.resolve = resolve;
            this.reject = reject;
        });
    }

    then(fn: () => void) {
        this.wait.then(() => fn());
    }
}

// This is currently operating by looking directly at the internal config. This
// should be done by listing OSC targets through svc, but there's no connection
// between the websocket handler and svc.
class OSCTargetSelect extends HTMLSelectElement {

    ready = new VoidPromise();

    constructor() {
        super();

        ctrl.ready().then(() => {
            ctrl.cfg().subscribe((newCfg) => this.update(newCfg));
            this.update(ctrl.cfg().last);
            this.ready.resolve();
        });
    }

    update(newCfg: intcfgpb.Config) {
        this.textContent = '';
        newCfg.oscConfig?.targets.forEach((target) => {
            let option = document.createElement('option');
            option.value = target.name;
            option.innerText = target.name;
            this.appendChild(option);
        });
    }
}
customElements.define('ak-osc-target-select', OSCTargetSelect, { extends: 'select' });

export { CfgUpdater, ControlPanel, OSCTargetSelect, UpdatingControlPanel, VoidPromise };