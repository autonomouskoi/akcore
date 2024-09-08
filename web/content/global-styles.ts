// https://eisenbergeffect.medium.com/using-global-styles-in-shadow-dom-5b80e802e89d
let globalSheets: CSSStyleSheet[] = [];

export function getGlobalStyleSheets() {
    if (globalSheets.length === 0) {
        globalSheets = Array.from(document.styleSheets)
            .map(x => {
                const sheet = new CSSStyleSheet();
                const css = Array.from(x.cssRules).map(rule => rule.cssText).join(' ');
                sheet.replaceSync(css);
                return sheet;
            });
    }

    return globalSheets;
}

export function addGlobalStylesToShadowRoot(shadowRoot: ShadowRoot) {
    shadowRoot.adoptedStyleSheets.push(
        ...getGlobalStyleSheets()
    );
}

export class GloballyStyledHTMLElement extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        addGlobalStylesToShadowRoot(this.shadowRoot);
    }
}