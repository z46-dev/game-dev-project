package shared

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/z46-dev/game-dev-project/util"
	"golang.org/x/image/colornames"
)

var possibleColors = []color.Color{colornames.Rosybrown, colornames.Lightblue, colornames.Lightcoral, colornames.Lightcyan, colornames.Lightgoldenrodyellow, colornames.Lightgray, colornames.Lightgreen, colornames.Lightpink, colornames.Lightsalmon, colornames.Lightseagreen, colornames.Lightskyblue, colornames.Lightslategray, colornames.Lightsteelblue, colornames.Lightyellow}

func CreateAssetForPolygon(poly []*util.Vector2D, size float64) (img *ebiten.Image) {
	img = ebiten.NewImage(int(size), int(size))
	img.Fill(color.Transparent)

	var path *vector.Path = &vector.Path{}
	var half float32 = float32(size) / 2
	for i, point := range poly {
		var x, y float32 = float32(point.X)*half + half, float32(point.Y)*half + half

		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}

	path.Close()

	var opts *vector.DrawPathOptions = &vector.DrawPathOptions{
		AntiAlias: true,
	}

	opts.ColorScale.ScaleWithColor(possibleColors[rand.Intn(len(possibleColors))])
	vector.FillPath(img, path, &vector.FillOptions{}, opts)
	return
}

func CreateShipAsset(poly []*util.Vector2D, size float64, fill color.Color, outline color.Color) (img *ebiten.Image) {
	img = ebiten.NewImage(int(size), int(size))
	img.Fill(color.Transparent)

	var path *vector.Path = &vector.Path{}
	var half float32 = float32(size) / 2
	for i, point := range poly {
		var x, y float32 = float32(point.X)*half + half, float32(point.Y)*half + half
		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	path.Close()

	var fillOpts *vector.DrawPathOptions = &vector.DrawPathOptions{AntiAlias: true}
	fillOpts.ColorScale.ScaleWithColor(fill)
	vector.FillPath(img, path, &vector.FillOptions{}, fillOpts)

	if outline != nil {
		var strokeOpts *vector.DrawPathOptions = &vector.DrawPathOptions{AntiAlias: true}
		strokeOpts.ColorScale.ScaleWithColor(outline)
		vector.StrokePath(img, path, &vector.StrokeOptions{Width: float32(size * 0.03)}, strokeOpts)
	}

	return
}

func CreateCircleAsset(radius float64, col color.Color) (img *ebiten.Image) {
	img = ebiten.NewImage(int(radius*2), int(radius*2))
	img.Fill(color.Transparent)

	var path *vector.Path = &vector.Path{}
	var r float32 = float32(radius)
	path.MoveTo(r*2, r)
	path.Arc(r, r, r, 0, 2*math.Pi, vector.Clockwise)
	path.Close()

	var opts *vector.DrawPathOptions = &vector.DrawPathOptions{
		AntiAlias: true,
	}

	opts.ColorScale.ScaleWithColor(col)
	vector.FillPath(img, path, &vector.FillOptions{}, opts)
	return
}

func CreateRoundedRectAsset(width, height, radius float64, fill color.Color, outline color.Color, outlineWidth float64) (img *ebiten.Image) {
	img = ebiten.NewImage(int(width), int(height))
	img.Fill(color.Transparent)

	var path *vector.Path = &vector.Path{}
	var w float32 = float32(width)
	var h float32 = float32(height)
	var r float32 = float32(radius)

	path.MoveTo(r, 0)
	path.LineTo(w-r, 0)
	path.Arc(w-r, r, r, -math.Pi/2, 0, vector.Clockwise)
	path.LineTo(w, h-r)
	path.Arc(w-r, h-r, r, 0, math.Pi/2, vector.Clockwise)
	path.LineTo(r, h)
	path.Arc(r, h-r, r, math.Pi/2, math.Pi, vector.Clockwise)
	path.LineTo(0, r)
	path.Arc(r, r, r, math.Pi, 3*math.Pi/2, vector.Clockwise)
	path.Close()

	var fillOpts *vector.DrawPathOptions = &vector.DrawPathOptions{AntiAlias: true}
	fillOpts.ColorScale.ScaleWithColor(fill)
	vector.FillPath(img, path, &vector.FillOptions{}, fillOpts)

	if outline != nil && outlineWidth > 0 {
		var strokeOpts *vector.DrawPathOptions = &vector.DrawPathOptions{AntiAlias: true}
		strokeOpts.ColorScale.ScaleWithColor(outline)
		vector.StrokePath(img, path, &vector.StrokeOptions{Width: float32(outlineWidth)}, strokeOpts)
	}

	return
}

func CreateRectAsset(width, height float64, col color.Color) (img *ebiten.Image) {
	img = ebiten.NewImage(int(width), int(height))
	img.Fill(col)
	return
}
