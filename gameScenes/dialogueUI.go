package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/dialogueData"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

type DialogUiState uint8

const (
	PrintingPlayerDialogue DialogUiState = iota
	PrintingNpcDialogue
	Completed
)

type DialogueUI struct {
	ui                    *ebitenui.UI
	TextPrinter           *TextPrinter
	statusBar             *ui.Menu
	face                  text.Face
	triggered             bool
	nextScene             bool
	triggerScene          sceneManager.SceneId
	ButtonEvent           bool
	PlayerDialogueTracker dialogueData.DialogueTracker
	NpcDialogueTracker    dialogueData.DialogueTracker
	StoryPoint            int
	State                 DialogUiState
	loaded                bool
}

func GenerateDialogueBarButton(d *DialogueUI) (button *widget.Button) {

	buttonImage := LoadStatusButtonImage()

	statusButton := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			d.ButtonEvent = true
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

func MakeDialogueUI(resolutionHeight int, resolutionWidth int) (*DialogueUI, error) {
	face, err := LoadFont(14, NovemberOutline)
	if err != nil {
		log.Fatal(err)
	}

	d := &DialogueUI{}
	d.face = face
	d.StoryPoint = 1

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
	d.TextPrinter.TextInput = ""
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
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(610, 200)),
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

func (d *DialogueUI) LoadDialogueUI(charName string) {
	playerDialogueTracker := dialogueData.DialogueTracker{
		CharName: charName,
		Index:    0,
	}

	npcDialogueTracker := dialogueData.DialogueTracker{
		CharName: charName,
		Index:    1,
	}
	d.triggered = true
	d.PlayerDialogueTracker = playerDialogueTracker
	d.NpcDialogueTracker = npcDialogueTracker
	d.TextPrinter.TextInput = dialogueData.GetResponse(d.NpcDialogueTracker.CharName, d.NpcDialogueTracker.Index)
	d.TextPrinter.NextMessage = true
	d.State = PrintingNpcDialogue
	d.loaded = true
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
func (d *DialogueUI) UpdateState() {
	playerResponse := dialogueData.GetPlayerResponse(d.PlayerDialogueTracker.CharName, d.StoryPoint, d.PlayerDialogueTracker.Index+1)
	npcResponse := dialogueData.GetResponse(d.NpcDialogueTracker.CharName, d.NpcDialogueTracker.Index+1)
	if playerResponse == "" && npcResponse == "" {
		println("completed dialogue")
		d.State = Completed
	}
	switch d.State {
	case PrintingNpcDialogue:
		d.State = PrintingPlayerDialogue
	case PrintingPlayerDialogue:
		d.State = PrintingNpcDialogue
	default:
		d.State = Completed
	}

}

func (d *DialogueUI) Reset() {
	d.loaded = false
	d.triggered = false
}

func (d *DialogueUI) Update() {
	if d.triggered {
		d.ui.Update()
		textInput := d.TextPrinter.TextInput

		if len(textInput) > 0 && d.TextPrinter.Counter%2 == 0 && d.TextPrinter.NextMessage {
			d.TextPrinter.CounterOn = true
			d.TextPrinter.MessageLoop()
		}

		if !d.TextPrinter.NextMessage && d.ButtonEvent {
			d.ButtonEvent = false
			d.UpdateState()
			d.TextPrinter.ResetTP()

			if d.State == PrintingNpcDialogue {
				println("entering npc dialogue")
				d.NpcDialogueTracker.Index++
				npcResponse := dialogueData.GetResponse(d.NpcDialogueTracker.CharName, d.NpcDialogueTracker.Index)
				if npcResponse != "" {
					d.TextPrinter.TextInput = npcResponse
					d.TextPrinter.NextMessage = true
				} else {
					println("no response from npc dialogue request")
					d.UpdateState()
				}

			}
			if d.State == PrintingPlayerDialogue {
				println("entering player dialogue")
				d.PlayerDialogueTracker.Index++
				playerResponse := dialogueData.GetPlayerResponse(d.PlayerDialogueTracker.CharName, d.StoryPoint, d.PlayerDialogueTracker.Index)
				if playerResponse != "" {
					d.TextPrinter.TextInput = playerResponse
					d.TextPrinter.NextMessage = true
				} else {
					println("no response from player dialogue request")
					d.UpdateState()
				}
			}

			if d.State == Completed {
				d.Reset()
				d.ButtonEvent = false
			}

		}
	}
}

func (d *DialogueUI) Draw(screen *ebiten.Image) error {
	if d.triggered {
		d.ui.Draw(screen)
	}
	return nil
}
