package gameScenes

import (
	"encoding/json"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/dataManagement"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
	"os"
)

type BattleScene struct {
	ui                                *ebitenui.UI
	resolutionWidth, resolutionHeight int
	inMenu                            bool
	audioPlayer                       *audioManagement.SFXAudioPlayer
	musicPlayer                       *audioManagement.DJ
	TextPrinter                       *TextPrinter
	SkillUsed                         string
	Cursor                            *BattleMenuCursorUpdater
	battle                            *battle.Battle
	playerBattleSprite                *gameObjects.BattleSprite
	enemyBattleSprite                 *gameObjects.BattleSprite
	dialogueMenu                      *ui.Menu
	statusBar                         *ui.Menu
	combatMenu                        *ui.Menu
	FrameCounter                      int
	FrameCounterOn                    bool
	trigger                           bool
	eventCountDown                    int
	events                            map[EventName]string
	currentEvent                      EventName
	graphicalEffectManager            GraphicalEffectManager
	turnTracker                       int
	loaded                            bool
	sceneChangeCountdown              int
	scene                             sceneManager.SceneId
	onScreenStatsUI                   *OnScreenStatsUI
}

type EventName uint8

const (
	MoveCursorToSkillMenu EventName = iota
	MoveCursorToCombatMenu
	MoveCursorToStatusBar
	HideSkillMenu
	ShowSkillMenu
	HideCombatMenu
	ShowCombatMenu
	EffectAnimationTriggered
	NoEvent
)

func (g *BattleScene) changeEvent(name EventName, timer int) {
	g.currentEvent = name
	g.eventCountDown = timer
}

func (g *BattleScene) TriggerEvent(name EventName) {

	if name == MoveCursorToSkillMenu {
		g.Cursor.MoveCursorToSkillMenu()
	}
	if name == MoveCursorToCombatMenu {
		g.Cursor.MoveCursorToCombatMenu()
	}
	if name == MoveCursorToStatusBar {
		g.Cursor.MoveCursorToStatusBar()
	}
	if name == HideSkillMenu {
		g.HideSkillMenu()
	}
	if name == ShowSkillMenu {
		g.ShowSkillMenu()
	}
	if name == HideCombatMenu {
		g.HideCombatMenu()
	}
	if name == ShowCombatMenu {
		g.ShowCombatMenu()
	}

}

func NewBattleScene() *BattleScene {
	bs := &BattleScene{
		resolutionWidth:  1512,
		resolutionHeight: 982,
		dialogueMenu:     &ui.Menu{},
		combatMenu:       &ui.Menu{},
		statusBar:        &ui.Menu{},
		inMenu:           false,
		loaded:           false,
	}

	return bs
}

func (g *BattleScene) IsLoaded() bool {
	return g.loaded
}

func (g *BattleScene) HideSkillMenu() {
	g.dialogueMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Hide
	//for _, button := range g.dialogueMenu.buttons {
	//button.GetWidget().Disabled = true
	//}
}

func (g *BattleScene) DisableSkillButtons() {

	for _, button := range g.dialogueMenu.Buttons {
		button.GetWidget().Disabled = true
	}
}

func (g *BattleScene) ShowSkillMenu() {
	g.dialogueMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
	for _, button := range g.dialogueMenu.Buttons {
		button.GetWidget().Disabled = false
	}

}

func (g *BattleScene) HideCombatMenu() {
	g.combatMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Hide
}

func (g *BattleScene) ShowCombatMenu() {
	g.combatMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
}

func (g *BattleScene) KeepCursorPressed() {
	g.Cursor.keepPressed(15)
}

func LoadPlayerBattleSprite() gameObjects.BattleSprite {

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/elyseBattleSprite.png")
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpritesheet(10, 7, 32, 48)

	playerBattleSprite, err := gameObjects.NewBattleSprite(playerImg, playerSpriteSheet, 700, 350, 4)
	if err != nil {
		log.Fatal(err)

	}
	return *playerBattleSprite

}

func LoadEnemyBattleSprite() gameObjects.BattleSprite {
	enemyImg, _, err := ebitenutil.NewImageFromFile("assets/images/sheriffBattleSprite.png")
	if err != nil {
		log.Fatal(err)
	}
	enemySpriteSheet := spritesheet.NewSpritesheet(10, 4, 32, 64)
	enemyBattleSprite, err := gameObjects.NewBattleSprite(enemyImg, enemySpriteSheet, 600, 50, 4)
	if err != nil {
		log.Fatal(err)
	}
	return *enemyBattleSprite
}

func DrawBattleSprite(sprite gameObjects.BattleSprite, screen *ebiten.Image, scale float64) {
	opts := ebiten.DrawImageOptions{}
	ani := sprite.GetAnimation()
	frame := ani.Frame()
	opts.GeoM.Scale(scale, scale)
	opts.GeoM.Translate(sprite.X, sprite.Y)
	screen.DrawImage(
		sprite.Img.SubImage(
			sprite.SpriteSheet.Rect(frame),
		).(*ebiten.Image),
		&opts)
}

func (g *BattleScene) playerTurn(turn *battle.Turn) {
	if turn.PlayerSkillUsed.SkillName == "reload" {
		g.audioPlayer.Play(audioManagement.Reload)
		g.playerBattleSprite.CombatButtonAnimationTrigger("reload")
		g.playerBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
	}

	if turn.PlayerSkillUsed.SkillName == "draw" && g.battle.GetPhase() == battle.Dialogue {
		g.enemyBattleSprite.DialogueButtonAnimationTrigger("draw")
		g.enemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
		g.playerBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
		g.playerBattleSprite.DialogueButtonAnimationTrigger("draw")
		g.audioPlayer.Play(audioManagement.PistolUnHolster)
		g.battle.UpdateBattlePhase()
		g.musicPlayer.Mix(audioManagement.BattleMusic)

	}

	turn.PlayerEventTriggered = true
	g.audioPlayer.ConfigureAttackResultSoundQueue(g.battle.GetTurn().DamageToEnemy, "NPC")
	g.graphicalEffectManager.PlayerEffects.ProcessPlayerTurnData(turn)
	g.graphicalEffectManager.PlayerEffects.TriggerEffectQueue()

	log.Printf("triggering playerBattleSprite effect\n")
	log.Printf("length of playerBattleSprite effect queue = %d\n", len(g.graphicalEffectManager.PlayerEffects.EffectQueue))

	if g.battle.GetPhase() == battle.Dialogue && g.playerBattleSprite.CurrentDialogueAnimation != gameObjects.NoDialogueSkill {
		g.playerBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
	}

	if g.battle.GetPhase() == battle.Shooting {
		g.battle.DamageEnemy()
		g.battle.UpdatePlayerAmmo()
		for _, effect := range turn.PlayerSkillUsed.Effects {
			if effect.EffectType == "shot" {
				g.onScreenStatsUI.ammoEffect.state = Triggered
			}
		}
		if g.battle.Enemy.DisplayStat(dataManagement.Health) <= 0 {
			g.battle.BattleWon = true
		}
		g.playerBattleSprite.UpdateState(gameObjects.UsingCombatSkill)

	}
}

func (g *BattleScene) enemyTurn(turn *battle.Turn) {

	turn.EnemyEventTriggered = true

	if turn.EnemySkillUsed.SkillName == "reload" {
		g.audioPlayer.Play(audioManagement.Reload)
	}

	if turn.EnemySkillUsed.SkillName == "draw" && g.battle.GetPhase() == battle.Dialogue {
		g.enemyBattleSprite.DialogueButtonAnimationTrigger("draw")
		g.enemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
		g.playerBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
		g.playerBattleSprite.DialogueButtonAnimationTrigger("draw")
		g.audioPlayer.Play(audioManagement.PistolUnHolster)
		g.battle.UpdateBattlePhase()
		g.musicPlayer.Mix(audioManagement.BattleMusic)

	}
	g.graphicalEffectManager.EnemyEffects.ProcessEnemyTurnData(turn)
	g.graphicalEffectManager.EnemyEffects.TriggerEffectQueue()

	if g.battle.GetPhase() == battle.Dialogue && g.battle.Turn > 0 {
		g.enemyBattleSprite.DialogueButtonAnimationTrigger(g.battle.GetTurn().EnemySkillUsed.SkillName)
		g.enemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
	}

	if g.battle.GetPhase() == battle.Shooting {
		g.battle.UpdateEnemyAmmo()
		g.battle.DamagePlayer()
		if g.battle.Player.DisplayStat(dataManagement.Health) <= 0 {
			g.battle.BattleLost = true
		}
		g.enemyBattleSprite.CombatButtonAnimationTrigger(g.battle.GetTurn().EnemySkillUsed.SkillName)
		g.enemyBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
		g.audioPlayer.ConfigureAttackResultSoundQueue(g.battle.GetTurn().DamageToPlayer, "Player")
	}
}

func (g *BattleScene) updateTurnLog() {
	turn := g.battle.GetTurn()
	jsonData, err := json.MarshalIndent(turn, "", "\t")
	if err != nil {
		panic(err)
	}

	// Write the JSON data to a file
	err = os.WriteFile("turnLogs/turnData.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}
}

func (g *BattleScene) CheckForWinner() {
	if g.battle.BattleWon {
		g.playerBattleSprite.CombatButtonAnimationTrigger("win")
		g.playerBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
		g.enemyBattleSprite.UpdateScale(8)
		g.enemyBattleSprite.CombatButtonAnimationTrigger("win")
		g.enemyBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
		victorySounds := []resource.AudioID{audioManagement.PistolUnHolster}
		g.audioPlayer.ConfigureSoundQueue(victorySounds)
		g.sceneChangeCountdown = 100
	}
}

func (g *BattleScene) UpdateSceneChangeCountdown() {
	if g.sceneChangeCountdown > 0 {
		g.sceneChangeCountdown--
	}
	if g.sceneChangeCountdown == 1 {
		if g.battle.BattleWon {
			g.scene = sceneManager.WinSceneID
		}
		if g.battle.BattleLost {
			g.scene = sceneManager.GameOverSceneID
		}
	}
}

func (g *BattleScene) characterReset() {
	g.battle.Enemy.ResetHealth()
	g.battle.Player.ResetHealth()
	g.battle.Player.ResetStatusStats()
	g.battle.Enemy.ResetStatusStats()
}
