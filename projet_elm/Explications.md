## 1. ParserTcTurtle.elm : L'analyseur syntaxique
- C'est la porte d'entrée des données utilisateur. Son rôle est de transformer une simple chaîne de caractères (du texte brut) en une structure de données que le langage Elm peut manipuler de manière sécurisée.
- Rôle technique : Il utilise la bibliothèque elm/parser. Il définit ce qu'est une "Instruction" (un type Union : Forward, Left, Right, Repeat).
- Fonctionnement : Si l'utilisateur écrit [Forward 50], le parser vérifie que la syntaxe est correcte (crochets, majuscules, nombres) et renvoie un objet Elm. S'il y a une erreur de frappe, il renvoie un message d'erreur explicite.

## 2. Drawing.elm : Le moteur de rendu géométrique
- C'est ici que réside la logique mathématique. Une fois que nous avons une liste d'instructions (ex: "Avance", "Tourne"), ce fichier calcule où se trouve la "tortue" sur le plan 2D.
- Rôle technique :
    - Simulation : Il parcourt la liste des instructions une par une en maintenant un état interne (position X, position Y et angle actuel).
    - Génération de segments : Chaque instruction Forward crée un segment de droite (coordonnées de départ et d'arrivée).
    - Calcul des limites (Bounds) : Il calcule automatiquement la taille minimale et maximale du dessin pour que le SVG généré soit toujours bien cadré (auto-zoom).
    - Conversion SVG : Il transforme ces données géométriques en balises <line> et <svg> pour l'affichage web.

## 3. Main.elm : Le chef d'orchestre (Interface utilisateur)
- C'est le point d'entrée de l'application. Il lie le parser et le rendu graphique à une interface interactive.
- Le Modèle (Model) : Stocke le texte écrit par l'utilisateur, le résultat de l'analyse, la couleur choisie et l'état de l'affichage.
- La Mise à jour (Update) : Gère les événements. Par exemple, quand vous cliquez sur le bouton "Carré", il injecte le texte Repeat 4 [Forward 80, Left 90] dans le champ de saisie.
- La Vue (View) : Définit l'aspect de la page (la zone de texte, les boutons de raccourcis comme "Cœur" ou "Étoile", et le sélecteur de couleur).

## 4. elm.json : La carte d'identité du projet
- C'est un fichier de configuration au format JSON indispensable à tout projet Elm.
- Rôle : Il liste les versions des bibliothèques externes nécessaires.
- Dépendances clés ici :
    - elm/parser pour analyser le texte.
    - elm/svg pour dessiner les formes à l'écran.
    - elm/browser et elm/html pour faire fonctionner l'application dans un navigateur.

## 5. main.js : Le pont vers le navigateur
- Le navigateur web ne sait pas lire directement le langage Elm.
- Rôle : C'est le résultat de la compilation. Le compilateur Elm a pris tous les fichiers .elm ci-dessus et les a traduits en JavaScript optimisé.
- Usage : C'est ce fichier qui est appelé dans une page HTML pour lancer l'application. Il contient également le "runtime" d'Elm qui gère la mise à jour ultra-rapide du DOM (Document Object Model).
