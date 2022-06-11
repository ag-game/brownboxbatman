package system

import (
	"math"

	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const CameraMoveSpeed = 0.132

type CameraSystem struct {
	Weapon   *component.Weapon
	Position *component.Position
}

func NewCameraSystem() *CameraSystem {
	s := &CameraSystem{}

	return s
}

func (s *CameraSystem) Update(e gohan.Entity) error {
	if !world.World.GameStarted || world.World.GameOver {
		return nil
	}

	world.World.CamMoving = world.World.CamY > 0
	if world.World.CamMoving {
		world.World.CamY -= CameraMoveSpeed
	} else {
		world.SetMessage("GAME OVER\n\nYOU  WIN!", math.MaxInt)
		world.World.GameOver = true
	}
	return nil
}

func (_ *CameraSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}
