import { calculerScoreMain } from './main.js';

export function majInterface(mains, etats, scores, auTourDe, donneur, arretes) {
    const tableEl = document.getElementById('table');
    const statusEl = document.getElementById('status-general');
    tableEl.innerHTML = "";
    
    const maxScore = Math.max(...scores);

    mains.forEach((main, i) => {
        const estActif = (i === auTourDe);
        const estElimine = !etats[i] && main.length === 0;
        const estArrete = arretes.includes(i);
        const estGele = !etats[i] && main.length > 0 && !estArrete;
        
        const sommeTour = calculerScoreMain(main);
        const iconeLeader = (scores[i] === maxScore && maxScore > 0) ? "👑" : "";
        const classeCible = (window.enAttenteDeCible && i !== auTourDe && etats[i]) ? 'cible-cliquable' : '';
        
        let div = document.createElement('div');
        div.className = `joueur-box ${estActif ? 'actif' : ''} ${estElimine ? 'elimine' : ''} ${estGele ? 'gele' : ''} ${estArrete ? 'stop' : ''} ${classeCible}`;
        
        div.onclick = () => { if (window.enAttenteDeCible) window.callbackCible(i); };

        div.innerHTML = `
            <div>${iconeLeader} <strong>Joueur ${i + 1}</strong> ${i === donneur ? "⭐" : ""}</div>
            <div class="score-total">Total: ${scores[i]}</div>
            <div class="score-tour">Tour: ${sommeTour}</div>
            <div class="main-cartes">
                ${main.map(c => `<div class="carte ${typeof c === 'number' ? '' : 'speciale'}">${c}</div>`).join('')}
            </div>
            <div class="statut-badge">
                ${estGele ? '❄️ GELÉ' : (estArrete ? '🛑 ARRÊTÉ' : '')}
            </div>
        `;
        tableEl.appendChild(div);
    });

    statusEl.innerText = window.enAttenteDeCible ? "🎯 CIBLEZ UN ADVERSAIRE" : `Tour du Joueur ${auTourDe + 1}`;
}

export function afficherPodium(scores) {
    const overlay = document.getElementById('overlay');
    const titre = document.getElementById('recap-titre');
    const liste = document.getElementById('podium-liste');
    
    const max = Math.max(...scores);
    titre.innerText = max >= 200 ? `🏆 VICTOIRE DU JOUEUR ${scores.indexOf(max) + 1} !` : "Fin de la manche";
    
    liste.innerHTML = scores
        .map((s, i) => ({ id: i + 1, score: s }))
        .sort((a, b) => b.score - a.score)
        .map((j, idx) => `<div class="podium-line"><span>${idx === 0 ? '🥇' : ''} Joueur ${j.id}</span><span>${j.score} pts</span></div>`)
        .join('');

    overlay.style.display = "flex";
}