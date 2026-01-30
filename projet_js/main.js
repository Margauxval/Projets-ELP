import { creerPaquet } from './deck.js';
import { majInterface, afficherPodium } from './int_graphique.js';

let nbJoueurs = parseInt(prompt("Combien de joueurs ?", "2")) || 2;
let scoresGlobaux = new Array(nbJoueurs).fill(0);
let mainsJoueurs, etatsJoueurs, paquet, donneurIndex, joueurQuiParle;
let joueursArretesVolontairement = []; // Pour différencier du gel

window.enAttenteDeCible = false;
window.callbackCible = null;

donneurIndex = Math.floor(Math.random() * nbJoueurs);

export function calculerScoreMain(main) {
    let nombres = main.filter(c => typeof c === 'number');
    let bonusFixes = main.filter(c => typeof c === 'string' && c.startsWith('+'));
    let aMultiplicateur = main.includes('x2');
    let totalNombres = nombres.reduce((a, b) => a + b, 0);
    let scoreBase = aMultiplicateur ? totalNombres * 2 : totalNombres;
    let totalBonus = bonusFixes.reduce((acc, b) => acc + parseInt(b.replace('+', '')), 0);
    return scoreBase + totalBonus;
}

function nouvelleManche() {
    paquet = creerPaquet();
    mainsJoueurs = Array.from({ length: nbJoueurs }, () => []);
    etatsJoueurs = new Array(nbJoueurs).fill(true);
    joueursArretesVolontairement = [];
    joueurQuiParle = (donneurIndex + 1) % nbJoueurs;
    
    const cartesInterdites = ["Freeze", "Flip 3", "Second Chance"];
    mainsJoueurs.forEach((m) => {
        let indexValide = paquet.findLastIndex(c => !cartesInterdites.includes(c));
        let carteInitiale = paquet.splice(indexValide, 1)[0];
        m.push(carteInitiale);
    });

    document.getElementById('overlay').style.display = "none";
    majInterface(mainsJoueurs, etatsJoueurs, scoresGlobaux, joueurQuiParle, donneurIndex, joueursArretesVolontairement);
}

function choisirCible(message) {
    return new Promise((resolve) => {
        window.enAttenteDeCible = true;
        majInterface(mainsJoueurs, etatsJoueurs, scoresGlobaux, joueurQuiParle, donneurIndex, joueursArretesVolontairement);
        alert(message);
        window.callbackCible = (indexCible) => {
            if (indexCible !== joueurQuiParle && etatsJoueurs[indexCible]) {
                window.enAttenteDeCible = false;
                resolve(indexCible);
            }
        };
    });
}

async function piocherCarte() {
    let c = paquet.pop();
    let mainActuelle = mainsJoueurs[joueurQuiParle];
    const adversairesActifs = etatsJoueurs.filter((e, i) => e && i !== joueurQuiParle).length;

    if (c === "Freeze" && adversairesActifs > 0) {
        mainActuelle.push(c);
        const cible = await choisirCible("FREEZE : Choisissez un adversaire à bloquer !");
        alert(`Le Joueur ${cible + 1} est gelé ❄️ !`);
        scoresGlobaux[cible] += calculerScoreMain(mainsJoueurs[cible]);
        etatsJoueurs[cible] = false; 
    } 
    else if (c === "Flip 3" && adversairesActifs > 0) {
        mainActuelle.push(c);
        const cible = await choisirCible("FLIP 3 : Choisissez la victime !");
        
        for (let i = 1; i <= 3; i++) {
            if (!etatsJoueurs[cible]) break; // Si déjà éliminé au tirage 1 ou 2
            
            let extra = paquet.pop();
            alert(`FLIP 3 (Carte ${i}/3) pour Joueur ${cible+1} : ${extra}`);
            
            // Vérification Doublon sur la victime
            if (typeof extra === 'number' && mainsJoueurs[cible].includes(extra)) {
                if (mainsJoueurs[cible].includes("Second Chance")) {
                    alert("Protégé par Second Chance !");
                    mainsJoueurs[cible].splice(mainsJoueurs[cible].indexOf("Second Chance"), 1);
                    mainsJoueurs[cible].push(extra);
                } else {
                    alert(`DOUBLON ! Joueur ${cible + 1} éliminé.`);
                    mainsJoueurs[cible] = [];
                    etatsJoueurs[cible] = false;
                }
            } else {
                mainsJoueurs[cible].push(extra);
            }
        }
    }
    else if (typeof c === 'number' && mainActuelle.includes(c)) {
        if (mainActuelle.includes("Second Chance")) {
            alert(`Doublon ${c} évité !`);
            mainActuelle.splice(mainActuelle.indexOf("Second Chance"), 1);
        } else {
            alert(`DOUBLON ${c} ! Joueur ${joueurQuiParle + 1} éliminé.`);
            mainsJoueurs[joueurQuiParle] = [];
            etatsJoueurs[joueurQuiParle] = false;
        }
    } 
    else {
        mainActuelle.push(c);
    }

    if (etatsJoueurs[joueurQuiParle] && mainActuelle.filter(card => typeof card === 'number').length >= 7) {
        alert("✨ FLIP 7 ! ✨");
        scoresGlobaux[joueurQuiParle] += (calculerScoreMain(mainActuelle) + 15);
        etatsJoueurs.fill(false);
    }
    
    finDeTour();
}

function finDeTour() {
    if (!etatsJoueurs.includes(true)) {
        afficherPodium(scoresGlobaux);
    } else {
        do {
            joueurQuiParle = (joueurQuiParle + 1) % nbJoueurs;
        } while (!etatsJoueurs[joueurQuiParle] && etatsJoueurs.includes(true));
        
        if (etatsJoueurs.includes(true)) {
            majInterface(mainsJoueurs, etatsJoueurs, scoresGlobaux, joueurQuiParle, donneurIndex, joueursArretesVolontairement);
        } else {
            afficherPodium(scoresGlobaux);
        }
    }
}

document.getElementById('btn-tirer').onclick = async () => {
    document.getElementById('btn-tirer').disabled = true;
    await piocherCarte();
    document.getElementById('btn-tirer').disabled = false;
};

document.getElementById('btn-stop').onclick = () => {
    scoresGlobaux[joueurQuiParle] += calculerScoreMain(mainsJoueurs[joueurQuiParle]);
    etatsJoueurs[joueurQuiParle] = false;
    joueursArretesVolontairement.push(joueurQuiParle);
    finDeTour();
};

document.getElementById('btn-continuer').onclick = () => {
    if (scoresGlobaux.some(s => s >= 200)) location.reload();
    else { donneurIndex = (donneurIndex + 1) % nbJoueurs; nouvelleManche(); }
};

nouvelleManche();