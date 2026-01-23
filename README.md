# Projets-ELP




- lien du discord : https://discord.gg/GeVZBCfj 
- lien du word partagé pour le CR de ELM: https://1drv.ms/w/c/1265599f2c973b15/IQBkCcqY6tWAR5_e4vzt7A3hAdIRe-JQZ7qYGWRfCDgiZhA?e=ceyja6
- lien du github pour JS : https://github.com/sfrenot/javascript/tree/master/projet4
## Projet GO : image filter 
Ne pas oublier d'évaluer les performances suivant la taille de l'image.
mettre le nombre de goroutine en parametre pour diviser le travail de l'image.
* indications du prof :client transferer un objet sans save en disque : gob dcoder img en local balancer dans un 
* attention si l'img est plus petite que la go routine (l 384) : verif si start y is in bound 

numcpu : nb de cpu ds l'ordi si pas param je prends osnumcpu attention aux divisions décimales 
on peut passer en go routine une référence et direct changer la photo au lieu d'en créer une nouvelle 

### Notre avancée / objectifs (GO)

- on a commencé à évaluer les perf : temps d'exe, on compte aussi tester sur des ordis avec un nb de coeurs différent (un avec 6 coeurs de 2 theards par coeur et un de 8 coeurs de simple threading), on sait pas encore le reste
- les filtres sont opérationnels
- on compte commencer le CR jeudi
- il y a encore bcp de chatgpt dans le code mais pareil on essaie de régler ça avant jeudi
- on compte faire un dossier "filters" qui contiendra les librairies de fichier pour raccourcir le fichier main.
- si on a le temps : améliorer les flags, finir modif couleur sélect et essayer distorsion. 

### Les filtres
- Noir et blanc (ok)
- Flou gaussien (ok)
- Filtre de couleur (ok)
- filtre thermique (ok)
- modification d'une couleur sélectionnée (par une autre couleur ou des pixel multicolors) (1 couleur suelement pour l'instant: jaune)


## Projet Elm 

On pensait plutot partir sur l'affichage d'un dessin que sur le word guesser, 
Pour l'instant, on a : 
- un fichier Page.elm qui code la mise en page et met tout ça sur un dossier main.js
- un fichier main.js qui contient le code en java script de la page
- un index.html pour interpréter et afficher le fichier js dans le navigateur

Pour exe le projet, en gros à chaque modif il faut actualiser le main.js --> taper elm make src/Main.elm --output=main.js
dans le terminal!!

### Des idées d'amélioration du projet : (ELM)

- Coder + intuitivement les directions : appuyer sur des boutons pour changer de direction, ou avancer par exemple
- des formes pré définies : cercle, rectangle, triangle, tortue par exemple : appuyer sur un bouton pour mettre le code qui les crée.
- Changement de couleur : pareil appui sur un bouton pour changer de couleur ??

- J'ai bcp vu sur internet que c'était récurrent d'importer le format svg et de créer des formes direct dans le fichier, et de les associer avec des fonctions, à voir !!

