package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"runtime"
	"sync"
	"time"

	"projet_gns/imagefilters"
)

func processChunk(src, dst *image.RGBA, bounds image.Rectangle, startY, endY int, filterName string, filter imagefilters.FilterFunc) {
	for y := startY; y < endY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			// Cas spécial : flou gaussien
			if filterName == "gaussien" {
				dst.Set(x, y, imagefilters.FlouGaussien(src, x, y, 30))
				continue
			}

			// Cas général : filtres pixel → pixel
			r16, g16, b16, a16 := src.At(x, y).RGBA()
			R := uint8(r16 >> 8)
			G := uint8(g16 >> 8)
			B := uint8(b16 >> 8)
			A := uint8(a16 >> 8)

			dst.Set(x, y, filter(R, G, B, A))
		}
	}
}

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

	bounds := img.Bounds()

	src := image.NewRGBA(bounds)
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			src.Set(x, y, img.At(x, y))
		}
	}

	filter, ok := imagefilters.Filters[filterName]
	if !ok && filterName != "gaussien" {
		fmt.Println("Filtre inconnu :", filterName)
		return
	}

	ncpu := runtime.NumCPU()
	ngoroutines := ncpu * 2

	height := bounds.Max.Y - bounds.Min.Y

	// Cas de résolution trop petite
	if ngoroutines > height {
		ngoroutines = height
	}

	packetSize := height / ngoroutines

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < ngoroutines; i++ {
		startY := bounds.Min.Y + i*packetSize
		endY := startY + packetSize
		if i == ngoroutines-1 {
			endY = bounds.Max.Y
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			processChunk(src, dst, bounds, s, e, filterName, filter)
		}(startY, endY)
	}

	wg.Wait()

	fmt.Println("Temps de traitement :", time.Since(start))

	out, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	jpeg.Encode(out, dst, nil)

	fmt.Println("Image traitée avec filtre", filterName, "et sauvegardée dans", outputPath)
}
