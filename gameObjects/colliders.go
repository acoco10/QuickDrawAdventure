package gameObjects

import (
	"image"
)

type Collider struct {
	*image.Rectangle
	*Trigger
}

func CheckCollisionHorizontal(playerSprite *Sprite, colliders []Collider, npcSprites map[string]*Character) {
	for _, collider := range colliders {
		//check if player is colliding with collider
		if collider.Overlaps(
			image.Rect(
				int(playerSprite.X),
				int(playerSprite.Y)+28,
				int(playerSprite.X)+16,
				int(playerSprite.Y)+32),
		) && (collider.Z == playerSprite.Z || collider.Z == 0) {
			if playerSprite.Dx > 0.0 { //check if player is going right
				//update player velocity
				playerSprite.X = float64(collider.Min.X) - 16
			} else if playerSprite.Dx < 0.0 { //check if player is going left
				playerSprite.X = float64(collider.Max.X)
			}

		}
	}

	for _, npcSprite := range npcSprites {
		playerCollider := image.Rect(
			int(playerSprite.X),
			int(playerSprite.Y)+28,
			int(playerSprite.X)+16,
			int(playerSprite.Y)+32)
		npcCollider :=
			image.Rect(
				int(npcSprite.X),
				int(npcSprite.Y)+28,
				int(npcSprite.X)+16,
				int(npcSprite.Y)+32,
			)
		if npcCollider.Overlaps(playerCollider) {
			if playerSprite.Dx > 0.0 { //check if player is going right
				//update player velocity
				playerSprite.X = float64(npcCollider.Min.X) - 16
			} else if playerSprite.Dx < 0.0 { //check if player is going left
				playerSprite.X = float64(npcCollider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(playerSprite *Sprite, colliders []Collider, npcSprites map[string]*Character) {
	for _, collider := range colliders {
		//check if player is colliding with collider
		if collider.Overlaps(
			image.Rect(
				int(playerSprite.X),
				int(playerSprite.Y)+28,
				int(playerSprite.X)+16,
				int(playerSprite.Y)+32),
		) && (collider.Z == playerSprite.Z || collider.Z == 0) {
			if playerSprite.Dy > 0.0 { //check if player is going down
				//update player position
				playerSprite.Y = float64(collider.Min.Y) - 32
			} else if playerSprite.Dy < 0.0 { //check if player is going up
				playerSprite.Y = float64(collider.Max.Y) - 28
			}

		}
	}
	for _, npcSprite := range npcSprites {
		playerCollider := image.Rect(
			int(playerSprite.X),
			int(playerSprite.Y)+28,
			int(playerSprite.X)+16,
			int(playerSprite.Y)+32)
		npcCollider :=
			image.Rect(
				int(npcSprite.X),
				int(npcSprite.Y)+28,
				int(npcSprite.X)+16,
				int(npcSprite.Y)+32,
			)
		if npcCollider.Overlaps(playerCollider) {
			if playerSprite.Dy > 0.0 { //check if player is going right
				//update player velocity
				playerSprite.Y = float64(npcCollider.Min.Y) - 32
			} else if playerSprite.Dy < 0.0 { //check if player is going left
				playerSprite.Y = float64(npcCollider.Max.Y) - 28
			}
		}
	}

}
