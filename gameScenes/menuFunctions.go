package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/audioManagement"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/ui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"image/color"
	"log"
	"strings"
)

func GenerateSkillButton(text string, g *BattleScene) (button *widget.Button) {

	// load gameScenes font, more fonts will be selectable later when we implement a resource manager
	face, err := assetManagement.LoadFont(20, assetManagement.November)
	buttonText := strings.ToUpper(string(text[0])) + text[1:]
	if err != nil {
		log.Fatal(err)
	}

	// loads a basic button image
	buttonImage := ui.LoadButtonImage()

	//make a new button with the name of each skill as text
	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			DialogueSkillButtonEvent(g, text, button)

		}),
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 102, G: 57, B: 48, A: 255},
			Disabled: color.NRGBA{R: 0, G: 100, B: 100, A: 255},
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

func MakeStatusContainer() *widget.Container {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/minorDialogue.png")
	if err != nil {
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{8, 66 - 16, 8}, [3]int{8, 32 - 16, 8})
	statusContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSliceImage),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(610, 200)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
		),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	return statusContainer
}

func StatusTextInput(textType string) *widget.TextInput {

	face, err := assetManagement.LoadFont(14, assetManagement.November)
	white := color.RGBA{R: 232, G: 225, B: 219, A: 255}
	faceColor := color.RGBA{R: 102, G: 57, B: 48, A: 255}
	if textType == "white" {
		faceColor = white
	}

	if textType == "playerName" {
		face, err = assetManagement.LoadFont(18, assetManagement.November)
		faceColor = color.RGBA{R: 153, G: 229, B: 80, A: 255}
	}
	if textType == "npcName" {
		face, err = assetManagement.LoadFont(18, assetManagement.November)
		faceColor = color.RGBA{R: 250, G: 218, B: 56, A: 255}
	}
	if err != nil {
		log.Fatal(err)
	}
	statusTextInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle: eimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 0}),
		}),
		widget.TextInputOpts.Face(face),

		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          faceColor,
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{R: 254, G: 255, B: 255, A: 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 0),
		),
		widget.TextInputOpts.Placeholder(""),
		widget.TextInputOpts.TabOrder(6),
	)

	return statusTextInput
}

func LoadDrawButtonImage() *widget.ButtonImage {
	/*hoverImg, _, err := ebitenutil.NewImageFromFile("buttonIdle.png")
	if err != nil {
	}*/

	//click imp prints a line through our button to indicate which option was selected
	clickImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/drawButtonClicked.png")

	if err != nil {
		log.Fatal("menuFunctions.go: 265", err)
	}

	//alpha can be changed for debugging purposes
	hover := eimage.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 0})
	idle := eimage.NewNineSliceColor(color.RGBA{R: 0, G: 100, B: 120, A: 0})
	pressed := eimage.NewNineSlice(clickImg, [3]int{54, 132 - 54 - 6, 6}, [3]int{2, 20, 2})
	disabled := eimage.NewNineSlice(clickImg, [3]int{54, 132 - 54 - 6, 6}, [3]int{2, 20, 2})

	return &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
}

func LoadStatusButtonImage() *widget.ButtonImage {
	/*hoverImg, _, err := ebitenutil.NewImageFromFile("buttonIdle.png")
	if err != nil {
	}*/

	statusButtonRaw, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/statusBarButton.png")
	if err != nil {
		log.Fatalf("button image file not loading")
	}

	statusButtonClicked := ebiten.NewImageFromImage(statusButtonRaw.SubImage(image.Rect(0, 0, 36, 23)))
	statusButton := ebiten.NewImageFromImage(statusButtonRaw.SubImage(image.Rect(36, 0, 72, 24)))

	//click imp prints a line through our button to indicate which option was selected

	//alpha can be changed for debugging purposes
	hover := eimage.NewNineSlice(statusButtonClicked, [3]int{3, 20, 14}, [3]int{7, 8, 7})
	idle := eimage.NewNineSlice(statusButton, [3]int{3, 20, 14}, [3]int{7, 8, 7})
	pressed := eimage.NewNineSlice(statusButtonClicked, [3]int{3, 20, 14}, [3]int{7, 8, 7})
	disabled := eimage.NewNineSliceColor(color.RGBA{R: 0, G: 100, B: 120, A: 0})

	return &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
}

func GenerateDrawButton(g *BattleScene) (button *widget.Button) {
	buttonText := "Draw"

	// load gameScenes font, more fonts will be selectable later when we implement a resource manager
	face, err := assetManagement.LoadFont(22, assetManagement.November)
	if err != nil {
		log.Fatal(err)
	}

	// loads a basic button image
	buttonImage := LoadDrawButtonImage()

	//make a new button with the name of each skill as text
	button = widget.NewButton(
		// specify the images to use

		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button

		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			DrawSkillButtonEvent(g, buttonText)
		}),

		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle: color.RGBA{R: 230, G: 10, B: 10, A: 255}}),

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

func CombatSkillBoxContainer(headerText string) *widget.Container {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/menuBackground.png")
	if err != nil {
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{12, 600 - 24, 12}, [3]int{12, 200 - 24, 12})
	//Container to vertically Dialogue with a header
	vertContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSliceImage),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(250, 50)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
		),
		),
	)

	face, err := assetManagement.LoadFont(24, assetManagement.November)
	if err != nil {
		log.Fatal(err)
	}

	//Container to orient Header Label

	headerContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 10)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	headerLbl := widget.NewText(
		widget.TextOpts.Text(headerText, face, color.RGBA{R: 102, G: 57, B: 48, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.TextOpts.Insets(widget.Insets{
			Left:   30,
			Right:  10,
			Top:    20,
			Bottom: 0,
		}),
	)

	//Horizontally organized container for skills buttons

	headerContainer.AddChild(headerLbl)
	vertContainer.AddChild(headerContainer)

	return vertContainer
}

func GenerateCombatSkillButtons(text string, g *BattleScene) (button *widget.Button) {

	// load gameScenes font, more fonts will be selectable later when we implement a resource manager
	face, err := assetManagement.LoadFont(20, assetManagement.November)
	buttonText := strings.ToUpper(string(text[0])) + text[1:]
	if err != nil {
		log.Fatal(err)
	}

	// loads a basic button image
	buttonImage := ui.LoadButtonImage()

	//make a new button with the name of each skill as text
	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			CombatSkillButtonEvent(g, text)
		}),
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle:     color.RGBA{R: 102, G: 57, B: 48, A: 255},
			Disabled: color.RGBA{R: 102, G: 100, B: 255, A: 255},
		}),

		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   84,
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

func GenerateStatusBarButton(g *BattleScene) (button *widget.Button) {

	buttonImage := LoadStatusButtonImage()

	statusButton := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			StatusEffectButtonEvent(g)
		}), widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(50, 50),
			//widget.WidgetOpts.CursorHovered("statusBar"),
			//widget.WidgetOpts.CursorPressed("statusBar"),
		),
		widget.ButtonOpts.TabOrder(1),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position:  widget.RowLayoutPositionEnd,
			MaxWidth:  36,
			MaxHeight: 24,
		},
		),
		),
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
		},
		),
	)
	return statusButton
}

func DialogueSkillButtonEvent(g *BattleScene, text string, button *widget.Button) {
	player := g.battle.CharacterBattleData[battle.Player]
	enemy := g.battle.CharacterBattleData[battle.Enemy]
	enemy.EventTriggered = false
	player.EventTriggered = false
	player.SkillUsed = g.battle.CharacterBattleData[battle.Player].DialogueSkills[text]
	g.battle.UpdateChar(player, enemy)
	g.changeEvent(HideSkillMenu, 15)
	g.inMenu = false
	g.KeepCursorPressed()
	if text == "comeBack" {
		button.GetWidget().Disabled = true
	}
}

func CombatSkillButtonEvent(g *BattleScene, text string) {
	g.battle.TakeCombatTurn(g.battle.CharacterBattleData[battle.Player].CombatSkills[text])
	g.changeEvent(HideCombatMenu, 15)
	g.inMenu = false
	g.KeepCursorPressed()
}

func DrawSkillButtonEvent(g *BattleScene, text string) {
	g.audioPlayer.Play(audioManagement.DrawButton)
	g.TextPrinter.ResetTP()
	g.battle.DialogueTurn(g.battle.CharacterBattleData[battle.Player].DialogueSkills["draw"])
	g.TextPrinter.NextMessage = true
	g.playerBattleSprite.DialogueButtonAnimationTrigger("draw")
	g.changeEvent(HideSkillMenu, 15)
	g.inMenu = false
	g.KeepCursorPressed()
}

func StatusEffectButtonEvent(g *BattleScene) {
	if g.TextPrinter.state == NotPrinting {
		g.StatusButtonEvent = true
	}
}

func GenerateGenericStatusBarButton(printer *TextPrinter) (button *widget.Button) {

	buttonImage := LoadStatusButtonImage()

	statusButton := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			GenericStatusEffectButtonEvent(printer)
		}), widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(50, 50),
			//widget.WidgetOpts.CursorHovered("statusBar"),
			//widget.WidgetOpts.CursorPressed("statusBar"),
		),
		widget.ButtonOpts.TabOrder(1),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position:  widget.RowLayoutPositionEnd,
			MaxWidth:  36,
			MaxHeight: 24,
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

func GenericStatusEffectButtonEvent(printer *TextPrinter) {

	if printer.NextMessage == false {

		println("resetting printer and moving cursor to Menu")

		printer.ResetTP()

	} else {
		//clear the last output
		printer.stringPosition = 1

		printer.StatusText[0].SetText("")
		printer.StatusText[1].SetText("")
		printer.StatusText[2].SetText("")

		//if there are more lines of the message trigger the printer again

		printer.NextMessage = true
	}
}

func MinorDialogueContainer() *widget.Container {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/minorDialogue.png")
	if err != nil {
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{8, 66 - 16, 8}, [3]int{8, 32 - 16, 8})
	statusContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(nineSliceImage),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 100)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
		),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	return statusContainer
}
