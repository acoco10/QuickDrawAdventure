package gameScenes

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/acoco10/QuickDrawAdventure/graphicEffects"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
)

type BattleScene struct {
	ui                                *ebitenui.UI
	resolutionWidth, resolutionHeight int
	inMenu                            bool
	gameLog                           *sceneManager.GameLog
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
	StatusButtonEvent                 bool
	FrameCounter                      int
	FrameCounterOn                    bool
	trigger                           bool
	eventCountDown                    int
	events                            map[EventName]string
	currentEvent                      EventName
	graphicalEffectManager            *graphicEffects.GraphicalEffectManager
	turnTracker                       int
	loaded                            bool
	sceneChangeCountdown              int
	scene                             sceneManager.SceneId
	statusMessage                     []string
	onScreenStatsUI                   *OnScreenStatsUI
	backGround                        ebiten.Image
	endTriggered                      bool
	gameEffects                       map[graphicEffects.EffectType]graphicEffects.GraphicEffect
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

func (g *BattleScene) HideStatusBar() {
	println("hiding status bar")
	g.statusBar.MenuContainer.GetWidget().Visibility = widget.Visibility_Hide
}

func (g *BattleScene) ShowStatusBar() {
	println("showing status bar")
	g.statusBar.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
}

func (g *BattleScene) DisableSkillButtons() {
	for _, button := range g.dialogueMenu.Buttons {
		button.GetWidget().Disabled = true
	}
}

func (g *BattleScene) HideCombatMenu() {
	g.combatMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Hide
}

func (g *BattleScene) DisableStatusButton() {
	for _, button := range g.statusBar.Buttons {
		button.GetWidget().Disabled = true
	}
}

func (g *BattleScene) ShowSkillMenu() {
	g.dialogueMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
	for _, button := range g.dialogueMenu.Buttons {
		button.GetWidget().Disabled = false
	}

}

func (g *BattleScene) ShowCombatMenu() {
	g.combatMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
}

func (g *BattleScene) KeepCursorPressed() {
	g.Cursor.keepPressed(15)
}

func LoadPlayerBattleSprite() gameObjects.BattleSprite {
	cAnimations := map[gameObjects.CAnimation]*animations.CyclicAnimation{
		gameObjects.AttackOne:   animations.NewCyclicAnimation(5, 25, 10, 14, 1),
		gameObjects.AttackTwo:   animations.NewCyclicAnimation(4, 34, 10, 14, 1),
		gameObjects.AttackThree: animations.NewCyclicAnimation(3, 23, 10, 14, 3),
		gameObjects.Win:         animations.NewCyclicAnimation(8, 68, 10, 14, 5),
		gameObjects.Reload:      animations.NewCyclicAnimation(9, 39, 10, 14, 1),
	}
	playerImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/characters/battleSprites/elyse/elyseBattleSprite.png")
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpritesheet(10, 7, 32, 48)

	playerBattleSprite, err := gameObjects.NewBattleSprite(playerImg, playerSpriteSheet, 700, 350, 5, cAnimations, nil)
	if err != nil {
		log.Fatal(err)

	}
	return *playerBattleSprite

}

func LoadEnemyBattleSprite(enemy battleStats.CharacterData) gameObjects.BattleSprite {
	humanCAnimations := map[gameObjects.CAnimation]*animations.CyclicAnimation{
		gameObjects.AttackOne:   animations.NewCyclicAnimation(5, 25, 10, 15, 1),
		gameObjects.AttackTwo:   animations.NewCyclicAnimation(4, 34, 10, 15, 1),
		gameObjects.AttackThree: animations.NewCyclicAnimation(3, 23, 10, 15, 3),
		gameObjects.Win:         animations.NewCyclicAnimation(8, 68, 10, 15, 5),
		gameObjects.Reload:      animations.NewCyclicAnimation(9, 39, 10, 15, 1),
	}
	animalCAnimations := map[gameObjects.CAnimation]*animations.CyclicAnimation{
		gameObjects.AttackOne: animations.NewCyclicAnimation(0, 2, 1, 15, 1),
	}
	enemyPath := fmt.Sprintf("images/characters/battleSprites/%s/%sBattleSprite.png", enemy.Name, enemy.Name)
	enemyImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, enemyPath)
	if err != nil {
		log.Fatal(err)
	}
	var enemyBs *gameObjects.BattleSprite
	if enemy.Name != "wolf" {
		enemySpriteSheet := spritesheet.NewSpritesheet(10, 4, 32, 64)
		enemyBs, err = gameObjects.NewBattleSprite(enemyImg, enemySpriteSheet, 600, 100, 3.2, humanCAnimations, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		idle := animations.NewAnimation(0, 0, 0, 10)
		enemySpriteSheet := spritesheet.NewSpritesheet(12, 1, 15, 19)
		enemyBs, err = gameObjects.NewBattleSprite(enemyImg, enemySpriteSheet, 600, 250, 4, animalCAnimations, idle)
	}
	if err != nil {
		log.Fatal(err)
	}
	return *enemyBs
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

func (g *BattleScene) incrementTextPrinter() {
	if len(g.statusMessage) > 0 {
		g.TextPrinter.ResetTP()
		g.TextPrinter.TextInput = g.statusMessage[0]
		g.statusMessage = g.statusMessage[1:]
	}
}

func (g *BattleScene) PlayerDialogueTurn(turn *battle.Turn) {
	if g.battle.State == battle.PlayerTurn && g.battle.BattlePhase == battle.Dialogue {
		if turn.PlayerIndex == 0 {
			if turn.PlayerMessage == nil {
				turn.PlayerTurnCompleted = true
			}
			g.statusMessage = g.battle.GetTurn().PlayerMessage
			g.battle.GetTurn().PlayerEventTriggered = true
			if len(g.statusMessage) > 0 {
				g.ShowStatusBar()
				g.incrementTextPrinter()
				g.battle.GetTurn().PlayerIndex++
				g.TextPrinter.NextMessage = true
			}
		}

		if turn.PlayerIndex == 1 && g.TextPrinter.state == NotPrinting && g.StatusButtonEvent {
			g.StatusButtonEvent = false
			g.ProcessTurnEffects(*g.playerBattleSprite, turn)
			g.graphicalEffectManager.PlayerEffects.TriggerEffectQueue()
			turn.PlayerEffectsTriggered = true
			turn.PlayerIndex++
			if g.graphicalEffectManager.PlayerEffects.EffectQueue[0].Type() == graphicEffects.Animated {
				g.HideStatusBar()
			}
			if g.graphicalEffectManager.PlayerEffects.EffectQueue[0].Type() == graphicEffects.Static {
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
			}
		}
		if turn.PlayerIndex == 2 {
			if len(g.graphicalEffectManager.PlayerEffects.EffectQueue) > 0 {
				if g.StatusButtonEvent && g.graphicalEffectManager.PlayerEffects.EffectQueue[0].Type() == graphicEffects.Static {
					g.graphicalEffectManager.PlayerEffects.EffectQueue[0].UnTrigger()
					g.StatusButtonEvent = false
				}
			}
			if g.graphicalEffectManager.PlayerEffects.GetState() == graphicEffects.NotTriggered {
				if turn.PlayerSkillUsed.SkillName == "draw" {
					g.DrawSkillUsed()
					turn.EnemyTurnCompleted = true
				}
				g.ShowStatusBar()
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
				turn.PlayerIndex++
				g.battle.EnactEffects(turn.PlayerSkillUsed, g.battle.Player, g.battle.Enemy, turn.PlayerRoll, turn.PlayerSecondaryRoll)
				g.onScreenStatsUI.ProcessTurn(*g, turn.DamageToEnemy, turn.PlayerSkillUsed.SkillName)
				g.battle.UpdateWinProbability(battle.DrawProb(g.battle.Player.DisplayStats(), g.battle.Enemy.DisplayStats()))
			}

		}
		if turn.PlayerIndex > 2 && g.StatusButtonEvent {
			g.StatusButtonEvent = false
			if len(g.statusMessage) > 0 {
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
			}
			if len(g.statusMessage) <= 0 {
				turn.PlayerTurnCompleted = true
				g.TextPrinter.ResetTP()

			}
		}
	}
}

func (g *BattleScene) PlayerShootingTurn(turn *battle.Turn) {
	if g.battle.State == battle.PlayerTurn && g.battle.BattlePhase == battle.Shooting {
		if turn.PlayerIndex == 0 {
			g.ShowStatusBar()
			if turn.PlayerMessage == nil {
				turn.PlayerTurnCompleted = true
			}

			g.statusMessage = g.battle.GetTurn().PlayerMessage
			for _, msg := range g.battle.GetTurn().PlayerMessage {
				println("Player msg:", msg)
			}
			g.battle.GetTurn().PlayerEventTriggered = true
			if len(g.statusMessage) > 0 {
				println("status msg[0]", g.statusMessage[0])
				g.incrementTextPrinter()
				g.battle.GetTurn().PlayerIndex++
				g.TextPrinter.NextMessage = true
			}
		}
		if turn.PlayerIndex == 1 && g.TextPrinter.state == NotPrinting {
			g.HideStatusBar()
			g.playerBattleSprite.CombatButtonAnimationTrigger(turn.PlayerSkillUsed.SkillName)
			g.playerBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
			g.StatusButtonEvent = false
			g.ProcessTurnEffects(*g.playerBattleSprite, turn)
			g.graphicalEffectManager.GameEffects.TriggerEffectQueue()
			g.graphicalEffectManager.PlayerEffects.TriggerEffectQueue()
			g.audioPlayer.ConfigureAttackResultSoundQueue(turn.DamageToEnemy, g.battle.Enemy.Name, g.battle.Player.Name)
			if turn.PlayerSkillUsed.SkillName == "reload" {
				g.audioPlayer.Play(audioManagement.Reload)
			}
			turn.PlayerEffectsTriggered = true
			turn.PlayerIndex++
			g.onScreenStatsUI.ProcessTurn(*g, turn.DamageToEnemy, turn.PlayerSkillUsed.SkillName)
			g.battle.DamageEnemy()
			g.battle.UpdatePlayerAmmo()
		}
		if turn.PlayerIndex == 2 && g.graphicalEffectManager.PlayerEffects.GetState() == graphicEffects.NotTriggered {
			g.ShowStatusBar()
			g.incrementTextPrinter()
			g.TextPrinter.NextMessage = true
			g.StatusButtonEvent = false
			turn.PlayerIndex++
		}
		if turn.PlayerIndex > 2 && g.StatusButtonEvent {
			g.CheckForWinner()
			if len(g.statusMessage) > 0 {
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
			}
			if len(g.statusMessage) <= 0 {
				turn.PlayerTurnCompleted = true
				g.TextPrinter.ResetTP()
			}

		}
	}
}

/*if turn.PlayerSkillUsed.SkillName == "reload" {
	g.audioPlayer.Play(audioManagement.Reload)
	g.playerBattleSprite.CombatButtonAnimationTrigger("reload")
	g.playerBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
}*/

func (g *BattleScene) DrawSkillUsed() {
	g.dialogueMenu.DisableButtons()
	g.battle.UpdateBattlePhase()
	g.playerBattleSprite.DialogueButtonAnimationTrigger("draw")
	g.enemyBattleSprite.DialogueButtonAnimationTrigger("draw")
	g.enemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
	g.playerBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
	g.musicPlayer.Mix(audioManagement.BattleMusic)
	soundList := []resource.AudioID{audioManagement.PistolUnHolster, audioManagement.PistolUnHolster}
	g.audioPlayer.ConfigureSoundQueue(soundList)
}

func (g *BattleScene) EnemyDialogueTurn(turn *battle.Turn) {
	if g.battle.State == battle.EnemyTurn && g.battle.BattlePhase == battle.Dialogue {
		if turn.EnemyMessage == nil {
			turn.EnemyTurnCompleted = true
		}
		g.ShowStatusBar()
		if turn.EnemyIndex == 0 {
			g.statusMessage = g.battle.GetTurn().EnemyMessage
			g.battle.GetTurn().EnemyEventTriggered = true
			if len(g.statusMessage) > 0 {
				g.ShowStatusBar()
				g.incrementTextPrinter()
				g.battle.GetTurn().EnemyIndex++
				g.TextPrinter.NextMessage = true
			}
		}
		if turn.EnemyIndex == 1 && g.TextPrinter.state == NotPrinting && g.StatusButtonEvent {
			g.StatusButtonEvent = false
			g.ProcessTurnEffects(*g.enemyBattleSprite, turn)
			g.graphicalEffectManager.EnemyEffects.TriggerEffectQueue()
			if turn.EnemySkillUsed.SkillName == "draw" {
				g.DrawSkillUsed()
				turn.PlayerTurnCompleted = true
			}

			turn.EnemyEffectsTriggered = true
			turn.EnemyIndex++

			if len(g.graphicalEffectManager.EnemyEffects.EffectQueue) > 0 {
				if g.graphicalEffectManager.EnemyEffects.EffectQueue[0].Type() == graphicEffects.Animated {
					g.HideStatusBar()
				}
				if g.graphicalEffectManager.EnemyEffects.EffectQueue[0].Type() == graphicEffects.Static {
					g.incrementTextPrinter()
					g.TextPrinter.NextMessage = true
				}
			}
		}
		if turn.EnemyIndex == 2 {
			g.battle.UpdateWinProbability(battle.DrawProb(g.battle.Player.DisplayStats(), g.battle.Enemy.DisplayStats()))
			if len(g.graphicalEffectManager.EnemyEffects.EffectQueue) > 0 {
				if g.StatusButtonEvent && g.graphicalEffectManager.EnemyEffects.EffectQueue[0].Type() == graphicEffects.Static {
					g.graphicalEffectManager.EnemyEffects.EffectQueue[0].UnTrigger()
					g.StatusButtonEvent = false
				}
			}
			if g.graphicalEffectManager.EnemyEffects.GetState() == graphicEffects.NotTriggered {
				g.ShowStatusBar()
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
				turn.EnemyIndex++
				g.battle.EnactEffects(turn.EnemySkillUsed, g.battle.Enemy, g.battle.Player, turn.EnemyRoll, turn.EnemySecondaryRoll)
				g.onScreenStatsUI.ProcessTurn(*g, turn.DamageToEnemy, turn.PlayerSkillUsed.SkillName)
				g.battle.UpdateWinProbability(battle.DrawProb(g.battle.Player.DisplayStats(), g.battle.Enemy.DisplayStats()))
			}
		}

		if turn.EnemyIndex > 2 && g.StatusButtonEvent {
			g.StatusButtonEvent = false
			if len(g.statusMessage) > 0 {
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
			}
			if len(g.statusMessage) <= 0 {
				turn.EnemyTurnCompleted = true
				g.TextPrinter.ResetTP()
			}
		}

	}
}

func (g *BattleScene) EnemyShootingTurn(turn *battle.Turn) {
	if g.battle.State == battle.EnemyTurn && g.battle.BattlePhase == battle.Shooting {
		if turn.EnemyMessage == nil {
			turn.EnemyTurnCompleted = true
		}
		if turn.EnemyIndex == 0 {
			g.ShowStatusBar()
			g.statusMessage = g.battle.GetTurn().EnemyMessage
			g.battle.GetTurn().EnemyEventTriggered = true
			if len(g.statusMessage) > 0 {
				g.incrementTextPrinter()
				g.battle.GetTurn().EnemyIndex++
				g.TextPrinter.NextMessage = true
			}
		}
		if turn.EnemyIndex == 1 && g.TextPrinter.state == NotPrinting {
			g.HideStatusBar()
			g.enemyBattleSprite.CombatButtonAnimationTrigger(turn.EnemySkillUsed.SkillName)
			g.enemyBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
			g.StatusButtonEvent = false
			g.ProcessTurnEffects(*g.enemyBattleSprite, turn)
			g.graphicalEffectManager.EnemyEffects.TriggerEffectQueue()
			g.graphicalEffectManager.GameEffects.TriggerEffectQueue()
			g.audioPlayer.ConfigureAttackResultSoundQueue(turn.DamageToPlayer, "Player", g.battle.Enemy.Name)
			if turn.EnemySkillUsed.SkillName == "reload" {
				g.audioPlayer.Play(audioManagement.Reload)
			}
			turn.EnemyEffectsTriggered = true
			turn.EnemyIndex++
			g.battle.UpdateEnemyAmmo()
			g.battle.DamagePlayer()
		}
		if turn.EnemyIndex == 2 {
			if g.graphicalEffectManager.EnemyEffects.GetState() == graphicEffects.NotTriggered {
				g.incrementTextPrinter()
				g.ShowStatusBar()
				g.TextPrinter.NextMessage = true
				turn.EnemyIndex++
			}
		}
		if turn.EnemyIndex > 2 {
			if g.StatusButtonEvent {
				g.StatusButtonEvent = false
				g.CheckForWinner()
				if len(g.statusMessage) > 0 {
					g.incrementTextPrinter()
					g.TextPrinter.NextMessage = true
				}
				if len(g.statusMessage) <= 0 {
					turn.EnemyTurnCompleted = true
					g.TextPrinter.ResetTP()
				}
			}
		}

	}
}

func (g *BattleScene) CheckAndEndBattle() {
	if g.battle.State == battle.Over && !g.endTriggered {
		g.endTriggered = true
		if g.battle.BattleWon {
			msg := fmt.Sprintf("You win! %s has been defeated!", g.battle.Enemy.Name)
			g.statusMessage = append(g.statusMessage, msg)
		}
		if g.battle.BattleLost {
			msg := fmt.Sprintf("You lost! Elyse has been defeated!")
			g.statusMessage = append(g.statusMessage, msg)
		}
		g.incrementTextPrinter()
	}
	if g.StatusButtonEvent && g.endTriggered {
		g.StatusButtonEvent = false
		println("changing scene")
		g.sceneChangeCountdown = 100
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
	if g.battle.Enemy.DisplayStat(battleStats.Health) <= 0 {
		g.battle.BattleWon = true
	}
	if g.battle.Player.DisplayStat(battleStats.Health) <= 0 {
		g.battle.BattleLost = true
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
			println("game over scene triggered")
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

func (g *BattleScene) DrawCharOutline(screen *ebiten.Image, sprite gameObjects.BattleSprite) {
	if sprite.EffectApplied == gameObjects.Outline {
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Scale(4.0, 4.0)
		opts.GeoM.Translate(sprite.X+2, sprite.Y+5)
		screen.DrawImage(
			sprite.Img.SubImage(
				sprite.SpriteSheet.Rect(31),
			).(*ebiten.Image),
			&opts)
	}
}

func (g *BattleScene) UpdateOutputDuringNonTurn() {
	if g.battle.State == battle.NextTurn {
		if g.battle.GetPhase() == battle.Dialogue && g.inMenu == false {
			g.HideStatusBar()
			println("next turn check")
			g.changeEvent(MoveCursorToSkillMenu, 20)
			g.inMenu = true
		}
		if g.battle.GetPhase() == battle.Shooting && g.inMenu == false {
			g.HideStatusBar()
			g.changeEvent(MoveCursorToCombatMenu, 20)
			g.inMenu = true
		}
	}

	if g.battle.State == battle.NotStarted {
		if g.StatusButtonEvent {
			g.StatusButtonEvent = false
			if len(g.statusMessage) > 0 {
				g.TextPrinter.ResetTP()
				g.TextPrinter.TextInput = g.statusMessage[0]
				g.statusMessage = g.statusMessage[1:]
				g.TextPrinter.NextMessage = true
			}
			if len(g.statusMessage) <= 0 && g.TextPrinter.state == NotPrinting {
				g.HideStatusBar()
				if g.battle.GetPhase() == battle.Dialogue {
					g.inMenu = true
					g.changeEvent(MoveCursorToSkillMenu, 20)
				}

				if g.battle.GetPhase() == battle.Shooting {
					g.inMenu = true
					g.changeEvent(MoveCursorToCombatMenu, 20)
				}
				g.TextPrinter.ResetTP()
			}
		}
	}
}

func (g *BattleScene) ProcessTurnEffects(bs gameObjects.BattleSprite, turn *battle.Turn) {
	var skillUsed battleStats.Skill
	var effectSequencer *graphicEffects.GraphicalEffectSequencer
	var dmg []int
	var weakness bool
	var roll bool
	var enemyBs gameObjects.BattleSprite
	gameEffectQueue := g.graphicalEffectManager.GameEffects.EffectQueue
	if g.battle.State == battle.PlayerTurn {
		weakness = turn.EnemyWeakness
		skillUsed = turn.PlayerSkillUsed
		effectSequencer = g.graphicalEffectManager.PlayerEffects
		dmg = turn.DamageToEnemy
		roll = turn.PlayerRoll
		enemyBs = *g.enemyBattleSprite
	}
	if g.battle.State == battle.EnemyTurn {
		skillUsed = turn.EnemySkillUsed
		effectSequencer = g.graphicalEffectManager.EnemyEffects
		dmg = turn.DamageToPlayer
		roll = turn.EnemyRoll
		enemyBs = *g.playerBattleSprite
	}
	if skillUsed.SkillName != "" {
		EffectQueue := make([]graphicEffects.GraphicEffect, 0)

		if skillUsed.SkillName == "stareDown" {
			EffectQueue = append(EffectQueue, bs.Effects[graphicEffects.StareEffect])
			effectSequencer.Counter = 1
		}

		if skillUsed.SkillName == "brag" {
			println("appending skill")
			if roll {
				EffectQueue = append(EffectQueue, bs.Effects[graphicEffects.SuccessfulEffect])
			}
			if !roll {
				EffectQueue = append(EffectQueue, bs.Effects[graphicEffects.UnsuccessfulEffect])
			}
			effectSequencer.Counter = 1
		}
		if skillUsed.SkillName == "insult" {
			if roll {
				EffectQueue = append(EffectQueue, bs.Effects[graphicEffects.SuccessfulEffect])
			}

			if !roll {
				EffectQueue = append(EffectQueue, bs.Effects[graphicEffects.UnsuccessfulEffect])
			}
			effectSequencer.Counter = 1
		}
		if skillUsed.SkillName == "draw" {
			if turn.EnemySkillUsed.SkillName != "draw" || turn.TurnInitiative == battle.Player {
				EffectQueue = append(EffectQueue, g.gameEffects[graphicEffects.DrawEffect])
				effectSequencer.Counter = 1
			}
		}

		if weakness {
			EffectQueue = append(EffectQueue, g.gameEffects[graphicEffects.WeaknessEffect])
			println("adding weakness effect to player effect queue")
		}

		Effect := skillUsed.Effects[0]
		if Effect.EffectType == "buff" && roll {
			if Effect.Stat == "fear" {
				amt := Effect.Amount
				if weakness {
					amt = amt * 2
				}
				eff := g.SetEffectCoord(g.gameEffects[graphicEffects.FearEffect], graphicEffects.FearEffect)
				eff = AddTextToImageEffect(g.gameEffects[graphicEffects.FearEffect], amt, graphicEffects.FearEffect)
				EffectQueue = append(EffectQueue, eff)
				effectSequencer.Counter = 1
				println("battleSceneFunctions.go:712 length of effect Queue =", len(EffectQueue), "\n")
			}

		}

		for _, result := range dmg {
			if result >= 0 {
				g.graphicalEffectManager.GameEffects.Counter = 14
				eff := g.SetEffectCoord(g.gameEffects[graphicEffects.MuzzleEffect], graphicEffects.MuzzleEffect)
				gameEffectQueue = append(gameEffectQueue, eff)
			}
			if result > 0 {
				effectSequencer.Counter = 36
				g.graphicalEffectManager.GameEffects.Counter = 36
				eff := AddTextToImageEffect(g.gameEffects[graphicEffects.HitSplatEffect], result, graphicEffects.HitSplatEffect)
				eff = g.SetEffectCoord(eff, graphicEffects.HitSplatEffect)
				gameEffectQueue = append(gameEffectQueue, eff)
				EffectQueue = append(EffectQueue, enemyBs.Effects[graphicEffects.TookDamageEffect])
			}

			if result == 0 {
				effectSequencer.Counter = 36
				g.graphicalEffectManager.GameEffects.Counter = 36
				eff := g.SetEffectCoord(g.gameEffects[graphicEffects.MissEffect], graphicEffects.MissEffect)
				gameEffectQueue = append(gameEffectQueue, eff)
				EffectQueue = append(EffectQueue, g.gameEffects[graphicEffects.FillerEffect])
			}
		}
		effectSequencer.EffectQueue = EffectQueue
		g.graphicalEffectManager.GameEffects.EffectQueue = gameEffectQueue
		effectSequencer.Configured = true

	}
}

func (g *BattleScene) LoadGameEffects() {
	drawImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/drawEffect.png")
	if err != nil {
		log.Fatalf("ebitenutil.NewImageFromFile file not found%s\n", err)
	}

	drawEffectSpriteSheet := spritesheet.NewSpritesheet(1, 4, 199, 125)
	drawEffect := graphicEffects.NewEffect(drawImg, drawEffectSpriteSheet, 450, 200, 3, 0, 1, 12, 4)

	missImage, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/missEffect.png")
	if err != nil {
		log.Printf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}
	missEffect := graphicEffects.NewStaticEffect(missImage, 520+float64(rand.IntN(80)), 200+float64(rand.IntN(80)), 12, 1)

	damageImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/damageEffect.png")
	if err != nil {
		log.Printf("graphicalEffectManager.go:107 ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	damagedEffect := graphicEffects.NewStaticEffect(damageImg, g.playerBattleSprite.X*g.playerBattleSprite.Scale, g.playerBattleSprite.Y*g.playerBattleSprite.Scale-float64(rand.IntN(80)), 50, 1)

	fearImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/fearEffect.png")
	if err != nil {
		log.Printf("graphicalEffectManager.go:84 ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	fearEffect := graphicEffects.NewStaticEffect(fearImg, g.playerBattleSprite.X, g.playerBattleSprite.Y-float64(rand.IntN(80)), 100, 1)

	weakImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/weaknessEffect.png")
	if err != nil {
		log.Fatalf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	weaknessEffect := graphicEffects.NewStaticEffect(weakImg, 600, 250, 100, 1)

	muzzleFlashImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/muzzleFlash+smoke.png")
	if err != nil {
		log.Fatalf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	muzzleSpriteSheet := spritesheet.NewSpritesheet(7, 1, 23, 32)
	muzzleFlashEffect := graphicEffects.NewEffect(muzzleFlashImg, muzzleSpriteSheet, 250, 100, 6, 0, 1, 2, 4)

	fillerImg := ebiten.NewImage(1, 1)
	fillerEffect := graphicEffects.NewStaticEffect(fillerImg, 600, 250, 12, 1)

	g.gameEffects = map[graphicEffects.EffectType]graphicEffects.GraphicEffect{
		graphicEffects.DrawEffect:     drawEffect,
		graphicEffects.FearEffect:     fearEffect,
		graphicEffects.HitSplatEffect: damagedEffect,
		graphicEffects.WeaknessEffect: weaknessEffect,
		graphicEffects.MissEffect:     missEffect,
		graphicEffects.FillerEffect:   fillerEffect,
		graphicEffects.MuzzleEffect:   muzzleFlashEffect,
	}
}
func (g *BattleScene) SetEffectCoord(effectInput graphicEffects.GraphicEffect, etype graphicEffects.EffectType) graphicEffects.GraphicEffect {
	var effect graphicEffects.GraphicEffect
	switch e := effectInput.(type) {
	case *graphicEffects.AnimatedEffect:
		effect = effectInput
	case *graphicEffects.StaticEffect:
		effect = graphicEffects.NewStaticEffect(e.AccessImage(), 0, 0, 10, 1)
	}
	if g.battle.State == battle.EnemyTurn {
		if etype == graphicEffects.MissEffect {

			check := rand.IntN(2)
			if check == 0 {
				effect.SetCoord(g.playerBattleSprite.X-20-rand.Float64()*70, g.playerBattleSprite.Y+50+rand.Float64()*100)
			} else {
				effect.SetCoord(g.playerBattleSprite.X+100+rand.Float64()*70, g.playerBattleSprite.Y+50+rand.Float64()*100)
			}
		}
		if etype == graphicEffects.HitSplatEffect {
			minX := float64(g.playerBattleSprite.Img.Bounds().Min.X) - 10
			minY := float64(g.playerBattleSprite.Img.Bounds().Min.Y) - 10
			effect.SetCoord(g.playerBattleSprite.X+rand.Float64()*minX, g.playerBattleSprite.Y+50+rand.Float64()*minY)
		}
		if etype == graphicEffects.FearEffect {
			effectInput.SetCoord(g.playerBattleSprite.X, g.playerBattleSprite.Y+100)
		}
	}
	if g.battle.State == battle.PlayerTurn {
		if etype == graphicEffects.MissEffect {
			check := rand.IntN(2)
			if check == 0 {
				effect.SetCoord(g.enemyBattleSprite.X-20-rand.Float64()*70, g.enemyBattleSprite.Y+50+rand.Float64()*100)
			} else {
				effect.SetCoord(g.enemyBattleSprite.X+100+rand.Float64()*70, g.enemyBattleSprite.Y+50+rand.Float64()*100)
			}
		}
		if etype == graphicEffects.HitSplatEffect {
			minX := float64(g.enemyBattleSprite.Img.Bounds().Min.X) - 10
			minY := float64(g.enemyBattleSprite.Img.Bounds().Min.Y) - 10
			effect.SetCoord(g.enemyBattleSprite.X+rand.Float64()*minX, g.enemyBattleSprite.Y+50+rand.Float64()*minY)
		}
		if etype == graphicEffects.FearEffect {
			effectInput.SetCoord(g.enemyBattleSprite.X, g.enemyBattleSprite.Y+100)
		}
		if etype == graphicEffects.MuzzleEffect {
			effectInput.SetCoord(g.playerBattleSprite.X+60, g.playerBattleSprite.Y)
		}
	}
	return effect
}

func AddTextToImageEffect(inputEffect graphicEffects.GraphicEffect, input int, effectType graphicEffects.EffectType) graphicEffects.GraphicEffect {
	img := ebiten.NewImageFromImage(inputEffect.AccessImage())
	effect := graphicEffects.NewStaticEffect(img, 0, 0, 10, 1)
	face, err := assetManagement.LoadFont(18, assetManagement.Lady)
	if err != nil {
		log.Fatal(err)
	}
	dopts := text.DrawOptions{}
	inputText := strconv.FormatInt(int64(input), 10)
	if effectType == graphicEffects.HitSplatEffect {
		dopts.GeoM.Translate(16, 10)
	}
	if effectType == graphicEffects.FearEffect {
		inputText += " Fear"
		dopts.GeoM.Translate(45, 10)
	}
	text.Draw(effect.AccessImage(), inputText, face, &dopts)
	return effect
}
