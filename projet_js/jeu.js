// --- CONFIGURATION ---
let nbJoueurs = parseInt(prompt("Combien de joueurs ?", "2"));
let scoresGlobaux = new Array(nbJoueurs).fill(0);
let etatsJoueurs = []; 
let mainsJoueurs = []; 

// Choix aléatoire du premier donneur [cite: 62]
let donneurIndex = Math.floor(Math.random() * nbJoueurs); 
let joueurQuiParle = (donneurIndex + 1) % nbJoueurs; 

let paquet = [];

function creerPaquet() {
    // Composition du paquet selon les règles [cite: 13, 40]
    paquet = [0]; // Une carte 0 qui vaut 0 point [cite: 40, 51]
    for (let i = 1; i <= 12; i++) {
        for (let j = 0; j < i; j++) {
            paquet.push(i); // Un 1, deux 2... douze 12 [cite: 13, 16]
        }
    }
    // Mélange de Fisher-Yates [cite: 62]
    for (let i = paquet.length - 1; i > 0; i--) {
        let j = Math.floor(Math.random() * (i + 1));
        [paquet[i], paquet[j]] = [paquet[j], paquet[i]];
    }
}

function initialiserManche() {
    creerPaquet();
    mainsJoueurs = Array.from({ length: nbJoueurs }, () => []);
    etatsJoueurs = new Array(nbJoueurs).fill(true); 
    
    // Affichage du donneur (on ajoute +1 pour l'humain)
    console.log(`%c LE DONNEUR EST LE JOUEUR ${donneurIndex + 1} `, "background: #222; color: #bada55");
    
    // Distribution initiale : une carte face visible à chaque joueur [cite: 63]
    for (let i = 0; i < nbJoueurs; i++) {
        mainsJoueurs[i].push(paquet.pop());
    }
    
    console.log("Distribution terminée. Début des tours de table.");
    afficherTableau();
}

function afficherTableau() {
    console.log("--- TABLEAU DES MAINS ---");
    mainsJoueurs.forEach((main, i) => {
        let statut = etatsJoueurs[i] ? "ACTIF" : "HORS-JEU";
        // On affiche i + 1 pour ne pas commencer à 0
        console.log(`Joueur ${i + 1} (${statut}): [${main.join(", ")}] | Somme: ${main.reduce((a, b) => a + b, 0)}`);
    });
    console.log(`C'est au tour du Joueur ${joueurQuiParle + 1}. (tirer() ou stop())`);
}

function tirer() {
    if (!etatsJoueurs[joueurQuiParle]) return;

    let carte = paquet.pop();
    console.log(`Joueur ${joueurQuiParle + 1} tire un ${carte}.`);

    // Si doublon, éliminé du tour et 0 point [cite: 11]
    if (mainsJoueurs[joueurQuiParle].includes(carte)) {
        console.log(`DOUBLON ! Joueur ${joueurQuiParle + 1} est éliminé.`);
        mainsJoueurs[joueurQuiParle] = []; 
        etatsJoueurs[joueurQuiParle] = false;
    } else {
        mainsJoueurs[joueurQuiParle].push(carte);
        
        // Bonus Flip 7 : 7 cartes différentes arrêtent le tour [cite: 10, 126]
        if (mainsJoueurs[joueurQuiParle].length === 7) {
            console.log("!!! FLIP 7 !!! Le tour s'arrête immédiatement !");
            let bonus = 15; // Bonus de 15 points [cite: 10, 138]
            scoresGlobaux[joueurQuiParle] += (mainsJoueurs[joueurQuiParle].reduce((a, b) => a + b, 0) + bonus);
            finDeMancheImmediate();
            return;
        }
    }
    passerAuSuivant();
}

function stop() {
    if (!etatsJoueurs[joueurQuiParle]) return;
    
    let points = mainsJoueurs[joueurQuiParle].reduce((a, b) => a + b, 0);
    scoresGlobaux[joueurQuiParle] += points; 
    console.log(`Joueur ${joueurQuiParle + 1} s'arrête avec ${points} pts.`);
    etatsJoueurs[joueurQuiParle] = false;
    passerAuSuivant();
}

function passerAuSuivant() {
    let joueursActifs = etatsJoueurs.filter(a => a === true).length;
    
    // Si tout le monde a fini [cite: 124]
    if (joueursActifs <= 0) {
        finDeMancheImmediate();
        return;
    }

    // Un seul tirage par joueur, puis passage au suivant 
    do {
        joueurQuiParle = (joueurQuiParle + 1) % nbJoueurs;
    } while (!etatsJoueurs[joueurQuiParle]);

    afficherTableau();
}

function finDeMancheImmediate() {
    console.log("--- RÉSULTATS DE LA MANCHE ---");
    
    // Création d'un objet lisible pour console.table sans l'index 0
    let affichageScores = {};
    scoresGlobaux.forEach((score, i) => {
        affichageScores[`Joueur ${i + 1}`] = { "Score Total": score };
    });
    console.table(affichageScores);
    
    // Fin de partie à 200 points [cite: 7, 150]
    if (scoresGlobaux.some(s => s >= 200)) {
        let gagnant = scoresGlobaux.indexOf(Math.max(...scoresGlobaux)) + 1;
        console.log(`%c LE JOUEUR ${gagnant} A GAGNÉ LA PARTIE ! `, "background: #gold; color: black; font-size: 20px");
    } else {
        // Le donneur passe à gauche pour le tour suivant 
        donneurIndex = (donneurIndex + 1) % nbJoueurs;
        joueurQuiParle = (donneurIndex + 1) % nbJoueurs;
        console.log("Tapez 'initialiserManche()' pour continuer.");
    }
}

initialiserManche();