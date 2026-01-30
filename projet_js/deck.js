export function creerPaquet() {
    let p = [0];
    for (let i = 1; i <= 12; i++) {
        for (let j = 0; j < i; j++) p.push(i);
    }
    const modifs = ["+2", "+4", "+6", "+8", "+10", "x2"];
    modifs.forEach(m => p.push(m));
    const actions = ["Second Chance", "Flip 3", "Freeze"];
    actions.forEach(a => { for (let k = 0; k < 3; k++) p.push(a); });
    return p.sort(() => Math.random() - 0.5);
}