package entity

import (
	"code.rocketnine.space/tslocum/brownboxbatman/asset"
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/gohan"
)

func NewPlayer() gohan.Entity {
	player := gohan.NewEntity()

	player.AddComponent(&component.Position{})

	player.AddComponent(&component.Velocity{})

	weapon := &component.Weapon{
		Damage:      1,
		FireRate:    144 / 16,
		BulletSpeed: 8,
	}
	player.AddComponent(weapon)

	player.AddComponent(&component.Sprite{
		Image: asset.ImgBat,
	})

	player.AddComponent(&component.Rail{})

	return player
}
