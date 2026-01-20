package imagefilters

import (
	"image"
	"image/color"
	"math/rand"
)

// ─────────────────────────────────────────────
//  UTILITAIRES COULEURS (HSV, FLUO…)
// ─────────────────────────────────────────────

func rgbToHSV(r, g, b uint8) (float64, float64, float64) {
	// convertit une couleur du format RGB (rouge, vert, bleu) vers le format HSV (teinte, saturation, valeur).
	// Paramètres :
	//   - r : composante rouge (0 à 255)
	//   - g : composante verte (0 à 255)
	//   - b : composante bleue (0 à 255)
	// Retour :
	//   - h : teinte (0 à 360 degrés)
	//   - s : saturation (0.0 à 1.0)
	//   - v : valeur (0.0 à 1.0)
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
	if delta != 0 {
		switch max { // switch : calcule le max, en fonction de ce que ça vaut fait on fait le cas adapté : ex : max == G --> on va pas calculer avec la méthode R et B, mais on fait que G, on switch.
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
	// détermine si une couleur RGB correspond à une teinte jaune ou orange en espace HSV.
	// Paramètres :
	//   - r : composante rouge (0 à 255)
	//   - g : composante verte (0 à 255)
	//   - b : composante bleue (0 à 255)
	// Retour :
	//   - bool : true si oui, false sinon

	h, s, v := rgbToHSV(r, g, b)
	return h >= 35 && h <= 65 && s >= 0.25 && v >= 0.2 // plage jaune-orange
}

func randomFluoColor() color.RGBA {
	// retourne une couleur fluorescente aléatoire parmi une palette prédéfinie.
	// Paramètres :
	//   - aucun
	// Retour :
	//   - color.RGBA : une couleur fluo choisie aléatoirement dans la palette

	palette := []color.RGBA{
		{255, 0, 255, 255},
		{0, 255, 255, 255},
		{255, 255, 0, 255},
		{0, 255, 0, 255},
		{255, 0, 0, 255},
		{0, 0, 255, 255},
	}
	return palette[rand.Intn(len(palette))]
}

func makeColorFilter(target color.RGBA) FilterFunc {
	// crée une fonction de filtre qui applique une coloration basée sur une couleur cible. (c'était pour éviter de faire des flags différents dans le main il me semble)
	// Paramètres :
	//   - target : couleur cible
	// Retour :
	//   - FilterFunc : fonction prenant en entrée des composantes RGBA et retournant une couleur recolorée selon la cible
	return func(r, g, b, a uint8) color.RGBA {
		return colorize(r, g, b, a, target)
	}
}

// ─────────────────────────────────────────────
//  THERMIQUE
// ─────────────────────────────────────────────

var thermalPalette = []color.RGBA{
	// palette de couleurs "thermique"
	{0, 0, 128, 255},
	{0, 0, 255, 255},
	{0, 255, 255, 255},
	{0, 255, 0, 255},
	{255, 255, 0, 255},
	{255, 128, 0, 255},
	{255, 0, 0, 255},
	{255, 255, 255, 255},
}

func luminance(r, g, b uint8) float64 {
	// calcule la luminosité perçue d'une couleur RGB selon le modèle de luminosité standard.
	// Paramètres :
	//   - r : composante rouge
	//   - g : composante verte
	//   - b : composante bleue
	// Retour :
	//   - float64 : luminosité normalisée entre 0.0 (noir) et 1.0 (blanc)
	return (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255.0
}

func thermalColor(r, g, b uint8) color.RGBA {
	// génère une couleur pseudo-thermique en fonction de la luminosité d'une couleur RGB.
	// Paramètres :
	//   - r : composante rouge
	//   - g : composante verte
	//   - b : composante bleue
	// Retour :
	//   - color.RGBA : couleur thermique associée
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

// ─────────────────────────────────────────────
//  FLOU BOX
// ─────────────────────────────────────────────

func clamp(x, min, max int) int {
	// limite une valeur entière à un intervalle donné [min, max].
	// Paramètres :
	//   - x : valeur à contraindre
	//   - min : borne inférieure de l'intervalle
	//   - max : borne supérieure de l'intervalle
	// Retour :
	//   - int : valeur bornée, comprise entre min et max
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func FlouBox(src *image.RGBA, x, y, radius int) color.RGBA {
	// applique un flou moyen sur un pixel donné en moyennant les couleurs de ses voisins dans un rayon défini.
	// Paramètres :
	//   - src : pointeur vers l'image source
	//   - x : coordonnée horizontale du pixel à traiter
	//   - y : coordonnée verticale du pixel à traiter
	//   - radius : rayon du voisinage à prendre en compte pour le flou
	// Retour :
	//   - color.RGBA : couleur floutée résultante du pixel (x, y)
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

// ─────────────────────────────────────────────
//  FILTRES PIXEL → PIXEL
// ─────────────────────────────────────────────

type FilterFunc func(r, g, b, a uint8) color.RGBA // définition d'un type : ici FilterFunc qui est n'importe quelle fonction qui prend 4 entiers (r, g, b, a) et retourne une couleur RGBA

func colorize(r, g, b, a uint8, target color.RGBA) color.RGBA {
	// aplique le filtre de couleurs
	return color.RGBA{
		R: uint8((int(r) + int(target.R)) / 2),
		G: uint8((int(g) + int(target.G)) / 2),
		B: uint8((int(b) + int(target.B)) / 2),
		A: a,
	}
}

func filterNoirBlanc(r, g, b, a uint8) color.RGBA {
	// applique le filtre noir et blanc
	gray := uint8((int(r) + int(g) + int(b)) / 3)
	return color.RGBA{gray, gray, gray, a}
}

func filterYellowOrangeFluo(r, g, b, a uint8) color.RGBA {
	// applique le filtre qui transforme chauqe pixel jaune / orange en pixel fluo
	if isYellowOrOrangeHSV(r, g, b) {
		return randomFluoColor()
	}
	return color.RGBA{r, g, b, a}
}

func filterThermal(r, g, b, a uint8) color.RGBA {
	// applique le filtre thermique
	return thermalColor(r, g, b)
}

// pour les flags des filtres

var colorTargets = map[string]color.RGBA{
	"rouge":  {255, 0, 0, 255},
	"orange": {255, 128, 0, 255},
	"jaune":  {255, 255, 0, 255},
	"vert":   {0, 255, 0, 255},
	"bleu":   {0, 0, 255, 255},
	"violet": {128, 0, 255, 255},
}

var Filters = map[string]FilterFunc{
	"floubox":    nil, // cas spécial
	"noirblanc":  filterNoirBlanc,
	"yellowfluo": filterYellowOrangeFluo,
	"thermique":  filterThermal,

	"rouge":  makeColorFilter(colorTargets["rouge"]),
	"orange": makeColorFilter(colorTargets["orange"]),
	"jaune":  makeColorFilter(colorTargets["jaune"]),
	"vert":   makeColorFilter(colorTargets["vert"]),
	"bleu":   makeColorFilter(colorTargets["bleu"]),
	"violet": makeColorFilter(colorTargets["violet"]),
}
