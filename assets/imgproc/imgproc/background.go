// internal/imgproc/background.go
package imgproc

import (
	"image"
	"math"
	"runtime"
	"sync"
)

type point struct{ x, y int }

// RemoveBackgroundParallel clears (alpha=0) any pixel whose colorDistance to (tgR, tgG, tgB) ≤ tol.
// Parallelized by rows.
func RemoveBackgroundParallel(img *image.RGBA, tgR, tgG, tgB, tol float64) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

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
				for x := 0; x < w; x++ {
					idx := img.PixOffset(x, y)
					r, g, b_ := float64(img.Pix[idx]), float64(img.Pix[idx+1]), float64(img.Pix[idx+2])
					if ColorDistance(r, g, b_, tgR, tgG, tgB) <= tol {
						img.Pix[idx+3] = 0
					}
				}
			}
		}(yStart, yEnd)
	}
	wg.Wait()
}

// RemoveBackgroundIterative does a flood-fill from any remaining pixel whose colorDistance ≤ tol.
// edgeCleanup layers beyond each border pixel are also cleared if they match tol*2.
func RemoveBackgroundIterative(img *image.RGBA, tgR, tgG, tgB, tol float64, edgeCleanup int) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	visited := make([]bool, w*h)

	isBg := func(x, y int) bool {
		idx := img.PixOffset(x, y)
		r, g, b_ := float64(img.Pix[idx]), float64(img.Pix[idx+1]), float64(img.Pix[idx+2])
		return ColorDistance(r, g, b_, tgR, tgG, tgB) <= tol
	}

	var floodFill func(int, int, *sync.Map) = func(sx, sy int, borderMap *sync.Map) {
		stack := []point{{sx, sy}}
		for len(stack) > 0 {
			p := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			ix := p.y*w + p.x
			if visited[ix] {
				continue
			}
			visited[ix] = true
			if !isBg(p.x, p.y) {
				continue
			}
			idx := img.PixOffset(p.x, p.y)
			img.Pix[idx+3] = 0 // transparent

			for _, d := range [][2]int{
				{-1, 0}, {+1, 0}, {0, -1}, {0, +1},
				{-1, -1}, {-1, +1}, {+1, -1}, {+1, +1},
			} {
				nx, ny := p.x+d[0], p.y+d[1]
				if nx < 0 || nx >= w || ny < 0 || ny >= h {
					continue
				}
				nIx := ny*w + nx
				if !isBg(nx, ny) {
					borderMap.Store(ix, struct{}{})
				} else if !visited[nIx] {
					stack = append(stack, point{nx, ny})
				}
			}
		}
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			ix := y*w + x
			if visited[ix] {
				continue
			}
			if isBg(x, y) {
				var borderMap sync.Map
				floodFill(x, y, &borderMap)
				borderMap.Range(func(key, _ interface{}) bool {
					bi := key.(int)
					px, py := bi%w, bi/w
					for dy := -edgeCleanup; dy <= edgeCleanup; dy++ {
						for dx := -edgeCleanup; dx <= edgeCleanup; dx++ {
							nx, ny := px+dx, py+dy
							if nx < 0 || nx >= w || ny < 0 || ny >= h {
								continue
							}
							nIdx := img.PixOffset(nx, ny)
							r, g, b_ := float64(img.Pix[nIdx]), float64(img.Pix[nIdx+1]), float64(img.Pix[nIdx+2])
							a := img.Pix[nIdx+3]
							if a != 0 && ColorDistance(r, g, b_, tgR, tgG, tgB) <= tol*2 {
								img.Pix[nIdx+3] = 0
							}
						}
					}
					return true
				})
			}
		}
	}
}

func Recolor(img *image.RGBA, r, g, b uint8, tolerance float64) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	var wg sync.WaitGroup
	workers := runtime.GOMAXPROCS(0)
	rowsPer := h / workers
	if rowsPer == 0 {
		rowsPer = 1
		workers = h
	}

	for i := range workers {
		yStart := i * rowsPer
		yEnd := yStart + rowsPer
		if i == workers-1 {
			yEnd = h
		}
		wg.Add(1)
		go func(ys, ye int) {
			defer wg.Done()
			for y := ys; y < ye; y++ {
				for x := range w {
					idx := img.PixOffset(x, y)

					R, G, B := MixColors(img.Pix[idx], img.Pix[idx+1], img.Pix[idx+2], r, g, b, tolerance)
					img.Pix[idx] = R
					img.Pix[idx+1] = G
					img.Pix[idx+2] = B
				}
			}
		}(yStart, yEnd)
	}
	wg.Wait()
}

// FuzzBorders thresholds α < 225 → 0, ≥ 225 → 255 in parallel, then
// any still‐opaque pixel with ≥2 transparent neighbors becomes transparent.
func FuzzBorders(img *image.RGBA) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	var wg sync.WaitGroup
	workers := runtime.GOMAXPROCS(0)
	rowsPer := h / workers
	if rowsPer == 0 {
		rowsPer = 1
		workers = h
	}

	// 1) Threshold α in parallel
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
				for x := 0; x < w; x++ {
					aIdx := img.PixOffset(x, y) + 3
					if img.Pix[aIdx] < 225 {
						img.Pix[aIdx] = 0
					} else {
						img.Pix[aIdx] = 255
					}
				}
			}
		}(yStart, yEnd)
	}
	wg.Wait()

	// 2) Collect “border pixels” to clear
	toClear := make([]int, 0)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			baseIdx := img.PixOffset(x, y)
			if img.Pix[baseIdx+3] == 0 {
				continue
			}
			countTrans := 0
			for _, d := range [][2]int{
				{-1, 0}, {+1, 0}, {0, -1}, {0, +1},
				{-1, -1}, {-1, +1}, {+1, -1}, {+1, +1},
			} {
				nx, ny := x+d[0], y+d[1]
				if nx < 0 || nx >= w || ny < 0 || ny >= h {
					continue
				}
				if img.Pix[img.PixOffset(nx, ny)+3] == 0 {
					countTrans++
				}
			}
			if countTrans >= 2 {
				toClear = append(toClear, baseIdx+3)
			}
		}
	}
	// 3) Clear them
	for _, aidx := range toClear {
		img.Pix[aidx] = 0
	}
}

// MaskMode selects how to build the initial boolean mask.
type MaskMode int

const (
	MaskRGB MaskMode = iota
	MaskHSV
	MaskEnsemble
)

// MaskConfig contains parameters for each masking strategy.
type MaskConfig struct {
	Mode MaskMode

	// For RGB thresholding:
	RGBTargetR, RGBTargetG, RGBTargetB float64
	RGBTolR, RGBTolG, RGBTolB          float64

	// For HSV thresholding:
	HSVTargetH, HSVTolH float64 // Hue in [0,360), tolerance in degrees
	HSVMinS, HSVMaxS    float64 // Saturation in [0,1]
	HSVMinV, HSVMaxV    float64 // Value in [0,1] (if you want to exclude too‐dark or too‐bright)
}

// RemoveBackgroundSmartWithConfig builds a mask according to cfg,
// applies a 1‐pixel morphological opening (erode→dilate) if desired,
// then flood‐fills from edges to clear the connected background.
// RemoveBackgroundSmartWithConfig builds a mask according to cfg,
// applies a 1‐pixel morphological opening (erode→dilate),
// then flood‐fills from edges to clear the connected background,
// and finally clears any remaining “interior” mask pixels (donut holes).
func RemoveBackgroundSmartWithConfig(
	img *image.RGBA,
	cfg MaskConfig,
	clearInterior bool,
) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	// 1) Build initial mask[]bool in parallel
	N := w * h
	mask := make([]bool, N)

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
				for x := 0; x < w; x++ {
					idx := img.PixOffset(x, y)
					R := float64(img.Pix[idx+0])
					G := float64(img.Pix[idx+1])
					B := float64(img.Pix[idx+2])

					var ok bool
					switch cfg.Mode {
					case MaskRGB:
						ok = maskPixelRGB(R, G, B,
							cfg.RGBTargetR, cfg.RGBTargetG, cfg.RGBTargetB,
							cfg.RGBTolR, cfg.RGBTolG, cfg.RGBTolB,
						)

					case MaskHSV:
						hh, ss, vv := rgbToHSV(R, G, B)
						ok = maskPixelHSV(hh, ss, vv,
							cfg.HSVTargetH, cfg.HSVTolH,
							cfg.HSVMinS, cfg.HSVMaxS,
							cfg.HSVMinV, cfg.HSVMaxV,
						)

					case MaskEnsemble:
						// Require BOTH RGB and HSV conditions:
						if maskPixelRGB(R, G, B,
							cfg.RGBTargetR, cfg.RGBTargetG, cfg.RGBTargetB,
							cfg.RGBTolR, cfg.RGBTolG, cfg.RGBTolB,
						) {
							hh, ss, vv := rgbToHSV(R, G, B)
							ok = maskPixelHSV(hh, ss, vv,
								cfg.HSVTargetH, cfg.HSVTolH,
								cfg.HSVMinS, cfg.HSVMaxS,
								cfg.HSVMinV, cfg.HSVMaxV,
							)
						} else {
							ok = false
						}
					}

					if ok {
						mask[y*w+x] = true
					}
				}
			}
		}(yStart, yEnd)
	}
	wg.Wait()

	// 2) Apply one‐pixel morphological opening: erode → dilate.
	mask = erodeMask(mask, w, h)
	mask = dilateMask(mask, w, h)

	// 3) Flood‐fill from any edge pixel where mask==true (single‐threaded)
	visited := make([]bool, N)
	queue := make([]point, 0, N)

	// enqueue top & bottom edges
	for x := 0; x < w; x++ {
		if mask[x] {
			queue = append(queue, point{x, 0})
		}
		if mask[(h-1)*w+x] {
			queue = append(queue, point{x, h - 1})
		}
	}
	// enqueue left & right edges (excluding corners)
	for y := 1; y < h-1; y++ {
		if mask[y*w] {
			queue = append(queue, point{0, y})
		}
		if mask[y*w+w-1] {
			queue = append(queue, point{w - 1, y})
		}
	}

	for head := 0; head < len(queue); head++ {
		p := queue[head]
		ix := p.y*w + p.x
		if visited[ix] {
			continue
		}
		visited[ix] = true
		// Clear alpha channel
		pixIdx := img.PixOffset(p.x, p.y)
		img.Pix[pixIdx+3] = 0

		// visit 8 neighbors
		for _, d := range [][2]int{
			{-1, 0}, {+1, 0}, {0, -1}, {0, +1},
			{-1, -1}, {-1, +1}, {+1, -1}, {+1, +1},
		} {
			nx, ny := p.x+d[0], p.y+d[1]
			if nx < 0 || nx >= w || ny < 0 || ny >= h {
				continue
			}
			nIdx := ny*w + nx
			if mask[nIdx] && !visited[nIdx] {
				queue = append(queue, point{nx, ny})
			}
		}
	}

	if clearInterior {
		// 4) Clear any “interior” mask pixels (donut holes):
		for y := 0; y < h; y++ {
			rowOff := y * w
			for x := 0; x < w; x++ {
				i := rowOff + x
				if mask[i] {
					pixIdx := img.PixOffset(x, y)
					img.Pix[pixIdx+3] = 0
				}
			}
		}
	}
}

// maskPixelRGB returns true if abs(R–tR) ≤ tolR, abs(G–tG) ≤ tolG, abs(B–tB) ≤ tolB.
func maskPixelRGB(
	R, G, B,
	tR, tG, tB float64,
	tolR, tolG, tolB float64,
) bool {
	if math.Abs(R-tR) > tolR {
		return false
	}
	if math.Abs(G-tG) > tolG {
		return false
	}
	if math.Abs(B-tB) > tolB {
		return false
	}
	return true
}

// maskPixelHSV returns true if hue ∈ [tH–tolH, tH+tolH], saturation ∈ [minS,maxS], value ∈ [minV,maxV].
func maskPixelHSV(
	H, S, V,
	tH, tolH float64,
	minS, maxS, minV, maxV float64,
) bool {
	// hue distance (circular)
	dh := math.Mod(math.Abs(H-tH), 360.0)
	if dh > 180 {
		dh = 360 - dh
	}
	if dh > tolH {
		return false
	}
	if S < minS || S > maxS {
		return false
	}
	if V < minV || V > maxV {
		return false
	}
	return true
}

// rgbToHSV converts R,G,B in [0..255] to H∈[0..360), S∈[0..1], V∈[0..1].
func rgbToHSV(r, g, b float64) (h, s, v float64) {
	R := r / 255.0
	G := g / 255.0
	B := b / 255.0

	max := math.Max(R, math.Max(G, B))
	min := math.Min(R, math.Min(G, B))
	delta := max - min

	// Hue
	switch {
	case delta == 0:
		h = 0
	case max == R:
		h = math.Mod(((G-B)/delta), 6) * 60
	case max == G:
		h = ((B-R)/delta + 2) * 60
	case max == B:
		h = ((R-G)/delta + 4) * 60
	}
	if h < 0 {
		h += 360
	}

	// Saturation
	if max == 0 {
		s = 0
	} else {
		s = delta / max
	}

	// Value
	v = max

	return
}

// erodeMask returns a new mask[] after one‐pixel erosion: only keep a pixel if
// all 8‐neighbors (and itself) are true in the original.
func erodeMask(mask []bool, w, h int) []bool {
	out := make([]bool, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := y*w + x
			if !mask[i] {
				continue
			}
			keep := true
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						keep = false
						break
					}
					if !mask[ny*w+nx] {
						keep = false
						break
					}
				}
				if !keep {
					break
				}
			}
			if keep {
				out[i] = true
			}
		}
	}
	return out
}

// dilateMask returns a new mask[] after one‐pixel dilation: if any neighbor (or itself) is true,
// then the output pixel becomes true.
func dilateMask(mask []bool, w, h int) []bool {
	out := make([]bool, w*h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := y*w + x
			any := false
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= w || ny < 0 || ny >= h {
						continue
					}
					if mask[ny*w+nx] {
						any = true
						break
					}
				}
				if any {
					break
				}
			}
			if any {
				out[i] = true
			}
		}
	}
	return out
}
