import { majInterface, afficherPodium } from './int_graphique.js';

const socket = io();  // ouvrir une co
let monIndex = null;
window.enAttenteDeCible = false;

// Logique de calcul du score (doit être identique au serveur)
export function calculerScoreMain(main) {
    let nombres = main.filter(c => typeof c === 'number'); //tableau de nombres uniquement
    let totalNombres = nombres.reduce((a, b) => a + b, 0); // additionne a et b, renvoie a et b partent de 0
    
    // Application du multiplicateur x2
    let scoreBase = main.includes('x2') ? totalNombres * 2 : totalNombres; 
    
    // Ajout des bonus fixes (+2, +4, etc.)
    let totalBonus = main.filter(c => typeof c === 'string' && c.startsWith('+'))
                         .reduce((acc, b) => acc + parseInt(b.replace('+', '')), 0);
    
    return scoreBase + totalBonus;
}

// Initialisation : on reçoit notre numéro de joueur
socket.on('init', (index) => {
    monIndex = index;
    console.log("Connecté en tant que joueur", index);
});

// Mise à jour globale de l'état du jeu
socket.on('update', (game) => {
    document.getElementById('overlay').style.display = "none";
    
    // On passe monIndex pour que l'interface sache qui est "VOUS"
    majInterface(game.mains, game.etats, game.scores, game.auTourDe, game.donneur, game.joueursArretes, monIndex);
    
    // Activation des boutons seulement si c'est mon tour et que je ne suis pas en train de cibler
    const estMonTour = (game.auTourDe === monIndex && game.etats[monIndex] && !window.enAttenteDeCible);
    document.getElementById('btn-tirer').disabled = !estMonTour;
    document.getElementById('btn-stop').disabled = !estMonTour;
});

// Gestion des cartes spéciales d'attaque
socket.on('demanderCible', (typeCarte) => {
    window.enAttenteDeCible = true;
    alert(`✨ ${typeCarte.toUpperCase()} ! Cliquez sur la boîte d'un adversaire pour l'attaquer.`);
    // On rafraîchit l'interface pour faire briller les cibles cliquables
    socket.emit('demanderRefresh'); 
});

// Fonction appelée quand on clique sur un joueur (via int_graphique.js)
window.callbackCible = (indexCible) => {
    if (window.enAttenteDeCible && indexCible !== monIndex) {
        window.enAttenteDeCible = false;
        socket.emit('cibleChoisie', indexCible);
    }
};

// Réception de messages (ex: résultat du Flip 3, Second Chance)
socket.on('message', (msg) => {
    alert(msg);
});

// Fin de manche classique
socket.on('finManche', (scores) => {
    afficherPodium(scores, false);
});

// Victoire finale (200 points)
socket.on('victoireFinale', (data) => {
    afficherPodium(data.scores, true, data.index);
});

// --- Événements Boutons ---
document.getElementById('btn-tirer').onclick = () => socket.emit('tirer');

document.getElementById('btn-stop').onclick = () => socket.emit('stop');

document.getElementById('btn-continuer').onclick = () => {
    const titre = document.getElementById('recap-titre').innerText;
    if (titre.includes("VICTOIRE FINALE")) {
        location.reload(); // Recommencer de zéro
    } else {
        socket.emit('continuer'); // Manche suivante
    }
};

// Touche "S" pour démarrer la partie au tout début
window.addEventListener('keydown', (e) => {
    if (e.key.toLowerCase() === 's') {
        socket.emit('demarrerPartie');
    }

});
