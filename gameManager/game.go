package gameManager

import (
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/gameScenes"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type Game struct {
	sceneMap      map[sceneManager.SceneId]sceneManager.Scene
	activeSceneId sceneManager.SceneId
	gameLog       *sceneManager.GameLog
}

func NewGame() *Game {
	sceneMap := map[sceneManager.SceneId]sceneManager.Scene{
		sceneManager.BattleSceneId:   gameScenes.NewBattleScene(),
		sceneManager.StartSceneId:    gameScenes.NewStartScene(),
		sceneManager.GameOverSceneID: gameScenes.NewGameOverScene(),
		sceneManager.WinSceneID:      gameScenes.NewWinScene(),
		sceneManager.TownSceneID:     gameScenes.NewTownScene(),
	}
	activeSceneId := sceneManager.StartSceneId

	game := &Game{
		sceneMap,
		activeSceneId,
		&sceneManager.GameLog{},
	}
	sceneMap[activeSceneId].FirstLoad(game.gameLog)
	return game
}

func (g *Game) Update() error {
	nextSceneId := g.sceneMap[g.activeSceneId].Update()
	// switched scenes
	if nextSceneId != g.activeSceneId {
		g.gameLog.PreviousScene = g.activeSceneId
		g.sceneMap[g.activeSceneId].OnExit()
		nextScene := g.sceneMap[nextSceneId]
		// if not loaded? then load in
		if !nextScene.IsLoaded() {
			nextScene.FirstLoad(g.gameLog)
		}
		nextScene.OnEnter()
	}
	g.activeSceneId = nextSceneId
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.sceneMap[g.activeSceneId].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1512, 982
}

func NewBattleTestGame() *Game {
	sceneMap := map[sceneManager.SceneId]sceneManager.Scene{
		sceneManager.BattleSceneId:   gameScenes.NewBattleScene(),
		sceneManager.GameOverSceneID: gameScenes.NewGameOverScene(),
		sceneManager.WinSceneID:      gameScenes.NewWinScene(),
		sceneManager.StartSceneId:    gameScenes.NewStartScene(),
	}
	elyseStats, err := battleStats.LoadSingleCharacter("elyse")
	if err != nil {
		log.Fatal(err)
	}
	activeSceneId := sceneManager.BattleSceneId
	glog := &sceneManager.GameLog{EnemyEncountered: battleStats.Sheriff,
		PlayerStats: &elyseStats}
	game := Game{
		sceneMap,
		activeSceneId,
		glog,
	}
	sceneMap[activeSceneId].FirstLoad(game.gameLog)
	return &game
}

func NewWebTestGame() *Game {
	sceneMap := map[sceneManager.SceneId]sceneManager.Scene{
		sceneManager.TestSceneID: gameScenes.NewTestScene(),
	}
	elyseStats, err := battleStats.LoadSingleCharacter("elyse")
	if err != nil {
		log.Fatal(err)
	}
	activeSceneId := sceneManager.TestSceneID
	glog := &sceneManager.GameLog{EnemyEncountered: battleStats.Sheriff,
		PlayerStats: &elyseStats}
	game := Game{
		sceneMap,
		activeSceneId,
		glog,
	}
	sceneMap[activeSceneId].FirstLoad(game.gameLog)
	return &game
}
