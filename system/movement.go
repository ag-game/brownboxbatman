package system

import (
	"image"
	"image/color"
	"math"
	"time"

	"code.rocketnine.space/tslocum/brownboxbatman/component"
	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const rewindThreshold = 1

type MovementSystem struct {
	ScreenW, ScreenH float64

	OnGround int
	OnLadder int

	Jumping  bool
	LastJump time.Time

	Dashing  bool
	LastDash time.Time

	collisionRects []image.Rectangle

	ladderRects []image.Rectangle

	fireRects []image.Rectangle

	debugCollisionRects []gohan.Entity
	debugLadderRects    []gohan.Entity
	debugFireRects      []gohan.Entity

	playerPositions     [][2]float64
	playerPosition      [2]float64
	playerPositionTicks int
	recordedPosition    bool
}

func NewMovementSystem() *MovementSystem {
	s := &MovementSystem{
		OnGround: -1,
		OnLadder: -1,
		ScreenW:  640,
		ScreenH:  480,
	}

	return s
}

func drawDebugRect(r image.Rectangle, c color.Color, overrideColorScale bool) gohan.Entity {
	rectEntity := ECS.NewEntity()

	rectImg := ebiten.NewImage(r.Dx(), r.Dy())
	rectImg.Fill(c)

	ECS.AddComponent(rectEntity, &component.PositionComponent{
		X: float64(r.Min.X),
		Y: float64(r.Min.Y),
	})

	ECS.AddComponent(rectEntity, &component.SpriteComponent{
		Image:              rectImg,
		OverrideColorScale: overrideColorScale,
	})

	return rectEntity
}

func (s *MovementSystem) removeDebugRects() {
	for _, e := range s.debugCollisionRects {
		ECS.RemoveEntity(e)
	}
	s.debugCollisionRects = nil

	for _, e := range s.debugLadderRects {
		ECS.RemoveEntity(e)
	}
	s.debugLadderRects = nil

	for _, e := range s.debugFireRects {
		ECS.RemoveEntity(e)
	}
	s.debugFireRects = nil
}

func (s *MovementSystem) addDebugCollisionRects() {
	s.removeDebugRects()

	for _, rect := range s.collisionRects {
		c := color.RGBA{200, 200, 200, 150}
		debugRect := drawDebugRect(rect, c, true)
		s.debugCollisionRects = append(s.debugCollisionRects, debugRect)
	}

	for _, rect := range s.ladderRects {
		c := color.RGBA{0, 0, 200, 150}
		debugRect := drawDebugRect(rect, c, true)
		s.debugLadderRects = append(s.debugLadderRects, debugRect)
	}

	for _, rect := range s.fireRects {
		c := color.RGBA{200, 0, 0, 150}
		debugRect := drawDebugRect(rect, c, false)
		s.debugFireRects = append(s.debugFireRects, debugRect)
	}
}

func (s *MovementSystem) UpdateDebugCollisionRects() {
	if world.World.Debug < 2 {
		s.removeDebugRects()
		return
	} else if len(s.debugCollisionRects) == 0 {
		s.addDebugCollisionRects()
	}

	for i, debugRect := range s.debugCollisionRects {
		sprite := ECS.Component(debugRect, component.SpriteComponentID).(*component.SpriteComponent)
		if s.OnGround == i {
			sprite.ColorScale = 1
		} else {
			sprite.ColorScale = 0.4
		}
	}

	for i, debugRect := range s.debugLadderRects {
		sprite := ECS.Component(debugRect, component.SpriteComponentID).(*component.SpriteComponent)
		if s.OnLadder == i {
			sprite.ColorScale = 1
		} else {
			sprite.ColorScale = 0.4
		}
	}
}

func (_ *MovementSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.VelocityComponentID,
	}
}

func (_ *MovementSystem) Uses() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.WeaponComponentID,
	}
}

func (s *MovementSystem) Update(ctx *gohan.Context) error {
	if world.World.MessageVisible {
		return nil
	}

	if world.World.GameOver && ctx.Entity == world.World.Player {
		return nil
	}

	position := component.Position(ctx)
	velocity := component.Velocity(ctx)

	vx, vy := velocity.X, velocity.Y
	if ctx.Entity == world.World.Player && (world.World.NoClip || world.World.Debug != 0) && ebiten.IsKeyPressed(ebiten.KeyShift) {
		vx, vy = vx*2, vy*2
	}

	position.X, position.Y = position.X+vx, position.Y+vy

	// Force player to remain within the screen bounds.
	if ctx.Entity == world.World.Player {
		screenX, screenY := s.levelCoordinatesToScreen(position.X, position.Y)
		if screenX < 0 {
			diff := screenX / world.World.CamScale
			position.X -= diff
		} else if screenX > float64(world.World.ScreenW)-world.World.PlayerWidth {
			diff := (float64(world.World.ScreenW) - world.World.PlayerWidth - screenX) / world.World.CamScale
			position.X += diff
		}
		if screenY < 0 {
			diff := screenY / world.World.CamScale
			position.Y -= diff
		} else if screenY > float64(world.World.ScreenH)-world.World.PlayerHeight {
			diff := (float64(world.World.ScreenH) - world.World.PlayerHeight - screenY) / world.World.CamScale
			position.Y += diff
		}

		playerRect := image.Rect(int(position.X), int(position.Y), int(position.X+world.World.PlayerWidth), int(position.Y+world.World.PlayerHeight))
		for _, r := range world.World.HazardRects {
			if playerRect.Overlaps(r) {
				world.World.SetGameOver()
				return nil
			}
		}
	} else if ctx.Entity == world.World.BrokenPieceA || ctx.Entity == world.World.BrokenPieceB {
		sprite := ECS.Component(ctx.Entity, component.SpriteComponentID).(*component.SpriteComponent)
		if ctx.Entity == world.World.BrokenPieceA {
			sprite.Angle -= 0.05
		} else {
			sprite.Angle += 0.05
		}
	}

	// TODO check bullet kill player

	return nil
}

func (s *MovementSystem) levelCoordinatesToScreen(x, y float64) (float64, float64) {
	return (x - world.World.CamX) * world.World.CamScale, (y - world.World.CamY) * world.World.CamScale
}

func (_ *MovementSystem) Draw(_ *gohan.Context, screen *ebiten.Image) error {
	return gohan.ErrSystemWithoutDraw
}

func (s *MovementSystem) RecordPosition(position *component.PositionComponent) {
	if math.Abs(position.X-s.playerPosition[0]) >= rewindThreshold || math.Abs(position.Y-s.playerPosition[1]) >= rewindThreshold {
		s.playerPosition[0], s.playerPosition[1] = position.X, position.Y
		s.playerPositions = append(s.playerPositions, s.playerPosition)
	}
}

func (s *MovementSystem) RemoveLastPosition() {
	if len(s.playerPositions) == 0 {
		return
	}

	s.playerPositions = s.playerPositions[:len(s.playerPositions)-1]
	if len(s.playerPositions) > 1 {
		s.playerPosition = s.playerPositions[len(s.playerPositions)-1]
	} else {
		s.playerPosition[0], s.playerPosition[1] = 0, 0
	}
}
