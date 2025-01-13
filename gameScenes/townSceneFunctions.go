package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"sort"
)

func (g *TownScene) UpdateDoors() {
	for _, object := range g.Objects {
		objectAnimation := object.ActiveAnimation(object.Status)
		if objectAnimation != nil {
			if object.Type == gameObjects.EntryDoor || object.Type == gameObjects.ExitDoor {
				if object.Status == gameObjects.Entering {
					println("updating door,", object.Name)
					objectAnimation.Update()

					if objectAnimation.Frame() == objectAnimation.LastF-3 {
						println("making player not visible for entering building effect")
						//remove sprite on last frame before they are shown inside the building
						g.Player.Visible = false
						objectAnimation.Update()
					}

					if objectAnimation.Frame() == objectAnimation.LastF {
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
						x, y := gameObjects.GetDoorCoord(g.MapData.EntryDoors, object.Name, "down")
						g.Player.X = x
						g.Player.Y = y
						g.Player.InAnimation = false
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

	opts := &ebiten.DrawImageOptions{}

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
		DrawCharacter(character, screen, *g.cam)
	}
}

func LockCursorForDialogue() {
	updater := CreateCursorUpdater()
	updater.MoveToLockedSpecificPosition(995, 710)
	input.SetCursorUpdater(updater)
}
