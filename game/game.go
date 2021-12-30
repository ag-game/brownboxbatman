package game

import (
	"image/color"
	"math/rand"
	"os"
	"sync"
	"time"

	"code.rocketnine.space/tslocum/brownboxbatman/asset"
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/brownboxbatman/entity"
	"code.rocketnine.space/tslocum/brownboxbatman/system"
	"code.rocketnine.space/tslocum/brownboxbatman/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

const sampleRate = 44100

// game is an isometric demo game.
type game struct {
	w, h int

	audioContext *audio.Context

	op *ebiten.DrawImageOptions

	disableEsc bool

	debugMode  bool
	cpuProfile *os.File

	movementSystem *system.MovementSystem
	renderSystem   *system.RenderSystem

	addedSystems bool

	sync.Mutex
}

// NewGame returns a new isometric demo game.
func NewGame() (*game, error) {
	g := &game{
		audioContext: audio.NewContext(sampleRate),
		op:           &ebiten.DrawImageOptions{},
	}

	err := g.loadAssets()
	if err != nil {
		panic(err)
	}

	const numEntities = 30000
	ECS.Preallocate(numEntities)

	return g, nil
}

func (g *game) tileToGameCoords(x, y int) (float64, float64) {
	return float64(x) * 32, float64(y) * 32
}

func (g *game) changeMap(filePath string) {
	world.LoadMap(filePath)

	if world.World.Player == 0 {
		world.World.Player = entity.NewPlayer()
	}

	const playerStartOffset = 128
	const camStartOffset = 480

	w := float64(world.World.Map.Width * world.World.Map.TileWidth)
	h := float64(world.World.Map.Height * world.World.Map.TileHeight)

	position := ECS.Component(world.World.Player, component.PositionComponentID).(*component.PositionComponent)
	position.X, position.Y = w/2, h-playerStartOffset

	world.World.CamX, world.World.CamY = 0, h-camStartOffset
}

// Layout is called when the game's layout changes.
func (g *game) Layout(w, h int) (int, int) {
	if !world.World.NativeResolution {
		w, h = 640, 480
	}
	if w != g.w || h != g.h {
		world.World.ScreenW, world.World.ScreenH = w, h
		g.w, g.h = w, h
	}
	return g.w, g.h
}

func (g *game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		g.Exit()
		return nil
	}

	if world.World.ResetGame {
		world.Reset()

		g.changeMap("map/m1.tmx")

		if !g.addedSystems {
			g.addSystems()

			if world.World.Debug == 0 {
				asset.SoundTitleMusic.Play()
			}

			g.addedSystems = true // TODO
		}

		rand.Seed(time.Now().UnixNano())

		world.World.ResetGame = false
		world.World.GameOver = false
	}

	err := ECS.Update()
	if err != nil {
		return err
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	err := ECS.Draw(screen)
	if err != nil {
		panic(err)
	}
}

func (g *game) addSystems() {
	ecs := ECS

	g.movementSystem = system.NewMovementSystem()
	ecs.AddSystem(system.NewPlayerMoveSystem(world.World.Player, g.movementSystem))
	ecs.AddSystem(system.NewplayerFireSystem())
	ecs.AddSystem(g.movementSystem)
	ecs.AddSystem(system.NewCreepSystem())
	ecs.AddSystem(system.NewCameraSystem())
	ecs.AddSystem(system.NewRailSystem())
	g.renderSystem = system.NewRenderSystem()
	ecs.AddSystem(g.renderSystem)
	ecs.AddSystem(system.NewRenderMessageSystem())
	ecs.AddSystem(system.NewRenderDebugTextSystem(world.World.Player))
	ecs.AddSystem(system.NewProfileSystem(world.World.Player))
}

func (g *game) loadAssets() error {
	asset.ImgWhiteSquare.Fill(color.White)
	asset.LoadSounds(g.audioContext)
	return nil
}

func (g *game) Exit() {
	os.Exit(0)
}
