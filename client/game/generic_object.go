package game

import (
	"image"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/z46-dev/game-dev-project/shared"
	"github.com/z46-dev/game-dev-project/util"
)

func genPolySides(n int) (sides []*util.Vector2D) {
	for i := range n {
		var angle float64 = 2 * math.Pi / float64(n) * float64(i)
		sides = append(sides, util.Vector(math.Cos(angle), math.Sin(angle)))
	}

	return
}

func genStarSides(n int, radMul float64) (sides []*util.Vector2D) {
	n *= 2
	for i := range n {
		var angle float64 = 2 * math.Pi / float64(n) * float64(i)
		var radius float64 = 1
		if i%2 == 0 {
			radius *= radMul
		}

		sides = append(sides, util.Vector(math.Cos(angle)*radius, math.Sin(angle)*radius))
	}

	return
}

var npcShapes [][]*util.Vector2D = [][]*util.Vector2D{
	genPolySides(3),
	genPolySides(4),
	genPolySides(5),
	genPolySides(6),
	genStarSides(3, .25),
	genStarSides(4, .5),
	genStarSides(5, .75),
}

func newGenericObject(game *Game) (o *GenericObject) {
	o = &GenericObject{
		game:           game,
		id:             game.next(),
		position:       util.Vector(0, 0),
		velocity:       util.Vector(0, 0),
		size:           32,
		frictionFactor: 0.95,
		density:        1,
		pushability:    1,
		elasticity:     0.2,
		rotation:       rand.Float64() * 2 * math.Pi,
		polygon:        util.NewPolygon(npcShapes[rand.IntN(len(npcShapes))], util.Vector(0, 0), 1, 0),
	}

	o.asset = shared.CreateAssetForPolygon(o.polygon, 1024)
	return
}

func (o *GenericObject) Spawn(position *util.Vector2D) *GenericObject {
	o.position = position.Copy()
	o.polygon.Transform(o.position, o.size/2, o.rotation)
	return o
}

func (o *GenericObject) SafelySpawn(posGetter func() *util.Vector2D, maxTries int) *GenericObject {
	for range maxTries {
		o.position = posGetter()
		o.polygon.Transform(o.position, o.size, o.rotation)

		var safe = true
		var potential []*GenericObject
		if potential = o.game.spatialHash.Retrieve(o.polygon.AABB); len(potential) != 0 {
			for _, obj := range potential {
				if util.TwoPolygonsIntersect(o.polygon, obj.polygon) {
					safe = false
					break
				}
			}
		}

		if safe {
			break
		}
	}

	return o
}

func (o *GenericObject) ID() (id uint64) {
	id = o.id
	return
}

func (o *GenericObject) Update() {
	o.velocity.Scale(o.frictionFactor)
	o.position.Add(o.velocity)

	o.Insert()
}

func (o *GenericObject) Insert() {
	o.polygon.Transform(o.position, o.size/2, o.rotation)
	o.game.spatialHash.Insert(o)
}

func (o *GenericObject) Collide() {
	var collisions []*GenericObject = o.game.spatialHash.Retrieve(o.polygon.AABB)
	for _, c := range collisions {
		if o.id == c.id || !util.TwoPolygonsIntersect(o.polygon, c.polygon) {
			continue
		}

		if o.pushability != 0 || c.pushability != 0 {
			o.resolveCollision(c)
		}
	}
}

func (o *GenericObject) resolveCollision(c *GenericObject) {
	var prevO *util.Vector2D = o.position.Copy().Subtract(o.velocity)
	var prevC *util.Vector2D = c.position.Copy().Subtract(c.velocity)
	if o.velocity.SquaredMagnitude() == 0 && c.velocity.SquaredMagnitude() == 0 {
		if resolution := util.ResolveTwoPolygons(o.polygon, c.polygon); resolution != nil {
			var normal *util.Vector2D = resolution.Copy().Normalize()
			const nudge float64 = 0.05
			if o.pushability > 0 {
				o.velocity.Add(normal.Copy().Scale(nudge))
			}
			if c.pushability > 0 {
				c.velocity.Subtract(normal.Copy().Scale(nudge))
			}
		}

		for i := 0; i < 4; i++ {
			if !util.TwoPolygonsIntersect(o.polygon, c.polygon) {
				break
			}
			o.applyMTVResolution(c)
		}
		return
	}

	if o.polygonsIntersectAt(c, prevO, prevC) {
		var curO *util.Vector2D = o.position.Copy()
		var curC *util.Vector2D = c.position.Copy()
		o.polygon.Transform(curO, o.size/2, o.rotation)
		c.polygon.Transform(curC, c.size/2, c.rotation)
		o.applyMTVResolution(c)
		return
	}

	var lo, hi float64 = 0, 1
	for range 12 {
		var mid float64 = (lo + hi) / 2
		var posO *util.Vector2D = prevO.Copy().Add(o.velocity.Copy().Scale(mid))
		var posC *util.Vector2D = prevC.Copy().Add(c.velocity.Copy().Scale(mid))

		if o.polygonsIntersectAt(c, posO, posC) {
			hi = mid
			continue
		}

		lo = mid
	}

	o.position = prevO.Copy().Add(o.velocity.Copy().Scale(hi))
	c.position = prevC.Copy().Add(c.velocity.Copy().Scale(hi))
	o.polygon.Transform(o.position, o.size/2, o.rotation)
	c.polygon.Transform(c.position, c.size/2, c.rotation)
	o.applyMTVResolution(c)
}

func (o *GenericObject) polygonsIntersectAt(c *GenericObject, posO, posC *util.Vector2D) (intersects bool) {
	o.polygon.Transform(posO, o.size/2, o.rotation)
	c.polygon.Transform(posC, c.size/2, c.rotation)
	intersects = util.TwoPolygonsIntersect(o.polygon, c.polygon)
	return
}

func (o *GenericObject) applyMTVResolution(c *GenericObject) {
	var (
		weightO float64 = o.pushWeight()
		weightC float64 = c.pushWeight()
		total   float64 = weightO + weightC
		lastMTV *util.Vector2D
	)

	if total == 0 {
		return
	}

	for i := 0; i < 4; i++ {
		var resolution *util.Vector2D = util.ResolveTwoPolygons(o.polygon, c.polygon)
		if resolution == nil {
			break
		}

		var moveO *util.Vector2D = resolution.Copy().Scale(weightO / total)
		var moveC *util.Vector2D = resolution.Copy().Scale(weightC / total)
		o.position.Add(moveO)
		c.position.Subtract(moveC)
		o.polygon.Transform(o.position, o.size/2, o.rotation)
		c.polygon.Transform(c.position, c.size/2, c.rotation)

		lastMTV = resolution
		if !util.TwoPolygonsIntersect(o.polygon, c.polygon) {
			break
		}
	}

	if lastMTV != nil {
		o.applyElasticity(c, lastMTV)
	}
}

func (o *GenericObject) pushWeight() (weight float64) {
	if o.pushability <= 0 {
		return 0
	}

	if o.density <= 0 {
		return 0
	}

	return o.pushability / o.density
}

func (o *GenericObject) invMass() (invMass float64) {
	if o.pushability <= 0 || o.density <= 0 {
		return 0
	}

	return o.pushability / o.density
}

func (o *GenericObject) applyElasticity(c *GenericObject, resolution *util.Vector2D) {
	var invMassO float64 = o.invMass()
	var invMassC float64 = c.invMass()
	var totalInvMass float64 = invMassO + invMassC
	if totalInvMass == 0 {
		return
	}

	var normal *util.Vector2D = resolution.Copy().Normalize()
	var relVel *util.Vector2D = o.velocity.Copy().Subtract(c.velocity)
	var velAlongNormal float64 = relVel.Dot(normal)
	if velAlongNormal >= 0 {
		return
	}

	var e float64 = math.Min(o.elasticity, c.elasticity)
	var impulse float64 = -(1 + e) * velAlongNormal / totalInvMass

	o.velocity.Add(normal.Copy().Scale(impulse * invMassO))
	c.velocity.Subtract(normal.Copy().Scale(impulse * invMassC))
}

func (o *GenericObject) Draw(screen *ebiten.Image) {
	if !o.game.Camera.IsInView(o.position, o.size) {
		return
	}

	var bounds image.Rectangle = o.asset.Bounds()
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

	screen.DrawImage(o.asset, options)
}

func (o *GenericObject) GetAABB() (aabb *util.AABB) {
	aabb = o.polygon.AABB
	return
}

func (o *GenericObject) Destroy() {
	// noop
}
