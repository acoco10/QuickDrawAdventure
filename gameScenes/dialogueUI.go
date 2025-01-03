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
			widget.WidgetOpts.MinSize(300, 300),
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
		if len(d.TextPrinter.TextInput) == d.TextPrinter.MessageIndex {

			println("resetting printer and moving cursor to Menu")

			d.TextPrinter.stringPosition = 1
			d.TextPrinter.MessageIndex = 0
			d.TextPrinter.StatusText[0].SetText("")
			d.TextPrinter.StatusText[1].SetText("")
			d.TextPrinter.StatusText[2].SetText("")
			d.TextPrinter.TextInput = []string{}
			d.TextPrinter.lines = []string{}
			d.TextPrinter.lineCounter = 0
			d.statusBar.DisableButtonVisibility()

		} else {
			println("triggering printer again, message index = ", d.TextPrinter.MessageIndex, "\n")
			//clear the last output
			d.TextPrinter.stringPosition = 1

			d.TextPrinter.StatusText[0].SetText("")
			d.TextPrinter.StatusText[1].SetText("")
			d.TextPrinter.StatusText[2].SetText("")

			//if there are more lines of the message trigger the printer again

			d.TextPrinter.NextMessage = true
		}
	}
}

func MakeDialogueUI(resolutionHeight int, resolutionWidth int) (*DialogueUI, error) {
	face, err := LoadFont(14)
	if err != nil {
		log.Fatal(err)
	}

	textInput := []string{"Owner of this Bar knows everyone in town"}
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
	d.TextPrinter = NewTextPrinter(textInput)

	//container for output menu
	statusContainer := MakeStatusContainer()
	d.statusBar.Buttons = append(d.statusBar.Buttons, GenerateDialogueBarButton(d))

	d.statusBar.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(
				widget.Insets{
					Top:    int(0.5 * float32(resolutionHeight)),
					Left:   int(0.5*float32(resolutionWidth)) - int(0.5*float32(600)),
					Right:  int(0.5*float32(resolutionWidth)) - int(0.5*float32(600)),
					Bottom: 400},
			),
		),
		),
	)

	//initialize empty lines for multi line text output
	statusText := StatusTextInput()
	statusTextLine2 := StatusTextInput()
	statusTextLine3 := StatusTextInput()

	//adding to container
	statusContainer.AddChild(statusText)
	statusContainer.AddChild(statusTextLine2)
	statusContainer.AddChild(statusTextLine3)
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
		d.TextPrinter.DialogueMessageLoop()
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

func (d *DialogueUI) UpdateDialogueUIText(text []string) {
	d.TextPrinter.TextInput = text
}
