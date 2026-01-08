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

Filtres disponibles :
    noirblanc
    thermique
    yellowfluo
    rouge, orange, jaune, vert, bleu, violet
    gaussien (flou)

──────────────────────────────────────────────────────────────
`)

	var inputPath string
	var outputPath string
	var filterName string

	fmt.Print("Chemin de l'image d'entrée (ex: input.jpg) : ")
	fmt.Scanln(&inputPath)

	fmt.Print("Chemin de l'image de sortie (ex: output.jpg) : ")
	fmt.Scanln(&outputPath)

	fmt.Print("Nom du filtre : ")
	fmt.Scanln(&filterName)

	if inputPath == "" {
		fmt.Println("Erreur : aucun fichier d'entrée fourni.")
		return
	}
	if outputPath == "" {
		outputPath = "output.jpg"
	}
	if filterName == "" {
		filterName = "noirblanc"
	}

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
