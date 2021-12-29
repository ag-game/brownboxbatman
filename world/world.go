package world

import (
	"image"
	"log"
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

const (
	SoundGunshot = iota
	SoundVampireDie1
	SoundVampireDie2
	SoundBat
	SoundPlayerHurt
	SoundPlayerDie
	SoundPickup
	SoundMunch
)

var World = &GameWorld{
	CamScale:     1,
	CamMoving:    true,
	PlayerWidth:  8,
	PlayerHeight: 32,
	TileImages:   make(map[uint32]*ebiten.Image),
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

	MessageVisible bool

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
}

func SetMessage(message string) {
	// TODO
}

func TileToGameCoords(x, y int) (float64, float64) {
	//return float64(x) * 32, float64(g.currentMap.Height*32) - float64(y)*32 - 32
	return float64(x) * 32, float64(y) * 32
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
				World.HazardRects = append(World.HazardRects, ObjectToRect(obj))
			}
		} else if grp.Name == "CREEPS" {
			for _, obj := range grp.Objects {
				c := NewCreep(component.CreepSnowblower, int64(obj.ID), float64(obj.X), float64(obj.Y))
				World.CreepRects = append(World.CreepRects, ObjectToRect(obj))
				World.CreepEntities = append(World.CreepEntities, c)
			}
		}
	}
}

func ObjectToRect(o *tiled.Object) image.Rectangle {
	x, y, w, h := int(o.X), int(o.Y), int(o.Width), int(o.Height)
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
		asset.SoundBatHit4.Play()
	} else {
		deathSound := rand.Intn(3)
		switch deathSound {
		case 0:
			asset.SoundBatHit1.Play()
		case 1:
			asset.SoundBatHit2.Play()
		case 2:
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

	w.BrokenPieceB = entity.NewCreepBullet(position.X, position.Y, xSpeedB, ySpeedB)
	pieceBSprite := &component.SpriteComponent{
		Image: asset.ImgBatBroken2,
	}
	ECS.AddComponent(w.BrokenPieceB, pieceBSprite)
}

// TODO move
func NewCreep(creepType int, creepID int64, x float64, y float64) gohan.Entity {
	creep := ECS.NewEntity()

	ECS.AddComponent(creep, &component.PositionComponent{
		X: x,
		Y: y,
	})

	if creepType == component.CreepSnowmanHead {
		ECS.AddComponent(creep, &component.VelocityComponent{})

		ECS.AddComponent(creep, &component.CreepComponent{
			Health:     128,
			FireAmount: 1,
			FireRate:   144 * 1.5,
			Rand:       rand.New(rand.NewSource(creepID)),
		})
	} else {
		ECS.AddComponent(creep, &component.CreepComponent{
			Health:     64,
			FireAmount: 8,
			FireRate:   144 / 4,
			Rand:       rand.New(rand.NewSource(creepID)),
		})
	}

	img := asset.ImgBat
	if creepType == component.CreepSnowmanHead {
		img = World.TileImages[8]
	} else if creepType == component.CreepSnowblower {
		img = World.TileImages[50]
	}
	ECS.AddComponent(creep, &component.SpriteComponent{
		Image: img,
	})

	return creep
}
