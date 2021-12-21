package component

import (
	. "code.rocketnine.space/tslocum/brownboxbatman/ecs"
	"code.rocketnine.space/tslocum/gohan"
)

type BulletComponent struct {
}

var BulletComponentID = ECS.NewComponentID()

func (p *BulletComponent) ComponentID() gohan.ComponentID {
	return BulletComponentID
}

func Bullet(ctx *gohan.Context) *BulletComponent {
	c, ok := ctx.Component(BulletComponentID).(*BulletComponent)
	if !ok {
		return nil
	}
	return c
}
