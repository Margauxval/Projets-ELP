import { creerPaquet } from './deck.js';
import { majInterface, afficherPodium } from './int_graphique.js';

let nbJoueurs = parseInt(prompt("Combien de joueurs ?", "2")) || 2;
let scoresGlobaux = new Array(nbJoueurs).fill(0);
let mainsJoueurs, etatsJoueurs, paquet, donneurIndex, joueurQuiParle;

donneurIndex = Math.floor(Math.random() * nbJoueurs);

function nouvelleManche() {
    paquet = creerPaquet();
    mainsJoueurs = Array.from({ length: nbJoueurs }, () => []);
    etatsJoueurs = new Array(nbJoueurs).fill(true);
    joueurQuiParle = (donneurIndex + 1) % nbJoueurs;
    
    mainsJoueurs.forEach(m => m.push(paquet.pop()));
    document.getElementById('overlay').style.display = "none";
    majInterface(mainsJoueurs, etatsJoueurs, scoresGlobaux, joueurQuiParle, donneurIndex);
}

document.getElementById('btn-tirer').onclick = () => {
    let c = paquet.pop();
    if (mainsJoueurs[joueurQuiParle].includes(c)) {
        alert(`DOUBLON ${c} ! J${joueurQuiParle + 1} éliminé.`);
        mainsJoueurs[joueurQuiParle] = [];
        etatsJoueurs[joueurQuiParle] = false;
        finDeTour();
    } else {
        mainsJoueurs[joueurQuiParle].push(c);
        if (mainsJoueurs[joueurQuiParle].length === 7) {
            scoresGlobaux[joueurQuiParle] += (mainsJoueurs[joueurQuiParle].reduce((a, b) => a + b, 0) + 15);
            afficherPodium(scoresGlobaux);
        } else { finDeTour(); }
    }
};

document.getElementById('btn-stop').onclick = () => {
    scoresGlobaux[joueurQuiParle] += mainsJoueurs[joueurQuiParle].reduce((a, b) => a + b, 0);
    etatsJoueurs[joueurQuiParle] = false;
    finDeTour();
};

function finDeTour() {
    if (!etatsJoueurs.includes(true)) {
        afficherPodium(scoresGlobaux);
    } else {
        do { joueurQuiParle = (joueurQuiParle + 1) % nbJoueurs; } while (!etatsJoueurs[joueurQuiParle]);
        majInterface(mainsJoueurs, etatsJoueurs, scoresGlobaux, joueurQuiParle, donneurIndex);
    }
}

document.getElementById('btn-continuer').onclick = () => {
    if (scoresGlobaux.some(s => s >= 200)) location.reload();
    else {
        donneurIndex = (donneurIndex + 1) % nbJoueurs;
        nouvelleManche();
    }
};

nouvelleManche();