package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/jpeg"
	"net"

	"project/filters"
	"project/processing"
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Serveur en écoute sur le port 9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erreur connexion :", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Lecture du filtre souhaité
	filterNameBuf := make([]byte, 32)
	conn.Read(filterNameBuf)
	filterName := string(filterNameBuf)
	filter := filters.Get(filterName)

	// taille et data de l'image
	var size uint32
	binary.Read(conn, binary.BigEndian, &size)
	imgData := make([]byte, size)
	conn.Read(imgData)

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		fmt.Println("Erreur decode image :", err)
		return
	}

	// lancement des go routines (traitement parallèle)
	result := processing.ApplyFilterParallel(img, filter)

	// Encoder et renvoyer (renvoi dans la commande)
	jpeg.Encode(conn, result, nil)
}
