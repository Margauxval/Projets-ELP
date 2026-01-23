export function majInterface(mains, etats, scores, auTourDe, donneur) {
    const tableEl = document.getElementById('table');
    const statusEl = document.getElementById('status-general');
    tableEl.innerHTML = "";
    let maxS = Math.max(...scores);

    mains.forEach((main, i) => {
        const estActif = (i === auTourDe);
        const estElimine = !etats[i] && main.length === 0;
        const somme = main.reduce((a, b) => a + b, 0);
        
        let div = document.createElement('div');
        div.className = `joueur-box ${estActif ? 'actif' : ''} ${estElimine ? 'elimine' : ''}`;
        
        let iconeLeader = (scores[i] === maxS && maxS > 0) ? "👑" : "";
        let iconeDonneur = (i === donneur) ? "⭐" : "";

        div.innerHTML = `
            <div>${iconeLeader} <strong>J${i + 1}</strong> ${iconeDonneur}</div>
            <div style="color:#ffeb3b">Total: ${scores[i]}</div>
            <div>Tour: ${somme}</div>
            <div class="main-cartes">
                ${main.map(c => `<div class="carte">${c}</div>`).join('')}
            </div>
        `;
        tableEl.appendChild(div);
    });
    statusEl.innerText = `Tour du Joueur ${auTourDe + 1}`;
}

export function afficherPodium(scores) {
    const overlay = document.getElementById('overlay');
    const liste = document.getElementById('podium-liste');
    let classement = scores.map((s, i) => ({ id: i + 1, score: s })).sort((a, b) => b.score - a.score);
    
    liste.innerHTML = "";
    classement.forEach((j, idx) => {
        let badge = (idx === 0) ? "🥇" : (idx === 1 ? "🥈" : (idx === 2 ? "🥉" : ""));
        liste.innerHTML += `<div class="podium-line"><span>${badge} J${j.id}</span><span>${j.score} pts</span></div>`;
    });
    overlay.style.display = "flex";
}