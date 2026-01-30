import { calculerScoreMain } from './client.js';

export function majInterface(mains, etats, scores, auTourDe, donneur, arretes, monIndex) {
    const tableEl = document.getElementById('table');
    const statusEl = document.getElementById('status-general');
    tableEl.innerHTML = "";
    
    const maxScore = Math.max(...scores);

    mains.forEach((main, i) => {
        const estActif = (i === auTourDe);
        const estArrete = arretes.includes(i);
        const estGele = !etats[i] && main.length > 0 && !estArrete;
  
        const estElimine = !etats[i] && main.length === 0; 
        
        const estLeader = scores[i] === maxScore && maxScore > 0;
        const classeCible = (window.enAttenteDeCible && i !== monIndex && etats[i]) ? 'cible-cliquable' : '';
        
        let div = document.createElement('div');
        div.className = `joueur-box ${estActif ? 'actif' : ''} ${estGele ? 'gele' : ''} ${estArrete ? 'stop' : ''} ${estElimine ? 'elimine' : ''} ${classeCible}`;
        div.onclick = () => { if (window.enAttenteDeCible) window.callbackCible(i); };

        div.innerHTML = `
            <div>
                ${estLeader ? '<span class="crown">👑</span>' : ''}
                <strong>Joueur ${i + 1} ${i === monIndex ? '(VOUS)' : ''}</strong> 
                ${i === donneur ? "⭐" : ""}
            </div>
            <div class="score-total">Total: ${scores[i]}</div>
            <div class="score-tour">Tour: ${calculerScoreMain(main)}</div>
            <div class="main-cartes">
                ${main.map(c => `<div class="carte ${typeof c === 'number' ? '' : 'speciale'}">${c}</div>`).join('')}
            </div>
        `;
        tableEl.appendChild(div); // ajoute le noeud div à la table
    });

    if (window.enAttenteDeCible) {
        statusEl.innerText = "🎯 CLIQUEZ SUR UNE VICTIME";
    } else {
        statusEl.innerText = (auTourDe === monIndex) ? "👉 À VOUS !" : `Attente du Joueur ${auTourDe + 1}...`;
    }
}

export function afficherPodium(scores, estVictoireFinale = false, indexGagnant = -1) {
    const overlay = document.getElementById('overlay');
    const titre = document.getElementById('recap-titre');
    const liste = document.getElementById('podium-liste');
    const btn = document.getElementById('btn-continuer');

    titre.innerText = estVictoireFinale ? `🏆 VICTOIRE FINALE : JOUEUR ${indexGagnant + 1} !` : "Fin de la manche";
    btn.innerText = estVictoireFinale ? "Recommencer une partie" : "Manche Suivante";

    liste.innerHTML = scores.map((s, i) => `
        <div class="podium-line">
            <span>Joueur ${i+1}</span>
            <span>${s} pts</span>
        </div>`).join('');
    
    overlay.style.display = "flex";

}
