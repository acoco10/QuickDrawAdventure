package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	ui2 "github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"math"
	"sort"
)

func (g *TownScene) NoTravelDoor(object *gameObjects.DoorObject) {
	objectAnimation := object.ActiveAnimation()
	if objectAnimation != nil {
		if objectAnimation.Frame() == objectAnimation.LastF {
			objectAnimation.Update()
			object.StopAnimation()
			objectAnimation.Reset()
		}
		if object.State == gameObjects.Entering {
			objectAnimation.Update()
		}
		if object.State == gameObjects.On {
			objectAnimation.Update()
		}
	}

}

func (g *TownScene) ExitDoor(fromDoor *gameObjects.DoorObject, toDoor *gameObjects.DoorObject) {
	objectAnimation := fromDoor.ActiveAnimation()
	if objectAnimation != nil {
		if fromDoor.State == gameObjects.Leaving {
			if objectAnimation.Frame() == objectAnimation.FirstF {
				g.cam.UpdateState(camera.Outside)
				x, y := gameObjects.GetDoorCoord(toDoor, fromDoor.Dir)
				g.Player.X = x
				g.Player.Y = y
				g.Player.ExitShadow()
				objectAnimation.Update()

			} else if objectAnimation.Frame() == objectAnimation.LastF {
				fromDoor.StopAnimation()
				fromDoor.Triggered = false
				fromDoor.State = gameObjects.NotTriggered
				objectAnimation.Reset()
			} else {
				objectAnimation.Update()
			}
		}
	}
}

func (g *TownScene) EnterDoor(fromDoor *gameObjects.DoorObject, toDoor *gameObjects.DoorObject) {
	fromDoorAnimation := fromDoor.ActiveAnimation()
	println("triggering function for enter door:", fromDoor.Name)
	if fromDoor.State == gameObjects.Entering {
		g.Player.InAnimation = true
		fromDoorAnimation.Update()
		if fromDoorAnimation.Frame() == fromDoorAnimation.LastF-3 {
			println("making player not visible for entering building effect")
			//remove sprite on last frame before they are shown inside the building
			g.Player.Visible = false
			fromDoorAnimation.Update()
		}
		if fromDoorAnimation.Frame() == fromDoorAnimation.LastF {
			g.cam.UpdateState(camera.Inside)
			g.cam.SetIndoorCameraBounds(g.MapData.CameraPoints[fromDoor.CameraPoint].Rect)
			if fromDoor.Name == "cave" {
				g.dark = true
			}
			println("changing player location to building interior")
			g.Player.Visible = true
			x, y := gameObjects.GetDoorCoord(toDoor, fromDoor.Dir)
			g.Player.X = x
			g.Player.Y = y
			fromDoorAnimation.Update()
			fromDoor.StopAnimation()
			fromDoorAnimation.Reset()
			g.Player.InAnimation = false
			fromDoor.Triggered = false
			fromDoor.State = gameObjects.NotTriggered
		}
	}
}
func (g *TownScene) InsideDoor(fromDoor *gameObjects.DoorObject, toDoor *gameObjects.DoorObject) {
	fromDoorAnimation := fromDoor.ActiveAnimation()
	x, y := gameObjects.GetDoorCoord(toDoor, fromDoor.Dir)
	if fromDoor.State == gameObjects.Leaving {
		g.cam.SetIndoorCameraBounds(g.MapData.CameraPoints[fromDoor.CameraPoint].Rect)
		if fromDoorAnimation.Frame() == fromDoorAnimation.FirstF {
			g.Player.X = x
			g.Player.Y = y
			fromDoor.StopAnimation()
			fromDoorAnimation.Reset()
			fromDoor.Triggered = false
			fromDoor.State = gameObjects.NotTriggered
		}
	}
}

func FindDoor(fromDoor *gameObjects.DoorObject, doors []*gameObjects.DoorObject) *gameObjects.DoorObject {
	for _, door := range doors {
		if door.Name == fromDoor.Name && door.Trigger != fromDoor.Trigger {
			return door
		}
	}
	return nil
}

func (g *TownScene) UpdateDoors() {
	for _, object := range g.MapData.Doors {
		if object.Triggered {
			switch object.Type {
			case gameObjects.EntryDoor:
				g.EnterDoor(object, FindDoor(object, g.MapData.Doors))
			case gameObjects.ExitDoor:
				g.ExitDoor(object, FindDoor(object, g.MapData.Doors))
			case gameObjects.ContextualObject:
				g.NoTravelDoor(object)
			case gameObjects.InsideDoor:
				g.InsideDoor(object, FindDoor(object, g.MapData.Doors))
			default:
				continue
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
		opts.ColorScale.Scale(0.5, 0.5, 0.5, 1)
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

func SetCursorForDialogue(ui DialogueUI, cursor *ui2.CursorUpdater) {
	if ui.DType == Dialogue {
		cursor.MoveToLockedSpecificPosition(1002, 306, 306)

	}
	if ui.DType == ShowDown {
		cursor.MoveToLockedSpecificPosition(1002, 720, 720)
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

func (g *TownScene) DrawItems(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	for _, item := range g.MapData.Items {
		if item.Name == "bridgeLeft" || item.Name == "bridgeRight" {
			itemRect := image.Rect(int(item.X), int(item.Y), int(item.X)+16, int(item.Y)+32)
			if itemRect.Overlaps(
				image.Rect(
					int(g.Player.X),
					int(g.Player.Y)+28,
					int(g.Player.X)+16,
					int(g.Player.Y)+32),
			) {
				opts.GeoM.Translate(item.X, item.Y+1)
				opts.GeoM.Translate(g.cam.X, g.cam.Y)
				opts.GeoM.Scale(4, 4)
				opts.GeoM.Rotate(1.5 * math.Pi / 360)
				screen.DrawImage(item.Img, opts)
				opts.GeoM.Reset()
				if item.State == "off" {
					if g.Player.Dx > 0 || item.Name == "bridgeRight" {
						g.Player.Y -= 0.05
					}
					if g.Player.Dx < 0 || item.Name == "bridgeRight" {
						g.Player.Y += 0.05
					}
					if g.Player.Dx < 0 || item.Name == "bridgeLeft" {
						g.Player.Y -= 0.05
					}
					if g.Player.Dx > 0 || item.Name == "bridgeLeft" {
						g.Player.Y += 0.05
					}
					item.State = "on"
				}
			} else {
				opts.GeoM.Translate(item.X, item.Y+2)
				opts.GeoM.Translate(g.cam.X, g.cam.Y)
				opts.GeoM.Scale(4, 4)
				screen.DrawImage(item.Img, opts)
				opts.GeoM.Reset()
			}
		} else {
			item.State = "off"
			opts.GeoM.Reset()
			opts.GeoM.Translate(item.X*4, item.Y*4-137)
			opts.GeoM.Translate(g.cam.X*4, g.cam.Y*4)
			screen.DrawImage(item.Img, opts)
			opts.GeoM.Reset()
		}
	}
}

func (g *TownScene) DrawObjects(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	for _, object := range g.MapData.Doors {
		opts.GeoM.Translate(object.X, object.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
		opts.GeoM.Scale(4, 4)

		objectFrame := 0
		objectAnimation := object.ActiveAnimation()

		if objectAnimation != nil {
			objectFrame = objectAnimation.Frame()
		}
		screen.DrawImage(
			object.Img.SubImage(
				object.SpriteSheet.Rect(objectFrame),
			).(*ebiten.Image),
			opts,
		)
		opts.GeoM.Reset()
	}
}

func (g *TownScene) DrawObjectsAbovePlayer(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	for _, object := range g.Objects {
		if object.DrawAbovePlayer && g.Player.Y+6 < object.Y {
			opts.GeoM.Translate(object.X, object.Y)
			opts.GeoM.Translate(g.cam.X, g.cam.Y)
			opts.GeoM.Scale(4, 4)

			objectFrame := 0
			objectAnimation := object.ActiveAnimation()

			if objectAnimation != nil {
				objectFrame = objectAnimation.Frame()
			}
			screen.DrawImage(
				object.Img.SubImage(
					object.SpriteSheet.Rect(objectFrame),
				).(*ebiten.Image),
				opts,
			)
		}
		opts.GeoM.Reset()
	}
}

func DistanceEq(x1, y1, x2, y2 float64) float64 {

	xdis := math.Abs(x1 - x2)
	ydis := math.Abs(y1 - y2)

	return math.Sqrt(xdis*xdis + ydis*ydis)
}

func (g *TownScene) ActivateTriggerCollidersVertical(trigger string) {
	gameObjects.CheckCollisionVertical(g.Player.Sprite, g.MapData.TriggerColliders[trigger], g.NPC)
}

func (g *TownScene) ActivateTriggerCollidersHorizontal(trigger string) {
	gameObjects.CheckCollisionHorizontal(g.Player.Sprite, g.MapData.TriggerColliders[trigger], g.NPC)
}

func (g *TownScene) UpdateTriggerColliders(way string) {
	for key, trig := range g.MapData.LayerTriggers {
		if trig.Triggered {
			if way == "vert" {
				g.ActivateTriggerCollidersVertical(key)
			} else if way == "horizontal" {
				g.ActivateTriggerCollidersHorizontal(key)
			}
		}
	}
}
