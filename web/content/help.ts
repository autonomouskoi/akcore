function SectionHelp(toggleNode: Node, detailsHTML: HTMLElement): HTMLElement {
    document.createTextNode
    let detailsDiv = document.createElement('div');
    detailsDiv.setAttribute('style', `
    display: none;
    padding: 0 0.5rem 0 0.5rem;
    border: solid 1px gray;
    border-radius: 5px;
`);
    let display = false;
    toggleNode.addEventListener('click', () => {
        display = !display;
        detailsDiv.style.setProperty('display', display ? 'block' : 'none');
    });

    detailsDiv.appendChild(detailsHTML);
    return detailsDiv;
}
export { SectionHelp };