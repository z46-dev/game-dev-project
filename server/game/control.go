package game

import (
	"math"

	"github.com/z46-dev/game-dev-project/util"
)

// Bound to a body. Assumes Game and Body are not nil
type Control struct {
	Game                *Game
	Body                *Ship
	Goal, PrimaryTarget *util.Vector2D
}

func NewControl(game *Game, body *Ship) (ctrl *Control) {
	ctrl = &Control{
		Game:          game,
		Body:          body,
		Goal:          util.Vector(0, 0),
		PrimaryTarget: nil,
	}

	return
}

func (c *Control) Update() {
	if c.Goal != nil && c.Goal.Magnitude() > 0 {
		var (
			angleToGoal      float64 = c.Goal.Direction()
			delta            float64 = wrapAngle(angleToGoal - c.Body.Rotation)
			speed, turnSpeed float64 = c.Body.Cfg.Speed / 120, c.Body.Cfg.TurnSpeed //math.Min(c.Body.Cfg.TurnSpeed, util.AngleDifference(c.Body.Rotation, angleToGoal))
		)

		delta = min(turnSpeed, max(-turnSpeed, delta))
		c.Body.Rotation += delta

		c.Body.Velocity.X += math.Cos(c.Body.Rotation) * speed
		c.Body.Velocity.Y += math.Sin(c.Body.Rotation) * speed
	}
}
