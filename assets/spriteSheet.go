package assets

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteSheet struct {
	SourceImage   *ebiten.Image
	Rows, Columns int
	Frames        []*ebiten.Image
}

func NewSpriteSheet(img *ebiten.Image, rows, columns int) (sheet *SpriteSheet) {
	sheet = &SpriteSheet{
		SourceImage: img,
		Rows:        rows,
		Columns:     columns,
	}

	var (
		frameWidth  = img.Bounds().Dx() / columns
		frameHeight = img.Bounds().Dy() / rows
	)

	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			rect := image.Rect(column*frameWidth, row*frameHeight, column*frameWidth+frameWidth, row*frameHeight+frameHeight)
			sub := img.SubImage(rect).(*ebiten.Image)
			frame := ebiten.NewImage(frameWidth, frameHeight)
			frame.DrawImage(sub, nil)
			sheet.Frames = append(sheet.Frames, frame)
		}
	}

	return
}
