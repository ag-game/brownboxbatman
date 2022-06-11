package system

import (
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/brownboxbatman/entity"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	fireSpeed = 1.5
)

type playerFireSystem struct {
	Position *component.Position
	Velocity *component.Velocity
	Weapon   *component.Weapon
	Sprite   *component.Sprite
}

func NewplayerFireSystem() *playerFireSystem {
	return &playerFireSystem{}
}

func (s *playerFireSystem) Update(e gohan.Entity) error {
	if !world.World.GameStarted || world.World.GameOver {
		return nil
	}

	position := s.Position
	weapon := s.Weapon

	if ebiten.IsKeyPressed(ebiten.KeyZ) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if weapon.NextFire == 0 {
			entity.NewPlayerBullet(position.X-8, position.Y-8, 0, -weapon.BulletSpeed)
			entity.NewPlayerBullet(position.X+8, position.Y-8, 0, -weapon.BulletSpeed)
			weapon.NextFire = weapon.FireRate
		}
	}
	if weapon.NextFire > 0 {
		weapon.NextFire--
	}
	return nil
}

func (s *playerFireSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
}
