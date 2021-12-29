package component

import (
	"math/rand"

	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type CreepComponent struct {
	Active bool

	Health     int
	FireAmount int
	FireRate   int // In ticks
	FireTicks  int //Ticks until next action

	Movement      int
	Movements     [][3]float64 // X, Y, pre-delay in ticks
	MovementTicks int          // Ticks until next action

	DamageTicks int

	Rand *rand.Rand
}

const (
	CreepSnowblower = iota + 1
	CreepSnowmanHead
)

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
