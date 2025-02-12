import * as controlpb from "/pb/modules/control_pb.js";
import * as manifestpb from "/pb/modules/manifest_pb.js";
import { GloballyStyledHTMLElement } from "./global-styles.js";
import { AKPanelListItem, AKPanel } from "./ak_panel.js";

class ModuleList extends GloballyStyledHTMLElement {
    onSelectModule = (id: string) => { };
    private _akPanelListItem = new AKPanelListItem();

    constructor() {
        super();

        this.shadowRoot.innerHTML = `
<style>
.module-name {
    display: flex;
    align-items: center;
    justify-content: center;
}
.state-started {
    font-weight: bold;
}
</style>
<div class="flex-column"></div>
`;
        this._akPanelListItem.addEventListener('click', () => this.onSelectModule('akctrl'));
    }

    set modules(entries: controlpb.ModuleListEntry[]) {
        let div = this.shadowRoot.querySelector('div');
        div.textContent = '';
        div.appendChild(this._akPanelListItem);
        entries.toSorted((a, b) => a.manifest.title.localeCompare(b.manifest.title))
            .forEach((entry) => {
                let item = new ModuleListItem(entry.manifest);
                item.id = `mod-${entry.manifest.id}`;
                if (entry.state.moduleState == controlpb.ModuleState.STARTED) {
                    item.classList.add('state-started');
                }
                item.addEventListener('click', () => this.onSelectModule(entry.manifest.id));
                div.appendChild(item);
            })
    }

    handleStateEvent(e: controlpb.ModuleCurrentStateEvent) {
        let item = this.shadowRoot.querySelector(`#mod-${e.moduleId}`);
        if (!item) {
            return;
        }
        if (e.moduleState === controlpb.ModuleState.STARTED) {
            item.classList.add('state-started');
        } else {
            item.classList.remove('state-started');
        }
    }
}
customElements.define('ak-module-list', ModuleList);

class ModuleListItem extends HTMLDivElement {
    constructor(manifest: manifestpb.Manifest) {
        super();

        this.classList.add('flex-row');

        this.innerHTML = `
<img src="/m/${manifest.id}/icon" width=48 height=48
        title="${manifest.title}: ${manifest.description}"
/>
<div class="module-name">${manifest.title}</div>
`;
    }
}
customElements.define('ak-module-list-item', ModuleListItem, { extends: "div" });

export { ModuleList };