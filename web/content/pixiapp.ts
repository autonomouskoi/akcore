import * as PIXI from "pixi.js";

const app = new PIXI.Application();
await app.init({
    width: 1280,
    height: 720,
    backgroundAlpha: 0,
 });
document.body.appendChild(app.canvas);

export { app };