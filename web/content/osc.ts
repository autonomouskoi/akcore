import { UpdatingControlPanel } from "./tk.js";
import * as intcfgpb from "/pb/svc/pb/svc_config_pb.js";
import * as oscpb from "/pb/svc/pb/osc_pb.js";
import { InternalConfig } from "./cfg_control.js";

let help = document.createElement('div');
help.innerHTML = `
<p><em>OSC</em> allows control of other systems like real time
visualizations over the network. <em>Targets</em> are address and port combinations representing
resources AK plugins can send OSC messages to.</p>
`;

class OSC extends UpdatingControlPanel<intcfgpb.Config> {
    private _targets: HTMLDivElement;
    private _edit: OSCEdit;

    constructor(cfg: InternalConfig) {
        super({ title: 'Open Sound Control', help, data: cfg });

        this.innerHTML = `
<h3 style="margin: 2px">Targets <button>&#10133;</button></h3>
<div id="targets" class="grid grid-4-col"></div>
`;

        this._targets = this.querySelector('#targets');
        this._edit = new OSCEdit();
        this.appendChild(this._edit);

        let newButton = this.querySelector('button');
        newButton.addEventListener('click', () => {
            this._edit.edit(new oscpb.OSCTarget(), (newTarget) => this._saveTarget('', newTarget));
        });

        this._resetTargets();
    }

    private _resetTargets() {
        this._targets.textContent = '';
        this._targets.innerHTML = `
<div class="column-header">Name</div>
<div class="column-header">Address/Port</div>
<div class="column-header"></div>
<div class="column-header"></div>
`;
    }

    update(cfg: intcfgpb.Config) {
        this._resetTargets();
        if (!cfg.oscConfig || !cfg.oscConfig.targets) {
            return
        }
        let oscCfg: oscpb.OSCConfig = cfg.oscConfig;
        oscCfg.targets.forEach((target) => {
            let name = document.createElement('div');
            name.innerText = target.name;
            this._targets.appendChild(name);

            let addr = document.createElement('div');
            addr.innerText = `${target.address}:${target.port}`;
            this._targets.appendChild(addr);

            let edit = document.createElement('button');
            edit.innerHTML = '&#9998;';
            edit.addEventListener('click',
                () => this._edit.edit(target,
                    (newTarget) => this._saveTarget(target.name, newTarget)
                )
            )
            this._targets.appendChild(edit);

            let del = document.createElement('button');
            del.innerHTML = '&#128465;';
            del.addEventListener('click', () => {
                if (confirm(`Delete target ${target.name}?`)) {
                    this._deleteTarget(target.name);
                }
            });
            this._targets.appendChild(del);
        });
    }

    private _saveTarget(origName: string, newT: oscpb.OSCTarget) {
        let cfg = this.last.clone();
        let oscCfg: oscpb.OSCConfig = cfg.oscConfig ? cfg.oscConfig : new oscpb.OSCConfig();

        if (!origName) {
            oscCfg.targets.push(newT);
        } else {
            for (let i = 0; i < oscCfg.targets.length; i++) {
                if (oscCfg.targets[i].name === origName) {
                    oscCfg.targets[i] = newT;
                    break;
                }
            }
        }

        cfg.oscConfig = oscCfg;
        this.save(cfg);
    }

    private _deleteTarget(name: string) {
        let cfg = this.last.clone();
        let oscCfg: oscpb.OSCConfig = cfg.oscConfig ? cfg.oscConfig : new oscpb.OSCConfig();
        oscCfg.targets = oscCfg.targets.filter((target) => target.name !== name);
        cfg.oscConfig = oscCfg;
        this.save(cfg);
    }
}
customElements.define('ak-config-osc', OSC, { extends: 'fieldset' });

class OSCEdit extends HTMLDialogElement {
    private _nameInput: HTMLInputElement;
    private _addrInput: HTMLInputElement;
    private _portInput: HTMLInputElement;

    private _onSave = (newTarget: oscpb.OSCTarget) => { };

    constructor() {
        super();

        this.innerHTML = `
<h3>Edit Target</h3>
<section class="grid grid-2-col">

<label for="name" title="A label to help you recognize this target">Name</label>
<input type="text" id="name" title="A label to help you recognize this target"" />

<label for="address" title="IP or hostname of the target">Address</label>
<input type="text" id="address" title="IP or hostname of the target"/>

<label for="port" title="Port for the the target">Port</label>
<input type="number" id="port" min="1" max="65535"
    title="Port for the target"/>

<button id="save">Save</button>
<button id="cancel">Cancel</button>

</section>
`;
        this._nameInput = this.querySelector('input#name');
        this._addrInput = this.querySelector('input#address');
        this._portInput = this.querySelector('input#port');

        this.querySelector('button#save').addEventListener('click', () => this._save());
        this.querySelector('button#cancel').addEventListener('click', () => {
            this._onSave = () => { };
            this.close()
        });

        this.addEventListener('keyup', (event) => {
            if (event.key === "Enter") {
                event.preventDefault();
                this._save();
            }
        });
    }

    private _save() {
        // maybe some validation
        let newTarget = new oscpb.OSCTarget({
            name: this._nameInput.value,
            address: this._addrInput.value,
            port: parseInt(this._portInput.value),
        });
        this._onSave(newTarget);
        this.close();
    }

    edit(target: oscpb.OSCTarget, onSave: (target: oscpb.OSCTarget) => void) {
        this._nameInput.value = target.name;
        this._addrInput.value = target.address;
        this._portInput.value = target.port.toString();
        this._onSave = onSave;
        this.showModal();
    }
}
customElements.define('ak-config-osc-edit', OSCEdit, { extends: 'dialog' });

export { OSC };