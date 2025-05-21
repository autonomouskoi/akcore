import { UpdatingControlPanel } from "./tk.js";
import * as intcfgpb from "/pb/svc/pb/config_pb.js";
import * as logpb from "/pb/svc/pb/log_pb.js";
import { InternalConfig } from "./cfg_control.js";

let help = document.createElement('div');
help.innerHTML = `
<p>Configure the level of logging AK does. <code>Debug</code> is the most detailed, <code>Error</code>
is the least. More detailed logs provide more detail for troubleshooting but may use more disk space.</p>
<p>On shutdown AK automatically deletes log files older than 30 days.</p>
`;

class Logging extends UpdatingControlPanel<intcfgpb.Config> {
    private _select: HTMLSelectElement;

    constructor(cfg: InternalConfig) {
        super({ title: 'Logging', help, data: cfg });

        this.innerHTML = `
<label for="level">Logging Level (requires restart):</label>
<select id="level">
    <option value="${logpb.LogLevel.ERROR}">Error</option>
    <option value="${logpb.LogLevel.WARN}">Warn</option>
    <option value="${logpb.LogLevel.INFO}">Info</option>
    <option value="${logpb.LogLevel.DEBUG}">Debug</option>
</select>
`;

        this._select = this.querySelector('select');
        this._select.addEventListener('change', () => this.onSelectChanged());
    }

    update(cfg: intcfgpb.Config) {
        this._select.value = (cfg.logLevel === undefined ? logpb.LogLevel.INFO : cfg.logLevel).toString();
    }

    private onSelectChanged() {
        let level = parseInt(this._select.value);
        if (level < 0 || level > 3) {
            return;
        }
        let cfg = this.last.clone();
        cfg.logLevel = level;
        this.save(cfg);
    }
}
customElements.define('ak-cfg-logging', Logging, { extends: 'fieldset' });

export { Logging };