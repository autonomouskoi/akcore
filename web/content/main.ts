import { Controller } from "./controller.js";
import { Status } from "./status.js";
import { ModuleStates } from "./module_state.js";

function start(container: HTMLElement) {
    container.style.setProperty('display', 'flex');
    container.style.setProperty('flex-direction', 'column');

    let h1 = document.createElement('h1');
    h1.innerText = 'AutonomousKoi';

    let ctrl = new Controller();

    container.appendChild(h1);
    container.appendChild(new Status(ctrl));
    container.appendChild(new ModuleStates());
}

export { start }