package gameObjects

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
)

type ObjectJSON struct {
	Name       string           `json:"name"`
	Height     float64          `json:"height"`
	Width      float64          `json:"width"`
	X          float64          `json:"x"`
	Y          float64          `json:"Y"`
	Class      string           `json:"class"`
	Type       string           `json:"type"`
	Properties []PropertiesJSON `json:"properties"`
}

type ObjectType uint8

const (
	EntryDoor ObjectType = iota
	ExitDoor
	InsideDoor
	StairTrigger
	ContextualObject
	Collider
	NpcSpawnObject
)

type Spawn struct {
	Name    string
	X, Y    float64
	spawned bool
}

type Trigger struct {
	Name        string
	Type        ObjectType
	Rect        image.Rectangle
	Triggered   bool
	Dir         Direction
	CameraPoint string
	Animation   string
	Sprite      string
}

func NewTrigger(json ObjectJSON) Trigger {
	newTrigger := new(Trigger)
	newTrigger.Name = json.Name
	rect := image.Rect(int(json.X), int(json.Y), int(json.X+json.Width), int(json.Y+json.Height))
	newTrigger.Rect = rect
	newTrigger.Triggered = false
	newTrigger.Rect.Min.Y = newTrigger.Rect.Min.Y - 32
	newTrigger.Rect.Max.Y = newTrigger.Rect.Max.Y - 32
	loadProperties(json.Properties, newTrigger)
	return *newTrigger
}

func loadProperties(props []PropertiesJSON, newTrigger *Trigger) {
	for _, prop := range props {
		switch prop.Name {
		case "direction":
			newTrigger.Dir = LoadDirection(prop.Value)
		case "cameraPoint":
			newTrigger.CameraPoint = prop.Value
		case "animation":
			newTrigger.Animation = prop.Value
		case "sprite":
			newTrigger.Sprite = prop.Value
		}
	}
}

type MapItem struct {
	State string
	Name  string
	X, Y  float64
	Img   *ebiten.Image
}

type MapObjectData struct {
	NpcSpawns        map[string]Spawn
	Colliders        []image.Rectangle
	Items            []MapItem
	InteractPoints   map[string]MapItem
	CameraPoints     map[string]Trigger
	ObjectSpawns     map[string]Spawn
	EnemySpawns      map[string]Trigger
	LayerTriggers    map[string]*Trigger
	TriggerColliders map[string][]image.Rectangle
	Triggers         []*Trigger
	Doors            []*DoorObject
}

// edit map load to store layer triggers and colldier in same place
type LayerTrigger struct {
	Trigger   *Trigger
	Colliders map[string][]image.Rectangle
}

func LoadDirection(dir string) Direction {
	switch dir {
	case "down":
		return Down
	case "up":
		return Up
	case "left":
		return Left
	case "right":
		return Right
	}
	return None
}

func LoadMapObjectData(tilemapJSON TilemapJSON) (MapObjectData, error) {

	colliders := []image.Rectangle{}
	var doors []*DoorObject
	npcSpawns := make(map[string]Spawn)
	contextualObjects := make(map[string]*Trigger)
	var items []MapItem
	layerTriggers := make(map[string]*Trigger)
	interactPoints := make(map[string]MapItem)
	cameraPoints := make(map[string]Trigger)
	objectSpawns := make(map[string]Spawn)
	enemies := make(map[string]Trigger)
	var allTrig []*Trigger
	triggerColliders := make(map[string][]image.Rectangle)

	for _, layer := range tilemapJSON.Layers {
		if layer.Name == "colliders" {
			for _, object := range layer.Objects {
				rect := image.Rect(
					int(object.X),
					int(object.Y)-32,
					int(object.Width+object.X),
					int(object.Y+object.Height)-32)
				colliders = append(colliders, rect)
			}
		}
		if layer.Name == "triggerColliders" {
			for _, object := range layer.Objects {
				rect := image.Rect(
					int(object.X),
					int(object.Y)-32,
					int(object.Width+object.X),
					int(object.Y+object.Height)-32)
				triggerColliders[object.Name] = append(triggerColliders[object.Name], rect)
				println("loaded Collision trigger:", object.Name)
			}
		}
		if layer.Name == "npcSpawns" {
			for _, object := range layer.Objects {
				println("loading npcSpawn:", object.Name)
				npcSpawn := Spawn{
					Name:    object.Name,
					X:       object.X,
					Y:       object.Y,
					spawned: true,
				}
				if object.Type == "conditionalNpcSpawn" {
					npcSpawn.spawned = false
				}
				npcSpawns[object.Name] = npcSpawn
			}
		}

		if layer.Name == "doors" {
			for _, object := range layer.Objects {
				door := LoadDoor(object)
				doors = append(doors, door)
			}
		}

		if layer.Type == "objectgroup" {
			for _, object := range layer.Objects {
				switch object.Type {
				case "objectSpawn":
					println("loading objectSpawn:", object.Name)
					objSpawn := Spawn{object.Name, object.X, object.Y - 32, true}
					objectSpawns[object.Name] = objSpawn
				case "contextualObject":
					println("loading contextualObject:", object.Name)
					contextualObject := NewTrigger(object)
					contextualObject.Type = ContextualObject
					contextualObjects[object.Name] = &contextualObject
					allTrig = append(allTrig, &contextualObject)
				case "itemSpawn":
					println("loading item:", object.Name)
					imgPath := fmt.Sprintf("images/items/%s.png", object.Name)
					img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
					if err != nil {
						log.Fatal(err)
					}
					item := MapItem{
						Name:  object.Name,
						X:     object.X,
						Y:     object.Y,
						Img:   img,
						State: "off",
					}
					items = append(items, item)
				case "interactPoint":
					println("loading interactPoint:", object.Name)
					var imgPath string
					for _, prop := range object.Properties {
						if prop.Name == "image" {
							imgPath = prop.Value
						}
					}
					img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
					if err != nil {
						println("interact image:", object.Name, "could not be loaded from", imgPath)
						img = ebiten.NewImage(10, 10)
					}
					item := MapItem{
						Name: object.Name,
						X:    object.X,
						Y:    object.Y - 32,
						Img:  img,
					}
					interactPoints[object.Name] = item
				case "cameraPoint":
					println("loading camera point:", object.Name)
					camPoint := NewTrigger(object)
					cameraPoints[object.Name] = camPoint
				case "enemySpawn":
					println("loading enemy spawn:", object.Name)
					enemySpawn := NewTrigger(object)
					enemies[object.Name] = enemySpawn
				case "layerTrigger":
					println("loading layer trigger", object.Name)
					layerTrigger := NewTrigger(object)
					layerTriggers[object.Name] = &layerTrigger
					allTrig = append(allTrig, &layerTrigger)
				}
			}
		}
	}

	mapObjects := MapObjectData{}
	mapObjects.NpcSpawns = npcSpawns
	mapObjects.Colliders = colliders
	mapObjects.Items = items
	mapObjects.InteractPoints = interactPoints
	mapObjects.CameraPoints = cameraPoints
	mapObjects.ObjectSpawns = objectSpawns
	mapObjects.EnemySpawns = enemies
	mapObjects.LayerTriggers = layerTriggers
	mapObjects.TriggerColliders = triggerColliders
	mapObjects.Triggers = allTrig
	mapObjects.Doors = doors
	return mapObjects, nil
}

func LoadDoor(object ObjectJSON) *DoorObject {
	door := NewTrigger(object)

	switch object.Type {
	case "entranceDoor":
		println("loading entranceDoor:", object.Name)
		door.Type = EntryDoor
	case "exitDoor":
		println("loading exitDoor:", object.Name)
		door.Type = ExitDoor
	case "insideDoor":
		println("loading insideDoor:", object.Name)
		door.Type = InsideDoor
	}

	spriteSheet, animation := GetAnimation(door.Animation)

	img := LoadDoorImage(door.Sprite)

	doorObject, err := NewObject(
		img,
		spriteSheet,
		animation,
		animation,
		&door,
	)

	if err != nil {
		log.Fatal(err)
	}

	return doorObject
}

func LoadDoorImage(doorSprite string) *ebiten.Image {
	standardDoorImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/buildings/sunriseInn/sideBuildingDoor.png")
	if err != nil {
		log.Fatal(err)
	}
	if doorSprite == "none" {
		return ebiten.NewImage(10, 10)
	}

	return standardDoorImg
}

func GetAnimation(doorType string) (spritesheet.SpriteSheet, *animations.Animation) {
	tavernDoorSpriteSheet := spritesheet.NewSpritesheet(2, 3, 20, 21)
	tavernDoorAnimation := animations.NewAnimation(0, 6, 1, 10.0)

	standardDoorAnimation := animations.NewAnimation(0, 3, 1, 10.0)
	standardDoorSpriteSheet := spritesheet.NewSpritesheet(3, 1, 18, 23)
	BigDoorSpriteSheet := spritesheet.NewSpritesheet(3, 1, 24, 38)

	switch doorType {
	case "big":
		return *BigDoorSpriteSheet, standardDoorAnimation
	case "tavern":
		return *tavernDoorSpriteSheet, tavernDoorAnimation
	default:
		return *standardDoorSpriteSheet, standardDoorAnimation
	}
}

func LoadMapObjects(entryDoors, exitDoors map[string]*Trigger, insideDoors []*Trigger, contextualObjects map[string]*Trigger, objSpawn map[string]Spawn) ([]*DoorObject, error) {

	objects := make([]*DoorObject, 0)
	//customDoorAnimations
	tavernDoorImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/buildings/tavernDoorSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	tavernDoorSpriteSheet := spritesheet.NewSpritesheet(2, 3, 20, 21)
	standardDoorSpriteSheet := spritesheet.NewSpritesheet(3, 1, 18, 23)
	standardDoorImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/buildings/sunriseInn/sideBuildingDoor.png")
	if err != nil {
		log.Fatal(err)
	}

	standardDoorAnimation := animations.NewAnimation(0, 3, 1, 10.0)

	BigDoorSpriteSheet := spritesheet.NewSpritesheet(3, 1, 24, 38)

	tavernDoorAnimation := animations.NewAnimation(0, 6, 1, 10.0)
	tavernDoorObject, _ := NewObject(
		tavernDoorImg,
		*tavernDoorSpriteSheet,
		tavernDoorAnimation,
		tavernDoorAnimation,
		entryDoors["rose"],
	)
	interiorDoorSpriteSheet := spritesheet.NewSpritesheet(3, 1, 18, 24)
	beadedCurtainImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/buildings/beadedCurtain.png")
	interiorDoorAnimation := animations.NewAnimation(0, 3, 1, 10.0)

	sunriseDoorImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/buildings/sunriseInn/sunriseInnDoor.png")
	if err != nil {
		log.Fatal(err)
	}

	sunriseDoor, err := NewObject(
		sunriseDoorImg,
		*BigDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		entryDoors["sunRise"],
	)

	if err != nil {
		log.Fatal(err)
	}

	beadedCurtain, err := NewObject(
		beadedCurtainImg,
		*interiorDoorSpriteSheet,
		interiorDoorAnimation,
		interiorDoorAnimation,
		contextualObjects["beadedCurtain1"],
	)

	beadedCurtain.X = objSpawn["beadedCurtain"].X
	beadedCurtain.Y = objSpawn["beadedCurtain"].Y

	beadedCurtain.DrawAbovePlayer = true
	if err != nil {
		log.Fatal(err)
	}

	sideDoor, err := NewObject(
		standardDoorImg,
		*standardDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		entryDoors["sideDoor"],
	)

	if err != nil {
		log.Fatal(err)
	}

	generalStoreDoor, err := NewObject(
		standardDoorImg,
		*standardDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		entryDoors["generalStore"],
	)

	farmCabinDoor, err := NewObject(
		standardDoorImg,
		*standardDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		entryDoors["farmCabin"],
	)

	caveDoor, err := NewObject(
		ebiten.NewImage(10, 10),
		*standardDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		entryDoors["cave"],
	)

	sunRiseRoof, err := NewObject(
		ebiten.NewImage(10, 10),
		*standardDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		entryDoors["sunRiseRoof"],
	)

	sunRiseOutsideBalcony, err := NewObject(
		ebiten.NewImage(10, 10),
		*standardDoorSpriteSheet,
		standardDoorAnimation,
		standardDoorAnimation,
		entryDoors["sunRiseOutsideBalcony"],
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, insideDoor := range insideDoors {
		door, err := NewObject(
			ebiten.NewImage(10, 10),
			*standardDoorSpriteSheet,
			standardDoorAnimation,
			standardDoorAnimation,
			insideDoor,
		)
		if err != nil {
			log.Fatal(err)
		}
		objects = append(objects, door)
	}

	for _, exitDoor := range exitDoors {
		door, err := NewObject(
			ebiten.NewImage(10, 10),
			*standardDoorSpriteSheet,
			standardDoorAnimation,
			standardDoorAnimation,
			exitDoor,
		)

		if err != nil {
			log.Fatal(err)
		}
		objects = append(objects, door)
	}

	objects = append(objects, tavernDoorObject, sunriseDoor, sideDoor, beadedCurtain, generalStoreDoor, farmCabinDoor, caveDoor, sunRiseRoof, sunRiseOutsideBalcony)
	return objects, err
}
