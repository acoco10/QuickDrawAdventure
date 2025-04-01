package ui

import (
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

type MainMenu struct {
	ui                 ebitenui.UI
	Triggered          bool
	SkillButtonPressed bool
	ButtonPressed      *widget.Button
	Player             *battleStats.CharacterData
	ResolutionHeight   int
	ResolutionWidth    int
	face               text.Face
	Cursor             *CursorUpdater
	DialogueSkillMenu  *DialogueSkillEquipMenu
	PrimaryMenu        *Menu
	SecondaryMenus     map[string]TriggerMenu
}

func NewMainMenu(resolutionHeight int, resolutionWidth int, playerData *battleStats.CharacterData, cursor *CursorUpdater) *MainMenu {
	m := MainMenu{
		ResolutionHeight: resolutionHeight,
		ResolutionWidth:  resolutionWidth,
		Player:           playerData,
		Cursor:           cursor,
	}

	return &m
}

func (m *MainMenu) Load() {
	m.PrimaryMenu = &Menu{}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	m.PrimaryMenu.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(
			widget.Insets{Top: 400, Left: m.ResolutionWidth/2 - 200, Right: m.ResolutionWidth / 2, Bottom: m.ResolutionHeight / 4 * 3},
		))),
	)

	OptionContainer := SkillsContainer()
	options := []string{"settings", "inventory", "skills"}

	for _, op := range options {
		//makes button with each skill name
		optionButton := GenerateMainMenuButton(op, m)
		OptionContainer.AddChild(optionButton)
		m.PrimaryMenu.Buttons = append(m.PrimaryMenu.Buttons, optionButton)
	}

	opContainer := SkillBoxContainerEquipUi("Main Menu")
	opContainer.AddChild(OptionContainer)
	m.PrimaryMenu.MenuContainer.AddChild(opContainer)

	rootContainer.AddChild(m.PrimaryMenu.MenuContainer)

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	m.ui = ui

	skills := NewDialogueEquipMenu(m.Cursor, m.ResolutionHeight, m.ResolutionWidth)
	skills.Load(m.ResolutionHeight, m.ResolutionWidth, m.Player)
	m.SecondaryMenus = make(map[string]TriggerMenu)
	m.SecondaryMenus["skills"] = skills
}

func (m *MainMenu) Trigger() {
	m.Triggered = true
	m.SetCursor()
}

func (m *MainMenu) UnTrigger() {
	m.Triggered = false
}

func (m *MainMenu) SetCursor() {
	m.Cursor.MoveCursorToMainMenu()
}

func (m *MainMenu) Update() {
	m.ui.Update()

	for _, menu := range m.SecondaryMenus {
		menu.Update()
	}
}

func (m *MainMenu) Draw(screen *ebiten.Image) {
	if m.Triggered == true {
		m.ui.Draw(screen)
	}

	for _, menu := range m.SecondaryMenus {
		menu.Draw(screen)
	}
}
func GenerateMainMenuButton(text string, menu *MainMenu) (button *widget.Button) {

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
			menu.SecondaryMenus[text].Trigger()
			menu.UnTrigger()
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
			}),
			widget.WidgetOpts.CursorExitHandler(func(args *widget.WidgetCursorExitEventArgs) {
			}),
		),
	)

	return button
}
