package gameScenes

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/battleStatsDataManagement"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

func (g *BattleScene) FirstLoad() {
	println("Battle Scene first load executing\n")

	characters, err := battleStatsDataManagement.LoadCharacters()

	if err != nil {
		log.Fatal("error loading characters.json error:", err)
	}

	elyse := characters[0]
	enemy := characters[1]

	println("elyse stats after loading=", elyse.DisplayStats()[battleStatsDataManagement.DrawSpeed])
	println("enemy stats after loading=", enemy.DisplayStats()[battleStatsDataManagement.DrawSpeed])

	var TextInput []string

	TextInput = []string{"you filthy Animal"}
	g.statusMessage = TextInput
	g.StatusButtonEvent = false
	g.audioPlayer = audioManagement.NewAudioPlayer()
	g.TextPrinter = NewTextPrinter()
	g.TextPrinter.TextInput = "Welcome to QuickDraw Adventure!"
	g.battle = battle.NewBattle(&elyse, &enemy)
	g.graphicalEffectManager = NewGraphicalEffectManager()
	g.musicPlayer = audioManagement.NewSongPlayer(audioManagement.DialogueMusic)
	g.scene = sceneManager.BattleSceneId

	playerBS := LoadPlayerBattleSprite()
	g.playerBattleSprite = &playerBS

	enemyBS := LoadEnemyBattleSprite()
	g.enemyBattleSprite = &enemyBS

	aEffect, err := LoadAmmoEffect()
	if err != nil {
		log.Fatal("error loading AmmoEffect error:", err)
	}

	g.onScreenStatsUI = &OnScreenStatsUI{
		ammoEffect:  aEffect,
		ammoCounter: 0,
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	g.dialogueMenu.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(
			widget.Insets{Top: g.resolutionHeight / 2, Left: g.resolutionWidth / 16, Right: g.resolutionWidth - g.resolutionWidth/16 - 600/2, Bottom: 0},
		))),
	)

	g.combatMenu.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(
			widget.Insets{Top: g.resolutionHeight / 2, Left: g.resolutionWidth / 16, Right: g.resolutionWidth - g.resolutionWidth/16 - 600/2, Bottom: 0},
		))),
	)

	g.statusBar.Buttons = append(g.statusBar.Buttons, GenerateStatusBarButton(g))
	g.statusBar.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(
				widget.Insets{
					Top:    int(0.57 * float32(g.resolutionHeight)),
					Left:   int(0.5*float32(g.resolutionWidth)) - int(0.5*float32(600)),
					Right:  int(0.5*float32(g.resolutionWidth)) - int(0.5*float32(600)),
					Bottom: 200},
			),
		),
		),
	)

	dSkillsLength := len(elyse.DialogueSkills)
	dialogueSkillNames := make([]string, dSkillsLength)
	for _, skill := range elyse.DialogueSkills {
		i := skill.Index
		dialogueSkillNames[i] = skill.SkillName
	}

	cSkillsLength := len(elyse.CombatSkills)
	combatSkillNames := make([]string, cSkillsLength)
	for _, skill := range elyse.CombatSkills {
		i := skill.Index
		combatSkillNames[i] = skill.SkillName
	}

	//Creating TextInput Container (status box prints events and dialogue
	statusContainer := MakeStatusContainer()

	//to dynamically update we need to create and ebitenUI textInput widget
	statusText := StatusTextInput("b")
	statusTextLine2 := StatusTextInput("b")
	statusTextLine3 := StatusTextInput("b")
	statusContainer.AddChild(statusText)
	statusContainer.AddChild(statusTextLine2)
	statusContainer.AddChild(statusTextLine3)
	statusContainer.AddChild(g.statusBar.Buttons[0])
	g.statusBar.MenuContainer.AddChild(statusContainer)

	rootContainer.AddChild(g.statusBar.MenuContainer)

	//creating menu for skill buttons that can be used by the playerBattleSprite
	dialogueSkillsContainer := SkillsContainer()
	drawButton := GenerateDrawButton(g)
	dialogueSkillsContainer.AddChild(drawButton)
	g.dialogueMenu.Buttons = append(g.dialogueMenu.Buttons, drawButton)

	for index, skillName := range dialogueSkillNames {
		//makes button with each skill name
		if skillName != "draw" {
			dialogueButton := GenerateSkillButtons(skillName, g)
			dialogueButton.Configure(widget.ButtonOpts.TabOrder(index))
			dialogueSkillsContainer.AddChild(dialogueButton)
			g.dialogueMenu.Buttons = append(g.dialogueMenu.Buttons, dialogueButton)
		}
	}

	dialogueContainer := SkillBoxContainer("Choose Skill")
	dialogueContainer.AddChild(dialogueSkillsContainer)
	g.dialogueMenu.MenuContainer.AddChild(dialogueContainer)
	rootContainer.AddChild(g.dialogueMenu.MenuContainer)
	//defining combat menu

	combatSkillsContainer := SkillsContainer()
	for index, skillName := range combatSkillNames {
		//makes button with each skill name
		combatButton := GenerateCombatSkillButtons(skillName, g)
		combatButton.Configure(widget.ButtonOpts.TabOrder(index))
		combatSkillsContainer.AddChild(combatButton)
		g.combatMenu.Buttons = append(g.combatMenu.Buttons, combatButton)
	}

	combatContainer := CombatSkillBoxContainer("Choose Combat Skill")
	combatContainer.AddChild(combatSkillsContainer)
	g.combatMenu.MenuContainer.AddChild(combatContainer)
	rootContainer.AddChild(g.combatMenu.MenuContainer)

	// construct the UI

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	g.ui = &ui
	g.TextPrinter.StatusText[0] = statusText
	g.TextPrinter.StatusText[1] = statusTextLine2
	g.TextPrinter.StatusText[2] = statusTextLine3

	//making input be controlled by arrowKeys through cursorHandling code
	g.Cursor = CreateCursorUpdater()
	input.SetCursorUpdater(g.Cursor)

	// Ebiten setup
	ebiten.SetWindowSize(g.resolutionWidth, g.resolutionHeight)
	ebiten.SetWindowTitle("Quick Draw Adventure")
	//Hiding mouse while we use custom cursor handling
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	g.HideSkillMenu()

	g.HideCombatMenu()

	g.loaded = true

}

func (g *BattleScene) OnEnter() {
}

func (g *BattleScene) OnExit() {
	g.characterReset()
	g.loaded = false
	g.musicPlayer.Stop()
}

// Layout implements gameScenes.
func (g *BattleScene) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update implements gameScenes.
func (g *BattleScene) Update() sceneManager.SceneId {
	g.battle.UpdateState()
	turn := g.battle.GetTurn()
	g.UpdateOutputDuringNonTurn()
	g.musicPlayer.Update()
	g.TextPrinter.countDownUpdate()
	// update the UI
	g.ui.Update()
	if g.turnTracker > 0 {
		g.playerBattleSprite.Update()
		g.enemyBattleSprite.Update()
		g.audioPlayer.Update()
	}

	if g.turnTracker < g.battle.Turn {
		g.turnTracker++
		log.Printf("Turn: %d proccessing turnTracker data for effect", g.battle.Turn)
		log.Printf("playerBattleSprite effect queue lengths = %d", len(g.graphicalEffectManager.PlayerEffects.EffectQueue))
		g.updateTurnLog()
	}

	g.graphicalEffectManager.EnemyEffects.Update()

	g.graphicalEffectManager.PlayerEffects.Update()

	if g.Cursor.countdown > 0 {
		g.Cursor.countdown--
	}

	if g.eventCountDown > 0 {
		println("event countdown:", g.eventCountDown)
		g.eventCountDown--
	}
	//have the menu handle its own events

	if g.TextPrinter.CounterOn {
		g.TextPrinter.UpdateCounter()
	}
	//
	g.TextPrinter.UpdateTPState()

	if len(g.TextPrinter.TextInput) > 0 && g.TextPrinter.Counter%2 == 0 && g.TextPrinter.NextMessage {

		g.TextPrinter.CounterOn = true
		g.TextPrinter.MessageLoop()
		if g.TextPrinter.Counter%4 == 0 {
			g.audioPlayer.Play(audioManagement.TextOutput)

		}
	}

	if g.TextPrinter.state == NotPrinting {
		g.statusBar.EnableButtonVisibility()
	}
	if g.TextPrinter.state == Printing {
		g.statusBar.DisableButtonVisibility()
	}

	if g.inMenu && g.battle.GetPhase() == battle.Dialogue {
		g.dialogueMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
		//g.statusBar.DisableButtonVisibility()
	}

	if g.inMenu && g.battle.GetPhase() == battle.Shooting {
		g.combatMenu.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
		g.statusBar.DisableButtonVisibility()
		g.HideSkillMenu()
	}

	if g.eventCountDown == 1 && g.currentEvent != NoEvent {
		g.TriggerEvent(g.currentEvent)
		println("event received:", g.currentEvent, "0= move to dialogue menu, 3 = move to skill menu")
		g.changeEvent(NoEvent, 0)
	}

	g.PlayerTurn(turn)
	g.enemyTurn(turn)

	err := g.onScreenStatsUI.Update(*turn)
	if err != nil {
		log.Fatal(err)
	}
	g.UpdateSceneChangeCountdown()

	return g.scene
}

// Draw implements Ebiten Draw method.
func (g *BattleScene) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{R: 205, G: 176, B: 109, A: 255})

	if g.turnTracker > 0 {
		PrintStatus(g, screen)
	}

	g.DrawCharOutline(screen, *g.playerBattleSprite)
	DrawBattleSprite(*g.playerBattleSprite, screen, g.playerBattleSprite.Scale)
	DrawBattleSprite(*g.enemyBattleSprite, screen, g.playerBattleSprite.Scale)
	g.onScreenStatsUI.Draw(screen)
	g.graphicalEffectManager.PlayerEffects.Draw(screen)
	g.graphicalEffectManager.EnemyEffects.Draw(screen)
	g.ui.Draw(screen)
	debugTextPrint(screen, g)

}

func PrintStatus(g *BattleScene, screen *ebiten.Image) {
	face, err := LoadFont(40, November)
	if err != nil {
		log.Fatal(err)
	}

	dp := g.battle.WinningProb

	playerAmmo := fmt.Sprintf("Player Ammo:%d", g.battle.PlayerAmmo)
	enemyAmmo := fmt.Sprintf("Enemy Ammo :%d", g.battle.EnemyAmmo)

	winningProbText := fmt.Sprintf("Probability of Winning Draw:%d", dp)
	playerHealth := fmt.Sprintf("Player Health:%d", g.battle.Player.DisplayStat(battleStatsDataManagement.Health))
	enemyHealth := fmt.Sprintf("Enemy Health:%d", g.battle.Enemy.DisplayStat(battleStatsDataManagement.Health))
	tensionMeter := fmt.Sprintf("Tension:%d", g.battle.Tension)

	dopts := text.DrawOptions{}
	dopts.DrawImageOptions.ColorScale.Scale(1, 0, 0, 255)
	dopts.GeoM.Translate(400, 50)
	text.Draw(screen, winningProbText, face, &dopts)
	dopts.GeoM.Reset()

	dopts.DrawImageOptions.ColorScale.Scale(1, 0, 0, 255)
	dopts.GeoM.Translate(400, 100)
	text.Draw(screen, tensionMeter, face, &dopts)

	face, err = LoadFont(16, November)
	if err != nil {
		log.Fatal(err)
	}
	dopts.GeoM.Translate(250, 200)
	text.Draw(screen, playerHealth, face, &dopts)
	dopts.GeoM.Reset()

	dopts.GeoM.Translate(250, 220)
	text.Draw(screen, playerAmmo, face, &dopts)
	dopts.GeoM.Reset()

	dopts.GeoM.Translate(800, 200)
	text.Draw(screen, enemyHealth, face, &dopts)
	dopts.GeoM.Reset()

	dopts.GeoM.Translate(800, 220)
	text.Draw(screen, enemyAmmo, face, &dopts)
	dopts.GeoM.Reset()

}

func debugTextPrint(screen *ebiten.Image, g *BattleScene) {
	face, err := LoadFont(40, November)
	if err != nil {
		log.Fatal(err)
	}

	dopts := text.DrawOptions{}

	var battleState string

	if g.battle.State == battle.PlayerTurn {
		battleState = "battle state = Player Turn"
	}
	if g.battle.State == battle.EnemyTurn {
		battleState = "battle state = Enemy Turn"
	}
	if g.battle.State == battle.NextTurn {
		battleState = "battle state = Next Turn"
	}
	if g.battle.State == battle.NotStarted {
		battleState = "battle state = Not Started"
	}

	dopts.GeoM.Translate(600, 500)
	text.Draw(screen, battleState, face, &dopts)
	dopts.GeoM.Reset()
}
