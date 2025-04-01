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
	"github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

type TownScene struct {

	//gameScenes elements
	Player              *gameObjects.Character
	NPC                 map[string]*gameObjects.Character
	tilemapJSON         *gameObjects.TilemapJSON
	tilesets            []gameObjects.Tileset
	cam                 *camera.Camera
	MapData             gameObjects.MapObjectData
	Objects             []*gameObjects.DoorObject
	debugCollisionMode  bool
	dialogueUi          *DialogueUI
	MainMenu            *ui.MainMenu
	ReadingMenu         *ui.TextBlockMenu
	loaded              bool
	cursor              *ui.CursorUpdater
	npcInProximity      gameObjects.Character
	interactInProximity *gameObjects.Trigger
	triggerInteraction  bool
	scene               sceneManager.SceneId
	gameLog             *sceneManager.GameLog
	enemyCountDown      int
	dark                bool
	InMenu              bool
	ResolutionHeight    int
	ResolutionWidth     int
}

func NewTownScene() *TownScene {
	ts := TownScene{}
	return &ts
}

func (g *TownScene) FirstLoad(gameLog *sceneManager.GameLog) {
	g.ResolutionWidth = 1512
	g.ResolutionHeight = 918

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

	martha, err := gameObjects.NewCharacter(npcSpawn["marthaFelten"], *marthaSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}

	player, err := gameObjects.NewCharacter(npcSpawn["elyse"], *charSpriteSheet, gameObjects.Player)
	if err != nil {
		log.Fatal(err)
	}

	antonio, err := gameObjects.NewCharacter(npcSpawn["antonio"], *charSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}
	bMSpriteSheet := spritesheet.NewSpritesheet(1, 1, 27, 35)
	boneMan, err := gameObjects.NewCharacter(npcSpawn["boneMan"], *bMSpriteSheet, gameObjects.NonPlayer)

	oMLSpriteSheet := spritesheet.NewSpritesheet(3, 1, 14, 18)
	oldManLandry, err := gameObjects.NewCharacter(npcSpawn["oldManLandry"], *oMLSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}

	georgeFelten, _ := gameObjects.NewCharacter(npcSpawn["georgeFelten"], *oMLSpriteSheet, gameObjects.NonPlayer)
	adahFelten, _ := gameObjects.NewCharacter(npcSpawn["adahFelten"], *oMLSpriteSheet, gameObjects.NonPlayer)

	g.Player = player

	g.gameLog.PlayerStats = player.BattleStats

	g.NPC = map[string]*gameObjects.Character{}

	g.NPC["bethAnne"] = bethAnne
	g.NPC["jarvis"] = jarvis
	g.NPC["zeph"] = zeph
	g.NPC["martha"] = martha
	g.NPC["antonio"] = antonio
	g.NPC["oldManLandry"] = oldManLandry
	g.NPC["boneMan"] = boneMan
	g.NPC["georgeFelten"] = georgeFelten
	g.NPC["adahFelten"] = adahFelten

	g.dialogueUi, err = MakeDialogueUI(1512, 918)
	if err != nil {
		log.Fatal(err)
	}

	g.dialogueUi.UpdateTriggerScene(sceneManager.TownSceneID)

	g.cursor = ui.CreateCursorUpdater(1512, 918)
	input.SetCursorUpdater(g.cursor)

	g.MainMenu = ui.NewMainMenu(g.ResolutionHeight, g.ResolutionWidth, g.Player.BattleStats, g.cursor)
	g.MainMenu.Load()
	g.MainMenu.SetCursor()
	g.ReadingMenu = &ui.TextBlockMenu{}
	g.ReadingMenu.Init()

}

func (g *TownScene) Update() sceneManager.SceneId {
	sceneUpdate, completed := g.dialogueUi.Update()

	if completed {
		g.InMenu = false
	}
	if sceneUpdate != sceneManager.TownSceneID {
		g.scene = sceneManager.BattleSceneId
	}

	g.Player.Dx = 0

	g.Player.Dy = 0

	//react to key presses by adding directional velocity
	if !g.Player.InAnimation && !g.InMenu {
		g.PlayerMovementInput()
	}

	g.MenuInput()

	//increase players position by their velocity every update
	g.Player.X += g.Player.Dx

	gameObjects.CheckCollisionHorizontal(g.Player.Sprite, g.MapData.Colliders, g.NPC)
	g.UpdateTriggerColliders("horizontal")

	g.Player.Y += g.Player.Dy

	gameObjects.CheckCollisionVertical(g.Player.Sprite, g.MapData.Colliders, g.NPC)
	g.UpdateTriggerColliders("vert")

	if !g.Player.InAnimation {
		gameObjects.CheckTriggers(g.Player, g.MapData.Triggers)
		gameObjects.CheckDoors(g.Player, g.MapData.Doors)
	}

	g.UpdateDoors()

	playerActiveAnimation := g.Player.ActiveAnimation(g.Player.Dx, g.Player.Dy)
	if playerActiveAnimation != nil {
		playerActiveAnimation.Update()
	}

	bethAnneAnimation := g.NPC["bethAnne"].ActiveAnimation(0, 0)
	bethAnneAnimation.Update()

	oldManActiveAnimation := g.NPC["oldManLandry"].ActiveAnimation(0, 0)
	oldManActiveAnimation.Update()

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

	g.interactInProximity = CheckInteractPopup(*g.Player, g.MapData.InteractPoints)
	g.CheckForNPCInteraction()

	enemyEncounter := battleStats.None
	if g.Player.Dx > 0 || g.Player.Dy > 0 {
		enemyEncounter = gameObjects.CheckEnemyTrigger(g.Player, g.MapData.EnemySpawns, g.enemyCountDown)
	}

	if enemyEncounter != battleStats.None {
		g.enemyCountDown = 0
		g.gameLog.EnemyEncountered = enemyEncounter
		enemyEncounter = battleStats.None
		g.scene = sceneManager.BattleSceneId
	}

	g.MainMenu.Update()
	g.ReadingMenu.Update()
	return g.scene

}

// Draw screen + sprites
func (g *TownScene) Draw(screen *ebiten.Image) {
	//map
	//loop through the tile map

	gameObjects.DrawMapBelowPlayer(*g.tilemapJSON, g.tilesets, *g.cam, screen, g.dark)

	//draw Player
	g.DrawItems(screen)
	g.DrawObjects(screen)
	g.DrawCharacters(screen)

	gameObjects.DrawMapAbovePlayer(*g.tilemapJSON, g.tilesets, *g.cam, screen, *g.Player, g.MapData.LayerTriggers, g.dark)
	g.DrawObjectsAbovePlayer(screen)

	if g.npcInProximity.Name != "" {
		DrawPopUp(screen, g.npcInProximity.X, g.npcInProximity.Y, float64(g.npcInProximity.SpriteSheet.SpriteWidth), g.cam)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	if g.dark {
		overlay := ebiten.NewImage(screenWidth, screenHeight)
		overlay.Fill(color.RGBA{0, 0, 10, 150}) // 50 = slight darkening
		screen.DrawImage(overlay, &ebiten.DrawImageOptions{})
	}

	if g.interactInProximity != nil && g.interactInProximity.Name != "" && !g.triggerInteraction {
		width := float64(g.interactInProximity.Rect.Bounds().Max.X - g.interactInProximity.Rect.Bounds().Min.X)
		DrawPopUp(screen, float64(g.interactInProximity.Rect.Min.X), float64(g.interactInProximity.Rect.Min.Y), width, g.cam)
	}
	gameObjects.DrawEvent(*g.tilemapJSON, g.tilesets, *g.cam, screen, *g.Player, g.MapData.LayerTriggers)
	g.MainMenu.Draw(screen)
	g.ReadingMenu.Draw(screen)
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
