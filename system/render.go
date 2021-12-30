package system

import (
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"code.rocketnine.space/tslocum/brownboxbatman/asset"

	"code.rocketnine.space/tslocum/brownboxbatman/component"
	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TileWidth = 32

	logoText      = "POWERED BY EBITEN"
	logoTextScale = 4.75
	logoTextWidth = 6.0 * float64(len(logoText)) * logoTextScale
	logoTime      = 144 * 3.5

	fadeInTime = 144 * 0.75
)

type RenderSystem struct {
	ScreenW int
	ScreenH int

	img *ebiten.Image
	op  *ebiten.DrawImageOptions

	camScale float64

	renderer gohan.Entity
}

func NewRenderSystem() *RenderSystem {
	s := &RenderSystem{
		renderer: ECS.NewEntity(),
		img:      ebiten.NewImage(320, 100),
		op:       &ebiten.DrawImageOptions{},
		camScale: 1,
		ScreenW:  640,
		ScreenH:  480,
	}

	return s
}

func (s *RenderSystem) Needs() []gohan.ComponentID {
	return []gohan.ComponentID{
		component.PositionComponentID,
		component.SpriteComponentID,
	}
}

func (s *RenderSystem) Uses() []gohan.ComponentID {
	return nil
}

func (s *RenderSystem) Update(_ *gohan.Context) error {
	return gohan.ErrSystemWithoutUpdate
}

func (s *RenderSystem) levelCoordinatesToScreen(x, y float64) (float64, float64) {
	px, py := world.World.CamX, world.World.CamY
	py *= -1
	return ((x - px) * s.camScale), ((y + py) * s.camScale)
}

// renderSprite renders a sprite on the screen.
func (s *RenderSystem) renderSprite(x float64, y float64, offsetx float64, offsety float64, angle float64, geoScale float64, colorScale float64, alpha float64, hFlip bool, vFlip bool, sprite *ebiten.Image, target *ebiten.Image) int {
	if alpha < .01 || colorScale < .01 {
		return 0
	}

	// Skip drawing off-screen tiles.
	drawX, drawY := s.levelCoordinatesToScreen(x, y)
	const padding = TileWidth * 4
	width, height := float64(TileWidth), float64(TileWidth)
	left := drawX
	right := drawX + width
	top := drawY
	bottom := drawY + height
	if (left < -padding || left > float64(s.ScreenW)+padding) || (top < -padding || top > float64(s.ScreenH)+padding) ||
		(right < -padding || right > float64(s.ScreenW)+padding) || (bottom < -padding || bottom > float64(s.ScreenH)+padding) {
		return 0
	}

	s.op.GeoM.Reset()

	if hFlip {
		s.op.GeoM.Scale(-1, 1)
		s.op.GeoM.Translate(TileWidth, 0)
	}
	if vFlip {
		s.op.GeoM.Scale(1, -1)
		s.op.GeoM.Translate(0, TileWidth)
	}

	s.op.GeoM.Scale(geoScale, geoScale)
	// Rotate
	s.op.GeoM.Translate(offsetx, offsety)
	s.op.GeoM.Rotate(angle)
	// Move to current isometric position.
	s.op.GeoM.Translate(x, y)
	// Translate camera position.
	s.op.GeoM.Translate(-world.World.CamX, -world.World.CamY)
	// Zoom.
	s.op.GeoM.Scale(s.camScale, s.camScale)
	// Center.
	//s.op.GeoM.Translate(float64(s.ScreenW/2.0), float64(s.ScreenH/2.0))

	s.op.ColorM.Scale(colorScale, colorScale, colorScale, alpha)

	target.DrawImage(sprite, s.op)

	s.op.ColorM.Reset()

	return 1
}

func (s *RenderSystem) Draw(ctx *gohan.Context, screen *ebiten.Image) error {
	if !world.World.GameStarted {
		if ctx.Entity == world.World.Player {
			world.World.GameStartedTicks++

			timeA := 144.0 * 2
			pct := float64(world.World.GameStartedTicks) / timeA
			if pct > 1 {
				pct = 1
			}
			s.op.GeoM.Reset()
			s.op.ColorM.Reset()
			s.op.ColorM.Scale(1, 1, 1, pct)
			screen.DrawImage(asset.ImgTitle1, s.op)
			s.op.ColorM.Reset()

			timeB := 144.0 * 2.75
			if world.World.GameStartedTicks > int(timeB) {
				pct = float64(world.World.GameStartedTicks-int(timeB)) / (144 * 0.5)
				if pct > 1 {
					pct = 1
				}
				s.op.ColorM.Reset()
				s.op.ColorM.Scale(1, 1, 1, pct)
				screen.DrawImage(asset.ImgTitle2, s.op)
			}

			timeC := 144.0 * 4.35
			if world.World.GameStartedTicks > int(timeC) {
				pct = float64(world.World.GameStartedTicks-int(timeC)) / (144 * 0.5)
				if pct > 1 {
					pct = 1
				}
				s.op.ColorM.Reset()
				s.op.ColorM.Scale(1, 1, 1, pct)
				screen.DrawImage(asset.ImgTitle3, s.op)

				if world.World.GameStartedTicks > 144*6 && world.World.GameStartedTicks%(144*2) < 144*1.5 {
					s.img.Clear()
					ebitenutil.DebugPrint(s.img, "PRESS <ENTER> OR <START> OR <LMB>")

					s.op.GeoM.Scale(2, 2)
					s.op.GeoM.Translate(120, 387)
					screen.DrawImage(s.img, s.op)
				}
			}
			s.op.ColorM.Reset()
		}
		return nil
	}

	position := component.Position(ctx)
	sprite := component.Sprite(ctx)

	if sprite.NumFrames > 0 && time.Since(sprite.LastFrame) > sprite.FrameTime {
		sprite.Frame++
		if sprite.Frame >= sprite.NumFrames {
			sprite.Frame = 0
		}
		sprite.Image = sprite.Frames[sprite.Frame]
		sprite.LastFrame = time.Now()
	}

	colorScale := 1.0
	if sprite.OverrideColorScale {
		colorScale = sprite.ColorScale
	}

	s.renderSprite(position.X, position.Y, 0, 0, sprite.Angle, 1.0, colorScale, 1.0, sprite.HorizontalFlip, sprite.VerticalFlip, sprite.Image, screen)
	return nil
}
