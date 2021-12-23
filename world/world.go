package world

import (
	"image"
	"log"
	"path/filepath"

	"code.rocketnine.space/tslocum/brownboxbatman/entity"

	"code.rocketnine.space/tslocum/brownboxbatman/asset"
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/gohan"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

var World = &GameWorld{
	GameStarted:  true, // TODO
	CamScale:     1,
	CamMoving:    true,
	PlayerWidth:  8,
	PlayerHeight: 32,
}

type GameWorld struct {
	*gohan.World

	Player gohan.Entity

	ScreenW, ScreenH int

	DisableEsc bool

	Debug  int
	NoClip bool

	GameStarted bool
	GameOver    bool

	MessageVisible bool

	CamX, CamY float64
	CamScale   float64
	CamMoving  bool

	PlayerWidth  float64
	PlayerHeight float64

	Map             *tiled.Map
	ObjectGroups    []*tiled.ObjectGroup
	HazardRects     []image.Rectangle
	TriggerEntities []gohan.Entity
	TriggerRects    []image.Rectangle
	TriggerNames    []string

	NativeResolution bool

	BrokenPieceA, BrokenPieceB gohan.Entity
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

	tileCache := make(map[uint32]*ebiten.Image)
	for i := uint32(0); i < uint32(tileset.TileCount); i++ {
		rect := tileset.GetTileRect(i)
		tileCache[i+tileset.FirstGID] = tilesetImg.SubImage(rect).(*ebiten.Image)
	}

	createTileEntity := func(t *tiled.LayerTile, x int, y int) gohan.Entity {
		tileX, tileY := TileToGameCoords(x, y)

		mapTile := ECS.NewEntity()
		ECS.AddComponent(mapTile, &component.PositionComponent{
			X: tileX,
			Y: tileY,
		})

		sprite := &component.SpriteComponent{
			Image:          tileCache[t.Tileset.FirstGID+t.ID],
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

				tileImg := tileCache[t.Tileset.FirstGID+t.ID]
				if tileImg == nil {
					continue
				}

				e := createTileEntity(t, x, y)
				if layer.Name == "CREEPS" {
					creep := &component.CreepComponent{
						Health:   1,
						FireRate: 144 / 2,
					}
					ECS.AddComponent(e, creep)
				}
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
					Image: tileCache[obj.GID],
				})

				World.TriggerNames = append(World.TriggerNames, obj.Name)
				World.TriggerEntities = append(World.TriggerEntities, mapTile)
				World.TriggerRects = append(World.TriggerRects, ObjectToRect(obj))
			}
		} else if grp.Name == "HAZARDS" {
			for _, obj := range grp.Objects {
				World.HazardRects = append(World.HazardRects, ObjectToRect(obj))
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

func (w *GameWorld) SetGameOver() {
	if w.GameOver {
		return
	}

	w.GameOver = true

	sprite := ECS.Component(w.Player, component.SpriteComponentID).(*component.SpriteComponent)
	sprite.Image = ebiten.NewImage(1, 1)

	position := ECS.Component(w.Player, component.PositionComponentID).(*component.PositionComponent)
	velocity := ECS.Component(w.Player, component.VelocityComponentID).(*component.VelocityComponent)

	xSpeedA := 1.5
	xSpeedB := -1.5
	ySpeedA := -1.5
	ySpeedB := -1.5
	if velocity.Y > 0 {
		ySpeedA = 1.5
		ySpeedB = 1.5
	} else if velocity.X < 0 {
		xSpeedA = -1.5
		xSpeedB = -1.5
		ySpeedA = -1.5
		ySpeedB = 1.5
	} else if velocity.X > 0 {
		xSpeedA = 1.5
		xSpeedB = 1.5
		ySpeedA = -1.5
		ySpeedB = 1.5
	}

	w.BrokenPieceA = entity.NewBullet(position.X, position.Y, xSpeedA, ySpeedA)
	pieceASprite := &component.SpriteComponent{
		Image: asset.ImgBatBroken1,
	}
	ECS.AddComponent(w.BrokenPieceA, pieceASprite)

	w.BrokenPieceB = entity.NewBullet(position.X, position.Y, xSpeedB, ySpeedB)
	pieceBSprite := &component.SpriteComponent{
		Image: asset.ImgBatBroken2,
	}
	ECS.AddComponent(w.BrokenPieceB, pieceBSprite)
}
