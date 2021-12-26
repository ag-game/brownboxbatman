package entity

import (
	"code.rocketnine.space/tslocum/brownboxbatman/asset"
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

func NewPlayer() gohan.Entity {
	player := ECS.NewEntity()

	ECS.AddComponent(player, &component.PositionComponent{})

	ECS.AddComponent(player, &component.VelocityComponent{})

	weapon := &component.WeaponComponent{
		Damage:      1,
		FireRate:    144 / 16,
		BulletSpeed: 8,
	}
	ECS.AddComponent(player, weapon)

	ECS.AddComponent(player, &component.SpriteComponent{
		Image: asset.ImgBat,
	})

	ECS.AddComponent(player, &component.RailComponent{})

	return player
}
