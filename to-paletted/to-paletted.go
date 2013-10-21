package main

// Demonstrates how to create a paletted image, e.g. GIF, from a full RGB image.
// The default encoding uses a pre-defined palette that creates non-optimal
// results.

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	// Open the original JPEG.
	f, err := os.Open("china.jpg")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	// Decode it.
	im, err := jpeg.Decode(f)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Do the conversion.
	pm := ImageToPaletted(im)

	// Create new, single frame GIF.
	g := &gif.GIF{
		Image: []*image.Paletted{pm},
		Delay: []int{0},
	}

	// Encode and save the gif.
	out, err := os.Create("china.gif")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer out.Close()
	gif.EncodeAll(out, g)
}

func ImageToPaletted(img image.Image) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, palette.Plan9)
	draw.FloydSteinberg.Draw(pm, b, img, image.ZP)
	return pm
}
