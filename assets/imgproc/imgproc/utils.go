package imgproc

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// ColorDistance computes Euclidean distance between two RGB triples.
func ColorDistance(r1, g1, b1, r2, g2, b2 float64) float64 {
	dr := r1 - r2
	dg := g1 - g2
	db := b1 - b2
	return math.Sqrt(dr*dr + dg*dg + db*db)
}

// CenterAndCrop finds the tight bounding box of all pixels with alpha>0,
// then returns a new square RGBA whose side = max(bboxW, bboxH), centering
// the cropped region with transparent padding.
func CenterAndCrop(img *image.RGBA) *image.RGBA {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	minX, minY := w, h
	maxX, maxY := 0, 0
	found := false

	// Find bounding box of alpha>0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if img.Pix[img.PixOffset(x, y)+3] != 0 {
				found = true
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// If fully transparent, return a 1×1 transparent image
	if !found {
		one := image.NewRGBA(image.Rect(0, 0, 1, 1))
		one.Pix[3] = 0
		return one
	}

	bboxW := maxX - minX + 1
	bboxH := maxY - minY + 1
	sz := bboxW
	if bboxH > sz {
		sz = bboxH
	}

	out := image.NewRGBA(image.Rect(0, 0, sz, sz))
	// Fill transparent
	draw.Draw(out, out.Bounds(), &image.Uniform{C: color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)

	// Offsets to center the cropped region
	offX, offY := 0, 0
	if bboxW < bboxH {
		offX = (bboxH - bboxW) / 2
	}
	if bboxH < bboxW {
		offY = (bboxW - bboxH) / 2
	}

	srcRect := image.Rect(minX, minY, minX+bboxW, minY+bboxH)
	draw.Draw(out,
		image.Rect(offX, offY, offX+bboxW, offY+bboxH),
		img, srcRect.Min, draw.Over)

	return out
}

func RGBToHue(r, g, b float64) float64 {
	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	delta := max - min

	if delta == 0 {
		return 0
	}

	var hue float64
	if max == r {
		hue = (g - b) / delta
	} else if max == g {
		hue = 2 + (b-r)/delta
	} else {
		hue = 4 + (r-g)/delta
	}

	hue *= 60
	if hue < 0 {
		hue += 360
	}

	return hue
}

func lerp(a, b float64, weight float64) float64 {
	return a + (b-a)*weight
}

func MixColors(r1, g1, b1, r2, g2, b2 uint8, weight float64) (uint8, uint8, uint8) {
	// return uint8(float64(r1) + float64(r2-r1)*weight), uint8(float64(g1) + float64(g2-g1)*weight), uint8(float64(b1) + float64(b2-b1)*weight)
	return uint8(lerp(float64(r1), float64(r2), weight)),
		uint8(lerp(float64(g1), float64(g2), weight)),
		uint8(lerp(float64(b1), float64(b2), weight))
}
