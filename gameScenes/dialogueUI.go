package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type DialogueUI struct {
	ui                *ebitenui.UI
	TextPrinter       *TextPrinter
	statusBar         *ui.Menu
	face              text.Face
	triggered         bool
	nextScene         bool
	triggerScene      sceneManager.SceneId
	playerPortrait    ebiten.Image
	characterPortrait ebiten.Image
}

func (d *DialogueUI) Trigger() {
	d.triggered = true
}

func (d *DialogueUI) TriggerScene() sceneManager.SceneId {
	if d.nextScene {
		return d.triggerScene
	}
	return sceneManager.TownSceneID
}

func (d *DialogueUI) UpdateTriggerScene(sceneId sceneManager.SceneId) {
	d.triggerScene = sceneId
}

func GenerateDialogueBarButton(d *DialogueUI) (button *widget.Button) {

	buttonImage := LoadStatusButtonImage()

	statusButton := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			DialogueStatusEffectButtonEvent(d)
		}), widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(100, 100),
			//widget.WidgetOpts.CursorHovered("statusBar"),
			//widget.WidgetOpts.CursorPressed("statusBar"),
		),
		widget.ButtonOpts.TabOrder(1),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position:  widget.RowLayoutPositionEnd,
			MaxWidth:  36, //36
			MaxHeight: 24, //24
		},
		),
		),
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			args.Button.GetWidget().Disabled = false
		},
		),
	)

	return statusButton
}

func DialogueStatusEffectButtonEvent(d *DialogueUI) {
	if len(d.TextPrinter.TextInput) == 0 {
		d.nextScene = true
	}
	if d.TextPrinter.NextMessage == false {

		println("resetting printer and moving cursor to Menu")

		d.TextPrinter.ResetTP()
		d.statusBar.DisableButtonVisibility()
	}
}

func MakeDialogueUI(resolutionHeight int, resolutionWidth int) (*DialogueUI, error) {
	face, err := LoadFont(14, November)
	if err != nil {
		log.Fatal(err)
	}

	textInput := []string{"Owner of this Bar knows everyone in town", "You should try the lemon daquiri"}
	d := &DialogueUI{}
	d.face = face

	//npc dialogue

	//clickable Player text options

	//initialDialogue button option

	//root container for ui
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	//empty menu to initialize dialogue output menu
	d.statusBar = &ui.Menu{}
	d.TextPrinter = NewTextPrinter()
	d.TextPrinter.TextInput = textInput[0]
	//container for output menu
	statusContainer := MinorDialogueContainer()
	d.statusBar.Buttons = append(d.statusBar.Buttons, GenerateDialogueBarButton(d))

	d.statusBar.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(
				widget.Insets{
					Top:    int(193 * 4),
					Left:   int(184 * 4),
					Right:  918 - (184 * 4),
					Bottom: 1512 - (184 * 4)},
			),
		),
		),
	)

	//initialize empty lines for multi line text output
	statusText := StatusTextInput("white")
	statusTextLine2 := StatusTextInput("white")
	statusTextLine3 := StatusTextInput("white")

	//adding to container
	statusContainer.AddChild(statusText)
	statusContainer.AddChild(statusTextLine2)
	statusContainer.AddChild(d.statusBar.Buttons[0])
	d.statusBar.MenuContainer.AddChild(statusContainer)

	rootContainer.AddChild(d.statusBar.MenuContainer)

	gUi := ebitenui.UI{
		Container: rootContainer,
	}
	d.nextScene = false
	d.ui = &gUi
	d.TextPrinter.StatusText[0] = statusText
	d.TextPrinter.StatusText[1] = statusTextLine2
	d.TextPrinter.StatusText[2] = statusTextLine3
	return d, nil
}

func (d *DialogueUI) UpdateDialogueUI() error {
	d.ui.Update()
	if len(d.TextPrinter.TextInput) > 0 && d.TextPrinter.Counter%2 == 0 && d.TextPrinter.NextMessage {
		d.TextPrinter.CounterOn = true
	}

	if d.TextPrinter.CounterOn {
		d.TextPrinter.UpdateCounter()
	}

	return nil
}

func (d *DialogueUI) Draw(screen *ebiten.Image) error {
	if d.triggered {
		d.ui.Draw(screen)
	}
	return nil
}

func (d *DialogueUI) UpdateDialogueUIText(text string) {
	d.TextPrinter.TextInput = text
}
