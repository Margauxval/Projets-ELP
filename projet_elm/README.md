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




