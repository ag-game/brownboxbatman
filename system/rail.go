package system

import (
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type RailSystem struct {
}

func NewRailSystem() *RailSystem {
	s := &RailSystem{}

	return s
}
func (_ *RailSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.RailComponentID,
		component.PositionComponentID,
	}
}

func (_ *RailSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *RailSystem) Update(ctx *gohan.Context) error {
	if world.World.MessageVisible || world.World.GameOver || !world.World.CamMoving {
		return nil
	}

	position := component.Position(ctx)
	position.Y -= CameraMoveSpeed
	return nil
}

func (_ *RailSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}
