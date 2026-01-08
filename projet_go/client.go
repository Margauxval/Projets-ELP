package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"os"
)

func main() {
	fmt.Println(`
──────────────────────────────────────────────────────────────
  Programme de traitement d'image — JP, Lilou & Margaux
──────────────────────────────────────────────────────────────

Utilisation :
    go run client.go <chemin_image_entree> <chemin_image_sortie> <nom_filtre>

Arguments :
    <chemin_image_entree>   : chemin vers l'image source (ex: input.jpg)
    <chemin_image_sortie>   : chemin vers l'image générée (ex: output.jpg)
    <nom_filtre>            : filtre à appliquer

Filtres disponibles :
    noirblanc
    thermique
    yellowfluo
    rouge, orange, jaune, vert, bleu, violet
    gaussien (flou)

Exemple :
    go run client.go photo.jpg resultat.jpg gaussien

──────────────────────────────────────────────────────────────
`)

	// Vérification des arguments
	if len(os.Args) < 4 {
		fmt.Println("Usage : go run client.go <inputPath> <outputPath> <filterName>")
		return
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]
	filterName := os.Args[3]

	// Valeurs par défaut si besoin
	if outputPath == "" {
		outputPath = "output.jpg"
	}
	if filterName == "" {
		filterName = "noirblanc"
	}

	// Lecture de l'image
	file, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		panic(err)
	}

	// Connexion au serveur
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	if err := encoder.Encode(filterName); err != nil {
		panic(err)
	}
	if err := encoder.Encode(buf.Bytes()); err != nil {
		panic(err)
	}

	var result []byte
	if err := decoder.Decode(&result); err != nil {
		panic(err)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	outImg, _, err := image.Decode(bytes.NewReader(result))
	if err != nil {
		panic(err)
	}

	if err := jpeg.Encode(outFile, outImg, nil); err != nil {
		panic(err)
	}

	fmt.Println("Image traitée avec filtre", filterName, "et sauvegardée dans", outputPath)
}


