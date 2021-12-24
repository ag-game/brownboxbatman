package system

import (
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/brownboxbatman/entity"
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
	if world.World.MessageVisible {
		return nil
	}

	creep := component.Creep(ctx)
	position := component.Position(ctx)

	// Skip inactive creeps.
	sx, sy := world.LevelCoordinatesToScreen(position.X, position.Y)
	if sx < 0 || sy < 0 || sx > 640 || sy > 480 {
		return nil
	}

	randSpeed := func() float64 {
		return 0.5 + creep.Rand.Float64()*0.5 + (0.5 - creep.Rand.Float64())
	}

	if creep.Ticks == 0 {
		for i := 0; i < 8; i++ {
			speedA := randSpeed()
			speedB := randSpeed()

			if creep.Rand.Intn(2) == 0 {
				speedA *= -1
			}
			if creep.Rand.Intn(2) == 0 {
				speedB *= -1
			}
			entity.NewBullet(position.X, position.Y, speedA, speedB)
		}
		creep.Ticks = creep.FireRate
	}

	creep.Ticks--
	return nil
}

func (_ *CreepSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
