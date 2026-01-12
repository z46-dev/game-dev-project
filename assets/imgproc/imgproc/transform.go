package imgproc

import (
	"image"
	"image/color"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

// PadToSquare returns a new RGBA image whose dimensions are square = max(width, height).
// The original is centered, and empty space is fully transparent.
func PadToSquare(src image.Image) *image.RGBA {
	srcBounds := src.Bounds()
	w, h := srcBounds.Dx(), srcBounds.Dy()
	size := w
	if h > size {
		size = h
	}
	dst := image.NewRGBA(image.Rect(0, 0, size, size))
	// Fill transparent
	draw.Draw(dst, dst.Bounds(), &image.Uniform{C: color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)

	offsetX := (size - w) / 2
	offsetY := (size - h) / 2
	draw.Draw(dst,
		image.Rect(offsetX, offsetY, offsetX+w, offsetY+h),
		src, srcBounds.Min, draw.Over)
	return dst
}

// ScaleImage scales src by factor f (e.g. 2.0, 4.0) using bilinear interpolation.
func ScaleImage(src *image.RGBA, f float64) *image.RGBA {
	srcW := src.Bounds().Dx()
	srcH := src.Bounds().Dy()
	dstW := int(float64(srcW) * f)
	dstH := int(float64(srcH) * f)
	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	xdraw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// ImageToRGBA copies any image.Image into a new *image.RGBA.
func ImageToRGBA(img image.Image) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, b, img, b.Min, draw.Src)
	return rgba
}
