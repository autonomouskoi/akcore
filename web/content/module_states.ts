import { GloballyStyledHTMLElement } from "./global-styles.js";
import { bus, enumName, Status } from "/bus.js";
import * as buspb from "/pb/bus/bus_pb.js";
import * as controlpb from "/pb/modules/control_pb.js";
import * as manifestpb from "/pb/modules/manifest_pb.js";

const TOPIC_COMMAND = enumName(controlpb.BusTopics, controlpb.BusTopics.MODULE_COMMAND);
const TOPIC_EVENT = enumName(controlpb.BusTopics, controlpb.BusTopics.MODULE_EVENT);
const TOPIC_REQUEST = enumName(controlpb.BusTopics, controlpb.BusTopics.MODULE_REQUEST);

const ICON_ACTION_UNKNOWN = '?';
const ICON_ACTION_START = 'START';
const ICON_ACTION_STOP = '️STOP';

const ICON_PATH_HELP = 'help.svg';
const ICON_PATH_OBS = 'OBS_Studio_Logo.svg';
const ICON_PATH_LINK = 'links-line.svg';
const ICON_PATH_CTRL = 'equalizer-line.svg';

class ModuleLink extends HTMLElement {
    constructor(webPath: manifestpb.ManifestWebPath) {
        super();

        let imgSrc = ICON_PATH_LINK;
        let typeDesc = 'general web link';
        switch (webPath.type) {
            case manifestpb.ManifestWebPathType.HELP:
                imgSrc = ICON_PATH_HELP;
                typeDesc = 'Help';
                break;
            case manifestpb.ManifestWebPathType.OBS_OVERLAY:
                imgSrc = ICON_PATH_OBS;
                typeDesc = 'An OBS Overlay';
                break;
            case manifestpb.ManifestWebPathType.EMBED_CONTROL:
                imgSrc = ICON_PATH_CTRL;
                typeDesc = 'Module controls';
                break;
        }
        this.innerHTML = `
<div>
        <a href="${webPath.path}" target="_blank" rel="noopener noreferrer" >
            <img src="${imgSrc}" alt="${typeDesc}" width="16" height="16" />
            ${webPath.description}
        </a>
</div>
`;
    }

    update(webPath: manifestpb.ManifestWebPath) {
    }
}
customElements.define('module-link', ModuleLink);

class ModuleDetails extends HTMLElement {
    constructor(manifest: manifestpb.Manifest) {
        super();

        this.innerHTML = `
<div class="row-flex">
    <fieldset><legend>Description</legend>${manifest?.description}</fieldset>
    <fieldset id="links"><legend>Links</legend></fieldset>
</div>
`;

        let links = this.querySelector('#links');
        manifest.webPaths.forEach((webPath) => {
            links.appendChild(new ModuleLink(webPath));
        });
    }
}
customElements.define('module-details', ModuleDetails);

class ModuleState {
    private _buttonState: HTMLButtonElement;
    private _checkAutoStart: HTMLInputElement;
    private _divState: HTMLDivElement;
    private _state = controlpb.ModuleState.UNSPECIFIED;

    private _elements: HTMLElement[] = [];

    setAutostart = (autostart: boolean) => { };
    setState = (state: controlpb.ModuleState) => { };

    constructor(mle: controlpb.ModuleListEntry) {
        let buttonExpand = document.createElement('button');
        buttonExpand.innerHTML = '&#x25B6;';
        this._elements.push(buttonExpand);

        let name = document.createElement('div');
        name.classList.add('name');
        name.innerText = mle.manifest.name;
        this._elements.push(name);

        let description = document.createElement('div');
        description.classList.add('description');
        description.innerText = mle.manifest.description;
        this._elements.push(description);

        this._divState = document.createElement('div');
        this._divState.innerText = '?';
        this._elements.push(this._divState);

        this._buttonState = document.createElement('button');
        this._buttonState.innerHTML = ICON_ACTION_UNKNOWN;
        this._buttonState.addEventListener('click', () => this._stateButtonClicked());
        this._elements.push(this._buttonState);

        let autoStartDiv = document.createElement('div');
        this._elements.push(autoStartDiv);

        let labelAutostart: HTMLLabelElement = document.createElement('label')
        labelAutostart.htmlFor = `autostart-${mle.manifest.id}`;
        labelAutostart.innerText = 'Auto Start:';
        autoStartDiv.appendChild(labelAutostart);

        this._checkAutoStart = document.createElement('input');
        this._checkAutoStart.type = 'checkbox';
        this._checkAutoStart.id = labelAutostart.htmlFor;
        this._checkAutoStart.addEventListener('change', () => this._autostartClicked());
        this.autoStart = mle.state.config.automaticStart;
        autoStartDiv.appendChild(this._checkAutoStart)

        let details = new ModuleDetails(mle.manifest);
        details.style.setProperty('display', 'none');
        this._elements.push(details);
        buttonExpand.addEventListener('click', () => {
            let hidden = details.style.getPropertyValue('display') === 'none';
            details.style.setProperty('display', hidden ? 'block' : 'none');
            buttonExpand.innerHTML = hidden ? '&#x25BC;' : '&#x25B6;';
        })

        this.state = mle.state.moduleState;
    }

    get elements(): HTMLElement[] {
        return this._elements;
    }

    set autoStart(value: boolean) {
        this._checkAutoStart.checked = value;
    }
    
    private _autostartClicked() {
        this.setAutostart(this._checkAutoStart.checked);
        if (this._checkAutoStart.checked && this._state !== controlpb.ModuleState.STARTED) {
            this._stateButtonClicked();
        }
    }

    set state(value: controlpb.ModuleState) {
        this._state = value;
        this._divState.innerText = enumName(controlpb.ModuleState, value);
        switch (value) {
            case controlpb.ModuleState.UNSTARTED:
                this._buttonState.innerHTML = ICON_ACTION_START;
                break;
            case controlpb.ModuleState.STARTED:
                this._buttonState.innerHTML = ICON_ACTION_STOP;
                break;
            case controlpb.ModuleState.STOPPED:
                this._buttonState.innerHTML = ICON_ACTION_START;
                break;
            case controlpb.ModuleState.FAILED:
                this._buttonState.innerHTML = ICON_ACTION_START;
                break;
            case controlpb.ModuleState.FINISHED:
                this._buttonState.innerHTML = ICON_ACTION_START;
                break;
            default:
                this._buttonState.innerHTML = ICON_ACTION_UNKNOWN;
        }
    }

    private _stateButtonClicked() {
        switch (this._state) {
            case  controlpb.ModuleState.UNSTARTED:
            case  controlpb.ModuleState.STOPPED:
            case  controlpb.ModuleState.FAILED:
            case  controlpb.ModuleState.FINISHED:
                this.setState(controlpb.ModuleState.STARTED);
                return
            case  controlpb.ModuleState.STARTED:
                this.setState(controlpb.ModuleState.STOPPED);
                return
        }
    }
}

class ModuleStates extends GloballyStyledHTMLElement {
    private _mainContainer: HTMLElement;
    private _modules: { [key: string]: ModuleState } = {};

    constructor() {
        super();

        this.shadowRoot!.innerHTML = `
<style>
#module-states {
    display: grid;
    grid-template-columns: [expand] 2rem [name] auto [description] 1fr [state] 6rem [state-button] 6rem [autostart] auto [end];
    row-gap: 10px;
    column-gap: 5px;
}
div.name {
    font-weight: bold;
}
div.description {
    overflow: hidden;
    white-space: no-wrap;
    text-overflow: elipsis;
}
module-details {
    grid-column-start: 1;
    grid-column-end: span end;
}
h3 {
    margin: 5px;
    padding-left: 3rem;
}
.row-flex {
    display: flex;
    flex-direction: row;
}
.row-flex > * {
    flex: 1;
}
</style>
</style>
<h3>Modules</h3>
<div id="module-states"></div>
`;
        this._mainContainer = this.shadowRoot.querySelector('div');
        this._mainContainer.style.setProperty('display', 'none');
        bus.addStatusListener((s: Status) => this._wsStatusChange(s));
        this._wsStatusChange(bus.getStatus());

        bus.waitForTopic(TOPIC_REQUEST, 5000)
            .then(() => {
                return bus.sendAnd(new buspb.BusMessage({
                    topic: TOPIC_REQUEST,
                    type: controlpb.MessageTypeRequest.MODULES_LIST_REQ,
                    message: new controlpb.ModulesListRequest().toBinary(),
                }))
            }).then((reply) => {
                let mlr = controlpb.ModulesListResponse.fromBinary(reply.message);
                mlr.entries.sort((a, b) => a.manifest.name.localeCompare(b.manifest.name))
                    .forEach((entry) => {
                        let ms = new ModuleState(entry);
                        ms.elements.forEach((elem) => this._mainContainer.appendChild(elem))
                        this._modules[entry.manifest.id] = ms;
                        ms.setAutostart = (autostart: boolean) => this._setAutostart(entry.manifest.id, autostart);
                        ms.setState = (state: controlpb.ModuleState) => this._setState(entry.manifest.id, state);
                    });
            }).catch((e) => { console.log(`error: ${e}`) });

        bus.subscribe(TOPIC_EVENT, (msg) => this._handleEvent(msg));
    }

    private _handleEvent(msg: buspb.BusMessage) {
        switch (msg.type) {
            case controlpb.MessageTypeEvent.MODULE_CURRENT_STATE:
                this._handleEventModuleCurrentState(msg);
                return;
            default:
                console.log(`unhandled event type ${msg.type}`);
        }
    }

    private _handleEventModuleCurrentState(msg: buspb.BusMessage) {
        let mcse = controlpb.ModuleCurrentStateEvent.fromBinary(msg.message);
        let ms = this._modules[mcse.moduleId];
        if (!ms) {
            console.log(`no such module: ${mcse.moduleId}`);
            return;
        }
        ms.autoStart = mcse.config.automaticStart;
        ms.state = mcse.moduleState;
    }

    private _setAutostart(moduleID: string, autostart: boolean) {
        bus.send(new buspb.BusMessage({
            topic: TOPIC_COMMAND,
            type: controlpb.MessageTypeCommand.MODULE_AUTOSTART_SET_REQ,
            message: new controlpb.ModuleAutostartSetRequest({
                moduleId: moduleID,
                autostart: autostart,
            }).toBinary(),
        }))
    }

    private _setState(moduleID: string, state: controlpb.ModuleState) {
        bus.send(new buspb.BusMessage({
            topic: TOPIC_COMMAND,
            type: controlpb.MessageTypeCommand.MODULE_STATE_SET_REQ,
            message: new controlpb.ModuleStateSetRequest({
                moduleId: moduleID,
                state: state,
            }).toBinary(),
        }));
    }

    private _wsStatusChange(s: Status) {
        switch(s) {
            case Status.Connecting:
            case Status.NotConnected:
                this._mainContainer.style.setProperty('display', 'none');
                return;
            case Status.Connected:
                this._mainContainer.style.removeProperty('display');
                return;
        }
    }
}
customElements.define('module-states', ModuleStates);

export { ModuleStates };