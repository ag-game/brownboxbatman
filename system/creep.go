package system

import (
	"code.rocketnine.space/tslocum/brownboxbatman/entity"

	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type CreepSystem struct {
}

func NewCreepSystem() *CreepSystem {
	s := &CreepSystem{}

	return s
}
func (_ *CreepSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.CreepComponentID,
		component.PositionComponentID,
	}
}

func (_ *CreepSystem) Uses() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.WeaponComponentID,
	}
}

func (s *CreepSystem) Update(ctx *gohan.Context) error {
	if world.World.MessageVisible || world.World.GameOver {
		return nil
	}

	creep := component.Creep(ctx)
	position := component.Position(ctx)

	sx, sy := world.LevelCoordinatesToScreen(position.X, position.Y)
	// TODO activate on visible
	_, _ = sx, sy

	if creep.Ticks == 0 {
		entity.NewBullet(position.X, position.Y, -0.5, 0)
		creep.Ticks = creep.FireRate
	}
	creep.Ticks--
	return nil
}

func (_ *CreepSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
