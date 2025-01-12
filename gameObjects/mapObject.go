package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
)

type ObjectJSON struct {
	Name   string  `json:"name"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
	X      float64 `json:"x"`
	Y      float64 `json:"Y"`
	Class  string  `json:"class"`
	Type   string  `json:"type"`
}
type ObjectType uint8

const (
	EntryDoor ObjectType = iota
	ExitDoor
	StairTrigger
	ContextualObject
	Collider
	NpcSpawnObject
)

type NpcSpawn struct {
	Name string
	X, Y float64
}

type Trigger struct {
	Name      string
	Type      ObjectType
	Rect      image.Rectangle
	Triggered bool
}

func NewTrigger(json ObjectJSON) Trigger {
	newTrigger := new(Trigger)
	newTrigger.Name = json.Name
	rect := image.Rect(int(json.X), int(json.Y), int(json.X+json.Width), int(json.Y+json.Height))
	newTrigger.Rect = rect
	newTrigger.Triggered = false
	return *newTrigger
}

type MapObjectData struct {
	EntryDoors        map[string]Trigger
	ExitDoors         map[string]Trigger
	NpcSpawns         map[string]Trigger
	Colliders         []image.Rectangle
	StairTriggers     map[string]*Trigger //pointer because has on/ off setting eg for balcony, not trigger switch
	ContextualObjects map[string]*Trigger
}

func LoadMapObjectData(tilemapJSON TilemapJSON) (MapObjectData, error) {

	colliders := []image.Rectangle{}
	entDoors := make(map[string]Trigger)
	exDoors := make(map[string]Trigger)
	npcSpawns := make(map[string]Trigger)
	stairTriggers := make(map[string]*Trigger)
	contextualObjects := make(map[string]*Trigger)

	for _, layer := range tilemapJSON.Layers {
		if layer.Type == "objectgroup" {
			for _, object := range layer.Objects {
				switch object.Type {
				case "entranceDoor":
					println("loading entranceDoor:", object.Name)
					door := NewTrigger(object)
					door.Type = EntryDoor
					door.Rect.Min.Y = door.Rect.Min.Y - 32
					door.Rect.Max.Y = door.Rect.Max.Y - 32
					entDoors[object.Name] = door
				case "exitDoor":
					println("loading exitDoor:", object.Name)
					door := NewTrigger(object)
					door.Type = ExitDoor
					exDoors[object.Name] = door
				case "npcSpawn":
					println("loading npcSpawn:", object.Name)
					npcSpawn := NewTrigger(object)
					npcSpawn.Type = NpcSpawnObject
					npcSpawns[object.Name] = npcSpawn
				case "collider":
					rect := image.Rect(
						int(object.X),
						int(object.Y)-32,
						int(object.Width+object.X),
						int(object.Y+object.Height)-32,
					)
					colliders = append(colliders, rect)
				case "stairTrigger":
					stair := NewTrigger(object)
					stair.Type = StairTrigger
					stairTriggers[object.Name] = &stair

				case "contextualObject":
					println("loading contextualObject:", object.Name)
					contextualObject := NewTrigger(object)
					contextualObject.Type = ContextualObject
					contextualObject.Rect.Min.Y = contextualObject.Rect.Min.Y - 32
					contextualObject.Rect.Max.Y = contextualObject.Rect.Max.Y - 32
					contextualObjects[object.Name] = &contextualObject

				}
			}
		}
	}

	mapObjects := MapObjectData{}
	mapObjects.EntryDoors = entDoors
	mapObjects.ExitDoors = exDoors
	mapObjects.NpcSpawns = npcSpawns
	mapObjects.StairTriggers = stairTriggers
	mapObjects.Colliders = colliders
	mapObjects.ContextualObjects = contextualObjects

	return mapObjects, nil
}

func LoadMapObjects(mapObjectData MapObjectData) ([]*DoorObject, error) {

	objects := make([]*DoorObject, 0)

	tavernDoorSpriteSheet := spritesheet.NewSpritesheet(2, 3, 20, 21)
	tavernDoorImg, _, err := ebitenutil.NewImageFromFile("assets/images/buildings/tavernDoorSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}

	tavernDoorAnimation := animations.NewAnimation(0, 6, 1, 10.0)
	tavernDoorObject, _ := NewObject(
		tavernDoorImg,
		*tavernDoorSpriteSheet,
		tavernDoorAnimation,
		tavernDoorAnimation,
		mapObjectData.EntryDoors["door1"],
	)

	sunRiseDoorSpriteSheet := spritesheet.NewSpritesheet(3, 1, 24, 38)
	sunriseDoorImg, _, err := ebitenutil.NewImageFromFile("assets/images/buildings/sunriseInn/sunriseInnDoor.png")
	if err != nil {
		log.Fatal(err)
	}

	standardDoorAnimation := animations.NewAnimation(0, 3, 1, 10.0)
	sunriseDoor, err := NewObject(
		sunriseDoorImg,
		*sunRiseDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		mapObjectData.EntryDoors["door2"],
	)
	if err != nil {
		log.Fatal(err)
	}
	sunRiseSideDoorSpriteSheet := spritesheet.NewSpritesheet(3, 1, 17, 23)
	sideDoorImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/buildings/sunriseInn/sideBuildingDoor.png")
	if err != nil {
		log.Fatal(err)
	}

	sideDoor, err := NewObject(
		sideDoorImg,
		*sunRiseSideDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		mapObjectData.EntryDoors["door3"],
	)
	if err != nil {
		log.Fatal(err)
	}

	beadedCurtainSpriteSheet := spritesheet.NewSpritesheet(3, 1, 18, 24)
	beadedCurtainImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/buildings/beadedCurtain.png")
	if err != nil {
		log.Fatal(err)
	}

	beadedCurtainAnimation := animations.NewAnimation(0, 3, 1, 10.0)

	beadedCurtain, err := NewObject(
		beadedCurtainImg,
		*beadedCurtainSpriteSheet,
		beadedCurtainAnimation,
		beadedCurtainAnimation,
		*mapObjectData.ContextualObjects["beadedCurtain1"],
	)
	beadedCurtain.DrawAbovePlayer = true
	if err != nil {
		log.Fatal(err)
	}
	objects = append(objects, tavernDoorObject, sunriseDoor, sideDoor, beadedCurtain)
	return objects, err
}
