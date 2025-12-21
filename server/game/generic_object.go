package game

import (
	"github.com/z46-dev/game-dev-project/util"
)

func NewGameObject(game *Game, position *util.Vector2D, f *Faction) (o *GenericObject) {
	o = &GenericObject{}

	// Important that IDs start at 1 because the protocol uses 0 as a terminator
	// 0 = no more objects coming, anything else = object ID
	game.nextID++
	o.ID = game.nextID

	o.Game = game
	o.Position = position
	o.Velocity = util.Vector(0, 0)
	o.Size = 32
	o.Rotation = 0
	o.Friction = 0.9
	o.Density = 1
	o.Pushability = 1
	o.Faction = f
	return
}

func (o *GenericObject) GetID() (id uint64) {
	id = o.ID
	return
}

func (o *GenericObject) GetAABB() (aabb *util.AABB) {
	aabb = nil
	return
}

func (o *GenericObject) Update() {
	o.Position.Add(o.Velocity)
	o.Velocity.Scale(o.Friction)

	o.Insert()
}

func (o *GenericObject) Insert() {
	// noop
}

func (o *GenericObject) Collide() {
	// noop
}

func (o *GenericObject) PushWeight() (weight float64) {
	if o.Pushability <= 0 || o.Density <= 0 {
		weight = 0
	} else {
		weight = o.Density / o.Pushability
	}

	return
}

func (o *GenericObject) invMass() (invMass float64) {
	if o.Pushability <= 0 || o.Density <= 0 {
		invMass = 0
	} else {
		invMass = o.Pushability / o.Density
	}

	return
}
