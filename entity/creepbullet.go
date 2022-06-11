package entity

import (
	"code.rocketnine.space/tslocum/brownboxbatman/asset"
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/gohan"
)

func NewCreepBullet(x, y, xSpeed, ySpeed float64) gohan.Entity {
	bullet := gohan.NewEntity()

	bullet.AddComponent(&component.Position{
		X: x,
		Y: y,
	})

	bullet.AddComponent(&component.Velocity{
		X: xSpeed,
		Y: ySpeed,
	})

	bullet.AddComponent(&component.Sprite{
		Image: asset.ImgWhiteSquare,
	})

	bullet.AddComponent(&component.CreepBullet{})

	return bullet
}
