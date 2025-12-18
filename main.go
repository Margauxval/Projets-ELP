package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"os"
	"sync"
	"time"
)

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
		if h < 0 {
			h += 360
		}
		if h >= 360 {
			h -= 360
		}
	}
	return h, s, v
}

func isYellowOrOrangeHSV(r, g, b uint8) bool {
	h, s, v := rgbToHSV(r, g, b)

	// Teinte entre orange (≈35°) et jaune (≈65°)
	isHueInRange := h >= 35 && h <= 65
	isSaturated := s >= 0.25
	isBrightEnough := v >= 0.2

	return isHueInRange && isSaturated && isBrightEnough
}

func randomFluoColor() color.RGBA {
	// Palette de couleurs fluos typiques
	palette := []color.RGBA{
		{R: 255, G: 0, B: 255, A: 255}, // rose fluo
		{R: 0, G: 255, B: 255, A: 255}, // cyan fluo
		{R: 255, G: 255, B: 0, A: 255}, // jaune fluo
		{R: 0, G: 255, B: 0, A: 255},   // vert fluo
		{R: 255, G: 0, B: 0, A: 255},   // rouge vif
		{R: 0, G: 0, B: 255, A: 255},   // bleu électrique
	}
	return palette[rand.Intn(len(palette))]
}

func processLine(rgba *image.RGBA, y int, bounds image.Rectangle, wg *sync.WaitGroup) {
	defer wg.Done()
	rand.Seed(time.Now().UnixNano())

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		r16, g16, b16, a16 := rgba.At(x, y).RGBA()
		R := uint8(r16 >> 8)
		G := uint8(g16 >> 8)
		B := uint8(b16 >> 8)
		A := uint8(a16 >> 8)

		if isYellowOrOrangeHSV(R, G, B) {
			// Remplacer par une couleur fluo aléatoire
			c := randomFluoColor()
			c.A = A
			rgba.Set(x, y, c)
		}
	}
}

func main() {
	// Ouvrir l'image
	file, err := os.Open("testphoto.jpg")
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	defer file.Close()

	src, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Erreur de décodage :", err)
		return
	}

	bounds := src.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, src.At(x, y))
		}
	}

	// Utiliser des goroutines pour traiter chaque ligne
	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go processLine(rgba, y, bounds, &wg)
	}
	wg.Wait()

	// Sauvegarder
	out, err := os.Create("phototest_filtre.jpg")
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	defer out.Close()

	jpeg.Encode(out, rgba, nil)
	fmt.Println("Image filtrée sauvegardée en phototest_filtre.jpg")
}
