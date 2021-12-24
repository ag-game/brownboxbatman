package component

import (
	"math/rand"

	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type CreepComponent struct {
	Health     int
	FireAmount int
	FireRate   int // In ticks
	Ticks      int // Ticks until next action
	Rand       *rand.Rand
}

var CreepComponentID = ECS.NewComponentID()

func (p *CreepComponent) ComponentID() gohan.ComponentID {
	return CreepComponentID
}

func Creep(ctx *gohan.Context) *CreepComponent {
	c, ok := ctx.Component(CreepComponentID).(*CreepComponent)
	if !ok {
		return nil
	}
	return c
}
