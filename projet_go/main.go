package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

var thermalPalette = []color.RGBA{
	{0, 0, 128, 255},     // bleu foncé
	{0, 0, 255, 255},     // bleu
	{0, 255, 255, 255},   // cyan
	{0, 255, 0, 255},     // vert
	{255, 255, 0, 255},   // jaune
	{255, 128, 0, 255},   // orange
	{255, 0, 0, 255},     // rouge
	{255, 255, 255, 255}, // jaune clair
}

//
// ────────────────────────────────────────────────────────────────
//  UTILITAIRES COULEURS (HSV, FLUO…)
// ────────────────────────────────────────────────────────────────
//

func rgbToHSV(r, g, b uint8) (float64, float64, float64) {
	R := float64(r) / 255.0
	G := float64(g) / 255.0
	B := float64(b) / 255.0

	max := R
	if G > max {
		max = G
	}
	if B > max {
		max = B
	}

	min := R
	if G < min {
		min = G
	}
	if B < min {
		min = B
	}

	v := max
	delta := max - min

	var s float64
	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	var h float64
	if delta == 0 {
		h = 0
	} else {
		switch max {
		case R:
			h = 60 * ((G - B) / delta)
			if h < 0 {
				h += 360
			}
		case G:
			h = 60 * (((B - R) / delta) + 2)
		case B:
			h = 60 * (((R - G) / delta) + 4)
		}
	}

	return h, s, v
}

func isYellowOrOrangeHSV(r, g, b uint8) bool {
	h, s, v := rgbToHSV(r, g, b)

	isHueInRange := h >= 35 && h <= 65
	isSaturated := s >= 0.25
	isBrightEnough := v >= 0.2

	return isHueInRange && isSaturated && isBrightEnough
}

func randomFluoColor() color.RGBA {
	palette := []color.RGBA{
		{255, 0, 255, 255}, // rose fluo
		{0, 255, 255, 255}, // cyan fluo
		{255, 255, 0, 255}, // jaune fluo
		{0, 255, 0, 255},   // vert fluo
		{255, 0, 0, 255},   // rouge vif
		{0, 0, 255, 255},   // bleu électrique
	}
	return palette[rand.Intn(len(palette))]
}

func makeColorFilter(target color.RGBA) FilterFunc {
	return func(r, g, b, a uint8) color.RGBA {
		return colorize(r, g, b, a, target)
	}
}

// pour filtre thermique :

// calcul de luminosité en %
func luminance(r, g, b uint8) float64 {
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255.0
}

// On prend la luminosité normalisée et on la mappe dans la palette (à comprendre car on a pas tt compris)
func thermalColor(r, g, b uint8) color.RGBA {
	lum := luminance(r, g, b)
	idx := lum * float64(len(thermalPalette)-1)

	i := int(idx)
	frac := idx - float64(i)

	if i >= len(thermalPalette)-1 {
		return thermalPalette[len(thermalPalette)-1]
	}

	c1 := thermalPalette[i]
	c2 := thermalPalette[i+1]

	return color.RGBA{
		R: uint8(float64(c1.R)*(1-frac) + float64(c2.R)*frac),
		G: uint8(float64(c1.G)*(1-frac) + float64(c2.G)*frac),
		B: uint8(float64(c1.B)*(1-frac) + float64(c2.B)*frac),
		A: 255,
	}
}

// flou gaussien : gérer les bords dde l'img : *
func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

//
// ────────────────────────────────────────────────────────────────
//  BIBLIOTHÈQUE DE FILTRES
// ────────────────────────────────────────────────────────────────
//

type FilterFunc func(r, g, b, a uint8) color.RGBA

func colorize(r, g, b, a uint8, target color.RGBA) color.RGBA {
	// Mélange simple : 50% pixel original, 50% couleur cible
	return color.RGBA{
		R: uint8((int(r) + int(target.R)) / 2),
		G: uint8((int(g) + int(target.G)) / 2),
		B: uint8((int(b) + int(target.B)) / 2),
		A: a,
	}
}
func flouGaussien(src *image.RGBA, x, y, radius int) color.RGBA {
	var rSum, gSum, bSum float64
	var count float64

	bounds := src.Bounds()

	for ky := -radius; ky <= radius; ky++ {
		for kx := -radius; kx <= radius; kx++ {

			px := clamp(x+kx, bounds.Min.X, bounds.Max.X-1)
			py := clamp(y+ky, bounds.Min.Y, bounds.Max.Y-1)

			r16, g16, b16, _ := src.At(px, py).RGBA()
			rSum += float64(uint8(r16 >> 8))
			gSum += float64(uint8(g16 >> 8))
			bSum += float64(uint8(b16 >> 8))
			count++
		}
	}

	return color.RGBA{
		R: uint8(rSum / count),
		G: uint8(gSum / count),
		B: uint8(bSum / count),
		A: 255,
	}
}

func filterNoirBlanc(r, g, b, a uint8) color.RGBA {
	gray := uint8((int(r) + int(g) + int(b)) / 3)
	return color.RGBA{gray, gray, gray, a}
}

func filterYellowOrangeFluo(r, g, b, a uint8) color.RGBA {
	if isYellowOrOrangeHSV(r, g, b) {
		return randomFluoColor()
	}
	return color.RGBA{r, g, b, a}
}

func filterThermal(r, g, b, a uint8) color.RGBA {
	return thermalColor(r, g, b)
}

// pour le filtre colorize
var colorTargets = map[string]color.RGBA{
	"rouge":  {255, 0, 0, 255},
	"orange": {255, 128, 0, 255},
	"jaune":  {255, 255, 0, 255},
	"vert":   {0, 255, 0, 255},
	"bleu":   {0, 0, 255, 255},
	"violet": {128, 0, 255, 255},
}

var filters = map[string]FilterFunc{
	"gaussien":   nil, // car n'appelle aucune fonction
	"noirblanc":  filterNoirBlanc,
	"yellowfluo": filterYellowOrangeFluo,
	"thermique":  filterThermal,

	// Filtres colorize
	"rouge":  makeColorFilter(colorTargets["rouge"]),
	"orange": makeColorFilter(colorTargets["orange"]),
	"jaune":  makeColorFilter(colorTargets["jaune"]),
	"vert":   makeColorFilter(colorTargets["vert"]),
	"bleu":   makeColorFilter(colorTargets["bleu"]),
	"violet": makeColorFilter(colorTargets["violet"]),
}

var selectedFilter = "gaussien" // en dur pour l’instant

//
// ────────────────────────────────────────────────────────────────
//  TRAITEMENT PARALLÈLE
// ────────────────────────────────────────────────────────────────
//

func processChunk(src, dst *image.RGBA, bounds image.Rectangle, startY, endY int, filterName string, filter FilterFunc) {
	for y := startY; y < endY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			// Cas spécial : flou gaussien
			if filterName == "gaussien" {
				dst.Set(x, y, flouGaussien(src, x, y, 30))
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

// evaluation des performances:
func logPerf(filter string, workers int, duration time.Duration, bounds image.Rectangle) {
    width := bounds.Max.X - bounds.Min.X
    height := bounds.Max.Y - bounds.Min.Y

    fmt.Println("────────────────────────────────────────")
    fmt.Println("Filtre utilisé      :", filter)
    fmt.Println("Goroutines utilisées:", workers)
    fmt.Println("Taille image        :", width, "x", height)
    fmt.Println("Temps de traitement :", duration)
    fmt.Println("────────────────────────────────────────")
}

//
// ────────────────────────────────────────────────────────────────
//  MAIN
// ────────────────────────────────────────────────────────────────
//

func main() {

	// ───────────────────────────────────────────────
	//  DOCSTRING AFFICHÉE AU LANCEMENT
	// ───────────────────────────────────────────────
	fmt.Println(`
──────────────────────────────────────────────────────────────
  Programme de traitement d'image — Lilou & Co
──────────────────────────────────────────────────────────────

Usage interactif :
    go run main.go
    → puis entrez les informations demandées

Filtres disponibles :
    noirblanc
    thermique
    yellowfluo
    rouge, orange, jaune, vert, bleu, violet
    gaussien (flou)

──────────────────────────────────────────────────────────────
`)

	// ───────────────────────────────────────────────
	//  MODE INTERACTIF : DEMANDE DES PARAMÈTRES
	// ───────────────────────────────────────────────
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

	// ───────────────────────────────────────────────
	//  CHARGEMENT DE L'IMAGE
	// ───────────────────────────────────────────────
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

	// src = image source
	src := image.NewRGBA(bounds)
	// dst = image destination
	dst := image.NewRGBA(bounds)

	// Copier l’image dans src
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			src.Set(x, y, img.At(x, y))
		}
	}

	// ───────────────────────────────────────────────
	//  CHOIX DU FILTRE
	// ───────────────────────────────────────────────
	// CHOIX DU FILTRE
	filter, ok := filters[filterName]
	if !ok && filterName != "gaussien" {
		fmt.Println("Filtre inconnu :", filterName)
		return
	}

	// PARALLÉLISATION
	ncpu := runtime.NumCPU()
	ngoroutines := ncpu * 2
	packetSize := (bounds.Max.Y - bounds.Min.Y) / ngoroutines

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

	duration := time.Since(start)
	logPerf(filterName, ngoroutines, duration, bounds)

	// SAUVEGARDE DE L'IMAGE
	out, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	jpeg.Encode(out, dst, nil)

	fmt.Println("Image traitée avec filtre", filterName, "et sauvegardée dans", outputPath)

}

