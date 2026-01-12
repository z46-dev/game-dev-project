package imgproc

import (
	"image"
	"runtime"
	"sync"
)

// SmoothEdges applies a box‐blur to the alpha channel of img.
// radius = how many pixels to blur in each direction (radius=1 ⇒ 3×3 kernel, radius=2 ⇒ 5×5, etc).
//
// Only the alpha channel is modified. The color channels remain untouched.
func SmoothEdges(img *image.RGBA, radius int) {
	if radius < 1 {
		return
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	N := w * h

	// Extract the original alpha into a flat []uint8
	orig := make([]uint8, N)
	for y := range h {
		rowOffset := y * w
		for x := 0; x < w; x++ {
			pixIdx := img.PixOffset(x, y) + 3
			orig[rowOffset+x] = img.Pix[pixIdx]
		}
	}

	// Prepare a buffer for blurred alphas
	blurred := make([]uint8, N)

	// We'll do a single‐step box blur (sum over a (2*radius+1)² window,
	// divide by count). For simplicity we do it in one pass.
	//
	// Parallelize by rows:
	var wg sync.WaitGroup
	workers := runtime.GOMAXPROCS(0)
	rowsPer := h / workers
	if rowsPer == 0 {
		rowsPer = 1
		workers = h
	}

	for i := 0; i < workers; i++ {
		yStart := i * rowsPer
		yEnd := yStart + rowsPer
		if i == workers-1 {
			yEnd = h
		}
		wg.Add(1)

		go func(ys, ye int) {
			defer wg.Done()
			for y := ys; y < ye; y++ {
				rowOffset := y * w
				for x := range w {
					// Compute box‐sum over [y-radius..y+radius] × [x-radius..x+radius]
					var sum int
					var count int

					y0 := y - radius
					if y0 < 0 {
						y0 = 0
					}
					y1 := y + radius
					if y1 >= h {
						y1 = h - 1
					}

					x0 := x - radius
					if x0 < 0 {
						x0 = 0
					}
					x1 := x + radius
					if x1 >= w {
						x1 = w - 1
					}

					for yy := y0; yy <= y1; yy++ {
						off := yy * w
						for xx := x0; xx <= x1; xx++ {
							sum += int(orig[off+xx])
							count++
						}
					}
					avg := sum / count
					blurred[rowOffset+x] = uint8(avg)
				}
			}
		}(yStart, yEnd)
	}
	wg.Wait()

	// Write blurred alphas back into img.Pix
	for y := 0; y < h; y++ {
		rowOffset := y * w
		for x := 0; x < w; x++ {
			pixIdx := img.PixOffset(x, y) + 3
			img.Pix[pixIdx] = blurred[rowOffset+x]
		}
	}
}

func SmoothEdgesWithColor(img *image.RGBA, radius int, R, G, B uint8) {
	if radius < 1 {
		return
	}

	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	N := w * h

	// Extract the original alpha into a flat []uint8
	orig := make([]uint8, N)
	for y := 0; y < h; y++ {
		rowOffset := y * w
		for x := 0; x < w; x++ {
			pixIdx := img.PixOffset(x, y) + 3
			orig[rowOffset+x] = img.Pix[pixIdx]
		}
	}

	// Prepare a buffer for blurred alphas
	blurred := make([]uint8, N)

	// We'll do a single‐step box blur (sum over a (2*radius+1)² window,
	// divide by count). For simplicity we do it in one pass.
	//
	// Parallelize by rows:
	var wg sync.WaitGroup
	workers := runtime.GOMAXPROCS(0)
	rowsPer := h / workers
	if rowsPer == 0 {
		rowsPer = 1
		workers = h
	}

	for i := 0; i < workers; i++ {
		yStart := i * rowsPer
		yEnd := yStart + rowsPer
		if i == workers-1 {
			yEnd = h
		}
		wg.Add(1)

		go func(ys, ye int) {
			defer wg.Done()
			for y := ys; y < ye; y++ {
				rowOffset := y * w
				for x := 0; x < w; x++ {
					// Compute box‐sum over [y-radius..y+radius] × [x-radius..x+radius]
					var sum int
					var count int

					y0 := y - radius
					if y0 < 0 {
						y0 = 0
					}
					y1 := y + radius
					if y1 >= h {
						y1 = h - 1
					}

					x0 := x - radius
					if x0 < 0 {
						x0 = 0
					}

					x1 := x + radius
					if x1 >= w {
						x1 = w - 1
					}

					for yy := y0; yy <= y1; yy++ {
						off := yy * w

						for xx := x0; xx <= x1; xx++ {
							sum += int(orig[off+xx])
							count++
						}
					}

					avg := sum / count
					blurred[rowOffset+x] = uint8(avg)
				}
			}
		}(yStart, yEnd)
	}
	wg.Wait()

	// Write blurred alphas back into img.Pix
	for y := 0; y < h; y++ {
		rowOffset := y * w
		for x := 0; x < w; x++ {
			a := blurred[rowOffset+x]

			if a < 200 {
				pixIdx := img.PixOffset(x, y)
				rr, gg, bb := MixColors(R, G, B, img.Pix[pixIdx], img.Pix[pixIdx+1], img.Pix[pixIdx+2], .2)
				img.Pix[pixIdx] = rr
				img.Pix[pixIdx+1] = gg
				img.Pix[pixIdx+2] = bb
				img.Pix[pixIdx+3] = a
			} else {
				pixIdx := img.PixOffset(x, y) + 3
				img.Pix[pixIdx] = a
			}
		}
	}
}
