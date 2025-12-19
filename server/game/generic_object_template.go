package game

import (
	"github.com/z46-dev/game-dev-project/util"
)

func (o *GenericObjectTemplate) GetID() (id uint64) {
	id = o.ID
	return
}

func (o *GenericObjectTemplate) GetAABB() (aabb *util.AABB) {
	aabb = nil
	return
}

func (o *GenericObjectTemplate) Insert() {
	// noop
}

func (o *GenericObjectTemplate) Collide() {
	// noop
}
