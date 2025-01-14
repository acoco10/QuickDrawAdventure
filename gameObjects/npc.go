package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/hajimehoshi/ebiten/v2"
)

type NPC struct {
	Name string
	*Sprite
	FollowsEnemy bool
	Animations   map[CharState]*animations.Animation
}

func (e *NPC) ActiveAnimation(dX, dY int) *animations.Animation {
	if dX > 0 {
		return e.Animations[Right]
	}
	if dX < 0 {
		return e.Animations[Left]
	}
	if dY > 0 {
		return e.Animations[Down]
	}
	if dY < 0 {
		return e.Animations[Up]
	}
	return nil
}

func NewNPC(eImg *ebiten.Image, npcSpawn Spawn) (*NPC, error) {
	npc := &NPC{
		Name: npcSpawn.Name,
		Sprite: &Sprite{
			Img:     eImg,
			X:       npcSpawn.X,
			Y:       npcSpawn.Y,
			Visible: true,
		},
		Animations: map[CharState]*animations.Animation{
			Down:  animations.NewAnimation(0, 4, 4, 22.0),
			Up:    animations.NewAnimation(2, 6, 4, 22.0),
			Left:  animations.NewAnimation(1, 10, 4, 11.0),
			Right: animations.NewAnimation(3, 11, 4, 11.0),
		},
	}

	return npc, nil
}
