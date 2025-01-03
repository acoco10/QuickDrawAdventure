package sceneManager

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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
	FirstLoad()
	OnEnter()
	OnExit()
	IsLoaded() bool
}
