# Projets-ELP
## Image filter



** lien du discord : https://discord.gg/GeVZBCfj **
Ne pas oublier d'évaluer les performances suivant la taille de l'image.
mettre le nombre de goroutine en parametre pour diviser le travail de limage.
Pas obligé client serveur.


* indications du prof :client transferer un objet sans save en disque : gob dcoder img en local balancer dans un 
* attention si l'img est plus petite que la go routine (l 384) : verif si start y is in bound 

numcpu : nb de cpu ds l'ordi si pas param je prends osnumcpu attention aux divisions décimales 
on peut passer en go routine une référence et direct changer la photo au lieu d'en créer une nouvelle 
flou gaussien : mmoyenne sur les px, attention aux bordures

## Notre avancée / objectifs

- on a commencé à évaluer les perf : temps d'exe, on compte aussi tester sur des ordis avec un nb de coeurs différent (un avec 6 coeurs de 2 theards par coeur et un de 8 coeurs de simple threading), on sait pas encore le reste
- les filtres sont opérationnels
- on compte commencer le CR jeudi
- il y a encore bcp de chatgpt dans le code mais pareil on essaie de régler ça avant jeudi
- on compte faire un dossier "filters" qui contiendra les librairies de fichier pour raccourcir le fichier main.
- si on a le temps : améliorer les flags, finir modif couleur sélect et essayer distorsion. 

## Les filtres
- Noir et blanc (ok)
- Flou gaussien (ok)
- Filtre de couleur (ok)
- filtre thermique (ok)
- (si on a le temps) filtre de distorsion
- modification d'une couleur sélectionnée (par une autre couleur ou des pixel multicolors) (1 couleur suelement pr l'instant)


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
- caulcul coeurs utilisables go routines

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


---

### 6. Les difficultés rencontrées  
Parlez honnêtement de :  
– bugs de concurrence  
– problèmes de protocole TCP  
– erreurs de design que vous avez corrigées  
– limites de votre solution  

### Pb à résoudre : 
-  flou gaussien, plage de couleurs thermique, sélect couleur et changer en une autre couleur : trop picky sur la sélection de couleur 
-  si la résqolution de l'img est trop basse : anticiper le cas
-  effets de bord flou gaussien
-  
- fonctions privées public "filters" on avait pas tt compris
- 


---

### 7. Conclusion  
En quelques lignes :  
– ce que vous avez appris sur Go  
– ce que vous auriez amélioré avec plus de temps  
– ce que ce projet vous a apporté

