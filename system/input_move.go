package system

import (
	"os"

	"code.rocketnine.space/tslocum/brownboxbatman/asset"
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
	Position *component.Position
	Velocity *component.Velocity
	Weapon   *component.Weapon
	Sprite   *component.Sprite

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

func (s *playerMoveSystem) Update(e gohan.Entity) error {
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
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyN) {
		world.World.NoClip = !world.World.NoClip
		return nil
	}

	if !world.World.GameStarted {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			world.StartGame()
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		if asset.SoundLevelMusic.IsPlaying() {
			asset.SoundLevelMusic.Pause()
		} else {
			asset.SoundLevelMusic.Play()
		}
	}

	if world.World.GameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			world.World.ResetGame = true
		}
		return nil
	}

	pressLeft := ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
	pressRight := ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD)
	pressUp := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	pressDown := ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS)

	if (pressLeft && !pressRight) ||
		(pressRight && !pressLeft) {
		if pressLeft {
			s.Velocity.X = -moveSpeed
		} else {
			s.Velocity.X = moveSpeed
		}
	} else {
		s.Velocity.X = 0
	}

	if (pressUp && !pressDown) ||
		(pressDown && !pressUp) {
		if pressUp {
			s.Velocity.Y = -moveSpeed
		} else {
			s.Velocity.Y = moveSpeed
		}
	} else {
		s.Velocity.Y = 0
	}
	return nil
}

func (s *playerMoveSystem) Draw(_ gohan.Entity, _ *ebiten.Image) error {
	return gohan.ErrUnregister
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
