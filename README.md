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

