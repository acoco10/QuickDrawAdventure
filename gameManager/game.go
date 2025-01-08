package gameManager

import (
	"github.com/acoco10/QuickDrawAdventure/gameScenes"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	sceneMap      map[sceneManager.SceneId]sceneManager.Scene
	activeSceneId sceneManager.SceneId
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
	sceneMap[activeSceneId].FirstLoad()
	return &Game{
		sceneMap,
		activeSceneId,
	}
}

func (g *Game) Update() error {
	nextSceneId := g.sceneMap[g.activeSceneId].Update()
	// switched scenes
	if nextSceneId != g.activeSceneId {
		g.sceneMap[g.activeSceneId].OnExit()
		nextScene := g.sceneMap[nextSceneId]
		// if not loaded? then load in
		if !nextScene.IsLoaded() {
			nextScene.FirstLoad()
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
	}
	activeSceneId := sceneManager.BattleSceneId
	sceneMap[activeSceneId].FirstLoad()
	return &Game{
		sceneMap,
		activeSceneId,
	}
}
