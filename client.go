package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"os"
)

const serverAddr = "localhost:9000" // adresse fixe

func main() {

	//  Flags : permettent à l'user d'utiliser des param juste en les écrivant
	imagePath := flag.String("image", "", "Chemin vers l'image à envoyer")
	filterName := flag.String("filter", "noirblanc", "Nom du filtre à appliquer")
	flag.Parse()

	if *imagePath == "" {
		fmt.Println("Erreur : vous devez spécifier --image=chemin")
		return
	}

	// Connexion au serveur via TCP - j'ai pas trop compris ça, à revoir
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Charger l'img
	file, err := os.Open(*imagePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// Encoder l'img : à quoi ca sert ?
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)
	data := buf.Bytes()

	// Envoyer filtre, taille et data img
	padded := make([]byte, 32)
	copy(padded, []byte(*filterName))
	conn.Write(padded)
	binary.Write(conn, binary.BigEndian, uint32(len(data)))
	conn.Write(data)

	// Recevoir et save l'img modifiée
	resultImg, _, err := image.Decode(conn)
	if err != nil {
		panic(err)
	}
	out, err := os.Create("resultat.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	jpeg.Encode(out, resultImg, nil)

	fmt.Println("Image modifiée reçue et sauvegardée sous resultat.jpg")
}
