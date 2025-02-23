package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"log"
	"sort"
)

func (g *TownScene) UpdateDoors() {
	for _, object := range g.Objects {
		objectAnimation := object.ActiveAnimation(object.Status)
		if objectAnimation != nil {
			if object.Type == gameObjects.EntryDoor || object.Type == gameObjects.ExitDoor {
				if object.Status == gameObjects.Entering {
					objectAnimation.Update()

					if objectAnimation.Frame() == objectAnimation.LastF-3 {
						println("making player not visible for entering building effect")
						//remove sprite on last frame before they are shown inside the building
						g.Player.Visible = false
						objectAnimation.Update()
					}

					if objectAnimation.Frame() == objectAnimation.LastF {
						g.cam.UpdateState(camera.Inside)
						g.cam.SetIndoorCameraBounds(g.MapData.CameraPoints[object.Name].Rect)
						println("changing player location to building interior")
						g.Player.Visible = true
						x, y := gameObjects.GetDoorCoord(g.MapData.ExitDoors, object.Name, "up")
						g.Player.X = x
						g.Player.Y = y
						objectAnimation.Update()
						object.StopAnimation()
						objectAnimation.Reset()
						g.Player.InAnimation = false
					}
				}

				if object.Status == gameObjects.Leaving {

					if objectAnimation.Frame() == objectAnimation.FirstF {
						g.cam.UpdateState(camera.Outside)
						x, y := gameObjects.GetDoorCoord(g.MapData.EntryDoors, object.Name, "down")
						g.Player.X = x
						g.Player.Y = y
						g.Player.InAnimation = false
						g.Player.ExitShadow()
						objectAnimation.Update()

					} else if objectAnimation.Frame() == objectAnimation.LastF {
						object.StopAnimation()
						objectAnimation.Reset()

					} else {

						objectAnimation.Update()
					}
				}
			}

			if object.Type == gameObjects.ContextualObject {
				if objectAnimation.Frame() == objectAnimation.LastF {
					objectAnimation.Update()
					object.StopAnimation()
					objectAnimation.Reset()
				}
				if object.Status == gameObjects.Entering {
					objectAnimation.Update()
				}
				if object.Status == gameObjects.On {
					objectAnimation.Update()
				}

			}
		}

	}
}

func (g *TownScene) SortCharacters() []*gameObjects.Character {
	characters := make([]*gameObjects.Character, 0)
	characters = append(characters, g.Player)
	for _, char := range g.NPC {
		characters = append(characters, char)
	}
	sort.Slice(characters, func(i, j int) bool {
		return characters[i].Y < characters[j].Y
	})

	return characters
}

func DrawCharacter(character *gameObjects.Character, screen *ebiten.Image, cam camera.Camera) {
	lDustEffect, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/characters/nonBattleCharacterAffects/walkingLeftDust.png")
	if err != nil {
		log.Fatal(err)
	}
	rDustEffect, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/characters/nonBattleCharacterAffects/walkingRightDust.png")

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(character.X, character.Y)
	opts.GeoM.Translate(cam.X, cam.Y)
	opts.GeoM.Scale(4, 4)
	if character.Dx > 0 {
		opts.GeoM.Translate(-34, 70)
		screen.DrawImage(lDustEffect, opts)

	}
	if character.Dy < 0 {
		opts.GeoM.Translate(0, 85)
		screen.DrawImage(rDustEffect, opts)
	}

	opts.GeoM.Reset()
	opts.GeoM.Translate(character.X, character.Y)
	opts.GeoM.Translate(cam.X, cam.Y)
	opts.GeoM.Scale(4, 4)
	characterFrame := 0
	characterActiveAnimation := character.ActiveAnimation(int(character.Dx), int(character.Dy))
	if characterActiveAnimation != nil {
		characterFrame = characterActiveAnimation.Frame()
	} else {
		if character.Direction == "U" {
			characterFrame = character.Animations[0].Frame()
		}
		if character.Direction == "D" {
			characterFrame = character.Animations[1].Frame()
		}
		if character.Direction == "R" {
			characterFrame = character.Animations[2].Frame()
		}
		if character.Direction == "L" {
			characterFrame = character.Animations[3].Frame()
		}

	}
	if character.Shadow {
		opts.ColorScale.Scale(0.5, 0.5, 0.5, 255)
	}
	if character.Visible {
		screen.DrawImage(
			//grab a subimage of the Spritesheet
			character.Img.SubImage(
				character.SpriteSheet.Rect(characterFrame)).(*ebiten.Image),
			opts,
		)
	}

	opts.GeoM.Reset()
}

func (g *TownScene) DrawCharacters(screen *ebiten.Image) {
	characters := g.SortCharacters()
	for _, character := range characters {
		if character.Spawned == true {
			DrawCharacter(character, screen, *g.cam)
		}
	}
}

func (g *TownScene) DrawColliders(screen *ebiten.Image) {
	for _, collider := range g.MapData.Colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)*4+float32(g.cam.X)*4,
			float32(collider.Min.Y)*4+float32(g.cam.Y)*4,
			float32(collider.Dx())*4,
			float32(collider.Dy())*4,
			1.0,
			color.RGBA{255, 0, 0, 255},
			false,
		)
	}
}

func LockCursorForDialogue(ui DialogueUI) {
	if ui.DType == Dialogue {
		updater := CreateCursorUpdater(10, 10)
		updater.MoveToLockedSpecificPosition(1002, 306)
		input.SetCursorUpdater(updater)
	}
	if ui.DType == ShowDown {
		updater := CreateCursorUpdater(10, 10)
		updater.MoveToLockedSpecificPosition(1002, 720)
		input.SetCursorUpdater(updater)
	}

}

func ProcessEvent() {

}

type Event interface {
	TriggerEvent()
}

type NpcSpawn struct {
	npcName string
}

func (n *NpcSpawn) TriggerEvent(npcName string) {

}
