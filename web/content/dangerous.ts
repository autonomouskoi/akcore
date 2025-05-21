import { UpdatingControlPanel } from "./tk.js";
import * as intcfgpb from "/pb/svc/pb/config_pb.js";
import { InternalConfig } from "./cfg_control.js";

let help = document.createElement('div');
help.innerHTML = `
<p><em>Network Accessible</em> allows AK to be accessed over the local network. This allows AK to
to be controlled and serve content such as OBS overlays over the local network. It also allows
anyone on the local network to control AK.</p>
`;

class Dangerous extends UpdatingControlPanel<intcfgpb.Config> {
    private _check: HTMLInputElement;

    constructor(cfg: InternalConfig) {
        super({ title: 'Dangerous Settings', help, data: cfg });

        this.innerHTML = `
<label for="network-accessible">Network Accessible (requires restart)</label>        
<input type="checkbox" id="network-accessible" />
`;

        this._check = this.querySelector('input');
        this._check.addEventListener('change', () => this.onListenChanged());
    }

    update(cfg: intcfgpb.Config) {
        this._check.checked = cfg.listenAddress === '0.0.0.0:8011';
    }

    private onListenChanged() {
        let cfg = this.last.clone()
        cfg.listenAddress = this._check.checked ? '0.0.0.0:8011' : '';
        this.save(cfg);
    }
}
customElements.define('ak-cfg-dangers', Dangerous, { extends: 'fieldset' });

export { Dangerous };