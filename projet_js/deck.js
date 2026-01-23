export function creerPaquet() {
    let p = [0];
    for (let i = 1; i <= 12; i++) {
        for (let j = 0; j < i; j++) p.push(i);
    }
    // Mélange
    return p.sort(() => Math.random() - 0.5);
}