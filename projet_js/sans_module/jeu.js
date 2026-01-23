document.addEventListener("DOMContentLoaded", () => {
    let nbJoueurs = parseInt(prompt("Combien de joueurs ?", "2")) || 2;
    let scoresGlobaux = new Array(nbJoueurs).fill(0);
    let mainsJoueurs = [], etatsJoueurs = [], paquet = [];
    
    let donneurIndex = Math.floor(Math.random() * nbJoueurs);
    let joueurQuiParle = (donneurIndex + 1) % nbJoueurs;

    const tableEl = document.getElementById('table');
    const statusEl = document.getElementById('status-general');
    const overlay = document.getElementById('overlay');

    function creerPaquet() {
        paquet = [0];
        for (let i = 1; i <= 12; i++) { for (let j = 0; j < i; j++) paquet.push(i); }
        for (let i = paquet.length - 1; i > 0; i--) {
            let j = Math.floor(Math.random() * (i + 1));
            [paquet[i], paquet[j]] = [paquet[j], paquet[i]];
        }
    }

    function initialiserManche() {
        creerPaquet();
        mainsJoueurs = Array.from({ length: nbJoueurs }, () => []);
        etatsJoueurs = new Array(nbJoueurs).fill(true);
        // Distribution initiale
        for (let i = 0; i < nbJoueurs; i++) mainsJoueurs[i].push(paquet.pop());
        overlay.style.display = "none";
        majInterface();
    }

    function majInterface() {
        tableEl.innerHTML = "";
        let maxScore = Math.max(...scoresGlobaux);

        mainsJoueurs.forEach((main, i) => {
            const estActif = (i === joueurQuiParle);
            const estElimine = !etatsJoueurs[i] && main.length === 0;
            const somme = main.reduce((a, b) => a + b, 0);
            
            let div = document.createElement('div');
            div.className = `joueur-box ${estActif ? 'actif' : ''} ${estElimine ? 'elimine' : ''}`;
            
            // Étoile pour le donneur, Couronne pour le premier au score
            let iconeDonneur = (i === donneurIndex) ? "⭐" : "";
            let iconeLeader = (scoresGlobaux[i] === maxScore && maxScore > 0) ? "<span class='crown'>👑</span>" : "";

            div.innerHTML = `
                <div>${iconeLeader} <strong>Joueur ${i + 1}</strong> ${iconeDonneur}</div>
                <div style="font-weight:bold; color:#ffeb3b">Total: ${scoresGlobaux[i]}</div>
                <div style="margin-top:5px">Somme Tour: ${somme}</div>
                <div class="main-cartes">
                    ${main.map(c => `<div class="carte">${c}</div>`).join('')}
                </div>
            `;
            tableEl.appendChild(div);
        });
        statusEl.innerText = `Tour du Joueur ${joueurQuiParle + 1}`;
    }

    document.getElementById('btn-tirer').onclick = () => {
        let carte = paquet.pop();
        if (mainsJoueurs[joueurQuiParle].includes(carte)) {
            alert(`DOUBLON DE ${carte} ! Le Joueur ${joueurQuiParle + 1} est éliminé !`);
            mainsJoueurs[joueurQuiParle] = []; 
            etatsJoueurs[joueurQuiParle] = false;
            passerAuSuivant();
        } else {
            mainsJoueurs[joueurQuiParle].push(carte);
            if (mainsJoueurs[joueurQuiParle].length === 7) {
                alert("FLIP 7 ! Bonus +15 pts !");
                scoresGlobaux[joueurQuiParle] += (mainsJoueurs[joueurQuiParle].reduce((a, b) => a + b, 0) + 15);
                afficherPodium("FLIP 7 : Manche Terminée !");
            } else {
                passerAuSuivant();
            }
        }
    };

    document.getElementById('btn-stop').onclick = () => {
        scoresGlobaux[joueurQuiParle] += mainsJoueurs[joueurQuiParle].reduce((a, b) => a + b, 0);
        etatsJoueurs[joueurQuiParle] = false;
        passerAuSuivant();
    };

    function passerAuSuivant() {
        if (!etatsJoueurs.includes(true)) {
            afficherPodium("Fin de la manche");
        } else {
            do {
                joueurQuiParle = (joueurQuiParle + 1) % nbJoueurs;
            } while (!etatsJoueurs[joueurQuiParle]);
            majInterface();
        }
    }

    function afficherPodium(titre) {
        let classement = scoresGlobaux.map((s, i) => ({ id: i + 1, score: s }))
                                     .sort((a, b) => b.score - a.score);

        document.getElementById('recap-titre').innerText = titre;
        const liste = document.getElementById('podium-liste');
        liste.innerHTML = "";
        
        classement.forEach((j, idx) => {
            let badge = (idx === 0) ? "🥇" : (idx === 1 ? "🥈" : (idx === 2 ? "🥉" : ""));
            liste.innerHTML += `
                <div class="podium-line">
                    <span>${badge} Joueur ${j.id}</span>
                    <span>${j.score} pts</span>
                </div>`;
        });

        overlay.style.display = "flex";
        
        document.getElementById('btn-continuer').onclick = () => {
            if (scoresGlobaux.some(s => s >= 200)) {
                alert("PARTIE FINIE ! Le navigateur va redémarrer.");
                location.reload();
            } else {
                donneurIndex = (donneurIndex + 1) % nbJoueurs;
                joueurQuiParle = (donneurIndex + 1) % nbJoueurs;
                initialiserManche();
            }
        };
    }

    initialiserManche();
});
