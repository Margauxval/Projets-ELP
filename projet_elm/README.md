# Magic's Draw — Interpréteur TcTurtle
## Interface graphique TC-Turtle 

![Interface graphique](interface_graphique.png)

## Structure du Projet
Voici les fichiers qui composent notre projet: 
- Elm
  - Main.elm : Cœur de l'application (Modèle, Update, Vue).
  - ParserTcTurtle.elm : Analyseur syntaxique utilisant le package elm/parser.
  - Drawing.elm : Moteur de rendu calculant les segments et les limites du dessin.
- un fichier ```elm.json``` : définit la manière dont le compilateur doit construire l'application (librairies à télécharger surtout)
- un fichier ```main.js``` qui est généré par le code ```Main.elm``` lors de la commande : ```elm make src/Main.elm --output=main.js```
- un fichier ```index.html``` qui affiche le code généré par le fichier ```main.js```

## Lancement du programme
### A la main 
- Tapez la commande ```elm make src/Main.elm --output=main.js``` puis ouvrez le fichier ```index.html``` dans votre navigateur.
### Avec elm reactor
- Tapez la commande ```elm reactor ``` : un lien du type ... va s'afficher, en appuyant sur la touche ```Ctrl``` cliquez sur ce lien et une interface prête à l'emploi s'affiche dans votre navigateur.

## Une fois dans l'application
- Amusez-vous à dessiner les formes que vous voulez, en tapant votre code dans la zone de texte, ou en utilisant directement les boutons proposés.
- N'hésitez pas à redimensionner les formes proposées par les boutons à l'aide de la zone de texte.
- Vous pouvez également changer la couleur de votre dessin.
- Attention, votre dessin ne s'affichera pas si vous n'appuyez pas sur le bouton "Tracer". Pour effacer votre dessin, il faut aussi rappuyer sur la touche "Tracer" pour que la page redevienne vierge.

### Enjoy Magic's Draw !! 
