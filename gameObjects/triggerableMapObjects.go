package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"image"
	"math/rand/v2"
)

type Door struct {
	Key   string
	Coord image.Rectangle
}

func NewDoor(obj ObjectJSON) (door Door) {
	door =
		Door{
			Key: obj.Name,
			Coord: image.Rect(
				int(obj.X),
				int(obj.Y)-20,
				int(obj.Width+obj.X),
				int(obj.Y)-10,
			),
		}

	return door
}

func CheckDoor(player *Character, entDoors map[string]Trigger) map[string]bool {
	returnMap := make(map[string]bool)
	for _, door := range entDoors {
		if door.Rect.Overlaps(
			image.Rect(
				int(player.X),
				int(player.Y)+28,
				int(player.X)+16,
				int(player.Y)+31),
		) {
			player.EnterShadow()
			returnMap[door.Name] = true
		} else {
			returnMap[door.Name] = false
		}
	}
	return returnMap
}

func CheckContextualTriggers(player *Character, contextTriggers map[string]*Trigger) map[string]ObjectState {
	playerRect := image.Rect(
		int(player.X),
		int(player.Y)+20,
		int(player.X)+16,
		int(player.Y)+31)
	returnMap := make(map[string]ObjectState)
	for _, trig := range contextTriggers {
		if trig.Rect.Overlaps(playerRect) {
			returnMap[trig.Name] = On
		}

		if !trig.Rect.Overlaps(playerRect) {
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
		if trig.Rect.Overlaps(playerRect) {
			countDown--
			if rand.IntN(1000-countDown) <= 2 {
				return battleStats.Wolf
			}
		}
	}
	return battleStats.None
}

func GetDoorCoord(doors map[string]Trigger, key string, direction string) (float64, float64) {
	x := 0.0
	y := 0.0

	if direction == "up" {
		x = float64(doors[key].Rect.Dx()/2+doors[key].Rect.Min.X) - 8
		y = float64(doors[key].Rect.Min.Y) - 32
	}
	if direction == "down" {
		x = float64(doors[key].Rect.Dx()/2+doors[key].Rect.Min.X) - 8
		y = float64(doors[key].Rect.Max.Y) - 16
	}

	return x, y

}

func CheckStairs(player *Character, stairTriggers map[string]*Trigger) {
	for _, stair := range stairTriggers {
		if stair.Rect.Overlaps(
			image.Rect(
				int(player.X),
				int(player.Y)+28,
				int(player.X)+16,
				int(player.Y)+31),
		) {
			if player.Dy < 0 {
				stair.Triggered = true
			}
			if player.Dy > 0 {
				stair.Triggered = false
			}
		}
	}
}
