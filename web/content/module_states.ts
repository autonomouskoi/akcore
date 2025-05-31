import { AKPanel } from "./ak_panel.js";
import { GloballyStyledHTMLElement } from "./global-styles.js";
import { bus, enumName, Status } from "/bus.js";
import { InternalConfig } from './cfg_control.js';
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

const ICON_TRIANGLE_RIGHT = '&#x25B6;';
const ICON_TRIANGLE_DOWN = '&#x25BC;';

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
            case manifestpb.ManifestWebPathType.CONTROL_PAGE:
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
}
customElements.define('module-link', ModuleLink);

function respondToVisibility(element: HTMLElement, callback: (visible: boolean) => void) {
    let options = { root: document.documentElement };
    let observer = new IntersectionObserver((entries, observer) => {
        entries.forEach((entry) => {
            callback(entry.intersectionRatio > 0);
        });
    }, options);
    observer.observe(element);
}

class ModuleLinks extends HTMLElement {
    constructor(manifest: manifestpb.Manifest) {
        super();

        this.style.setProperty('padding', '0.4em');
        manifest.webPaths
            .filter((webPath) =>
                webPath.type != manifestpb.ManifestWebPathType.EMBED_CONTROL
            ).forEach((webPath) => this.appendChild(new ModuleLink(webPath)));
    }
}
customElements.define('ak-module-links', ModuleLinks);

class ModuleDetails extends HTMLElement {
    private _divCtrl: HTMLDivElement;
    private _loading = false;
    private _ctrlLink: manifestpb.ManifestWebPath;

    constructor(manifest: manifestpb.Manifest) {
        super();

        this.innerHTML = `
<div class="column-flex">
    <div class="row-flex">
        <fieldset><legend>Description</legend>${manifest?.description}</fieldset>
        <fieldset id="links"><legend>Links</legend></fieldset>
    </div>
    <div id="embedded-ctrl"></div>
</div>
`;

        let links = this.querySelector('#links');
        manifest.webPaths
            .filter((webPath) =>
                webPath.type !== manifestpb.ManifestWebPathType.EMBED_CONTROL
            ).forEach((webPath) => {
                links.appendChild(new ModuleLink(webPath));
            });
        this._divCtrl = this.querySelector('#embedded-ctrl');
        manifest.webPaths.forEach((wp) => {
            if (wp.type === manifestpb.ManifestWebPathType.EMBED_CONTROL) {
                this._ctrlLink = wp;
            }
        });

        respondToVisibility(this, (visible) => this._handleVisibility(visible, manifest.name));
    }

    private _handleVisibility(visible: boolean, name: string) {
        if (!visible || this._loading || !this._ctrlLink) {
            return;
        }
        this._loading = true;
        import(this._ctrlLink.path)
            .then((mod) => mod.start(this._divCtrl))
            .catch((e) => { console.log(`failed to load ${name}: ${e}`) });
        return;
    }
}
customElements.define('module-details', ModuleDetails);

class ModuleState {
    private _buttonExpand: HTMLButtonElement;
    private _buttonState: HTMLButtonElement;
    private _checkAutoStart: HTMLInputElement;
    private _divDetails: ModuleDetails;
    private _divState: HTMLDivElement;
    private _state = controlpb.ModuleState.UNSPECIFIED;

    private _elements: HTMLElement[] = [];

    setAutostart = (autostart: boolean) => { };
    setState = (state: controlpb.ModuleState) => { };

    constructor(mle: controlpb.ModuleListEntry) {
        this._buttonExpand = document.createElement('button');
        this._buttonExpand.innerHTML = ICON_TRIANGLE_RIGHT;
        this._buttonExpand.disabled = true;
        this._elements.push(this._buttonExpand);

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

        this._divDetails = new ModuleDetails(mle.manifest);
        this._divDetails.style.setProperty('display', 'none');
        this._elements.push(this._divDetails);
        this._buttonExpand.addEventListener('click', () => {
            this._expanded = !this._divDetails.checkVisibility();
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

    private set _expanded(v: boolean) {
        this._divDetails.style.setProperty('display', v ? 'block' : 'none');
        this._buttonExpand.innerHTML = v ? ICON_TRIANGLE_DOWN : ICON_TRIANGLE_RIGHT;
    }

    set state(value: controlpb.ModuleState) {
        this._state = value;
        this._divState.innerText = enumName(controlpb.ModuleState, value);
        if (value === controlpb.ModuleState.STARTED) {
            this._buttonExpand.disabled = false;
        } else {
            this._buttonExpand.disabled = true;
            this._expanded = false;
        }
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
            case controlpb.ModuleState.UNSTARTED:
            case controlpb.ModuleState.STOPPED:
            case controlpb.ModuleState.FAILED:
            case controlpb.ModuleState.FINISHED:
                this.setState(controlpb.ModuleState.STARTED);
                return
            case controlpb.ModuleState.STARTED:
                this.setState(controlpb.ModuleState.STOPPED);
                return
        }
    }
}

class ModuleStateIndicator extends HTMLElement {
    constructor() {
        super();
    }

    set state(state: controlpb.ModuleState) {
        this.innerHTML = `Status: ${enumName(controlpb.ModuleState, state)}`;
    }
}
customElements.define('ak-module-state-indicator', ModuleStateIndicator);

class ModuleStateButton extends HTMLButtonElement {
    private _state: controlpb.ModuleState;

    private _changeState = (newState: controlpb.ModuleState) => { };

    constructor(changeState: (newState: controlpb.ModuleState) => void) {
        super();
        this.style.setProperty('width', '8rem');
        this.state = controlpb.ModuleState.UNSTARTED;

        this._changeState = changeState;
        this.addEventListener('click', () => { this._onClick() });
    }

    set state(state: controlpb.ModuleState) {
        if (this._state === state) {
            return;
        }
        this._state = state;

        switch (state) {
            case controlpb.ModuleState.UNSTARTED:
                this.innerHTML = ICON_ACTION_START;
                break;
            case controlpb.ModuleState.STARTED:
                this.innerHTML = ICON_ACTION_STOP;
                break;
            case controlpb.ModuleState.STOPPED:
                this.innerHTML = ICON_ACTION_START;
                break;
            case controlpb.ModuleState.FAILED:
                this.innerHTML = ICON_ACTION_START;
                break;
            case controlpb.ModuleState.FINISHED:
                this.innerHTML = ICON_ACTION_START;
                break;
            default:
                this.innerHTML = ICON_ACTION_UNKNOWN;
        }
    }

    private _onClick() {
        switch (this._state) {
            case controlpb.ModuleState.UNSTARTED: // these fall through
            case controlpb.ModuleState.STOPPED:
            case controlpb.ModuleState.FAILED:
            case controlpb.ModuleState.FINISHED:
                this._changeState(controlpb.ModuleState.STARTED);
                return
            case controlpb.ModuleState.STARTED:
                this._changeState(controlpb.ModuleState.STOPPED);
                return
        }
    }
}
customElements.define('ak-module-state-button', ModuleStateButton, { extends: 'button' });

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
    padding: 12px 5px 12px 5px;
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
        bus.addStatusListener((s: Status) => this._wsStatusChange(s));
        this._wsStatusChange(bus.getStatus());

        bus.subscribe(TOPIC_EVENT, (msg) => this._handleEvent(msg));
    }

    private _populate() {
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
        switch (s) {
            case Status.Connecting:
            case Status.NotConnected:
                this._modules = {};
                this._mainContainer.textContent = '';
                return;
            case Status.Connected:
                this._populate();
                return;
        }
    }
}
customElements.define('module-states', ModuleStates);

class ModuleAutostart extends HTMLElement {
    constructor(setAutostart: (autostart: boolean) => void) {
        super();

        this.innerHTML = `
<div>
    <label for="check-autostart">Autostart</label>
    <input type="checkbox" id="check-autostart" />
</div>
`;
        let input: HTMLInputElement = this.querySelector('input');
        input.addEventListener('change', () => setAutostart(input.checked));
    }

    set autostart(autostart: boolean) {
        let input: HTMLInputElement = this.querySelector('input')
        input.checked = autostart;
    }
}
customElements.define('ak-module-autostart', ModuleAutostart);

class ModulePanel extends GloballyStyledHTMLElement {
    private _manifest: manifestpb.Manifest;
    private _state: controlpb.ModuleState;

    constructor(ctrl: Controller, manifest: manifestpb.Manifest) {
        super();
        this._manifest = manifest;
        this.shadowRoot.innerHTML = `
<style>
section {
    gap: 10px;
}
.state > * {
    flex-grow: 1;
}
</style>
<section id="heading" class="flex-row">
<div style="flex-basis: fit-content">
    <img src="/m/${this._manifest.id}/icon" width="96" height="96" />
</div>
<div style="flex-grow: 1">
    <h2>${this._manifest.title}${this._manifest.version ? ' - ' + this._manifest.version : ''}</h2>        
    <p>${this._manifest.description}</p>
</div>
</section>

<section class="flex-row state">
<div>
    <ak-module-state-indicator></ak-module-state-indicator>
</div>
<div>
    <button>BONK</button>
</div>
<div id="autostart">
    <label for="check-autostart">Autostart</label>
    <input type="checkbox" id="check-autostart" />
</div>
</section>

<section id="embed-ctrls"></section>
`;
        let placeholderButton = this.shadowRoot.querySelector('button');
        let msb = new ModuleStateButton((newState: controlpb.ModuleState) => ctrl.changeState(this._manifest.id, newState));
        placeholderButton.parentElement.replaceChild(msb, placeholderButton);

        let placeholderAutostart = this.shadowRoot.querySelector('#autostart');
        let autostart = new ModuleAutostart((autostart: boolean) => ctrl.setAutostart(this._manifest.id, autostart));
        autostart.id = 'autostart';
        placeholderAutostart.parentElement.replaceChild(autostart, placeholderAutostart);

        this.shadowRoot.querySelector('section#heading').appendChild(new ModuleLinks(this._manifest));
    }

    set state(state: controlpb.ModuleState) {
        if (state === this._state) {
            return;
        }
        this._state = state;

        let stateIndicator: ModuleStateIndicator = this.shadowRoot.querySelector('ak-module-state-indicator');
        stateIndicator.state = state;
        let stateButton = this.shadowRoot.querySelector('button') as ModuleStateButton;
        stateButton.state = state;

        let embedCtrls = this.shadowRoot.querySelector('#embed-ctrls');

        if (state !== controlpb.ModuleState.STARTED) {
            embedCtrls.textContent = '';
            return;
        }

        this._manifest.webPaths
            .filter((wp) => wp.type === manifestpb.ManifestWebPathType.EMBED_CONTROL)
            .forEach((wp) => {
                let ctrl: HTMLDivElement = document.createElement('div');
                ctrl.classList.add('embed-ctrl');
                import(wp.path)
                    .then((mod) => mod.start(ctrl))
                    .then(() => { embedCtrls.appendChild(ctrl) })
                    .catch((e) => { console.log(`failed to load ${this._manifest.name}: ${e}`) });
            });
    }

    set autostart(autostart: boolean) {
        let mas: ModuleAutostart = this.shadowRoot.querySelector('#autostart');
        mas.autostart = autostart;
    }
}
customElements.define('ak-module-panel', ModulePanel);

class SingleDisplayPanel extends HTMLElement {
    constructor() {
        super();
    }

    add(element: HTMLElement) {
        element.style.setProperty('display', 'none');
        this.appendChild(element);
    }

    display(id: string) {
        Array.from(this.children)
            .filter((elem) => elem instanceof HTMLElement)
            .forEach((elem: HTMLElement) => {
                if (elem.id === id) {
                    elem.style.removeProperty('display');
                } else {
                    elem.style.setProperty('display', 'none');
                }
            });
    }
}

interface Controller {
    changeState(moduleID: string, newState: controlpb.ModuleState): void;
    setAutostart(moduleID: string, autostart: boolean): void;
    ready(): Promise<void>;
    cfg(): InternalConfig;
}

class ModulesPanel extends SingleDisplayPanel {
    private _ctrl: Controller;
    private _akPanel: AKPanel;

    constructor(ctrl: Controller) {
        super();
        this._ctrl = ctrl;

        this._akPanel = new AKPanel(ctrl);
        this._akPanel.id = 'mod-akctrl';
    }

    display(id: string) {
        super.display(`mod-${id}`);
    }

    set modules(entries: controlpb.ModuleListEntry[]) {
        this.textContent = '';
        this.add(this._akPanel);
        entries.toSorted((a, b) => a.manifest.name.localeCompare(b.manifest.name))
            .forEach((entry) => {
                let panel = new ModulePanel(this._ctrl, entry.manifest);
                panel.id = `mod-${entry.manifest.id}`;
                panel.state = entry.state.moduleState;
                panel.autostart = entry.state.config.automaticStart;
                this.add(panel);
            })
    }

    handleStateEvent(e: controlpb.ModuleCurrentStateEvent) {
        let panel: ModulePanel = this.querySelector(`#mod-${e.moduleId}`);
        if (!panel) {
            return;
        }
        panel.state = e.moduleState;
        panel.autostart = e.config.automaticStart;
    }
}
customElements.define('ak-modules-panel', ModulesPanel);

export { ModuleStates, ModulesPanel };