import { Controller } from "./cfg_control.js";
import { StatusContainer } from "./status.js";
import { ModuleStates } from "./module_states.js";

function start(container: HTMLElement) {
    container.style.setProperty('display', 'flex');
    container.style.setProperty('flex-direction', 'column');

    let ctrl = new Controller();

    container.appendChild(new StatusContainer(ctrl));
    container.appendChild(new ModuleStates());
}

export { start }