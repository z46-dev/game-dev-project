package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/z46-dev/game-dev-project/shared/definitions"
	"golang.org/x/image/colornames"
)

const projectileAssetSize float64 = 64

func projectileAssetForID(id definitions.ProjectileID) *ebiten.Image {
	switch id {
	case definitions.PROJECTILE_LIGHT_LASER:
		return lightLaserAsset
	case definitions.PROJECTILE_HEAVY_LASER:
		return heavyLaserAsset
	case definitions.PROJECTILE_PULSE_EMITTER:
		return pulseEmitterAsset
	case definitions.PROJECTILE_LIGHT_MISSILE:
		return lightMissileAsset
	default:
		return fallbackProjectileAsset
	}
}

var (
	lightLaserAsset         *ebiten.Image = createLaserAsset(projectileAssetSize, 0.9, 0.12, colornames.Lightcyan, colornames.White)
	heavyLaserAsset         *ebiten.Image = createLaserAsset(projectileAssetSize, 0.85, 0.2, colornames.Orange, colornames.Mistyrose)
	pulseEmitterAsset       *ebiten.Image = createPulseAsset(projectileAssetSize, colornames.Gold, colornames.Orange)
	lightMissileAsset       *ebiten.Image = createMissileAsset(projectileAssetSize, colornames.Slategray, colornames.Lightgray, colornames.Orangered)
	fallbackProjectileAsset *ebiten.Image = createPulseAsset(projectileAssetSize, colornames.Lightyellow, colornames.Gold)
)

func createLaserAsset(size, lengthRatio, thicknessRatio float64, glow color.Color, core color.Color) *ebiten.Image {
	img := ebiten.NewImage(int(size), int(size))
	img.Fill(color.Transparent)

	length := size * lengthRatio
	thickness := size * thicknessRatio
	x := (size - length) / 2
	y := (size - thickness) / 2
	vector.FillRect(img, float32(x), float32(y), float32(length), float32(thickness), glow, true)

	coreLength := length * 0.85
	coreThickness := thickness * 0.5
	cx := (size - coreLength) / 2
	cy := (size - coreThickness) / 2
	vector.FillRect(img, float32(cx), float32(cy), float32(coreLength), float32(coreThickness), core, true)

	return img
}

func createPulseAsset(size float64, glow color.Color, core color.Color) *ebiten.Image {
	img := ebiten.NewImage(int(size), int(size))
	img.Fill(color.Transparent)

	drawCircle(img, size, size*0.45, glow)
	drawCircle(img, size, size*0.25, core)
	return img
}

func createMissileAsset(size float64, body color.Color, trim color.Color, noseColor color.Color) *ebiten.Image {
	img := ebiten.NewImage(int(size), int(size))
	img.Fill(color.Transparent)

	bodyLength := size * 0.6
	bodyHeight := size * 0.18
	bodyX := (size - bodyLength) / 2
	bodyY := (size - bodyHeight) / 2
	vector.FillRect(img, float32(bodyX), float32(bodyY), float32(bodyLength), float32(bodyHeight), body, true)

	finWidth := size * 0.1
	finHeight := size * 0.08
	finX := bodyX + bodyLength*0.15
	vector.FillRect(img, float32(finX), float32(bodyY-finHeight), float32(finWidth), float32(finHeight), trim, true)
	vector.FillRect(img, float32(finX), float32(bodyY+bodyHeight), float32(finWidth), float32(finHeight), trim, true)

	nosePath := &vector.Path{}
	noseX := bodyX + bodyLength
	noseY := size / 2
	noseSize := size * 0.15
	nosePath.MoveTo(float32(noseX), float32(noseY-noseSize))
	nosePath.LineTo(float32(noseX+noseSize), float32(noseY))
	nosePath.LineTo(float32(noseX), float32(noseY+noseSize))
	nosePath.Close()

	noseOpts := &vector.DrawPathOptions{AntiAlias: true}
	noseOpts.ColorScale.ScaleWithColor(noseColor)
	vector.FillPath(img, nosePath, &vector.FillOptions{}, noseOpts)

	return img
}

func drawCircle(img *ebiten.Image, size float64, radius float64, col color.Color) {
	path := &vector.Path{}
	r := float32(radius)
	center := float32(size / 2)
	path.MoveTo(center+r, center)
	path.Arc(center, center, r, 0, 2*math.Pi, vector.Clockwise)
	path.Close()

	opts := &vector.DrawPathOptions{AntiAlias: true}
	opts.ColorScale.ScaleWithColor(col)
	vector.FillPath(img, path, &vector.FillOptions{}, opts)
}
