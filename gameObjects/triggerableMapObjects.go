package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math/rand/v2"
)

type Door struct {
	Key   string
	Coord image.Rectangle
}

func CheckContextualTriggers(player *Character, contextTriggers map[string]*Trigger) map[string]ObjectState {
	playerRect := image.Rect(
		int(player.X),
		int(player.Y)+20,
		int(player.X)+16,
		int(player.Y)+31)
	returnMap := make(map[string]ObjectState)
	for _, trig := range contextTriggers {
		if trig.Overlaps(playerRect) {
			returnMap[trig.Name] = On
		}

		if !trig.Overlaps(playerRect) {
			returnMap[trig.Name] = NotTriggered
		}
	}
	return returnMap
}

func CheckEnemyTrigger(player *Character, enemySpawn map[string]Trigger, countDown int) battleStats.CharacterName {
	playerRect := image.Rect(
		int(player.X),
		int(player.Y)+20,
		int(player.X)+16,
		int(player.Y)+31)

	for _, trig := range enemySpawn {
		if trig.Overlaps(playerRect) {
			countDown--
			if rand.IntN(1000-countDown) <= 2 {
				return battleStats.Wolf
			}
		}
	}
	return battleStats.None
}

func GetDoorCoord(door *DoorObject, direction Direction) (float64, float64) {
	//direction determines if player is walking north or south coming out of the destination door
	x := 0.0
	y := 0.0

	if direction == Up {
		x = float64(door.Rectangle.Dx()/2+door.Min.X) - 8
		y = float64(door.Min.Y) - 32
	}
	if direction == Down {
		x = float64(door.Rectangle.Dx()/2+door.Min.X) - 8
		y = float64(door.Max.Y) - 16
	}

	return x, y

}

func CheckDoors(player *Character, doors []*DoorObject) {
	for _, door := range doors {
		CheckTrigger(player, door.Trigger)
	}
}

func CheckTrigger(player *Character, trigger *Trigger) {
	if trigger.Overlaps(
		image.Rect(
			int(player.X),
			int(player.Y)+28,
			int(player.X)+16,
			int(player.Y)+31),
	) {
		yTriggerCheck(player, trigger)
		xTriggerCheck(player, trigger)
	} else {
		if trigger.Triggered {
			if trigger.Z > player.Z {
				player.AddToFuncQueue(player.IncreaseZ)
			}

		}
		trigger.Triggered = false
	}
}

func yTriggerCheck(player *Character, trigger *Trigger) {
	if player.Dy < 0 {
		if trigger.Dir == Up {
			trigger.Triggered = true
		}
	} else if player.Dy > 0 {
		if trigger.Dir == Down {
			trigger.Triggered = true
		}
	}
}

func xTriggerCheck(player *Character, trigger *Trigger) {
	if player.Dx < 0 {
		if trigger.Dir == Left {
			trigger.Triggered = true
		} else {
			trigger.Triggered = false
		}
	} else if player.Dx > 0 {
		if trigger.Dir == Right {
			trigger.Triggered = true
			return
		} else {
			trigger.Triggered = false
		}
	}
}

func CheckTriggers(player *Character, triggers []*Trigger) {
	for _, trigger := range triggers {
		if trigger.Auto {
			CheckTrigger(player, trigger)
		}
	}
}

func DrawTriggers(triggers []*Trigger, screen *ebiten.Image, cam *camera.Camera) {
	for _, trig := range triggers {

		scale := float32(4)

		x := float64(trig.Min.X)
		y := float64(trig.Min.Y)

		w, h := trig.Dx(), trig.Dy()

		screenX := float32(x+cam.X) * scale
		screenY := float32(y+cam.Y) * scale
		width := float32(w) * scale
		height := float32(h) * scale

		vector.StrokeRect(
			screen,
			screenX,
			screenY,
			width,
			height,
			1.0,
			color.RGBA{100, 255, 100, 255},
			false,
		)
	}
}
