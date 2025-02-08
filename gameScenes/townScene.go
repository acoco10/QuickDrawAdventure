package gameScenes

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
	"log"
)

type TownScene struct {

	//gameScenes elements
	Player              *gameObjects.Character
	NPC                 map[string]*gameObjects.Character
	tilemapJSON         *gameObjects.TilemapJSON
	tilesets            []gameObjects.Tileset
	cam                 *camera.Camera
	colliders           []image.Rectangle
	MapData             gameObjects.MapObjectData
	Objects             []*gameObjects.DoorObject
	action              bool
	debugCollisionMode  bool
	dialogueUi          *DialogueUI
	loaded              bool
	cursor              BattleMenuCursorUpdater
	npcInProximity      gameObjects.Character
	interactInProximity gameObjects.MapItem
	triggerInteraction  bool
	dustEffect          *ebiten.Image
	scene               sceneManager.SceneId
	gameLog             *sceneManager.GameLog
}

func NewTownScene() *TownScene {
	ts := TownScene{}
	return &ts
}

func (g *TownScene) FirstLoad(gameLog *sceneManager.GameLog) {

	tileMapFile, err := assets.Map.ReadFile("map/town1Map.json")
	if err != nil {
		log.Fatal("tilemap file not loading", err)
	}

	tilemapJSON, err := gameObjects.NewTilemapJSON(tileMapFile)
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	tilesets, err := tilemapJSON.GenTileSets()

	if err != nil {
		log.Fatal(err)
	}
	g.gameLog = gameLog
	g.tilemapJSON = tilemapJSON
	g.tilesets = tilesets

	g.MapData, err = gameObjects.LoadMapObjectData(*tilemapJSON)
	g.Objects, err = gameObjects.LoadMapObjects(g.MapData)
	g.npcInProximity = gameObjects.Character{}
	g.scene = sceneManager.TownSceneID
	if err != nil {
		log.Fatal(err)
	}

	g.loaded = true

	g.cam = camera.NewCamera(0.0, 0.0)
	g.cam.UpdateState(camera.Outside)

	charSpriteSheet := spritesheet.NewSpritesheet(4, 6, 16, 32)

	npcSpawn := g.MapData.NpcSpawns

	jarvis, err := gameObjects.NewCharacter(npcSpawn["jarvis"], *charSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}

	bethSpriteSheet := spritesheet.NewSpritesheet(3, 1, 18, 25)

	bethAnne, err := gameObjects.NewCharacter(npcSpawn["bethAnne"], *bethSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}

	zephSpriteSheet := spritesheet.NewSpritesheet(3, 1, 15, 29)

	zeph, err := gameObjects.NewCharacter(npcSpawn["zeph"], *zephSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}
	marthaSpriteSheet := spritesheet.NewSpritesheet(1, 1, 15, 27)

	martha, err := gameObjects.NewCharacter(npcSpawn["marthaJean"], *marthaSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}

	player, err := gameObjects.NewCharacter(npcSpawn["elyse"], *charSpriteSheet, gameObjects.Player)
	if err != nil {
		log.Fatal(err)
	}

	antonio, err := gameObjects.NewCharacter(npcSpawn["antonio"], *charSpriteSheet, gameObjects.Player)
	if err != nil {
		log.Fatal(err)
	}

	g.Player = player

	g.gameLog.PlayerStats = player.BattleStats

	g.NPC = map[string]*gameObjects.Character{}

	g.NPC["bethAnne"] = bethAnne
	g.NPC["jarvis"] = jarvis
	g.NPC["zeph"] = zeph
	g.NPC["martha"] = martha
	g.NPC["antonio"] = antonio

	g.dialogueUi, err = MakeDialogueUI(1512, 918)
	if err != nil {
		log.Fatal(err)
	}

	g.dialogueUi.UpdateTriggerScene(sceneManager.TownSceneID)

}

func (g *TownScene) Update() sceneManager.SceneId {

	sceneUpdate := g.dialogueUi.Update()
	if sceneUpdate != sceneManager.TownSceneID {
		g.scene = sceneManager.BattleSceneId
	}

	g.Player.Dx = 0

	g.Player.Dy = 0

	//react to key presses by adding directional velocity
	if !g.Player.InAnimation {
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.Player.Dx = 1.5
			g.Player.Direction = "L"
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.Player.Dx = -1.5
			g.Player.Direction = "R"
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			g.Player.Dy = 1.5
			g.Player.Direction = "U"
		}
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			g.Player.Dy = -1.5
			g.Player.Direction = "D"
		}

	}

	if ebiten.IsKeyPressed(ebiten.KeyE) && g.npcInProximity.Name != "" {
		g.dialogueUi.LoadUI(g.npcInProximity.Name)
		if g.npcInProximity.Name == "marthaJean" {
			g.NPC["antonio"].Spawned = true
		}
		if g.npcInProximity.Name == "antonio" {
			g.dialogueUi.UpdateTriggerScene(sceneManager.BattleSceneId)
			g.gameLog.EnemyEncountered = battleStats.Sheriff
		}
		LockCursorForDialogue()
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) && g.interactInProximity.Name != "" {
		g.Player.Visible = false
		g.triggerInteraction = true
	}

	//increase players position by their velocity every update
	g.Player.X += g.Player.Dx

	gameObjects.CheckCollisionHorizontal(g.Player.Sprite, g.MapData.Colliders, g.NPC)

	g.Player.Y += g.Player.Dy

	gameObjects.CheckCollisionVertical(g.Player.Sprite, g.MapData.Colliders, g.NPC)

	playerOnEntDoor := make(map[string]bool)
	playerOnExDoor := make(map[string]bool)
	playerOnContextTrigger := make(map[string]gameObjects.ObjectState)

	if !g.Player.InAnimation {
		playerOnEntDoor = gameObjects.CheckDoor(g.Player, g.MapData.EntryDoors)
	}

	if !g.Player.InAnimation {
		playerOnExDoor = gameObjects.CheckDoor(g.Player, g.MapData.ExitDoors)
	}

	playerOnContextTrigger = gameObjects.CheckContextualTriggers(g.Player, g.MapData.ContextualObjects)

	if !g.Player.InAnimation {
		gameObjects.CheckStairs(g.Player, g.MapData.StairTriggers)
	}

	if !g.Player.InAnimation {
		gameObjects.CheckContextualTriggers(g.Player, g.MapData.ContextualObjects)
	}

	playerActiveAnimation := g.Player.ActiveAnimation(int(g.Player.Dx), int(g.Player.Dy))
	if playerActiveAnimation != nil {
		playerActiveAnimation.Update()
	}

	bethAnneAnimation := g.NPC["bethAnne"].ActiveAnimation(int(g.Player.Dx), int(g.Player.Dy))
	bethAnneAnimation.Update()

	g.UpdateDoors()

	//updating camera to Player position
	g.cam.FollowTarget(g.Player.X-16, g.Player.Y, 320, 180)

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

		if playerOnExDoor[object.Name] && object.Status == gameObjects.NotTriggered {
			println("playerOnExDoor:", object.Name)
			g.Player.InAnimation = true
			object.Status = gameObjects.Leaving
		}
		if playerOnEntDoor[object.Name] && object.Status == gameObjects.NotTriggered {
			println("playerOnEntDoor:", object.Name)
			g.Player.InAnimation = true
			object.Status = gameObjects.Entering
		}
		if object.Type == gameObjects.ContextualObject {
			object.Status = playerOnContextTrigger[object.Name]
		}
	}

	//custom script animation for tavern door (swings forward on entrance)
	g.npcInProximity = CheckDialoguePopup(*g.Player, g.NPC)
	g.interactInProximity = CheckInteractPopup(*g.Player, g.MapData.InteractPoints)
	enemyEncounter := battleStats.None
	if g.Player.Dx > 0 || g.Player.Dy > 0 {
		enemyEncounter = gameObjects.CheckEnemyTrigger(g.Player, g.MapData.EnemySpawns)
	}

	if enemyEncounter != battleStats.None {
		g.gameLog.EnemyEncountered = enemyEncounter
		enemyEncounter = battleStats.None
		g.scene = sceneManager.BattleSceneId
	}
	return g.scene

}

// Draw screen + sprites
func (g *TownScene) Draw(screen *ebiten.Image) {

	opts := ebiten.DrawImageOptions{}

	//map
	//loop through the tile map

	gameObjects.DrawMapBelowPlayer(*g.tilemapJSON, g.tilesets, *g.cam, screen, g.MapData.StairTriggers)

	//draw Player

	for _, item := range g.MapData.Items {
		opts.GeoM.Reset()
		opts.GeoM.Translate(item.X*4, item.Y*4-137)
		opts.GeoM.Translate(g.cam.X*4, g.cam.Y*4)
		screen.DrawImage(item.Img, &opts)
		opts.GeoM.Reset()
	}
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
	for _, object := range g.Objects {
		if object.DrawAbovePlayer && g.Player.Y+6 < object.Y {
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

		}

		opts.GeoM.Reset()
	}

	if g.triggerInteraction {
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
		opts.GeoM.Translate(g.interactInProximity.X, g.interactInProximity.Y)
		opts.GeoM.Scale(4, 4)
		screen.DrawImage(g.interactInProximity.Img, &opts)
		g.Player.InAnimation = true
	}

	gameObjects.DrawMapAbovePlayer(*g.tilemapJSON, g.tilesets, *g.cam, screen, *g.Player, g.MapData.StairTriggers)

	/*	for _, door := range g.MapData.EntryDoors {
			vector.StrokeRect(
				screen,
				float32(door.Rect.Min.X)*4+float32(g.cam.X)*4,
				float32(door.Rect.Min.Y)*4+float32(g.cam.Y)*4,
				float32(door.Rect.Dx())*4,
				float32(door.Rect.Dy())*4,
				1.0,
				color.RGBA{255, 0, 0, 255},
				false,
			)
		}

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
		}*/
	if g.npcInProximity.Name != "" {
		DrawPopUp(screen, g.npcInProximity.X, g.npcInProximity.Y, float64(g.npcInProximity.SpriteSheet.SpriteWidth), g.cam)
	}

	if g.interactInProximity.Name != "" && !g.triggerInteraction {
		width := float64(g.interactInProximity.Img.Bounds().Max.X - g.interactInProximity.Img.Bounds().Min.X)
		DrawPopUp(screen, g.interactInProximity.X, g.interactInProximity.Y, width, g.cam)
	}

	err := g.dialogueUi.Draw(screen)
	if err != nil {
		log.Fatal(err)
	}

}

func (g *TownScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *TownScene) OnEnter() {
	if g.gameLog.PreviousScene == sceneManager.StartSceneId {
		g.FirstLoad(g.gameLog)
	}
}
func (g *TownScene) OnExit() {
	g.scene = sceneManager.TownSceneID

}
func (g *TownScene) IsLoaded() bool {
	return g.loaded

}

func DebugCamera(cam camera.Camera, player gameObjects.Character, screen *ebiten.Image) {
	face, err := assetManagement.LoadFont(40, assetManagement.November)
	if err != nil {
		log.Fatal(err)
	}
	cameraDebugText := fmt.Sprintf("X:%f Y:%f", cam.X, cam.Y)
	playerPosition := fmt.Sprintf("player X:%f player Y:%f", player.X, player.Y)
	dopts := text.DrawOptions{}
	dopts.DrawImageOptions.ColorScale.Scale(1, 0, 0, 255)
	dopts.GeoM.Translate(player.X*4+cam.X*4, player.Y*4+cam.Y*4)
	text.Draw(screen, cameraDebugText, face, &dopts)
	dopts.GeoM.Translate(0, 50)
	text.Draw(screen, playerPosition, face, &dopts)

}
