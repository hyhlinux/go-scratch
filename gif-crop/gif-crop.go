package main

// Demonstrates a bug in the gif encoder when an image with non-(0,0) bounds is
// saved.  See https://code.google.com/p/go/issues/detail?id=6635

import (
	"image"
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
		im.Image[index] = frame.SubImage(dstBounds).(*image.Paletted)
	}

	gout, err := os.Create("actual.gif")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer gout.Close()
	gif.EncodeAll(gout, im)

	// Save the frame as a jpeg.
	jout, err := os.Create("expected.jpg")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer jout.Close()
	jpeg.Encode(jout, im.Image[0], &jpeg.Options{90})
}
