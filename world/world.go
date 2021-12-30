package world

import (
	"image"
	"log"
	"math"
	"math/rand"
	"path/filepath"

	"code.rocketnine.space/tslocum/brownboxbatman/asset"
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/brownboxbatman/entity"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

var World = &GameWorld{
	CamScale:     1,
	CamMoving:    true,
	PlayerWidth:  8,
	PlayerHeight: 32,
	TileImages:   make(map[uint32]*ebiten.Image),
	ResetGame:    true,
}

type GameWorld struct {
	*gohan.World

	Player gohan.Entity

	ScreenW, ScreenH int

	DisableEsc bool

	Debug  int
	NoClip bool

	GameStarted      bool
	GameStartedTicks int
	GameOver         bool

	MessageVisible  bool
	MessageTicks    int
	MessageDuration int
	MessageUpdated  bool
	MessageText     string

	PlayerX, PlayerY float64

	CamX, CamY float64
	CamScale   float64
	CamMoving  bool

	PlayerWidth  float64
	PlayerHeight float64

	Map             *tiled.Map
	ObjectGroups    []*tiled.ObjectGroup
	HazardRects     []image.Rectangle
	CreepRects      []image.Rectangle
	CreepEntities   []gohan.Entity
	TriggerEntities []gohan.Entity
	TriggerRects    []image.Rectangle
	TriggerNames    []string

	NativeResolution bool

	BrokenPieceA, BrokenPieceB gohan.Entity

	TileImages map[uint32]*ebiten.Image

	ResetGame bool

	resetTipShown bool
}

func TileToGameCoords(x, y int) (float64, float64) {
	//return float64(x) * 32, float64(g.currentMap.Height*32) - float64(y)*32 - 32
	return float64(x) * 32, float64(y) * 32
}

func Reset() {
	for _, e := range ECS.Entities() {
		ECS.RemoveEntity(e)
	}
	World.Player = 0

	World.ObjectGroups = nil
	World.HazardRects = nil
	World.CreepRects = nil
	World.CreepEntities = nil
	World.TriggerEntities = nil
	World.TriggerRects = nil
	World.TriggerNames = nil

	World.MessageVisible = false
}

func LoadMap(filePath string) {
	loader := tiled.Loader{
		FileSystem: asset.FS,
	}

	// Parse .tmx file.
	m, err := loader.LoadFromFile(filepath.FromSlash(filePath))
	if err != nil {
		log.Fatalf("error parsing world: %+v", err)
	}

	// Load tileset.

	tileset := m.Tilesets[0]

	imgPath := filepath.Join("./map/", tileset.Image.Source)
	f, err := asset.FS.Open(filepath.FromSlash(imgPath))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	tilesetImg := ebiten.NewImageFromImage(img)

	// Load tiles.

	for i := uint32(0); i < uint32(tileset.TileCount); i++ {
		rect := tileset.GetTileRect(i)
		World.TileImages[i+tileset.FirstGID] = tilesetImg.SubImage(rect).(*ebiten.Image)
	}

	createTileEntity := func(t *tiled.LayerTile, x int, y int) gohan.Entity {
		tileX, tileY := TileToGameCoords(x, y)

		mapTile := ECS.NewEntity()
		ECS.AddComponent(mapTile, &component.PositionComponent{
			X: tileX,
			Y: tileY,
		})

		sprite := &component.SpriteComponent{
			Image:          World.TileImages[t.Tileset.FirstGID+t.ID],
			HorizontalFlip: t.HorizontalFlip,
			VerticalFlip:   t.VerticalFlip,
			DiagonalFlip:   t.DiagonalFlip,
		}
		ECS.AddComponent(mapTile, sprite)

		return mapTile
	}

	var t *tiled.LayerTile
	for _, layer := range m.Layers {
		for y := 0; y < m.Height; y++ {
			for x := 0; x < m.Width; x++ {
				t = layer.Tiles[y*m.Width+x]
				if t == nil || t.Nil {
					continue // No tile at this position.
				}

				tileImg := World.TileImages[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					continue
				}
				createTileEntity(t, x, y)
			}
		}
	}

	// Load ObjectGroups.

	var objects []*tiled.ObjectGroup
	var loadObjects func(grp *tiled.Group)
	loadObjects = func(grp *tiled.Group) {
		for _, subGrp := range grp.Groups {
			loadObjects(subGrp)
		}
		for _, objGrp := range grp.ObjectGroups {
			objects = append(objects, objGrp)
		}
	}
	for _, grp := range m.Groups {
		loadObjects(grp)
	}
	for _, objGrp := range m.ObjectGroups {
		objects = append(objects, objGrp)
	}

	World.Map = m
	World.ObjectGroups = objects

	for _, grp := range World.ObjectGroups {
		if grp.Name == "TRIGGERS" {
			for _, obj := range grp.Objects {
				mapTile := ECS.NewEntity()
				ECS.AddComponent(mapTile, &component.PositionComponent{
					X: obj.X,
					Y: obj.Y - 32,
				})
				ECS.AddComponent(mapTile, &component.SpriteComponent{
					Image: World.TileImages[obj.GID],
				})

				World.TriggerNames = append(World.TriggerNames, obj.Name)
				World.TriggerEntities = append(World.TriggerEntities, mapTile)
				World.TriggerRects = append(World.TriggerRects, ObjectToRect(obj))
			}
		} else if grp.Name == "HAZARDS" {
			for _, obj := range grp.Objects {
				r := ObjectToRect(obj)
				r.Min.Y += 32
				r.Max.Y += 32
				World.HazardRects = append(World.HazardRects, r)
			}
		} else if grp.Name == "CREEPS" {
			for _, obj := range grp.Objects {
				creepType := component.CreepSnowblower
				switch obj.GID {
				case 9:
					creepType = component.CreepSmallRock
				case 18:
					creepType = component.CreepMediumRock
				case 23:
					creepType = component.CreepLargeRock
				}
				r := ObjectToRect(obj)
				c := NewCreep(creepType, int64(obj.ID), float64(r.Min.X), float64(r.Min.Y))
				World.CreepRects = append(World.CreepRects, r)
				World.CreepEntities = append(World.CreepEntities, c)
			}
		}
	}
}

func ObjectToRect(o *tiled.Object) image.Rectangle {
	x, y, w, h := int(o.X), int(o.Y), int(o.Width), int(o.Height)
	y -= 32
	return image.Rect(x, y, x+w, y+h)
}

func LevelCoordinatesToScreen(x, y float64) (float64, float64) {
	return (x - World.CamX) * World.CamScale, (y - World.CamY) * World.CamScale
}

func (w *GameWorld) SetGameOver(vx, vy float64) {
	if w.GameOver {
		return
	}

	w.GameOver = true

	if rand.Intn(100) == 7 {
		asset.SoundBatHit4.Rewind()
		asset.SoundBatHit4.Play()
	} else {
		deathSound := rand.Intn(3)
		switch deathSound {
		case 0:
			asset.SoundBatHit1.Rewind()
			asset.SoundBatHit1.Play()
		case 1:
			asset.SoundBatHit2.Rewind()
			asset.SoundBatHit2.Play()
		case 2:
			asset.SoundBatHit3.Rewind()
			asset.SoundBatHit3.Play()
		}
	}

	sprite := ECS.Component(w.Player, component.SpriteComponentID).(*component.SpriteComponent)
	sprite.Image = ebiten.NewImage(1, 1)

	position := ECS.Component(w.Player, component.PositionComponentID).(*component.PositionComponent)

	if vx == 0 && vy == 0 {
		velocity := ECS.Component(w.Player, component.VelocityComponentID).(*component.VelocityComponent)
		vx, vy = velocity.X, velocity.Y
	}

	xSpeedA := 1.5
	xSpeedB := -1.5
	ySpeedA := -1.5
	ySpeedB := -1.5
	if vy > 0 {
		ySpeedA = 1.5
		ySpeedB = 1.5
	} else if vx < 0 {
		xSpeedA = -1.5
		xSpeedB = -1.5
		ySpeedA = -1.5
		ySpeedB = 1.5
	} else if vx > 0 {
		xSpeedA = 1.5
		xSpeedB = 1.5
		ySpeedA = -1.5
		ySpeedB = 1.5
	}

	w.BrokenPieceA = entity.NewCreepBullet(position.X, position.Y, xSpeedA, ySpeedA)
	pieceASprite := &component.SpriteComponent{
		Image: asset.ImgBatBroken1,
	}
	ECS.AddComponent(w.BrokenPieceA, pieceASprite)
	ECS.AddComponent(w.BrokenPieceA, &component.CreepBulletComponent{
		Invulnerable: true,
	})

	w.BrokenPieceB = entity.NewCreepBullet(position.X, position.Y, xSpeedB, ySpeedB)
	pieceBSprite := &component.SpriteComponent{
		Image: asset.ImgBatBroken2,
	}
	ECS.AddComponent(w.BrokenPieceB, pieceBSprite)
	ECS.AddComponent(w.BrokenPieceB, &component.CreepBulletComponent{
		Invulnerable: true,
	})

	if !World.resetTipShown {
		SetMessage("  GAME  OVER\n\nRESET: <ENTER>", math.MaxInt)
		World.resetTipShown = true
	} else {
		SetMessage("GAME OVER", math.MaxInt)
	}
}

// TODO move
func NewCreep(creepType int, creepID int64, x float64, y float64) gohan.Entity {
	creep := ECS.NewEntity()

	ECS.AddComponent(creep, &component.PositionComponent{
		X: x,
		Y: y,
	})

	if creepType == component.CreepSmallRock {
		ECS.AddComponent(creep, &component.VelocityComponent{})

		ECS.AddComponent(creep, &component.CreepComponent{
			Type:       creepType,
			Health:     32,
			FireAmount: 2,
			FireRate:   144 * 1,
			Rand:       rand.New(rand.NewSource(creepID)),
		})
	} else if creepType == component.CreepMediumRock {
		ECS.AddComponent(creep, &component.VelocityComponent{})

		ECS.AddComponent(creep, &component.CreepComponent{
			Type:       creepType,
			Health:     64,
			FireAmount: 4,
			FireRate:   144 * 1,
			Rand:       rand.New(rand.NewSource(creepID)),
		})
	} else if creepType == component.CreepLargeRock {
		ECS.AddComponent(creep, &component.VelocityComponent{})

		ECS.AddComponent(creep, &component.CreepComponent{
			Type:       creepType,
			Health:     96,
			FireAmount: 8,
			FireRate:   144,
			Rand:       rand.New(rand.NewSource(creepID)),
		})
	} else { // CreepSnowblower
		ECS.AddComponent(creep, &component.CreepComponent{
			Type:       creepType,
			Health:     64,
			FireAmount: 8,
			FireRate:   144 / 4,
			Rand:       rand.New(rand.NewSource(creepID)),
		})
	}

	// TODO handle flipped creep
	var img *ebiten.Image
	if creepType == component.CreepSmallRock {
		img = World.TileImages[9]
	} else if creepType == component.CreepMediumRock {
		img = World.TileImages[18]
	} else if creepType == component.CreepLargeRock {
		img = World.TileImages[23]
	} else { // CreepSnowblower
		img = World.TileImages[50]
	}
	ECS.AddComponent(creep, &component.SpriteComponent{
		Image: img,
	})

	return creep
}

func StartGame() {
	if World.GameStarted {
		return
	}
	World.GameStarted = true

	asset.SoundTitleMusic.Pause()

	asset.SoundLevelMusic.Play()
}

func SetMessage(message string, duration int) {
	World.MessageText = message
	World.MessageVisible = true
	World.MessageUpdated = true
	World.MessageDuration = duration
	World.MessageTicks = 0
}
