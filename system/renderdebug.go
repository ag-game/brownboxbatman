package system

import (
	"fmt"
	"image/color"
	_ "image/png"

	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type RenderDebugTextSystem struct {
	Position *component.Position
	Velocity *component.Velocity
	Weapon   *component.Weapon

	player   gohan.Entity
	op       *ebiten.DrawImageOptions
	debugImg *ebiten.Image
}

func NewRenderDebugTextSystem(player gohan.Entity) *RenderDebugTextSystem {
	s := &RenderDebugTextSystem{
		player:   player,
		op:       &ebiten.DrawImageOptions{},
		debugImg: ebiten.NewImage(94, 114),
	}

	return s
}

func (s *RenderDebugTextSystem) Update(_ gohan.Entity) error {
	return gohan.ErrUnregister
}

func (s *RenderDebugTextSystem) Draw(e gohan.Entity, screen *ebiten.Image) error {
	if world.World.Debug <= 0 {
		return nil
	}

	position := s.Position
	velocity := s.Velocity

	s.debugImg.Fill(color.RGBA{0, 0, 0, 80})
	ebitenutil.DebugPrint(s.debugImg, fmt.Sprintf("POS %.0f,%.0f\nVEL %.2f,%.2f\nENT %d\nUPD %d\nDRA %d\nTPS %0.0f\nFPS %0.0f", position.X, position.Y, velocity.X, velocity.Y, gohan.CurrentEntities(), gohan.CurrentUpdates(), gohan.CurrentDraws(), ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	screen.DrawImage(s.debugImg, nil)
	return nil
}
