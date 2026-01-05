# Projets-ELP
## Image filter


* lien du replit : https://replit.com/@lilounadlerr/Image-FilterGo *
Ne pas oublier d'évaluer les performances suivant la taille de l'image.
mettre le nombre de goroutine en parametre pour diviser le travail de limage.
Pas obligé client serveur.


* indications du prof : mieux vaut se pencher sur l'évaluation des perf et la connexion tcp
attention pas opti car une gr par pixel, idem pour ligne par ligne ou colonne par colonne
optimal = nb cœur *2, mettre en param nb go routine

numcpu : nb de cpu ds l'ordi si pas param je prends osnumcpu attention aux divisions décimales 
on peut passer en go routine une référence et direct changer la photo au lieu d'en créer une nouvelle 
flou gaussien : mmoyenne sur les px, attention aux bordures

## Notre avancée / objectifs

- commencer à regarder les perfs du programme
- faire une "bilbiothèque" de filtres utilisables
- établir une co TCP serv/ client avec des tags pour exe le programme selon ce que veut le client. (commencé)

## Les filtres
- Noir et blanc
- Flou gaussien
- Filtre de couleur
- modification d'une couleur sélectionnée (par une autre couleur ou des pixel multicolors)


## CR 

### 1. Le problème choisi  
Explique en quelques phrases :  
– quel problème vous avez décidé de résoudre (ex : multiplication de matrices, Levenshtein, random walks, etc.)  
– pourquoi ce problème est intéressant pour la concurrence  
– quelles sont les entrées et les sorties de votre programme

---

### 2. L’architecture générale  
Décris votre application :  
– version locale simple  
– version locale concurrente  
– version serveur TCP concurrent  
– worker pool si vous en avez un  

Explique comment les différentes parties communiquent : goroutines, channels, workers, etc.

---

### 3. Le design concurrent  
C’est la partie la plus importante.  
Explique :  
– quelles goroutines existent  
– ce qu’elles font  
– comment vous synchronisez (waitgroups, channels)  
– pourquoi ce design est efficace pour votre problème  
– comment vous évitez les blocages / deadlocks

---

### 4. Le serveur TCP  
Décris :  
– comment vous acceptez les connexions  
– comment vous gérez plusieurs clients en parallèle  
– comment vous découpez les messages (lecture bloquante, fin de ligne, protocole simple)  
– comment vous renvoyez les résultats

---

### 5. Les performances  
Explique ce que vous avez mesuré :  
– temps d’exécution en séquentiel  
– temps d’exécution en concurrent  
– impact du nombre de workers  
– limites rencontrées (CPU, I/O, taille des données, etc.)  

Pas besoin de graphiques compliqués, juste montrer que vous avez testé.

---

### 6. Les difficultés rencontrées  
Parlez honnêtement de :  
– bugs de concurrence  
– problèmes de protocole TCP  
– erreurs de design que vous avez corrigées  
– limites de votre solution  

Pierre aime quand on montre qu’on a appris quelque chose.

---

### 7. Conclusion  
En quelques lignes :  
– ce que vous avez appris sur Go  
– ce que vous auriez amélioré avec plus de temps  
– ce que ce projet vous a apporté

