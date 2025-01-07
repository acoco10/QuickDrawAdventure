package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"os"
)

type TownScene struct {

	//gameScenes elements
	Player             *gameObjects.Character
	NPC                []*gameObjects.Character
	npcSpawns          map[string]gameObjects.NPCspawn
	tilemapJSON        *gameObjects.TilemapJSON
	tilesets           []gameObjects.Tileset
	cam                *camera.Camera
	colliders          []image.Rectangle
	Objects            []*gameObjects.Object
	EntranceDoors      map[string]gameObjects.Door
	ExitDoors          map[string]gameObjects.Door
	action             bool
	debugCollisionMode bool
	dialoguei          *DialogueUI
	loaded             bool
	cursor             BattleMenuCursorUpdater
}

func NewTownScene() *TownScene {
	ts := TownScene{}
	return &ts
}

func (g *TownScene) FirstLoad() {
	g.loaded = true

	tavernDoorSpriteSheet := spritesheet.NewSpritesheet(2, 2, 20, 21)

	tavernDoorImg, _, err := ebitenutil.NewImageFromFile("assets/images/buildings/tavernDoorSpriteSheet.png")

	if err != nil {
		log.Fatal(err)
	}
	tavernDoorAnimation := animations.NewAnimation(0, 6, 1, 10.0)
	tavernDoorObject, _ := gameObjects.NewObject(tavernDoorImg, 167.18, 158.76, *tavernDoorSpriteSheet, tavernDoorAnimation, tavernDoorAnimation, "door1")

	sunRiseDoorSpriteSheet := spritesheet.NewSpritesheet(3, 1, 24, 38)
	sunriseDoorImg, _, err := ebitenutil.NewImageFromFile("assets/images/buildings/sunriseInn/sunriseInnDoor.png")

	if err != nil {
		log.Fatal(err)
	}

	sunriseDoorAnimation := animations.NewAnimation(0, 3, 1, 10.0)
	sunriseDoor, _ := gameObjects.NewObject(sunriseDoorImg, 557, 145, *sunRiseDoorSpriteSheet, sunriseDoorAnimation, sunriseDoorAnimation, "door2")

	g.Objects = append(g.Objects, tavernDoorObject, sunriseDoor)

	tileMapFile, err := os.ReadFile("assets/map/town1Map.json")
	if err != nil {
		log.Fatal("tilemap file not loading", err)
	}
	tilemapJSON, err := gameObjects.NewTilemapJSON(tileMapFile)
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	tileset, err := tilemapJSON.GenTileSets()
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	g.tilesets = tileset
	g.tilemapJSON = tilemapJSON

	g.cam = camera.NewCamera(0.0, 0.0)

	colliders, entranceDoors, exitDoors, npcSpawns := gameObjects.StoreMapObjects(*tilemapJSON)

	g.colliders = colliders
	g.EntranceDoors = entranceDoors
	g.ExitDoors = exitDoors
	g.npcSpawns = npcSpawns

	charSpriteSheet := spritesheet.NewSpritesheet(4, 6, 16, 32)

	playerimg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/elyse/elyseSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	player, _ := gameObjects.NewCharacter(playerimg, 75, 75, *charSpriteSheet, "elyse")
	g.Player = player

	NPCImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/npc/TownPerson1.png")
	if err != nil {
		log.Fatal(err)
	}

	npc, err := gameObjects.NewCharacter(NPCImg, g.npcSpawns["Jarvis"].X, g.npcSpawns["Jarvis"].Y-64, *charSpriteSheet, "Jarvis")
	if err != nil {
		log.Fatal(err)
	}

	g.NPC = append(g.NPC, npc)

	g.dialoguei, err = MakeDialogueUI(1512, 918)
	if err != nil {
		log.Fatal(err)
	}
	g.dialoguei.UpdateTriggerScene(sceneManager.BattleSceneId)

}

func (g *TownScene) Update() sceneManager.SceneId {
	err := g.dialoguei.UpdateDialogueUI()
	if err != nil {
		log.Fatal(err)
	}
	g.Player.Dx = 0
	g.Player.Dy = 0
	//react to key presses by adding directional velocity
	if !g.Player.InAnimation {
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.Player.Dx = 1.8
			g.Player.Direction = "L"
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.Player.Dx = -1.8
			g.Player.Direction = "R"
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.Player.Dy = 1.8
			g.Player.Direction = "U"
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.Player.Dy = -1.8
			g.Player.Direction = "D"
		}

	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.dialoguei.Trigger()
		//LockCursorForDialogue()
	}

	//increase players position by their velocity every update
	g.Player.X += g.Player.Dx

	gameObjects.CheckCollisionHorizontal(g.Player.Sprite, g.colliders, g.NPC)

	g.Player.Y += g.Player.Dy

	gameObjects.CheckCollisionVertical(g.Player.Sprite, g.colliders, g.NPC)

	playerOnEntDoor := make(map[string]bool)
	playerOnExDoor := make(map[string]bool)
	if !g.Player.InAnimation {
		playerOnEntDoor = gameObjects.CheckEntDoor(g.Player, g.EntranceDoors, g.ExitDoors)
	}
	if !g.Player.InAnimation {
		playerOnExDoor = gameObjects.CheckEntDoor(g.Player, g.ExitDoors, g.EntranceDoors)
	}

	/* //for _, enemy := range g.enemies {

	enemy.Dx = 0.0

	enemy.Dy = 0.0

	if enemy.FollowsPlayer {
		if enemy.X < g.Player.X {
			enemy.Dx += 1

		} else if enemy.X > g.Player.X {
			enemy.Dx -= 1
		}
		if enemy.Y < g.Player.Y {
			enemy.Dy += 1

		} else if enemy.Y > g.Player.Y {
			enemy.Dy -= 1
		}
	}
	enemy.X += enemy.Dx */

	//gameObjects.CheckCollisionHorizontal(enemy.Sprite, g.colliders)

	//enemy.Y += enemy.Dy

	//gameObjects.CheckCollisionVertical(enemy.Sprite, g.colliders)

	//checking active Player animation
	playerActiveAnimation := g.Player.ActiveAnimation(int(g.Player.Dx), int(g.Player.Dy))
	if playerActiveAnimation != nil {
		playerActiveAnimation.Update()
	}

	g.UpdateDoors()

	//updating camera to Player position
	g.cam.FollowTarget(g.Player.X+16, g.Player.Y+16, 320, 240)

	//when Player hits the edge of the map the camera does not follow
	//need to update this logic for interiors, new map?
	g.cam.Constrain(
		//width of maps from map JSON * tile size
		float64(g.tilemapJSON.Layers[0].Width)*16,
		float64(g.tilemapJSON.Layers[0].Height)*16,
		//screen resolution
		320,
		240,
	)

	//check if Player has entered a door and update door object eventually this will need to be a loop for all object animations
	for _, object := range g.Objects {
		if playerOnExDoor[object.Name] && object.Status == "" {
			g.Player.InAnimation = true
			object.Status = "leaving"
		}

		if playerOnEntDoor[object.Name] && object.Status == "" {
			g.Player.InAnimation = true
			object.Status = "entering"
		}
	}

	//custom script animation for tavern door (swings forward on entrance)

	return g.dialoguei.TriggerScene()
}

// Draw screen + sprites
func (g *TownScene) Draw(screen *ebiten.Image) {

	opts := ebiten.DrawImageOptions{}

	//map
	//loop through the tile map
	gameObjects.DrawMapBelowPlayer(*g.tilemapJSON, g.tilesets, *g.cam, screen)
	//draw Player
	for _, object := range g.Objects {
		opts.GeoM.Translate(object.X, object.Y)
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
		opts.GeoM.Scale(4, 4)

		objectFrame := 0
		objectAnimation := object.ActiveAnimation(object.Status)

		if objectAnimation != nil {
			objectFrame = objectAnimation.Frame()
		}

		screen.DrawImage(
			object.Img.SubImage(
				object.SpriteSheet.Rect(objectFrame),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	g.DrawCharacters(screen)
	gameObjects.DrawMapAbovePlayer(*g.tilemapJSON, g.tilesets, *g.cam, screen, *g.Player)

	opts.GeoM.Scale(4, 4)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)
	for _, door := range g.EntranceDoors {
		vector.StrokeRect(
			screen,
			float32(door.Coord.Min.X)*4+float32(g.cam.X)*4,
			float32(door.Coord.Min.Y)*4+float32(g.cam.Y)*4,
			float32(door.Coord.Dx())*4,
			float32(door.Coord.Dy())*4,
			1.0,
			color.RGBA{255, 0, 0, 255},
			false,
		)
	}

	for _, collider := range g.colliders {
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

	vector.StrokeRect(
		screen,
		float32(g.NPC[0].X)*4+float32(g.cam.X)*4,
		float32(g.NPC[0].Y)*4+float32(g.cam.Y)*4+28*4,
		16*4,
		4*4,
		1.0,
		color.RGBA{255, 0, 0, 255},
		false,
	)
	err := g.dialoguei.Draw(screen)
	if err != nil {
		log.Fatal(err)
	}
	for _, npc := range g.NPC {
		DrawDialoguePopup(*g.Player, *npc, screen, *g.cam)
	}
	return

}

func (g *TownScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *TownScene) OnEnter() {

}
func (g *TownScene) OnExit() {

}
func (g *TownScene) IsLoaded() bool {
	return g.loaded

}
