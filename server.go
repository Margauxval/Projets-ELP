package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"net"
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
			fmt.Println("Erreur de connexion :", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Lire l'image envoyée
	img, _, err := image.Decode(conn)
	if err != nil {
		fmt.Println("Erreur de décodage image :", err)
		return
	}

	// Traiter l'image (ex: filtre bleuté)
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			R := uint8(r >> 8)
			G := uint8(g >> 8)
			B := uint8(b >> 8)
			A := uint8(a >> 8)

			// Filtre bleuté simple
			newR := R / 2
			newG := G / 2
			newB := uint8(min(int(B)+80, 255))

			rgba.Set(x, y, color.RGBA{R: newR, G: newG, B: uint8(newB), A: A})
		}
	}

	// Renvoyer l'image modifiée
	jpeg.Encode(conn, rgba, nil)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
