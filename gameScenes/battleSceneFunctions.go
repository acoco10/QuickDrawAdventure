package gameScenes

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/gameObjects"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	StatusButtonEvent                 bool
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
	statusMessage                     []string
	onScreenStatsUI                   *OnScreenStatsUI
	backGround                        ebiten.Image
	endTriggered                      bool
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
		gameObjects.AttackOne:   animations.NewCyclicAnimation(5, 25, 10, 15, 1),
		gameObjects.AttackTwo:   animations.NewCyclicAnimation(4, 34, 10, 15, 1),
		gameObjects.AttackThree: animations.NewCyclicAnimation(3, 23, 10, 15, 3),
		gameObjects.Win:         animations.NewCyclicAnimation(8, 68, 10, 15, 5),
		gameObjects.Reload:      animations.NewCyclicAnimation(9, 39, 10, 15, 1),
	}
	playerImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/characters/elyse/elyseBattleSprite.png")
	if err != nil {
		log.Fatal(err)
	}

	playerSpriteSheet := spritesheet.NewSpritesheet(10, 7, 32, 48)

	playerBattleSprite, err := gameObjects.NewBattleSprite(playerImg, playerSpriteSheet, 700, 350, 5, cAnimations)
	if err != nil {
		log.Fatal(err)

	}
	return *playerBattleSprite

}

func LoadEnemyBattleSprite(enemy battleStats.Character) gameObjects.BattleSprite {
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
	enemyPath := fmt.Sprintf("images/characters/npc/battleSprites/%sBattleSprite.png", enemy.Name)
	enemyImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, enemyPath)
	if err != nil {
		log.Fatal(err)
	}
	var enemyBs *gameObjects.BattleSprite
	if enemy.Name != "wolf" {
		enemySpriteSheet := spritesheet.NewSpritesheet(10, 4, 32, 64)
		enemyBs, err = gameObjects.NewBattleSprite(enemyImg, enemySpriteSheet, 600, 100, 3.2, humanCAnimations)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		enemySpriteSheet := spritesheet.NewSpritesheet(3, 1, 15, 19)
		enemyBs, err = gameObjects.NewBattleSprite(enemyImg, enemySpriteSheet, 600, 150, 4, animalCAnimations)
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
			g.graphicalEffectManager.PlayerEffects.ProcessPlayerTurnData(turn)
			g.graphicalEffectManager.PlayerEffects.TriggerEffectQueue()
			turn.PlayerEffectsTriggered = true
			turn.PlayerIndex++
			if g.graphicalEffectManager.PlayerEffects.EffectQueue[0].Type() == Animated {
				g.HideStatusBar()
			}
			if g.graphicalEffectManager.PlayerEffects.EffectQueue[0].Type() == Static {
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
			}
		}
		if turn.PlayerIndex == 2 {
			if len(g.graphicalEffectManager.PlayerEffects.EffectQueue) > 0 {
				if g.StatusButtonEvent && g.graphicalEffectManager.PlayerEffects.EffectQueue[0].Type() == Static {
					g.graphicalEffectManager.PlayerEffects.EffectQueue[0].UnTrigger()
					g.StatusButtonEvent = false
				}
			}
			if g.graphicalEffectManager.PlayerEffects.state == NotTriggered {
				if turn.PlayerSkillUsed.SkillName == "draw" {
					g.dialogueMenu.DisableButtons()
					g.battle.UpdateBattlePhase()
					g.playerBattleSprite.DialogueButtonAnimationTrigger("draw")
					g.enemyBattleSprite.DialogueButtonAnimationTrigger("draw")
					g.enemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
					g.playerBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
					g.onScreenStatsUI.ammoEffect.MakeVisible()
					g.musicPlayer.Mix(audioManagement.BattleMusic)
					turn.EnemyTurnCompleted = true
				}
				g.ShowStatusBar()
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
				turn.PlayerIndex++
				g.battle.EnactEffects(turn.PlayerSkillUsed, g.battle.Player, g.battle.Enemy, turn.PlayerRoll, turn.PlayerSecondaryRoll)
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
			g.graphicalEffectManager.PlayerEffects.ProcessPlayerTurnData(turn)
			g.graphicalEffectManager.PlayerEffects.TriggerEffectQueue()
			g.audioPlayer.ConfigureAttackResultSoundQueue(turn.DamageToEnemy, g.battle.Enemy.Name)
			turn.PlayerEffectsTriggered = true
			turn.PlayerIndex++
			g.onScreenStatsUI.ProcessTurn(turn.DamageToEnemy, turn.PlayerSkillUsed.SkillName)
			g.battle.DamageEnemy()
			g.battle.UpdatePlayerAmmo()
		}
		if turn.PlayerIndex == 2 && g.graphicalEffectManager.PlayerEffects.state == NotTriggered {
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
			g.graphicalEffectManager.EnemyEffects.ProcessEnemyTurnData(turn)
			g.graphicalEffectManager.EnemyEffects.TriggerEffectQueue()
			if turn.EnemySkillUsed.SkillName == "draw" {
				g.battle.UpdateBattlePhase()
				g.playerBattleSprite.DialogueButtonAnimationTrigger("draw")
				g.enemyBattleSprite.DialogueButtonAnimationTrigger("draw")
				g.enemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
				g.playerBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
				turn.PlayerTurnCompleted = true
			}

			turn.EnemyEffectsTriggered = true
			turn.EnemyIndex++

			if len(g.graphicalEffectManager.EnemyEffects.EffectQueue) > 0 {
				if g.graphicalEffectManager.EnemyEffects.EffectQueue[0].Type() == Animated {
					g.HideStatusBar()
				}
				if g.graphicalEffectManager.EnemyEffects.EffectQueue[0].Type() == Static {
					g.incrementTextPrinter()
					g.TextPrinter.NextMessage = true
				}
			}
		}
		if turn.EnemyIndex == 2 {
			g.battle.UpdateWinProbability(battle.DrawProb(g.battle.Player.DisplayStats(), g.battle.Enemy.DisplayStats()))
			if len(g.graphicalEffectManager.EnemyEffects.EffectQueue) > 0 {
				if g.StatusButtonEvent && g.graphicalEffectManager.EnemyEffects.EffectQueue[0].Type() == Static {
					g.graphicalEffectManager.EnemyEffects.EffectQueue[0].UnTrigger()
					g.StatusButtonEvent = false
				}
			}
			if g.graphicalEffectManager.EnemyEffects.state == NotTriggered {
				g.ShowStatusBar()
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
				turn.EnemyIndex++
				g.battle.EnactEffects(turn.EnemySkillUsed, g.battle.Enemy, g.battle.Player, turn.EnemyRoll, turn.EnemySecondaryRoll)
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

		/*if turn.EnemySkillUsed.SkillName == "reload" {
			g.audioEnemy.Play(audioManagement.Reload)
			g.EnemyBattleSprite.CombatButtonAnimationTrigger("reload")
			g.EnemyBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
		}

		if turn.EnemySkillUsed.SkillName == "draw" && g.battle.GetPhase() == battle.Dialogue {
			turn.EnemyEventTriggered = true
			g.enemyBattleSprite.DialogueButtonAnimationTrigger("draw")
			g.enemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
			g.EnemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
			g.EnemyBattleSprite.DialogueButtonAnimationTrigger("draw")
			g.audioEnemy.Play(audioManagement.PistolUnHolster)
			g.battle.UpdateBattlePhase()
			g.musicEnemy.Mix(audioManagement.BattleMusic)

		}
		if turn.EnemySkillUsed.SkillName == "stare down" {
			g.audioEnemy.Play(audioManagement.StareDownEffect)
			g.EnemyBattleSprite.UpdateCharEffect(gameObjects.Outline, 150)
		}*/
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
			g.graphicalEffectManager.EnemyEffects.ProcessEnemyTurnData(turn)
			g.graphicalEffectManager.EnemyEffects.TriggerEffectQueue()
			g.audioPlayer.ConfigureAttackResultSoundQueue(g.battle.GetTurn().DamageToPlayer, "Player")
			turn.EnemyEffectsTriggered = true
			turn.EnemyIndex++
			g.battle.UpdateEnemyAmmo()
			g.battle.DamagePlayer()
		}
		if turn.EnemyIndex == 2 {
			if g.graphicalEffectManager.EnemyEffects.state == NotTriggered {
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
		g.sceneChangeCountdown = 10
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
		//g.sceneChangeCountdown = 100
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
