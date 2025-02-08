package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type TextState uint8

const (
	PrintingText TextState = iota
	TextCompleted
)

type TextUI struct {
	ui          *ebitenui.UI
	TextPrinter *TextPrinter
	statusBar   *ui.Menu
	face        text.Face
	triggered   bool
	ButtonEvent bool
	State       TextState
	loaded      bool
	inputText   []string
	index       int
}

func MakeTextUI(resolutionHeight int, resolutionWidth int) (*TextUI, error) {
	face, err := assetManagement.LoadFont(14, assetManagement.NovemberOutline)
	if err != nil {
		log.Fatal(err)
	}

	t := &TextUI{}
	t.face = face

	//root container for ui
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
	)

	//empty menu to initialize dialogue output menu
	t.statusBar = &ui.Menu{}
	t.TextPrinter = NewTextPrinter()
	t.TextPrinter.TextInput = ""
	//container for output menu
	statusContainer := MinorDialogueContainer()
	t.statusBar.Buttons = append(t.statusBar.Buttons, GenerateMenuButton(t))

	t.statusBar.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(
				widget.Insets{
					Top:    screenHeight / 4,
					Left:   screenWidth / 2,
					Right:  screenWidth / 2,
					Bottom: screenHeight - screenHeight/4},
			),
		),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(610, 200)),
	)

	//initialize empty lines for multi line text output
	statusText := StatusTextInput("white")
	statusTextLine2 := StatusTextInput("white")
	statusTextLine3 := StatusTextInput("white")

	//adding to container
	statusContainer.AddChild(statusText)
	statusContainer.AddChild(statusTextLine2)
	statusContainer.AddChild(t.statusBar.Buttons[0])
	t.statusBar.MenuContainer.AddChild(statusContainer)

	rootContainer.AddChild(t.statusBar.MenuContainer)

	gUi := ebitenui.UI{
		Container: rootContainer,
	}
	t.ui = &gUi
	t.TextPrinter.StatusText[0] = statusText
	t.TextPrinter.StatusText[1] = statusTextLine2
	t.TextPrinter.StatusText[2] = statusTextLine3
	return t, nil
}

func (t *TextUI) LoadTextUI(text []string) {
	t.triggered = true
	t.TextPrinter.TextInput = text[0]
	t.index = 0
	t.TextPrinter.NextMessage = true
	t.State = PrintingText
	t.loaded = true
}

func (t *TextUI) TriggerButton() {
	t.triggered = true
}

func (t *TextUI) UpdateState() {

}

func (t *TextUI) Reset() {
	t.loaded = false
	t.triggered = false
}

func (t *TextUI) Update() {
	if t.triggered {
		t.ui.Update()

	}
}

func (t *TextUI) Draw(screen *ebiten.Image) error {
	if t.triggered {
		t.ui.Draw(screen)
	}
	return nil
}
