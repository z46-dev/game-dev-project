package main

import (
	"flag"
	"fmt"
	"log"

	"imgfixr/imgproc"
)

var (
	inputPath  string
	outputPath string
)

func init() {
	flag.StringVar(&inputPath, "in", "", "Path to input image (PNG, JPEG, etc.)")
	flag.StringVar(&outputPath, "out", "", "Path to output PNG")
}

var (
	TargetWhiteBG = imgproc.MaskConfig{
		Mode: imgproc.MaskHSV,

		HSVTargetH: 0,
		HSVTolH:    10,
		HSVMinS:    0,
		HSVMaxS:    1,
		HSVMinV:    0,
		HSVMaxV:    1,
	}

	TargetMagentaBG = imgproc.MaskConfig{
		Mode: imgproc.MaskHSV,

		HSVTargetH: 300,
		HSVTolH:    30,
		HSVMinS:    .15,
		HSVMaxS:    1,
		HSVMinV:    .15,
		HSVMaxV:    1,
	}

	TargetGreenBG = imgproc.MaskConfig{
		Mode: imgproc.MaskHSV,

		HSVTargetH: 120,
		HSVTolH:    40,
		HSVMinS:    .1,
		HSVMaxS:    1,
		HSVMinV:    .1,
		HSVMaxV:    1,
	}
)

func main() {
	flag.Parse()
	if inputPath == "" || outputPath == "" {
		log.Fatal("You must specify both -in and -out paths.")
	}

	// 1) Load input image
	srcImg, err := imgproc.LoadImage(inputPath)
	if err != nil {
		log.Fatalf("Failed to load input image: %v", err)
	}

	// 2) Pad to square, centering original
	square := imgproc.PadToSquare(srcImg)

	// 3) Scale up (e.g. factor = 4)
	scaled := imgproc.ScaleImage(square, 4)

	// Convert to RGBA for pixel‐level editing
	rgba := imgproc.ImageToRGBA(scaled)

	imgproc.RemoveBackgroundSmartWithConfig(rgba, TargetMagentaBG, true)

	// 6) Fuzz border pixels
	imgproc.FuzzBorders(rgba)

	// 7) Center & crop to tight square
	final := imgproc.CenterAndCrop(rgba)

	imgproc.FuzzBorders(rgba)
	imgproc.SmoothEdges(final, 2)

	// 8) Save as PNG
	if err := imgproc.SavePNG(outputPath, imgproc.ScaleImage(final, .5)); err != nil {
		log.Fatalf("Failed to save output: %v", err)
	}

	fmt.Println("Done!")
}
