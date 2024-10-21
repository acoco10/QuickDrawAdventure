package mapobjects

import (
	"ShootEmUpAdventure/entities"
	"image"
)

func CheckEnterDoor(player *entities.Player, entdoors map[string]Door, exdoors map[string]Door) {
	for key, door := range entdoors {
		if door.Coord.Overlaps(
			image.Rect(
				int(player.X),
				int(player.Y)+28,
				int(player.X)+16,
				int(player.Y)+31),
		) {
			player.X = float64(exdoors[key].Coord.Max.X - (exdoors[key].Coord.Max.X-exdoors[key].Coord.Min.X)/2)
			player.Y = float64(exdoors[key].Coord.Max.Y) - 60
		}
	}
}

func CheckExitDoor(player *entities.Player, entdoors map[string]Door, exdoors map[string]Door) {
	for key, door := range exdoors {
		if door.Coord.Overlaps(
			image.Rect(
				int(player.X),
				int(player.Y)+28,
				int(player.X)+16,
				int(player.Y)+31),
		) {
			player.X = float64(entdoors[key].Coord.Min.X + (entdoors[key].Coord.Min.X-entdoors[key].Coord.Max.X)/2)
			player.Y = float64(entdoors[key].Coord.Min.Y + 20)
		}

	}
}
