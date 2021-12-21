package component

import (
	"time"

	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type WeaponComponent struct {
	Equipped bool

	Damage int

	FireRate time.Duration
	LastFire time.Time

	BulletSpeed float64
}

var WeaponComponentID = ECS.NewComponentID()

func (p *WeaponComponent) ComponentID() gohan.ComponentID {
	return WeaponComponentID
}

func Weapon(ctx *gohan.Context) *WeaponComponent {
	c, ok := ctx.Component(WeaponComponentID).(*WeaponComponent)
	if !ok {
		return nil
	}
	return c
}
