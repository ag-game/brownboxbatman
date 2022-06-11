package entity

import (
	"code.rocketnine.space/tslocum/brownboxbatman/asset"
	"code.rocketnine.space/tslocum/brownboxbatman/component"
	"code.rocketnine.space/tslocum/gohan"
)

func NewPlayerBullet(x, y, xSpeed, ySpeed float64) gohan.Entity {
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
		Image: asset.ImgBlackSquare,
	})

	bullet.AddComponent(&component.PlayerBullet{})

	return bullet
}
