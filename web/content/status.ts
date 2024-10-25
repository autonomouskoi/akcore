import { GloballyStyledHTMLElement } from "./global-styles.js";
import { Controller } from "./controller.js";

class Status extends GloballyStyledHTMLElement {
    constructor(ctrl: Controller) {
        super();

        this.shadowRoot.innerHTML = `
<style>
div {
    display: flex;
    padding: 5px;
    justify-content: flex-end;
}
</style>
<div></div>
`

        let div = this.shadowRoot.querySelector('div');

        ctrl.ready().then(() => {
            div.appendChild(new Listen(ctrl));
        })
    }
}
customElements.define('core-status-main-unused', Status);

class Listen extends HTMLElement {
    constructor(ctrl: Controller) { 
        super();

        this.innerHTML = `
<label for="check"
        title="Allow others on the local network to control AK. Could be dangerous!"
>Network Accessible (requires restart)</label>
<input type="checkbox" id="check" />
`;
        let check = this.querySelector('input') as HTMLInputElement;
        check.checked = ctrl.listenAddress.last;
        ctrl.listenAddress.subscribe((v) => {
            check.checked = v;
        })
        check.addEventListener('change', () => {
            ctrl.listenAddress.save(check.checked);
        });
    }
}
customElements.define('core-status-listen-unused', Listen);

export { Status };
