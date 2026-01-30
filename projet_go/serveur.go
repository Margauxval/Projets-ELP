package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"projet_go/imagefilters"
	"runtime"
	"sync"
	"time"
)

func processChunk(src, dst *image.RGBA, bounds image.Rectangle, startY, endY int, filterName string, filter imagefilters.FilterFunc) {
	// applique un filtre à une bande horizontale comprise entre startY et endY de l'image source src
	// Paramètres :
	//   - src : pointeur vers l'image source
	//   - dst : pointeur vers l'image destination
	//   - bounds : coordonnées de l'image à traiter
	//   - startY : ligne de début (incluse)
	//   - endY : ligne de fin (exclue)
	//   - filterName : nom du filtre à appliquer
	//Retour : pas de retour

	for y := startY; y < endY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if filterName == "floubox" {
				dst.Set(x, y, imagefilters.FlouBox(src, x, y, 30)) // modifiez 30 pour intensifier ou réduire le flou
				continue // passe à la boucle suivante
			}
			r16, g16, b16, a16 := src.At(x, y).RGBA()
			R := uint8(r16 >> 8) // .RGBA() renvoie une couleur sur 16 bits, on la remet sur 8 bits pr avoir entre 0 et 255
			G := uint8(g16 >> 8)
			B := uint8(b16 >> 8)
			A := uint8(a16 >> 8)
			dst.Set(x, y, filter(R, G, B, A))
		}
	}
}

func main() {
	// fonction main qui gère et appelle toutes les fonctions
	// connexion TCP au client
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Serveur en écoute sur le port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Erreur de connexion :", err)
			continue
		}
		go handleConnection(conn) //go routines envoyées pr gérer client
	}
}

func logPerf(filter string, workers int, duration time.Duration, bounds image.Rectangle) {
	// affiche les informations de performance du traitement d'image dans la console.
	// Paramètres :
	//   - filter : nom du filtre appliqué
	//   - workers : nombre de goroutines utilisées pour le traitement
	//   - duration : durée totale du traitement
	//   - bounds : dimensions de l'image traitée (largeur et hauteur)
	// Retour : pas de retour

	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	fmt.Println("────────────────────────────────────────")
	fmt.Println("Filtre utilisé      :", filter)
	fmt.Println("Goroutines utilisées:", workers)
	fmt.Println("Taille image        :", width, "x", height)
	fmt.Println("Temps de traitement :", duration)
	fmt.Println("────────────────────────────────────────")
}

func handleConnection(conn net.Conn) {
	// gère une connexion entrante, applique un filtre à l'image reçue, puis renvoie l'image modifiée.
	// Paramètres :
	//   - conn : connexion réseau TCP avec le client
	// Retour : pas de retour
	defer conn.Close() // si fin imprévue du programme cette ligne sera tapée à la toute fin

	//création de pipes	
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	var filterName string
	var imgBytes []byte

	if err := decoder.Decode(&filterName); err != nil {
		fmt.Println("Erreur lecture filtre :", err)
		return
	}
	if err := decoder.Decode(&imgBytes); err != nil {
		fmt.Println("Erreur lecture image :", err)
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		fmt.Println("Erreur décodage image :", err)
		return
	}

	bounds := img.Bounds()
	src := image.NewRGBA(bounds) // crée une nv img de taille bounds ?
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			src.Set(x, y, img.At(x, y))
		}
	}
	filter, ok := imagefilters.Filters[filterName]
	if !ok && filterName != "floubox" {
		fmt.Println("Filtre inconnu :", filterName)
		return
	}

	ncpu := runtime.NumCPU()
	ngoroutines := ncpu * 2
	height := bounds.Max.Y - bounds.Min.Y
	if ngoroutines > height {
		ngoroutines = height
	}
	packetSize := height / ngoroutines

	start := time.Now()

	var wg sync.WaitGroup // pr vérifier que toutes les go routines sont synchronisées avant de lancer
	for i := 0; i < ngoroutines; i++ {
		startY := bounds.Min.Y + i*packetSize
		endY := startY + packetSize
		if i == ngoroutines-1 {
			endY = bounds.Max.Y
		}
		wg.Add(1)
		go func(s, e int) { //start et end
			defer wg.Done()
			processChunk(src, dst, bounds, s, e, filterName, filter)
		}(startY, endY)
	}
	wg.Wait()

	duration := time.Since(start)
	logPerf(filterName, ngoroutines, duration, bounds)

	var outBuf bytes.Buffer // en gros crée un buffer pour simuler le fichier dans la mémoire vive 
	if err := jpeg.Encode(&outBuf, dst, nil); err != nil {
		fmt.Println("Erreur encodage JPEG :", err)
		return
	}

	if err := encoder.Encode(outBuf.Bytes()); err != nil {
		fmt.Println("Erreur envoi image :", err)
		return
	}
}


