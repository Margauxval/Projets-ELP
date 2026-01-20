# Interface graphique TC-Turtle 


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


