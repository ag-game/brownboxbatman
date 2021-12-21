package system

import (
	"os"

	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	moveSpeed = 1.5
)

type playerMoveSystem struct {
	player       gohan.Entity
	movement     *MovementSystem
	lastWalkDirL bool

	rewindTicks    int
	nextRewindTick int
}

func NewPlayerMoveSystem(player gohan.Entity, m *MovementSystem) *playerMoveSystem {
	return &playerMoveSystem{
		player:   player,
		movement: m,
	}
}

func (_ *playerMoveSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
		component.WeaponComponentID,
		component.SpriteComponentID,
	}
}

func (_ *playerMoveSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *playerMoveSystem) Update(ctx *gohan.Context) error {
	velocity := component.Velocity(ctx)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) && !world.World.DisableEsc {
		os.Exit(0)
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
		v := 1
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			v = 2
		}
		if world.World.Debug == v {
			world.World.Debug = 0
		} else {
			world.World.Debug = v
		}
		s.movement.UpdateDebugCollisionRects()
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyN) {
		world.World.NoClip = !world.World.NoClip
		return nil
	}

	/*if world.World.MessageVisible {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			world.World.MessageVisible = false
		}
		return nil
	}*/

	if (ebiten.IsKeyPressed(ebiten.KeyA) && !ebiten.IsKeyPressed(ebiten.KeyD)) ||
		(ebiten.IsKeyPressed(ebiten.KeyD) && !ebiten.IsKeyPressed(ebiten.KeyA)) {
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			velocity.X = -moveSpeed
		} else {
			velocity.X = moveSpeed
		}
	} else {
		velocity.X = 0
	}

	if (ebiten.IsKeyPressed(ebiten.KeyW) && !ebiten.IsKeyPressed(ebiten.KeyS)) ||
		(ebiten.IsKeyPressed(ebiten.KeyS) && !ebiten.IsKeyPressed(ebiten.KeyW)) {
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			velocity.Y = -moveSpeed
		} else {
			velocity.Y = moveSpeed
		}
	} else {
		velocity.Y = 0
	}
	return nil
}

func (s *playerMoveSystem) Draw(_ *gohan.Context, _ *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}

func deltaXY(x1, y1, x2, y2 float64) (dx float64, dy float64) {
	dx, dy = x1-x2, y1-y2
	if dx < 0 {
		dx *= -1
	}
	if dy < 0 {
		dy *= -1
	}
	return dx, dy
}
