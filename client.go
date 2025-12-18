package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Ouvrir et décoder l'image
	file, err := os.Open("testphoto.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	// Envoyer l'image au serveur
	err = jpeg.Encode(conn, img, nil)
	if err != nil {
		panic(err)
	}

	// Recevoir l'image modifiée
	outImg, _, err := image.Decode(conn)
	if err != nil {
		panic(err)
	}

	// Sauvegarder l'image reçue
	out, err := os.Create("photo_modifiee.jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	err = jpeg.Encode(out, outImg, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("Image modifiée reçue et sauvegardée.")
}
