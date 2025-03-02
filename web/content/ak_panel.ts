import { GloballyStyledHTMLElement } from "./global-styles.js";
import * as status from "./status.js";
import { ControlPanel } from "./tk.js";
import { bus, Status } from "/bus.js";

class AKPanel extends GloballyStyledHTMLElement {
    constructor(ctrl: status.Controller) {
        super();

        this.classList.add('flex-column');
        this.shadowRoot.appendChild(new status.StatusContainer(ctrl));
        this.shadowRoot.appendChild(new Dangerous(ctrl));
    }
}
customElements.define('ak-panel', AKPanel);

class AKPanelListItem extends HTMLDivElement {
    constructor() {
        super();

        this.classList.add('flex-row');

        this.innerHTML = `
<svg width="24" height="18" xmlns="http://www.w3.org/2000/svg"
    style="position: absolute";
>
    <circle cx="12" cy="9" r="8" stroke="black" fill="white" />
</svg>
<img src="/ak_logo.svg" width=48 height=48
    title="AutonomousKoi controls"
/>
<div class="module-name">AK Controls</div>
`;
        this._update(bus.getStatus());
        bus.addStatusListener((s) => this._update(s));
    }

    private _update(s: Status) {
        let color = s === Status.NotConnected ? 'red' :
            s === Status.Connecting ? 'yellow' :
                s === Status.Connected ? 'green' : 'white';
        let circle = this.querySelector('circle') as SVGCircleElement;
        circle.style.fill = color;
    }
}
customElements.define('ak-panel-list-item', AKPanelListItem, { extends: 'div' });

class Dangerous extends ControlPanel {
    constructor(ctrl: status.Controller) {
        let help = document.createElement('div');
        help.textContent = `Dangerous settings are dangerous`;

        super({ title: 'Dangerous Settings', help });

        let listen = new status.Listen(ctrl);
        this.appendChild(listen);
    }
}
customElements.define('ak-cfg-dangers', Dangerous, { extends: 'fieldset' });

export { AKPanel, AKPanelListItem };