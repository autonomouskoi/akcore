import { bus, enumName } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as controlpb from "/pb/modules/control_pb.js";
import * as manifestpb from "/pb/modules/manifest_pb.js";

const ICON_ACTION_UNKNOWN = '?';
const ICON_ACTION_START = 'START';
const ICON_ACTION_STOP = '️STOP';

const ICON_PATH_OBS = 'OBS_Studio_Logo.svg';
const ICON_PATH_LINK = 'links-line.svg';
const ICON_PATH_CTRL = 'equalizer-line.svg';

class ModuleLink extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
    }

    update(webPath: manifestpb.ManifestWebPath) {
        let imgSrc = ICON_PATH_LINK;
        let typeDesc = 'general web link';
        switch (webPath.type) {
            case manifestpb.ManifestWebPathType.OBS_OVERLAY:
                imgSrc = ICON_PATH_OBS;
                typeDesc = 'An OBS Overlay';
                break;
            case manifestpb.ManifestWebPathType.EMBED_CONTROL:
                imgSrc = ICON_PATH_CTRL;
                typeDesc = 'Module controls';
                break;
        }
        this.shadowRoot!.innerHTML = `
<div>
        <a href="${webPath.path}" target="_blank" rel="noopener noreferrer" >
            <img src="${imgSrc}" alt="${typeDesc}" width="16" height="16" />
            ${webPath.description}
        </a>
</div>
`;
    }
}

customElements.define('module-link', ModuleLink);

class ModuleState extends HTMLElement {
    private _expanded: boolean;
    private _id: string;
    private _name: string = '';
    private _state: string;
    private _stateButtonDisabled: boolean;
    private _icon: string;
    private _manifest: manifestpb.Manifest;
    private _autostart: boolean;

    constructor() {
        super();
        this._icon = ICON_ACTION_UNKNOWN;
        this.attachShadow({ mode: 'open' });
    }

    update() {
        this.shadowRoot!.innerHTML = `
<style>
.module {
    border: 1px solid black;
    border-radius: 5px;
    padding: 1rem;
}
.controls {
    display: flex;
    flex-direction: row;
}
.module-block {
    display: inline-block;
}
#module-autostart {
    padding-left: 1rem;
}
#module-name {
    width: 20rem;
}
#module-state {
    width: 10rem;
}
#state-button {
    font-family: sans-serif;
    width: 4rem;
}
.column-flex {
    display: flex;
    flex-direction: column;
}
.row-flex {
    display: flex;
    flex-direction: row;
}
.row-flex > * {
    flex: 1;
}
</style>
<div class="module">
    <div class="controls">
        <button id="expand-button" class="module-block">${this._expanded ? "-" : "+"}</button>
        <div id="module-name" class="module-block">${this._name}</div>
        <div id="module-state" class="module-block">${this._state}</div>
        <div id="module-controls" class="module-block">
            <button id="state-button" ${this._stateButtonDisabled ? "disabled" : ""}>${this._icon}</button>
        </div>
        <div id="module-autostart" class="module-block">
            <label for="autostart-check">Auto Start:</label>
            <input type="checkbox" id="autostart-check" ${this._autostart ? "checked" : ""}></input>
        </div>
    </div>
    <div id="expanded-box" class="column-flex" style="display: ${this._expanded ? "block" : "none"}">
        <div class="row-flex">
            <fieldset><legend>Description</legend>${this._manifest?.description}</fieldset>
            <fieldset id="links"><legend>Links</legend></fieldset>
        </div>
        <div id="controls"><em>Inline controls go here</em></div>
    </div>
</div>
`;
        (this.shadowRoot.querySelector("#state-button") as HTMLButtonElement).onclick = () => this._stateButtonClicked();
        (this.shadowRoot.querySelector("#expand-button") as HTMLButtonElement).onclick = () => this._expandButtonClicked();
        (this.shadowRoot.querySelector("#autostart-check") as HTMLImageElement).onchange = (ev) => this._autostartClicked(ev);
        if (!this._manifest || !this._manifest.webPaths) {
            return;
        }
        let links = this.shadowRoot.querySelector("#links") as HTMLElement;
        this._manifest.webPaths.forEach((webPath) => {
            let pathLink = document.createElement('module-link') as ModuleLink;
            pathLink.update(webPath);
            links.appendChild(pathLink);
        })
    }

    set id(value: string) {
        this._id = value;
    }

    get name() {
        return this._name;
    }

    set name(value: string) {
        this._name = value;
    }

    set autostart(value: boolean) {
        this._autostart = value;
        this.update();
    }

    set state(value: string) {
        this._state = value;
        let disabled = false;
        switch (value) {
            case enumName(controlpb.ModuleState, controlpb.ModuleState.UNSTARTED):
                this._stateButtonDisabled = false;
                this._icon = ICON_ACTION_START;
                break;
            case enumName(controlpb.ModuleState, controlpb.ModuleState.STARTED):
                this._stateButtonDisabled = false;
                this._icon = ICON_ACTION_STOP;
                break;
            case enumName(controlpb.ModuleState, controlpb.ModuleState.STOPPED):
                this._stateButtonDisabled = false;
                this._icon = ICON_ACTION_START;
                break;
            case enumName(controlpb.ModuleState, controlpb.ModuleState.FAILED):
                this._stateButtonDisabled = false;
                this._icon = ICON_ACTION_START;
                break;
            case enumName(controlpb.ModuleState, controlpb.ModuleState.FINISHED):
                this._stateButtonDisabled = false;
                this._icon = ICON_ACTION_START;
                break;
            default:
                this._stateButtonDisabled = true;
                this._icon = ICON_ACTION_UNKNOWN;
        }
        this.update();
    }
    set manifest(manifest: manifestpb.Manifest) {
        this._manifest = manifest;
        this.update();
    }

    private _stateButtonClicked() {
        this._stateButtonDisabled = true; // disable after click
        this.update();
        let cms = new controlpb.ChangeModuleState();
        cms.moduleId = this._id;
        switch (this._state) {
            case enumName(controlpb.ModuleState, controlpb.ModuleState.UNSTARTED):
                cms.moduleState = controlpb.ModuleState.STARTED;
                break;
            case enumName(controlpb.ModuleState, controlpb.ModuleState.STARTED):
                cms.moduleState = controlpb.ModuleState.STOPPED;
                break;
            case enumName(controlpb.ModuleState, controlpb.ModuleState.STOPPED):
                cms.moduleState = controlpb.ModuleState.STARTED;
                break;
            case enumName(controlpb.ModuleState, controlpb.ModuleState.FAILED):
                cms.moduleState = controlpb.ModuleState.STARTED;
                break;
            case enumName(controlpb.ModuleState, controlpb.ModuleState.FINISHED):
                cms.moduleState = controlpb.ModuleState.STARTED;
                break;
        }
        let msg = new buspb.BusMessage()
        msg.topic = enumName(controlpb.BusTopics, controlpb.BusTopics.CONTROL);
        msg.type = controlpb.MessageType.TYPE_CHANGE_STATE;
        msg.message = cms.toBinary();
        bus.send(msg);
    }

    private _expandButtonClicked() {
        this._expanded = !this._expanded;
        this.update();
    }

    private _autostartClicked(event: Event) {
        let target = event.target as HTMLInputElement;
        let msg = new buspb.BusMessage();
        msg.topic = enumName(controlpb.BusTopics, controlpb.BusTopics.CONTROL);
        msg.type = controlpb.MessageType.TYPE_CHANGE_MODULE_AUTOSTART;
        let cma = new controlpb.ChangeModuleAutostart();
        cma.moduleId = this._id;
        cma.autostart = target.checked;
        msg.message = cma.toBinary();
        bus.send(msg);
    }
}

customElements.define('module-state', ModuleState);

class ModuleStates extends HTMLElement {
    private _modules: { [key: string]: ModuleState };
    private _mainContainer: HTMLElement;

    constructor() {
        super();
        this._modules = {};
        this.attachShadow({ mode: 'open' });
        this.shadowRoot!.innerHTML = `
<div id="module-states"></div>
`;
        this._mainContainer = this.shadowRoot.querySelector("#module-states") as HTMLElement;
    }

    update() {
        this._mainContainer.textContent = '';
        Object.keys(this._modules)
            .map((key: string) => this._modules[key])
            .toSorted((a: ModuleState, b: ModuleState) => a.name.localeCompare(b.name))
            .forEach((ms: ModuleState) => this._mainContainer.appendChild(ms));
    }

    connectedCallback() {
        bus.subscribe(
            enumName(controlpb.BusTopics, controlpb.BusTopics.STATE),
            (msg: buspb.BusMessage) => this.handleBusMessage(msg),
        );
        // get all current states
        let msg = new buspb.BusMessage();
        msg.topic = enumName(controlpb.BusTopics, controlpb.BusTopics.CONTROL);
        msg.type = controlpb.MessageType.TYPE_GET_CURRENT_STATES;
        setTimeout(() => bus.send(msg), 500);
        //bus.send(msg);
    }
    disconnectedCallback() {
        bus.unsubscribe(enumName(controlpb.BusTopics, controlpb.BusTopics.STATE));
    }

    handleBusMessage(msg: buspb.BusMessage) {
        switch (msg.type) {
            case controlpb.MessageType.TYPE_CURRENT_STATE:
                let state = controlpb.CurrentModuleState.fromBinary(msg.message);
                let modState = this._modules[state.moduleId];
                if (!modState) {
                    modState = document.createElement('module-state') as ModuleState;
                    modState.id = state.moduleId;

                    let gmr = new controlpb.GetManifestRequest();
                    gmr.moduleId = state.moduleId;
                    let msg = new buspb.BusMessage();
                    msg.topic = enumName(controlpb.BusTopics, controlpb.BusTopics.CONTROL);
                    msg.type = controlpb.MessageType.TYPE_GET_MANIFEST_REQ;
                    msg.message = gmr.toBinary();
                    bus.sendWithReply(msg, (resp: buspb.BusMessage) => {
                        // TODO: handle error
                        if (resp.type != controlpb.MessageType.TYPE_GET_MANIFEST_RESP) {
                            return;
                        }
                        let gmResp = controlpb.GetManifestResponse.fromBinary(resp.message);
                        modState.name = gmResp.manifest.name;
                        modState.manifest = gmResp.manifest;
                        this.update();
                    });

                    this._modules[state.moduleId] = modState;
                    this.update();
                }
                modState.state = enumName(controlpb.ModuleState, state.moduleState);
                modState.autostart = state.config.automaticStart;
                break;
            default:
                console.log(`unhandled control message type: ${msg.type}`);
        }
    }
}

customElements.define('module-states', ModuleStates);

let ms = document.createElement('module-states');
document.querySelector('#ui').appendChild(ms);