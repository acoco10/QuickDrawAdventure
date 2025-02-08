package sceneManager

import (
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameLog struct {
	PlayerLocation   string
	PlayerStats      *battleStats.CharacterData
	EnemyEncountered battleStats.CharacterName
	PreviousScene    SceneId
}

type SceneId uint

const (
	BattleSceneId SceneId = iota
	StartSceneId
	GameOverSceneID
	WinSceneID
	TownSceneID
)

type Scene interface {
	Update() SceneId
	Draw(screen *ebiten.Image)
	FirstLoad(log *GameLog)
	OnEnter()
	OnExit()
	IsLoaded() bool
}
