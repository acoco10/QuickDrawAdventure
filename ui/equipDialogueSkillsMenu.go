package ui

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
	"strings"
)

type DialogueSkillEquipMenu struct {
	ui                 ebitenui.UI
	Triggered          bool
	ButtonText         string
	SkillButtonPressed bool
	ButtonPressed      *widget.Button
	equippedMenu       Menu
	selectMenu         Menu
	hoverMenu          map[string]*widget.Container
	Player             *battleStats.CharacterData
	ResolutionHeight   int
	ResolutionWidth    int
	face               text.Face
	cursor             *CursorUpdater
	timer              int
	triggerFunc        func()
}

func NewDialogueEquipMenu(updater *CursorUpdater, resolutionHeight, resolutionWidth int) *DialogueSkillEquipMenu {
	face, err := assetManagement.LoadFont(16, assetManagement.November)
	if err != nil {
		log.Fatal(err)
	}

	d := DialogueSkillEquipMenu{}

	d.face = face
	d.ResolutionHeight = resolutionHeight
	d.ResolutionWidth = resolutionWidth
	d.cursor = updater

	return &d
}

func (d *DialogueSkillEquipMenu) Load(resolutionHeight int, resolutionWidth int, char *battleStats.CharacterData) {

	dialogueMenu := Menu{}
	equipMenu := Menu{}
	d.hoverMenu = make(map[string]*widget.Container, len(char.LearnedDialogueSkills))

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	equipMenu.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(
			widget.Insets{Top: resolutionHeight / 4, Left: 100, Right: 900, Bottom: resolutionHeight / 4 * 3},
		))),
	)

	dialogueMenu.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(
			widget.Insets{Top: resolutionHeight / 4 * 2, Left: 100, Right: 900, Bottom: 300},
		))),
	)

	dialogueSkillsContainer := SkillsContainer()

	Options := char.LearnedDialogueSkills
	for _, skill := range Options {
		//makes button with each skill name
		if skill.SkillName != "draw" {
			dialogueButton := GenerateSkillButton(skill.SkillName, d)
			dialogueSkillsContainer.AddChild(dialogueButton)
			dialogueMenu.Buttons = append(dialogueMenu.Buttons, dialogueButton)
			d.hoverMenu[skill.SkillName] = d.MakeHoverMenuForSkill(skill, d.face)
		}
	}

	dialogueContainer := SkillBoxContainerEquipUi("Dialogue Options")
	dialogueContainer.AddChild(dialogueSkillsContainer)
	dialogueMenu.MenuContainer.AddChild(dialogueContainer)
	rootContainer.AddChild(dialogueMenu.MenuContainer)

	equippedSkillOptionsContainer := SkillsContainer()
	for slot := range char.DialogueSlots {
		buttonText := fmt.Sprintf("Slot %d:", slot)
		dialogueButton := GenerateSlotButton(buttonText, d, 20)
		equippedSkillOptionsContainer.AddChild(dialogueButton)
		equipMenu.Buttons = append(equipMenu.Buttons, dialogueButton)
	}

	equippedSkillsContainer := SkillBoxContainerEquipUi("Equipped Dialogue Skills")
	equippedSkillsContainer.AddChild(equippedSkillOptionsContainer)
	equipMenu.MenuContainer.AddChild(equippedSkillsContainer)
	rootContainer.AddChild(equipMenu.MenuContainer)

	for _, skill := range d.hoverMenu {
		rootContainer.AddChild(skill)
		skill.GetWidget().Visibility = widget.Visibility_Hide
	}

	ui := ebitenui.UI{
		Container: rootContainer,
	}
	d.ui = ui
	d.Triggered = false

	d.Player = char

	d.equippedMenu = equipMenu
	d.selectMenu = dialogueMenu
	d.SkillButtonPressed = false

}

func (d *DialogueSkillEquipMenu) Update() {
	if d.Triggered == true {
		d.ui.Update()
	}
	if d.timer == 1 {
		if d.SkillButtonPressed {
			d.cursor.SetSkillMenuEquip()
		} else {
			d.cursor.SetSkillMenuSelect()
		}
	}
	if d.timer > 0 {
		d.timer--
	}
}

func (d *DialogueSkillEquipMenu) MakeHoverMenuForSkill(skill battleStats.Skill, face text.Face) *widget.Container {
	hoverMenu := Menu{}
	MenuContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(
			widget.Insets{Top: d.ResolutionHeight / 4 * 2, Left: 500, Right: 100, Bottom: 300},
		))),
	)

	hoverMenu.MenuContainer = SkillBoxContainerEquipUi("Skill Attributes")
	attributeContainer := SkillsContainer()

	nameTag := fmt.Sprintf("Name: %s", skill.SkillName)
	textTag := fmt.Sprintf("Dialogue: %s", skill.Text)
	var effects []string
	for i, effect := range skill.Effects {
		effectTag := fmt.Sprintf("Effect %d: %s, Amount: %d ", i, effect.Stat, effect.Amount)
		effects = append(effects, effectTag)
	}

	allAttributes := []string{nameTag, textTag, effects[0]}
	if len(effects) > 1 {
		allAttributes = append(allAttributes, effects[1])
	}

	for _, attribute := range allAttributes {
		button := GenerateSlotButton(attribute, d, 16)
		hoverMenu.Buttons = append(hoverMenu.Buttons, button)
		attributeContainer.AddChild(button)
	}

	hoverMenu.MenuContainer.AddChild(attributeContainer)
	MenuContainer.AddChild(hoverMenu.MenuContainer)
	return MenuContainer
}

func (d *DialogueSkillEquipMenu) Draw(screen *ebiten.Image) {
	if d.Triggered == true {
		d.ui.Draw(screen)
	}
}

func (d *DialogueSkillEquipMenu) Trigger() {
	d.Triggered = true
	d.cursor.SetSkillMenuSelect()
}

func (d *DialogueSkillEquipMenu) UnTrigger() {
	d.Triggered = false
}

func GenerateSkillButton(text string, menu *DialogueSkillEquipMenu) (button *widget.Button) {

	// load gameScenes font, more fonts will be selectable later when we implement a resource manager
	face, err := assetManagement.LoadFont(20, assetManagement.November)
	buttonText := strings.ToUpper(string(text[0])) + text[1:]
	if err != nil {
		log.Fatal(err)
	}

	// loads a basic button image
	buttonImage := LoadButtonImage()

	//make a new button with the name of each skill as text
	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			menu.SkillSelect(text, button)
			button.ToggleMode = true
		}),

		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle: color.RGBA{R: 102, G: 57, B: 48, A: 255},
		}),

		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   64,
			Right:  16,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 10),
			widget.WidgetOpts.CursorHovered("buttonHover"),
			widget.WidgetOpts.CursorPressed("buttonPressed"),
		),
		widget.ButtonOpts.TabOrder(5),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.CursorEnterHandler(func(args *widget.WidgetCursorEnterEventArgs) {
				menu.hoverMenu[text].GetWidget().Visibility = widget.Visibility_Show
			}),
			widget.WidgetOpts.CursorExitHandler(func(args *widget.WidgetCursorExitEventArgs) {
				menu.hoverMenu[text].GetWidget().Visibility = widget.Visibility_Hide
			}),
		),
	)

	return button
}

func GenerateSlotButton(text string, menu *DialogueSkillEquipMenu, fontSize float64) (button *widget.Button) {

	// load gameScenes font, more fonts will be selectable later when we implement a resource manager
	face, err := assetManagement.LoadFont(fontSize, assetManagement.November)
	buttonText := strings.ToUpper(string(text[0])) + text[1:]
	if err != nil {
		log.Fatal(err)
	}

	// loads a basic button image
	buttonImage := LoadButtonImage()

	//make a new button with the name of each skill as text
	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			menu.EquipSkill(button)
		}),
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle: color.RGBA{R: 102, G: 57, B: 48, A: 255},
		}),

		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   64,
			Right:  100,
			Top:    5,
			Bottom: 5,
		}),
		widget.ButtonOpts.TextPosition(widget.TextPositionStart, widget.TextPositionStart),

		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 10),
			widget.WidgetOpts.CursorHovered("buttonHover"),
			widget.WidgetOpts.CursorPressed("buttonPressed"),
		),
		widget.ButtonOpts.TabOrder(5),
	)

	return button
}

func (d *DialogueSkillEquipMenu) SkillSelect(text string, button *widget.Button) {
	d.ButtonText = text
	d.SkillButtonPressed = true
	d.ButtonPressed = button
	d.timer = 15
	d.cursor.KeepPressed(15)
}

func (d *DialogueSkillEquipMenu) EquipSkill(button *widget.Button) {
	if d.SkillButtonPressed {
		button.Text().Label = d.ButtonText
		d.SkillButtonPressed = false
		d.timer = 15
		d.cursor.KeepPressed(15)
	}
	d.Player.EquippedDialogueSkills[d.ButtonText] = d.Player.LearnedDialogueSkills[d.ButtonText]
	d.ButtonPressed.ToggleMode = false
}
