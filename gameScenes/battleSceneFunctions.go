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
	g.onScreenStatsUI.tensionMeter.MakeNoyVisible()
}

func (g *BattleScene) ShowStatusBar() {
	g.statusBar.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
	if g.battle.BattlePhase == battle.Dialogue {
		g.onScreenStatsUI.tensionMeter.MakeVisible()
	}
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

	playerBattleSprite, err := gameObjects.NewBattleSprite(playerImg, playerSpriteSheet, 750, 350, 5, cAnimations, nil)
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
		enemyBs, err = gameObjects.NewBattleSprite(enemyImg, enemySpriteSheet, 630, 140, 3.2, humanCAnimations, nil)
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

/*if turn.PlayerSkillUsed.SkillName == "reload" {
	g.audioPlayer.Play(audioManagement.Reload)
	g.playerBattleSprite.CombatButtonAnimationTrigger("reload")
	g.playerBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
}*/

func (g *BattleScene) DrawSkillAnimationTrigger() {
	g.dialogueMenu.DisableButtons()
	g.playerBattleSprite.DialogueButtonAnimationTrigger("draw")
	g.enemyBattleSprite.DialogueButtonAnimationTrigger("draw")
	g.enemyBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
	g.playerBattleSprite.UpdateState(gameObjects.UsingDialogueSkill)
	g.musicPlayer.Mix(audioManagement.BattleMusic)
	soundList := []resource.AudioID{audioManagement.PistolUnHolster, audioManagement.PistolUnHolster}
	g.audioPlayer.ConfigureSoundQueue(soundList)

}

func (g *BattleScene) UpdateBattleAfterDraw(enemy *battle.CharacterBattleData) {
	g.battle.UpdateBattlePhase()
	enemy.Completed = true
	enemy.Message = []string{}
	g.battle.SetDrawBonus()
	g.battle.UpdateState()
	g.battle.UpdateBattlePhase()

}

func (g *BattleScene) EnemyDialogueTurn(turn *battle.Turn) {
	enemy := g.battle.CharacterBattleData[battle.Enemy]
	if g.battle.State == battle.EnemyTurn && g.battle.BattlePhase == battle.Dialogue {
		if enemy.Message == nil {
			turn.EnemyTurnCompleted = true
		}
		g.ShowStatusBar()
		if turn.EnemyIndex == 0 {
			g.statusMessage = enemy.Message
			enemy.EventTriggered = true
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
			if enemy.SkillUsed.SkillName == "draw" {
				g.DrawSkillAnimationTrigger()
			}

			enemy.EffectsTriggered = true
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
			g.battle.UpdateWinProbability(battle.DrawProb(g.battle.CharacterBattleData[battle.Player].DisplayStats(), enemy.DisplayStats()))
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
				g.onScreenStatsUI.ProcessTurn(*g, enemy.DamageOutput, g.battle.CharacterBattleData[battle.Player].SkillUsed.SkillName)
				g.battle.UpdateWinProbability(battle.DrawProb(g.battle.CharacterBattleData[battle.Player].DisplayStats(), enemy.DisplayStats()))
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
				if enemy.SkillUsed.SkillName == "draw" {
					g.battle.SetDrawBonus()
					g.battle.UpdateState()
					g.battle.UpdateBattlePhase()
				}

			}
		}

	}
}

func (g *BattleScene) UpdateTurn() {
	player := g.battle.CharacterBattleData[battle.Player]
	enemy := g.battle.CharacterBattleData[battle.Enemy]
	if g.battle.State == battle.PlayerTurn && g.battle.BattlePhase == battle.Dialogue {
		g.DialogueTurn(g.battle.CharacterBattleData[battle.Player], g.battle.CharacterBattleData[battle.Enemy])
	}
	if g.battle.State == battle.PlayerTurn && g.battle.BattlePhase == battle.Shooting {
		g.ShootingTurn(g.battle.CharacterBattleData[battle.Player], g.battle.CharacterBattleData[battle.Enemy])
	}
	if g.battle.State == battle.EnemyTurn && g.battle.BattlePhase == battle.Dialogue {
		if !enemy.EventTriggered {
			enemySkillUsed, err := battle.EnemyChooseSkill(*g.battle, enemy.DialogueSkills)
			if err != nil {
				log.Fatal(err)
			}
			enemy.SkillUsed = enemySkillUsed
			g.battle.UpdateChar(enemy, player)
		}
		println("Triggering enemy turn")
		g.DialogueTurn(g.battle.CharacterBattleData[battle.Enemy], g.battle.CharacterBattleData[battle.Player])

	}
	if g.battle.State == battle.EnemyTurn && g.battle.BattlePhase == battle.Shooting {
		g.ShootingTurn(g.battle.CharacterBattleData[battle.Enemy], g.battle.CharacterBattleData[battle.Player])
	}
}

func (g *BattleScene) DialogueTurn(char *battle.CharacterBattleData, enemy *battle.CharacterBattleData) {
	if char.Message == nil {
		char.Completed = true
	}
	g.ShowStatusBar()
	if char.Index == 0 {
		g.statusMessage = char.Message
		char.EventTriggered = true
		if len(g.statusMessage) > 0 {
			g.ShowStatusBar()
			g.incrementTextPrinter()
			char.Index++
			g.TextPrinter.NextMessage = true
		}
	}
	if char.Index == 1 && g.TextPrinter.state == NotPrinting && g.StatusButtonEvent {
		if len(g.statusMessage) == 0 {
			char.Completed = true
			g.TextPrinter.ResetTP()

		}
		g.StatusButtonEvent = false
		//g.ProcessTurnEffects(*g.enemyBattleSprite, g.battle.GetTurn())
		//g.graphicalEffectManager.EnemyEffects.TriggerEffectQueue()
		if char.SkillUsed.SkillName == "draw" {
			g.DrawSkillAnimationTrigger()
		}

		char.EffectsTriggered = true
		char.Index++

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
	if char.Index == 2 {
		if len(g.statusMessage) == 0 {
			char.Completed = true
			g.TextPrinter.ResetTP()

		}
		g.battle.UpdateWinProbability(battle.DrawProb(g.battle.CharacterBattleData[battle.Player].DisplayStats(), enemy.DisplayStats()))
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
			char.Index++
			g.battle.EnactEffects(char.SkillUsed, char, enemy, char.Roll, char.SecondaryRoll)
			g.onScreenStatsUI.ProcessTurn(*g, enemy.DamageOutput, g.battle.CharacterBattleData[battle.Player].SkillUsed.SkillName)
			g.battle.UpdateWinProbability(battle.DrawProb(g.battle.CharacterBattleData[battle.Player].DisplayStats(), enemy.DisplayStats()))
		}
	}

	if char.Index > 2 && g.StatusButtonEvent {
		g.StatusButtonEvent = false
		if len(g.statusMessage) > 0 {
			g.incrementTextPrinter()
			g.TextPrinter.NextMessage = true
		}
		if len(g.statusMessage) == 0 {
			if char.SkillUsed.SkillName == "draw" {
				g.UpdateBattleAfterDraw(enemy)
			}
			char.Completed = true
			g.TextPrinter.ResetTP()

		}
	}
}

func (g *BattleScene) ShootingTurn(char *battle.CharacterBattleData, enemy *battle.CharacterBattleData) {
	if char.Message == nil {
		char.Completed = true
	}
	if char.Index == 0 {
		g.ShowStatusBar()
		g.statusMessage = char.Message
		char.EventTriggered = true
		if len(g.statusMessage) > 0 {
			g.incrementTextPrinter()
			char.Index++
			g.TextPrinter.NextMessage = true
		}
	}
	if char.Index == 1 && g.TextPrinter.state == NotPrinting {
		g.HideStatusBar()
		//g.enemyBattleSprite.CombatButtonAnimationTrigger(char.SkillUsed.SkillName)
		//g.enemyBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
		g.StatusButtonEvent = false
		//g.ProcessTurnEffects(*g.enemyBattleSprite, g.battle.GetTurn())
		//g.graphicalEffectManager.EnemyEffects.TriggerEffectQueue()
		//g.graphicalEffectManager.GameEffects.TriggerEffectQueue()
		//g.audioPlayer.ConfigureAttackResultSoundQueue(char.DamageOutput, "Player", g.battle.CharacterBattleData[battle.Enemy].Name)
		if char.SkillUsed.SkillName == "reload" {
			g.audioPlayer.Play(audioManagement.Reload)
		}
		char.EffectsTriggered = true
		char.Index++
		char.UpdateAmmo()
		if char.Name == "elyse" {
			g.battle.DamageCharacter(battle.Enemy, battle.Player)
			g.onScreenStatsUI.ProcessTurn(*g, char.DamageOutput, char.SkillUsed.SkillName)
		} else {
			g.battle.DamageCharacter(battle.Player, battle.Enemy)
		}

	}
	if char.Index == 2 {
		//if g.graphicalEffectManager.EnemyEffects.GetState() == graphicEffects.NotTriggered
		g.incrementTextPrinter()
		g.ShowStatusBar()
		g.TextPrinter.NextMessage = true
		char.Index++
	}
	if char.Index > 2 {
		if g.StatusButtonEvent {
			g.StatusButtonEvent = false
			g.CheckForWinner()
			if len(g.statusMessage) > 0 {
				g.incrementTextPrinter()
				g.TextPrinter.NextMessage = true
			}
			if len(g.statusMessage) <= 0 {
				if char.DrawBonus {
					char.DrawBonus = false
					char.Completed = true
					enemy.Completed = true
				}
				char.Completed = true
				g.TextPrinter.ResetTP()
			}
		}
	}
}

func (g *BattleScene) EnemyShootingTurn(turn *battle.Turn) {
	enemy := g.battle.CharacterBattleData[battle.Enemy]
	if g.battle.State == battle.EnemyTurn && g.battle.BattlePhase == battle.Shooting {
		if enemy.Message == nil {
			turn.EnemyTurnCompleted = true
		}
		if turn.EnemyIndex == 0 {
			g.ShowStatusBar()
			g.statusMessage = enemy.Message
			enemy.EventTriggered = true
			if len(g.statusMessage) > 0 {
				g.incrementTextPrinter()
				g.battle.GetTurn().EnemyIndex++
				g.TextPrinter.NextMessage = true
			}
		}
		if turn.EnemyIndex == 1 && g.TextPrinter.state == NotPrinting {
			g.HideStatusBar()
			g.enemyBattleSprite.CombatButtonAnimationTrigger(enemy.SkillUsed.SkillName)
			g.enemyBattleSprite.UpdateState(gameObjects.UsingCombatSkill)
			g.StatusButtonEvent = false
			g.ProcessTurnEffects(*g.enemyBattleSprite, turn)
			g.graphicalEffectManager.EnemyEffects.TriggerEffectQueue()
			g.graphicalEffectManager.GameEffects.TriggerEffectQueue()
			g.audioPlayer.ConfigureAttackResultSoundQueue(enemy.DamageOutput, "Player", g.battle.CharacterBattleData[battle.Enemy].Name)
			if enemy.SkillUsed.SkillName == "reload" {
				g.audioPlayer.Play(audioManagement.Reload)
			}
			enemy.EffectsTriggered = true
			turn.EnemyIndex++
			g.battle.UpdateAmmo()
			g.battle.DamageCharacter(battle.Player, battle.Enemy)
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
					if enemy.DrawBonus {
						enemy.DrawBonus = false
					}
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
			msg := fmt.Sprintf("You win! %s has been defeated!", g.battle.CharacterBattleData[battle.Enemy].Name)
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

	// Write the JSON playerData to a file
	err = os.WriteFile("turnLogs/turnData.json", jsonData, 0644)
	if err != nil {
		panic(err)
	}
}

func (g *BattleScene) CheckForWinner() {
	if g.battle.CharacterBattleData[battle.Enemy].DisplayStat(battleStats.Health) <= 0 {
		g.battle.BattleWon = true
	}
	if g.battle.CharacterBattleData[battle.Player].DisplayStat(battleStats.Health) <= 0 {
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
	player := g.battle.CharacterBattleData[battle.Player]
	enemy := g.battle.CharacterBattleData[battle.Enemy]
	enemy.ResetHealth()
	enemy.ResetStatusStats()
	player.ResetHealth()
	player.ResetStatusStats()
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
			if g.battle.CharacterBattleData[battle.Enemy].DrawBonus {
				g.battle.TakeCombatTurn(g.battle.CharacterBattleData[battle.Player].CombatSkills["reload"])
			} else {
				g.changeEvent(MoveCursorToCombatMenu, 20)
				g.inMenu = true
			}
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
	var char *battle.CharacterBattleData
	var enemy *battle.CharacterBattleData
	if g.battle.State == battle.PlayerTurn {
		char = g.battle.CharacterBattleData[battle.Player]
		enemy = g.battle.CharacterBattleData[battle.Enemy]
		enemyBs = *g.enemyBattleSprite
	} else {
		char = g.battle.CharacterBattleData[battle.Enemy]
		enemy = g.battle.CharacterBattleData[battle.Player]
		enemyBs = *g.playerBattleSprite
	}
	weakness = char.WeaknessTargeted
	skillUsed = char.SkillUsed
	effectSequencer = g.graphicalEffectManager.PlayerEffects
	dmg = char.DamageOutput
	roll = char.Roll

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
			if enemy.SkillUsed.SkillName != "draw" || turn.TurnInitiative == battle.Player {
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
			amt := Effect.Amount
			if weakness {
				amt = amt * 2
			}
			if Effect.Stat == "fear" {
				eff := AddTextToImageEffect(g.gameEffects[graphicEffects.FearEffect], amt, graphicEffects.FearEffect)
				eff = g.SetEffectCoord(eff, graphicEffects.FearEffect)
				EffectQueue = append(EffectQueue, eff)
				effectSequencer.Counter = 1
				println("battleSceneFunctions.go:712 length of effect Queue =", len(EffectQueue), "\n")
			}
			if Effect.Stat == "anger" {
				eff := AddTextToImageEffect(g.gameEffects[graphicEffects.AngerEffect], amt, graphicEffects.AngerEffect)
				eff = g.SetEffectCoord(eff, graphicEffects.AngerEffect)
				EffectQueue = append(EffectQueue, eff)
				effectSequencer.Counter = 1
			}

		}
		if char.DrawBonus || enemy.DrawBonus && g.battle.BattlePhase == battle.Shooting {
			gameEffectQueue = append(gameEffectQueue, g.gameEffects[graphicEffects.DrawWinEffect])
		}

		for _, result := range dmg {
			if result >= 0 {
				if !char.DrawBonus || !enemy.DrawBonus {
					g.graphicalEffectManager.GameEffects.Counter = 14
				}
				eff := g.SetEffectCoord(g.gameEffects[graphicEffects.MuzzleEffect], graphicEffects.MuzzleEffect)
				gameEffectQueue = append(gameEffectQueue, eff)
				g.graphicalEffectManager.GameEffects.Counter = 14
			}
			if result > 0 {
				effectSequencer.Counter = 18
				eff := AddTextToImageEffect(g.gameEffects[graphicEffects.HitSplatEffect], result, graphicEffects.HitSplatEffect)
				eff = g.SetEffectCoord(eff, graphicEffects.HitSplatEffect)
				gameEffectQueue = append(gameEffectQueue, eff)
				EffectQueue = append(EffectQueue, enemyBs.Effects[graphicEffects.TookDamageEffect])
				g.gameEffects[graphicEffects.FillerEffect].SetDuration(44)
				EffectQueue = append(EffectQueue, g.gameEffects[graphicEffects.FillerEffect])
				//EffectQueue = append(EffectQueue, g.gameEffects[graphicEffects.FillerEffect])
			}

			if result == 0 {
				effectSequencer.Counter = 28
				eff := g.SetEffectCoord(g.gameEffects[graphicEffects.MissEffect], graphicEffects.MissEffect)
				gameEffectQueue = append(gameEffectQueue, eff)
				g.gameEffects[graphicEffects.FillerEffect].SetDuration(50)
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
	missEffect := graphicEffects.NewStaticEffect(missImage, 520+float64(rand.IntN(80)), 200+float64(rand.IntN(80)), 20, 1, graphicEffects.MissEffect)

	damageImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/damageEffect.png")
	if err != nil {
		log.Printf("graphicalEffectManager.go:107 ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	hitSplatEffect := graphicEffects.NewStaticEffect(damageImg, g.playerBattleSprite.X*g.playerBattleSprite.Scale, g.playerBattleSprite.Y*g.playerBattleSprite.Scale-float64(rand.IntN(80)), 50, 1, graphicEffects.HitSplatEffect)
	hitSplatEffect.SetDepth(1)
	fearImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/fearEffect.png")
	if err != nil {
		log.Printf("graphicalEffectManager.go:84 ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	fearEffect := graphicEffects.NewStaticEffect(fearImg, g.playerBattleSprite.X, g.playerBattleSprite.Y-float64(rand.IntN(80)), 100, 1, graphicEffects.FearEffect)

	weakImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/weaknessEffect.png")
	if err != nil {
		log.Fatalf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	weaknessEffect := graphicEffects.NewStaticEffect(weakImg, 600, 250, 100, 1, graphicEffects.WeaknessEffect)

	angerImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/angerEffect.png")
	if err != nil {
		log.Fatalf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	angerEffect := graphicEffects.NewStaticEffect(angerImg, 600, 250, 100, 1, graphicEffects.AngerEffect)

	muzzleFlashImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/muzzleFlash+smoke.png")
	if err != nil {
		log.Fatalf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}

	muzzleSpriteSheet := spritesheet.NewSpritesheet(7, 1, 29, 32)
	muzzleFlashEffect := graphicEffects.NewEffect(muzzleFlashImg, muzzleSpriteSheet, 250, 100, 6, 0, 1, 2, 4)

	drawWinEffectimg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/effectAssets/drawWinEffect.png")
	if err != nil {
		log.Fatalf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
	}
	drawWinEffect := graphicEffects.NewStaticEffect(drawWinEffectimg, 300, 170, 50, 2, graphicEffects.DrawWinEffect)

	fillerImg := ebiten.NewImage(1, 1)
	fillerEffect := graphicEffects.NewStaticEffect(fillerImg, 600, 250, 20, 1, graphicEffects.FillerEffect)

	g.gameEffects = map[graphicEffects.EffectType]graphicEffects.GraphicEffect{
		graphicEffects.DrawEffect:     drawEffect,
		graphicEffects.FearEffect:     fearEffect,
		graphicEffects.HitSplatEffect: hitSplatEffect,
		graphicEffects.WeaknessEffect: weaknessEffect,
		graphicEffects.MissEffect:     missEffect,
		graphicEffects.FillerEffect:   fillerEffect,
		graphicEffects.MuzzleEffect:   muzzleFlashEffect,
		graphicEffects.AngerEffect:    angerEffect,
		graphicEffects.DrawWinEffect:  drawWinEffect,
	}
}
func (g *BattleScene) SetEffectCoord(effectInput graphicEffects.GraphicEffect, etype graphicEffects.EffectType) graphicEffects.GraphicEffect {
	var effect graphicEffects.GraphicEffect
	switch e := effectInput.(type) {
	case *graphicEffects.AnimatedEffect:
		effect = effectInput
	case *graphicEffects.StaticEffect:
		effect = graphicEffects.NewStaticEffect(e.AccessImage(), 0, 0, 50, 1, etype)
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
			minX := float64(g.playerBattleSprite.X) + 50
			minY := float64(g.playerBattleSprite.Y) + 50
			randFactorX := 10 * rand.Float64()
			randFactorY := 100 * rand.Float64()
			effect.SetCoord(minX+randFactorX, minY+randFactorY)
			effect.SetDepth(1)
		}
		if etype == graphicEffects.FearEffect || etype == graphicEffects.AngerEffect {
			effect.SetCoord(g.playerBattleSprite.X+30, g.playerBattleSprite.Y+100)
		}
		if etype == graphicEffects.MuzzleEffect {
			effect.SetCoord(g.enemyBattleSprite.X, g.enemyBattleSprite.Y)
			effect.SetDepth(1)
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
			minX := float64(g.enemyBattleSprite.X) + 20
			minY := float64(g.enemyBattleSprite.Y) + 50
			randFactorX := 10 * rand.Float64()
			randFactorY := 100 * rand.Float64()
			effect.SetCoord(minX+randFactorX, minY+randFactorY)
			effect.SetDepth(1)
		}
		if etype == graphicEffects.FearEffect || etype == graphicEffects.AngerEffect {
			effect.SetCoord(g.enemyBattleSprite.X, g.enemyBattleSprite.Y+100)
		}
		if etype == graphicEffects.MuzzleEffect {
			effect.SetCoord(g.playerBattleSprite.X+60, g.playerBattleSprite.Y)
			effect.SetDepth(0)
		}
	}
	return effect
}

func AddTextToImageEffect(inputEffect graphicEffects.GraphicEffect, input int, effectType graphicEffects.EffectType) graphicEffects.GraphicEffect {
	img := ebiten.NewImageFromImage(inputEffect.AccessImage())
	effect := graphicEffects.NewStaticEffect(img, 0, 0, 10, 1, effectType)
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
	if effectType == graphicEffects.AngerEffect {
		inputText += " Anger"
		dopts.GeoM.Translate(45, 10)
	}

	text.Draw(effect.AccessImage(), inputText, face, &dopts)
	return effect
}
