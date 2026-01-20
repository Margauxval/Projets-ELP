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

1. Lancer le Serveur
Dans le premier terminal, allez dans le dossier contenant serveur.go et lancez-le :
```go run serveur.go```
Le serveur va afficher "Serveur en écoute sur le port 8080...". Laissez cette fenêtre ouverte.

2. Lancer le Client
Dans un second terminal, utilisez la commande suivante en donnant les bons arguments :
```go run client.go <image_entree> <image_sortie> <nom_filtre>```
Liste des arguments:

Si vous oubliez quels arguments mettre, pas de panique ! Lancez le client sans argument et il vous rappèlera quoi mettre. 

