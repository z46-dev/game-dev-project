package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/game-dev-project/util"
	"golang.org/x/image/colornames"
)

var genericBoxImage *ebiten.Image = ebiten.NewImage(16, 16)

func newGenericObject(game *Game, position *util.Vector2D) (o *GenericObject) {
	o = &GenericObject{
		game:           game,
		id:             game.next(),
		position:       position,
		velocity:       util.Vector(1, 0),
		size:           32,
		frictionFactor: 0.99,
		aabb:           &util.AABB{},
	}

	return
}

func (o *GenericObject) ID() (id uint64) {
	id = o.id
	return
}

func (o *GenericObject) Update() {
	o.velocity.Scale(o.frictionFactor)
	o.position.Add(o.velocity)

	o.aabb.X1 = o.position.X - o.size/2
	o.aabb.Y1 = o.position.Y - o.size/2
	o.aabb.X2 = o.position.X + o.size/2
	o.aabb.Y2 = o.position.Y + o.size/2

	o.game.spatialHash.Insert(o)
}

func (o *GenericObject) Collide() {
	var collisions []*GenericObject = o.game.spatialHash.Retrieve(o.aabb)
	for _, c := range collisions {
		if o.id == c.id {
			continue
		}

		// Right now, two squares will always collide if they are touching
		var angle float64 = util.AngleBetween(o.position, c.position)
		o.velocity.Add(util.Vector(0.25, 0).Rotate(angle + math.Pi))
	}
}

func (o *GenericObject) Draw(screen *ebiten.Image) {
	if !o.game.Camera.IsInView(o.position, o.size) {
		return
	}

	genericBoxImage.Fill(colornames.Blue)

	var bounds image.Rectangle = genericBoxImage.Bounds()
	var dx, dy float64 = float64(bounds.Dx()), float64(bounds.Dy())
	var width, height float64 = o.size / dx, o.size / dy

	var options *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}

	// Object transformations
	options.GeoM.Translate(-dx/2, -dy/2)
	options.GeoM.Scale(width, height)
	options.GeoM.Rotate(o.rotation)
	options.GeoM.Translate(o.position.X, o.position.Y)

	// Camera transformations
	options.GeoM.Scale(o.game.Camera.Zoom, o.game.Camera.Zoom)
	options.GeoM.Translate(o.game.Camera.Width/2, o.game.Camera.Height/2)
	options.GeoM.Translate(-o.game.Camera.Position.X*o.game.Camera.Zoom, -o.game.Camera.Position.Y*o.game.Camera.Zoom)

	// Graphical improvements
	options.Filter = ebiten.FilterLinear
	options.DisableMipmaps = false

	screen.DrawImage(genericBoxImage, options)
}

func (o *GenericObject) GetAABB() (aabb *util.AABB) {
	aabb = o.aabb
	return
}

func (o *GenericObject) Destroy() {
	// noop
}
