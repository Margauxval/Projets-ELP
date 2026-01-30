const express = require('express');
const http = require('http');
const { Server } = require('socket.io');

const app = express();
const server = http.createServer(app);
const io = new Server(server);

app.use(express.static(__dirname));

let game = {
    mains: [], etats: [], scores: [], joueursArretes: [],
    auTourDe: 0, donneur: 0, paquet: [], sockets: [],
    enAttenteCible: false, carteEnAttente: null
};

function creerPaquet() {
    let p = [0];
    for (let i = 1; i <= 12; i++) { for (let j = 0; j < i; j++) p.push(i); }
    ["+2", "+4", "+6", "+8", "+10", "x2"].forEach(m => p.push(m));
    ["Second Chance", "Flip 3", "Freeze"].forEach(a => { for (let k = 0; k < 3; k++) p.push(a); });
    return p.sort(() => Math.random() - 0.5);
}

function calculerScore(main) {
    let nombres = main.filter(c => typeof c === 'number');
    let bonus = main.filter(c => typeof c === 'string' && c.startsWith('+'));
    let total = nombres.reduce((a, b) => a + b, 0);
    if (main.includes('x2')) total *= 2;
    return total + bonus.reduce((acc, b) => acc + parseInt(b), 0);
}

function passerTour() {
    if (!game.etats.includes(true)) {
        io.emit('finManche', game.scores);
    } else {
        do {
            game.auTourDe = (game.auTourDe + 1) % game.sockets.length;
        } while (!game.etats[game.auTourDe] && game.etats.includes(true));
        io.emit('update', game);
    }
}

io.on('connection', (socket) => {
    if (game.sockets.length >= 4) return socket.disconnect();
    game.sockets.push(socket.id);
    socket.emit('init', game.sockets.length - 1);

    socket.on('demarrerPartie', () => {
        game.paquet = creerPaquet();
        game.mains = game.sockets.map(() => []);
        game.etats = game.sockets.map(() => true);
        game.scores = game.scores.length ? game.scores : game.sockets.map(() => 0);
        game.joueursArretes = [];
        
        game.mains.forEach(m => {
            let idx = game.paquet.findLastIndex(c => typeof c === 'number');
            m.push(game.paquet.splice(idx, 1)[0]);
        });
        io.emit('update', game);
    });

   socket.on('tirer', () => {
        if (socket.id !== game.sockets[game.auTourDe] || game.enAttenteCible) return;
        let c = game.paquet.pop();
        let main = game.mains[game.auTourDe];

        if (c === "Freeze" || c === "Flip 3") {
            game.enAttenteCible = true;
            game.carteEnAttente = c;
            socket.emit('demanderCible', c);
            io.emit('update', game); 
        } else if (typeof c === 'number' && main.includes(c)) {
            if (main.includes("Second Chance")) {
                main.splice(main.indexOf("Second Chance"), 1);
                main.push(c);
                io.emit('message', `Sauvé par Second Chance !`);
                passerTour(); // Ne pas oublier de passer le tour après sauvetage
            } else {
                game.mains[game.auTourDe] = []; // Main vide = doublon
                game.etats[game.auTourDe] = false;
                
                // --- AJOUT : Vérification si tout le monde a perdu ---
                if (!game.etats.includes(true)) {
                    io.emit('update', game); // Update pour montrer la couleur rouge
                    setTimeout(() => verifierFinDePartie(), 1000);
                } else {
                    passerTour();
                }
            }
        } else {
            main.push(c);
            if (main.filter(x => typeof x === 'number').length >= 7) {
                game.scores[game.auTourDe] += (calculerScore(main) + 15);
                game.etats.fill(false);
                verifierFinDePartie();
                return;
            }
            passerTour();
        }
        io.emit('update', game);
    });

    socket.on('cibleChoisie', (indexCible) => {
        if (socket.id !== game.sockets[game.auTourDe]) return;
        
        let typeAttaque = game.carteEnAttente;
        // On ajoute la carte d'attaque (Flip 3 ou Freeze) à la main de l'attaquant
        game.mains[game.auTourDe].push(typeAttaque);

        if (typeAttaque === "Freeze") {
            game.scores[indexCible] += calculerScore(game.mains[indexCible]);
            game.etats[indexCible] = false;
        } 
        else if (typeAttaque === "Flip 3") {
            let tirees = [];
            for (let i = 0; i < 3; i++) {
                let extra = game.paquet.pop();
                tirees.push(extra);
                
                // On ajoute la carte à la main de la victime IMMEDIATEMENT pour l'affichage
                if (extra !== "Freeze") {
                    game.mains[indexCible].push(extra);
                }

                // Logique des effets
                if (extra === "Freeze") {
                    game.mains[indexCible].push(extra); // On l'ajoute quand même pour le style
                    game.scores[indexCible] += calculerScore(game.mains[indexCible]);
                    game.etats[indexCible] = false;
                    break; 
                } 
                else if (typeof extra === 'number' && checkDoublon(game.mains[indexCible], extra)) {
                    if (game.mains[indexCible].includes("Second Chance")) {
                        // Consomme la seconde chance
                        game.mains[indexCible].splice(game.mains[indexCible].indexOf("Second Chance"), 1);
                        io.emit('message', `Joueur ${indexCible + 1} sauvé par Second Chance pendant le Flip 3 !`);
                    } else {
                        // Doublon fatal : on vide la main
                        game.mains[indexCible] = []; 
                        game.etats[indexCible] = false; 
                        break;
                    }
                }
                // Optionnel : émettre un update ici si tu veux voir les cartes apparaître une par une
                // io.emit('update', game); 
            }
            io.emit('message', `Flip 3 terminé sur Joueur ${indexCible + 1} ! Cartes tirées : ${tirees.join(', ')}`);
        }

        game.enAttenteCible = false;
        game.carteEnAttente = null;

        // Vérification si tout le monde est hors-jeu
        if (!game.etats.includes(true)) {
            io.emit('update', game);
            setTimeout(() => verifierFinDePartie(), 1200);
        } else {
            passerTour();
        }
    });

    function checkDoublon(main, nouvelleCarte) {
        return main.filter(c => c === nouvelleCarte).length > 1;
    }

    function verifierFinDePartie() {
        const gagnant = game.scores.findIndex(s => s >= 200);
        if (gagnant !== -1) {
            io.emit('victoireFinale', { index: gagnant, scores: game.scores });
            game.scores = []; // Reset pour la prochaine
        } else if (!game.etats.includes(true)) {
            io.emit('finManche', game.scores);
        }
    }

    socket.on('stop', () => {
        if (socket.id !== game.sockets[game.auTourDe]) return;
        game.scores[game.auTourDe] += calculerScore(game.mains[game.auTourDe]);
        game.etats[game.auTourDe] = false;
        game.joueursArretes.push(game.auTourDe);
        passerTour();
    });

    socket.on('continuer', () => {
        game.donneur = (game.donneur + 1) % game.sockets.length;
        socket.emit('demarrerPartie');
    });
});

server.listen(3000, '0.0.0.0', () => console.log("Serveur Flip 7 pret sur port 3000"));