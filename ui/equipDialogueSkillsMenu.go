package ui

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tidwall/gjson"
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
	Player             battleStats.CharacterData
}

func (d *DialogueSkillEquipMenu) Load(resolutionWidth int, resolutionHeight int, char battleStats.CharacterData) {

	dialogueMenu := Menu{}
	equipMenu := Menu{}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	equipMenu.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(
			widget.Insets{Top: resolutionHeight / 4, Left: 100, Right: 0, Bottom: resolutionHeight / 4 * 3},
		))),
	)

	dialogueMenu.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(
			widget.Insets{Top: resolutionHeight / 4 * 3, Left: 100, Right: 0, Bottom: 300},
		))),
	)

	skillOptions := char.LearnedInsults
	dialogueSkillsContainer := SkillsContainer()

	for _, skill := range skillOptions {
		query := fmt.Sprintf("insults.#(id = %d).insult", skill)
		results := gjson.Get(char.DialogueData, query)
		response := results.String()
		//makes button with each skill name
		dialogueButton := GenerateSkillButton(response, d)
		dialogueSkillsContainer.AddChild(dialogueButton)
		dialogueMenu.Buttons = append(dialogueMenu.Buttons, dialogueButton)
	}

	dialogueContainer := SkillBoxContainer("Choose Insults to Equip")
	dialogueContainer.AddChild(dialogueSkillsContainer)
	dialogueMenu.MenuContainer.AddChild(dialogueContainer)

	bragOptions := char.LearnedBrags
	for _, skill := range bragOptions {
		query := fmt.Sprintf("brags.#(id = %d).name", skill)
		results := gjson.Get(char.DialogueData, query)
		response := results.String()
		//makes button with each skill name
		dialogueButton := GenerateSkillButton(response, d)
		dialogueSkillsContainer.AddChild(dialogueButton)
		dialogueMenu.Buttons = append(dialogueMenu.Buttons, dialogueButton)
	}

	rootContainer.AddChild(dialogueMenu.MenuContainer)

	equippedSkillOptionsContainer := SkillsContainer()

	for slot := range char.DialogueSlots {
		text := fmt.Sprintf("Slot %d:", slot)
		dialogueButton := GenerateSlotButton(text, d)
		equippedSkillOptionsContainer.AddChild(dialogueButton)
		equipMenu.Buttons = append(equipMenu.Buttons, dialogueButton)
	}

	equippedSkillsContainer := SkillBoxContainer("Equipped Dialogue Skills")
	equippedSkillsContainer.AddChild(equippedSkillOptionsContainer)
	equipMenu.MenuContainer.AddChild(equippedSkillsContainer)
	rootContainer.AddChild(equipMenu.MenuContainer)

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
}

func (d *DialogueSkillEquipMenu) Draw(screen *ebiten.Image) {
	if d.Triggered == true {
		d.ui.Draw(screen)
	}
}

func (d *DialogueSkillEquipMenu) Trigger() {
	d.Triggered = true
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
	)

	return button
}

func GenerateSlotButton(text string, menu *DialogueSkillEquipMenu) (button *widget.Button) {

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
			menu.EquipSkill(button)
		}),
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle: color.RGBA{R: 102, G: 57, B: 48, A: 255},
		}),

		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   10,
			Right:  16,
			Top:    5,
			Bottom: 5,
		}),

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
}

func (d *DialogueSkillEquipMenu) EquipSkill(button *widget.Button) {
	if d.SkillButtonPressed {
		button.Text().Label = d.ButtonText
		d.SkillButtonPressed = false
	}
	d.ButtonPressed.ToggleMode = false
}
