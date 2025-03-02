import { SectionHelp } from './help.js';
import { ValueSubscriber, ValueUpdater } from './vu.js';

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

export { CfgUpdater, ControlPanel, UpdatingControlPanel };