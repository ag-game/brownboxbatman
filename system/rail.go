package system

import (
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

type RailSystem struct {
	Rail     *component.Rail
	Position *component.Position
}

func NewRailSystem() *RailSystem {
	s := &RailSystem{}

	return s
}

func (s *RailSystem) Update(e gohan.Entity) error {
	if !world.World.GameStarted || world.World.GameOver || !world.World.CamMoving {
		return nil
	}

	s.Position.Y -= CameraMoveSpeed
	return nil
}

func (_ *RailSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}
