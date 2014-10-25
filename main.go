package main

import (
	"flag"
	"fmt"
	"github.com/claudetech/smartcrop"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"os"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%s", err.Error())
	os.Exit(1)
}

var (
	width  = flag.Int("width", 128.0, "cropped image width")
	height = flag.Int("height", 128.0, "cropped image height")
	input  = flag.String("input", "input.jpg", "path for the input file")
	output = flag.String("output", "output.jpg", "path for the output file")
)

func main() {
	flag.Parse()
	fi, err := os.Open(*input)
	if err != nil {
		fail(err)
	}
	defer fi.Close()

	img, _, err := image.Decode(fi)
	if err != nil {
		fail(err)
	}
	crop, err := smartcrop.SmartCrop(&img, *width, *height)
	if err != nil {
		fail(err)
	}
	sub, ok := img.(SubImager)
	if ok {
		cropImage := sub.SubImage(image.Rect(crop.X, crop.Y,
			crop.Width+crop.X, crop.Height+crop.Y))
		m := resize.Resize(uint(*width), uint(*height), cropImage, resize.Lanczos3)
		toimg, _ := os.Create(*output)
		defer toimg.Close()
		err = jpeg.Encode(toimg, m, &jpeg.Options{jpeg.DefaultQuality})
	} else {
		err = fmt.Errorf("subimage not supported")
	}
	if err != nil {
		fail(err)
	}
}
