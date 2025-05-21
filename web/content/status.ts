import { GloballyStyledHTMLElement } from "./global-styles.js";
import { InternalConfig } from './cfg_control.js';
import { bus, Status } from "/bus.js";

interface Controller {
    ready(): Promise<void>;
    cfg(): InternalConfig;
}

class StatusContainer extends GloballyStyledHTMLElement {
    constructor(ctrl: Controller) {
        super();

        this.shadowRoot.innerHTML = `
<style>
div {
    display: flex;
    flex-direction: row;
    margin: 5px;
    align-items: center;
}
div > * {
    flex-grow: 1;
}
.title {
    font-size: x-large;
    font-weight: bolder;
}
</style>
<div></div>
`

        let div = this.shadowRoot.querySelector('div');
        let title = document.createElement('div');
        title.classList.add('title')

        fetch('./build.json')
            .then((resp) => {
                return resp.json();
            }).then((js) => {
                title.innerHTML = `<a href="https://autonomouskoi.org" target="_blank">AutonomousKoi ${js.Build}</a>`;
            })


        ctrl.ready().then(() => {
            div.appendChild(title);
            div.appendChild(new BusConnection());
        })
    }
}
customElements.define('core-status-main-unused', StatusContainer);

class BusConnection extends GloballyStyledHTMLElement {
    constructor() {
        super();

        this._update(bus.getStatus());
        bus.addStatusListener((s) => { this._update(s) });
    }

    private _update(s: Status) {
        let color = s === Status.NotConnected ? 'red' :
            s === Status.Connecting ? 'yellow' :
                s === Status.Connected ? 'green' : 'white';
        this.shadowRoot.innerHTML = `
<style>
#outer {
    display: flex;
    flex-direction: row;
    align-items: center;
}
</style>
<div id="outer">
    <svg width="24" height="18" xmlns="http://www.w3.org/2000/svg">
    <circle cx="12" cy="9" r="8" stroke="black" fill="${color}" />
    </svg>
    <div>${s}</div>
</div>
`;
    }
}
customElements.define('core-status-busconn-unused', BusConnection);

export { Controller, StatusContainer };
