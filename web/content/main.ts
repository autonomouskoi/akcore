import * as controlpb from "/pb/modules/control_pb.js";
import { Controller } from "./cfg_control.js";
import { ModulesPanel } from "./module_states.js";
import { ModuleList } from "./module_list.js";
import { GloballyStyledHTMLElement } from "./global-styles.js";
import { bus, Status } from "/bus.js";

function start(container: HTMLElement) {
    container.classList.add('flex-column');

    let ctrl = new Controller();

    //container.appendChild(new StatusContainer(ctrl));
    container.appendChild(new AKContainer(ctrl));
}

class AKContainer extends GloballyStyledHTMLElement {
    private _moduleListEntries: controlpb.ModuleListEntry[] = [];

    private _modList = new ModuleList();
    private _modPanel: ModulesPanel;
    private _ctrl: Controller;

    constructor(ctrl: Controller) {
        super();
        this._ctrl = ctrl;

        this.shadowRoot.innerHTML = `
<style>
ak-modules-panel {
    flex-grow: 1;
}
.flex-row {
    width: 100%;
    gap: 1rem;
}
</style>
<div class="flex-row"></div>
`;
        this._modPanel = new ModulesPanel(ctrl);

        this._modList.onSelectModule = (id: string) => { this._modPanel.display(id)};
        ctrl.subscribeModuleState((e) => this._modList.handleStateEvent(e));
        ctrl.subscribeModuleState((e) => this._modPanel.handleStateEvent(e));

        let container = this.shadowRoot.querySelector('div');
        container.appendChild(this._modList);
        container.appendChild(this._modPanel);

        this._wsStatus(bus.getStatus());
        bus.addStatusListener((s) => this._wsStatus(s));
    }

    private _wsStatus(s: Status) {
        if (s === Status.Connected) {
            this._ctrl.list_modules().then((entries) => this._updateModules(entries));
        } else {
            this._updateModules([]);
        }
    }

    private _updateModules(entries: controlpb.ModuleListEntry[]) {
        let getIDs = (entry: controlpb.ModuleListEntry) => entry.manifest.id;
        let current = new Set(this._moduleListEntries.map(getIDs));
        let updated = new Set(entries.map(getIDs));
        if (current === updated) {
            return;
        }
        this._modList.modules = entries;
        this._modPanel.modules = entries;
    }
}
customElements.define('ak-container', AKContainer);

export { start }