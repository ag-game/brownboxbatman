package entity

import (
	"time"

	"code.rocketnine.space/tslocum/brownboxbatman/asset"
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

const (
	spawnX = -100
	spawnY = 5
)

func NewPlayer() gohan.Entity {
	player := ECS.NewEntity()

	ECS.AddComponent(player, &component.PositionComponent{})

	ECS.AddComponent(player, &component.VelocityComponent{})

	weapon := &component.WeaponComponent{
		Damage:      1,
		FireRate:    100 * time.Millisecond,
		BulletSpeed: 15,
	}
	ECS.AddComponent(player, weapon)

	ECS.AddComponent(player, &component.SpriteComponent{
		Image: asset.ImgBrownBat,
	})

	return player
}
