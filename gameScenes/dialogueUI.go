package gameScenes

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/dialogueData"
	"github.com/acoco10/QuickDrawAdventure/sceneManager"
	"github.com/acoco10/QuickDrawAdventure/ui"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
	"strings"
	"unicode"
)

type TextPopup interface {
	TriggerButton()
}
type DialogUiState uint8

const (
	PrintingPlayerDialogue DialogUiState = iota
	PrintingNpcDialogue
	Completed
)

const (
	screenWidth  = int(1512)
	screenHeight = int(918)
)

type DialogueType uint8

const (
	ShowDown DialogueType = iota
	Dialogue
)

type DialogueUI struct {
	ui                        *ebitenui.UI
	TextPrinter               *TextPrinter
	statusBar                 *ui.Menu
	ShowDownStatusBar         *ui.Menu
	face                      text.Face
	triggered                 bool
	nextScene                 bool
	ScreenWidth, ScreenHeight int
	triggerScene              sceneManager.SceneId
	ButtonEvent               bool
	PlayerDialogueTracker     dialogueData.DialogueTracker
	NpcDialogueTracker        dialogueData.DialogueTracker
	index                     int
	StoryPoint                int
	PlayerImg                 *ebiten.Image
	CharImg                   *ebiten.Image
	State                     DialogUiState
	DType                     DialogueType
	loaded                    bool
	nameTag                   *widget.TextInput
	playerData                string
	npcData                   string
}

func GenerateMenuButton(popup TextPopup) (button *widget.Button) {

	buttonImage := LoadStatusButtonImage()

	statusButton := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			popup.TriggerButton()
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
	face, err := assetManagement.LoadFont(14, assetManagement.NovemberOutline)
	if err != nil {
		log.Fatal(err)
	}

	d := &DialogueUI{}
	d.face = face
	d.StoryPoint = 1
	d.ScreenHeight = 918
	d.ScreenWidth = 1512

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
	d.statusBar.Buttons = append(d.statusBar.Buttons, GenerateMenuButton(d))
	d.statusBar.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(
				widget.Insets{
					Top:    d.ScreenHeight / 4,
					Left:   d.ScreenWidth / 2,
					Right:  d.ScreenWidth / 2,
					Bottom: d.ScreenHeight - d.ScreenHeight/4},
			),
		),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(610, 200)),
	)

	d.ShowDownStatusBar = &ui.Menu{}
	d.ShowDownStatusBar.Buttons = append(d.statusBar.Buttons, GenerateMenuButton(d))
	d.ShowDownStatusBar.MenuContainer = widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(
				widget.Insets{
					Top:    int(0.75 * float32(d.ScreenHeight)),
					Left:   int(0.5*float32(d.ScreenWidth)) - int(0.5*float32(600)),
					Right:  int(0.5*float32(d.ScreenWidth)) - int(0.5*float32(600)),
					Bottom: int(0.35 * float32(d.ScreenHeight)),
				},
			),
		),
		),
	)

	//initialize empty lines for multi line text output
	name := StatusTextInput("playerName")
	statusText := StatusTextInput("white")
	statusTextLine2 := StatusTextInput("white")
	statusTextLine3 := StatusTextInput("white")

	//adding to container
	statusContainer.AddChild(name)
	statusContainer.AddChild(statusText)
	statusContainer.AddChild(statusTextLine2)
	statusContainer.AddChild(d.statusBar.Buttons[0])
	d.statusBar.MenuContainer.AddChild(statusContainer)
	d.ShowDownStatusBar.MenuContainer.AddChild(statusContainer)

	rootContainer.AddChild(d.statusBar.MenuContainer)
	rootContainer.AddChild(d.ShowDownStatusBar.MenuContainer)

	gUi := ebitenui.UI{
		Container: rootContainer,
	}
	d.nextScene = false
	d.ui = &gUi
	d.TextPrinter.StatusText[0] = statusText
	d.TextPrinter.StatusText[1] = statusTextLine2
	d.TextPrinter.StatusText[2] = statusTextLine3
	d.nameTag = name
	d.PlayerImg, _, err = ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/characters/battleSprites/elyse/elyseDialogue.png")
	if err != nil {
		log.Fatal(err)
	}
	d.LoadDialogueJson()
	return d, nil
}
func (d *DialogueUI) TriggerButton() {
	d.ButtonEvent = true
}

func FormatJsonName(name string) string {
	name = strings.ToUpper(string(name[0])) + name[1:]
	for i := 1; i < len(name); i++ {

		if unicode.IsUpper(rune(name[i])) == true {
			name = name[:i] + " " + name[i:]
			break
		}
	}
	return name
}

func (d *DialogueUI) Load(charName string, dType DialogueType) {
	d.index = 1
	playerDialogueTracker := dialogueData.DialogueTracker{
		CharName: charName,
		Index:    0,
	}

	npcDialogueTracker := dialogueData.DialogueTracker{
		CharName: charName,
		Index:    0,
	}

	d.DType = dType
	d.triggered = true
	d.PlayerDialogueTracker = playerDialogueTracker
	d.NpcDialogueTracker = npcDialogueTracker

	if d.DType == ShowDown {
		d.statusBar.MenuContainer.GetWidget().Visibility = widget.Visibility_Hide
		d.ShowDownStatusBar.MenuContainer.GetWidget().Visibility = widget.Visibility_Show

		enemyImgPath := "images/characters/battleSprites/" + charName + "/" + charName + "Dialogue.png"

		charImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, enemyImgPath)
		if err != nil {
			log.Fatal(err)
		}

		d.CharImg = charImg
	}
	if d.DType == Dialogue {
		d.statusBar.MenuContainer.GetWidget().Visibility = widget.Visibility_Show
		d.ShowDownStatusBar.MenuContainer.GetWidget().Visibility = widget.Visibility_Hide
	}

	playerFirst := dialogueData.TalkFirst(charName, d.StoryPoint, d.playerData)

	if playerFirst {
		d.nameTag.SetText("Elyse")
		d.TextPrinter.TextInput = dialogueData.GetResponse(d.PlayerDialogueTracker.CharName, d.index, d.playerData)
		d.State = PrintingPlayerDialogue
	} else {
		d.nameTag.SetText(FormatJsonName(d.NpcDialogueTracker.CharName))
		d.TextPrinter.TextInput = dialogueData.GetResponse(d.NpcDialogueTracker.CharName, d.index, d.npcData)
		d.State = PrintingNpcDialogue
	}

	d.TextPrinter.NextMessage = true
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
	playerResponse := dialogueData.GetResponse(d.PlayerDialogueTracker.CharName, d.index+1, d.playerData)
	npcResponse := dialogueData.GetResponse(d.NpcDialogueTracker.CharName, d.index+1, d.npcData)
	if playerResponse == "" && npcResponse == "" {
		println("completed dialogue")
		d.State = Completed
		return
	}
	if d.State == PrintingNpcDialogue {
		if playerResponse != "" {
			d.State = PrintingPlayerDialogue
			return
		}
	}
	if d.State == PrintingPlayerDialogue {
		if npcResponse != "" {
			d.State = PrintingNpcDialogue
		}
	}
}

func (d *DialogueUI) Reset() {
	d.loaded = false
	d.triggered = false
}

func (d *DialogueUI) Update() sceneManager.SceneId {
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
				d.nameTag.SetText(FormatJsonName(d.NpcDialogueTracker.CharName))
				fmt.Printf("entering npc: %s\n", d.NpcDialogueTracker.CharName)
				d.index++
				npcResponse := dialogueData.GetResponse(d.NpcDialogueTracker.CharName, d.index, d.npcData)
				if npcResponse != "" {
					d.TextPrinter.TextInput = npcResponse
					d.TextPrinter.NextMessage = true
				} else {
					println("no response from npc dialogue request")
					d.UpdateState()
				}

			}
			if d.State == PrintingPlayerDialogue {
				d.nameTag.SetText("Elyse")
				println("entering player dialogue")
				d.index++
				playerResponse := dialogueData.GetResponse(d.PlayerDialogueTracker.CharName, d.index, d.playerData)
				if playerResponse != "" {
					println(playerResponse)
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
				return d.triggerScene
			}

		}
	}
	return sceneManager.TownSceneID
}

func (d *DialogueUI) Draw(screen *ebiten.Image) error {
	if d.triggered {
		if d.DType == ShowDown {
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Scale(5, 5)
			if d.State == PrintingPlayerDialogue {
				screen.DrawImage(d.PlayerImg, opts)
			}
			if d.State == PrintingNpcDialogue {
				screen.DrawImage(d.CharImg, opts)
			}

		}

		d.ui.Draw(screen)
	}
	return nil
}

func (d *DialogueUI) LoadDialogueJson() {
	data, err := assets.Dialogue.ReadFile("dialogueData/elyseDialogue.json")
	if err != nil {
		log.Fatal(err)
	}

	jsonString := string(data)
	d.playerData = jsonString

	data, err = assets.Dialogue.ReadFile("dialogueData/townDialogue.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonString = string(data)
	d.npcData = jsonString

}
