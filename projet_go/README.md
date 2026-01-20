# Ce que fait le code

## Le Serveur (serveur.go) :
- Reste en écoute de connexions entrantes.
- Découpe l'image reçue en "chunks" (morceaux).
- Applique le filtre choisi (noir et blanc, thermique, flou, etc.) en utilisant toutes les ressources du processeur disponible.
- Mesure les performances (temps de traitement) et renvoie le résultat au client.
  
## Le Client (client.go) :
- Lit une image JPEG sur votre disque.
- Se connecte au serveur via le protocole TCP (port 8080).
- Envoie le nom du filtre souhaité et les données de l'image.
- Récupère l'image transformée et l'enregistre localement.


# Comment lancer le code :

Pour faire fonctionner ce programme, vous devez ouvrir deux terminaux différents.

### 1. Lancer le Serveur
Dans le premier terminal, allez dans le dossier contenant serveur.go et lancez-le :

```go run serveur.go```

Le serveur va afficher "Serveur en écoute sur le port 8080...". Laissez cette fenêtre ouverte.

### 2. Lancer le Client
Dans un second terminal, utilisez la commande suivante en donnant les bons arguments :

```go run client.go <chemin_image_entree> <chemin_image_sortie> <nom_filtre>```

Liste des arguments:
- <chemin_image_entree>   : chemin vers l'image source (ex: input.jpg)
- <chemin_image_sortie>   : chemin vers l'image générée (ex: output.jpg) vous pouvez choisir le nom de votre nouvelle image !
- <nom_filtre>            : filtre à appliquer

Filtres disponibles :
- noirblanc -> applique un filtre noir et blanc
- thermique -> applique un filtre thermique 
- yellowfluo -> remplace les pixels jaunes par des pixels fluos random
- rouge, orange, jaune, vert, bleu, violet -> applique un filtre de couleur sur toute l'image
- floubox -> floute l'entiereté de l'image. Rq : en modifiant le code dans la fonction processchunk (env. l.30 dans serveur.go), vous pouvez modifier l'intensité du flou.

Si vous oubliez quels arguments mettre, pas de panique ! Lancez le client sans argument et il vous rappèlera quoi mettre. 

### 3. Arrêter le serveur
Une fois toutes les photos modifiées à votre guise, vous devez arrêter le serveur. Afin de pouvoir en lancer un nouveau derrière, n'oubliez pas de taper la commande :

``` kill -9 <PID> ```

Le PID est accessible avec la commande : 

``` ps aux | grep go ``` 

 Le PID correspondra à celui en face de serveur.go.

 ## Amusez-vous bien avec les filtres fous !!


