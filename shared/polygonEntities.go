package shared

import (
	"image/color"
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
