package main

// Shows a possible workaround for the gif encoder bug where images with bounds
// that don't start at (0,0) are incorrectly encoded.

import (
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	f, err := os.Open("lucha.gif")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	im, err := gif.DecodeAll(f)
	if err != nil {
		log.Fatal(err.Error())
	}

	firstFrame := im.Image[0]
	srcBounds := firstFrame.Bounds()

	// Create a crop region equal to the middle half.
	dstBounds := image.Rect(
		srcBounds.Min.X,
		srcBounds.Min.Y+srcBounds.Dy()/4,
		srcBounds.Max.X,
		srcBounds.Max.Y-srcBounds.Dy()/4)

	// Should only be one frame.
	for index, frame := range im.Image {
		si := frame.SubImage(dstBounds)

		// Creates a new image bounded at (0,0) and copies in the SubImage.
		b := image.Rect(0, 0, dstBounds.Dx(), dstBounds.Dy())
		pm := image.NewPaletted(b, frame.Palette)
		draw.Draw(pm, b, si, dstBounds.Min, draw.Src)
		im.Image[index] = pm
	}

	gout, err := os.Create("actual.wa.gif")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer gout.Close()
	gif.EncodeAll(gout, im)

	// Save the frame as a jpeg.
	jout, err := os.Create("expected.wa.jpg")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer jout.Close()
	jpeg.Encode(jout, im.Image[0], &jpeg.Options{90})
}
