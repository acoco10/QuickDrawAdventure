package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
	"os"
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
	interactInProximity gameObjects.Item
	triggerInteraction  bool
	dustEffect          *ebiten.Image
}

func NewTownScene() *TownScene {
	ts := TownScene{}
	return &ts
}

func (g *TownScene) FirstLoad() {

	tileMapFile, err := os.ReadFile("assets/map/town1Map.json")
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

	g.tilemapJSON = tilemapJSON
	g.tilesets = tilesets

	g.MapData, err = gameObjects.LoadMapObjectData(*tilemapJSON)
	g.Objects, err = gameObjects.LoadMapObjects(g.MapData)
	g.npcInProximity = gameObjects.Character{}

	if err != nil {
		log.Fatal(err)
	}

	g.loaded = true

	g.cam = camera.NewCamera(0.0, 0.0)

	charSpriteSheet := spritesheet.NewSpritesheet(4, 6, 16, 32)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/elyse/elyseSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}

	bethAnneImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/characters/npc/townFolk/bethAnne.png")
	if err != nil {
		log.Fatal(err)
	}

	jarvisImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/npc/townFolk/jarvis.png")
	if err != nil {
		log.Fatal(err)
	}

	zephImg, _, err := ebitenutil.NewImageFromFile("assets/images/characters/npc/townFolk/zeph.png")
	if err != nil {
		log.Fatal(err)
	}

	npcSpawn := g.MapData.NpcSpawns

	jarvis, err := gameObjects.NewCharacter(jarvisImg, npcSpawn["jarvis"], *charSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}

	bethSpriteSheet := spritesheet.NewSpritesheet(3, 1, 18, 25)

	bethAnne, err := gameObjects.NewCharacter(bethAnneImg, npcSpawn["bethAnne"], *bethSpriteSheet, gameObjects.NonPlayer)
	if err != nil {
		log.Fatal(err)
	}

	zephSpriteSheet := spritesheet.NewSpritesheet(3, 1, 15, 29)

	zeph, err := gameObjects.NewCharacter(zephImg, npcSpawn["zeph"], *zephSpriteSheet, gameObjects.NonPlayer)

	player, err := gameObjects.NewCharacter(playerImg, npcSpawn["player"], *charSpriteSheet, gameObjects.Player)
	if err != nil {
		log.Fatal(err)
	}

	g.Player = player

	g.NPC = map[string]*gameObjects.Character{}

	g.NPC["bethAnne"] = bethAnne
	g.NPC["jarvis"] = jarvis
	g.NPC["zeph"] = zeph

	g.dialogueUi, err = MakeDialogueUI(1512, 918)
	if err != nil {
		log.Fatal(err)
	}
	g.dialogueUi.UpdateTriggerScene(sceneManager.BattleSceneId)

}

func (g *TownScene) Update() sceneManager.SceneId {
	g.dialogueUi.Update()

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

	if ebiten.IsKeyPressed(ebiten.KeyE) && g.npcInProximity.Name != "" {
		g.dialogueUi.LoadDialogueUI(g.npcInProximity.Name)
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
	return g.dialogueUi.TriggerScene()

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
		if object.DrawAbovePlayer {
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

	}
	gameObjects.DrawMapAbovePlayer(*g.tilemapJSON, g.tilesets, *g.cam, screen, *g.Player, g.MapData.StairTriggers)

	opts.GeoM.Scale(4, 4)
	opts.GeoM.Translate(g.cam.X, g.cam.Y)

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

	if g.triggerInteraction {
		opts.GeoM.Translate(g.cam.X, g.cam.Y)
		opts.GeoM.Translate(g.interactInProximity.X, g.interactInProximity.Y)
		opts.GeoM.Scale(4, 4)
		screen.DrawImage(g.interactInProximity.Img, &opts)
		g.Player.InAnimation = true
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

}
func (g *TownScene) OnExit() {

}
func (g *TownScene) IsLoaded() bool {
	return g.loaded

}
