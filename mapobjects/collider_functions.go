package mapobjects

import (
	"QuickDrawAdventure/entities"
	"image"
)

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		//check if player is colliding with collider
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y)+27,
				int(sprite.X)+14,
				int(sprite.Y)+31),
		) {
			if sprite.Dx > 0.0 { //check if player is going right
				//update player velocity
				sprite.X = float64(collider.Min.X) - 14
			} else if sprite.Dx < 0.0 { //check if player is going left
				sprite.X = float64(collider.Max.X)
			}

		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		//check if player is colliding with collider
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y)+27,
				int(sprite.X)+14,
				int(sprite.Y)+31),
		) {
			if sprite.Dy > 0.0 { //check if player is going down
				//update player position
				sprite.Y = float64(collider.Min.Y) - 31
			} else if sprite.Dy < 0.0 { //check if player is going up
				sprite.Y = float64(collider.Max.Y) - 27
			}

		}
	}
}
