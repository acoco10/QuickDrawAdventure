package ui

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
)

type MenuState uint8

const (
	Cover MenuState = iota
	Reading
	Completed
)

type TextBlockMenu struct {
	ui        ebitenui.UI
	triggered bool
	cover     *ebiten.Image
	text      []string
	index     int
	textArea  *widget.TextArea
	State     MenuState
	SkillId   int
}

func (t *TextBlockMenu) Init() {

	coverImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/items/topNotchWesternCover.png")

	if err != nil {
		log.Fatal(err)
	}

	t.cover = coverImg

	text := []string{"Top Notch Westerns Volume 1. \n\nAs they stared each other down, Shiela looked Billy in the eye and said,\n\n\"Your so dumb if your brains was dynamite you wouldn't have enough to blow your nose\"",
		"Elyse learned a new skill! Dynamite for Brains!"}
	t.text = text

	face, err := assetManagement.LoadFont(24, assetManagement.November)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/menuBackground.png")
	if err != nil {
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{12, 600 - 24, 12}, [3]int{12, 200 - 24, 12})

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	// construct a textarea
	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				//Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionCenter,
					MaxWidth:  600,
					MaxHeight: 600,
				}),
				//Set the minimum size for the widget
				widget.WidgetOpts.MinSize(600, 600),
			),
		),
		//Set gap between scrollbar and text
		widget.TextAreaOpts.ControlWidgetSpacing(2),
		//Tell the textarea to display bbcodes
		widget.TextAreaOpts.ProcessBBCode(true),
		//Set the font color
		widget.TextAreaOpts.FontColor(color.Black),
		//Set the font face (size) to use
		widget.TextAreaOpts.FontFace(face),
		//Set the initial text for the textarea
		//It will automatically line wrap and process newlines characters
		//If ProcessBBCode is true it will parse out bbcode
		widget.TextAreaOpts.Text(text[0]),
		//Tell the TextArea to show the vertical scrollbar
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(30)),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: nineSliceImage,
				Mask: eimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{
					Idle:  eimage.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
					Hover: eimage.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle:    eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Hover:   eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Pressed: eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				},
			),
		),
	)
	t.textArea = textarea
	//Add text to the end of the textarea
	//textarea.AppendText("\nLast Row")
	//Add text to the beginning of the textarea
	//textarea.PrependText("First Row\n")
	//Replace the current text with the new value
	//textarea.SetText("New Value!")
	//Retrieve the current value of the text area text
	fmt.Println(textarea.GetText())
	// add the textarea as a child of the container
	rootContainer.AddChild(textarea)

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	t.ui = ui
	t.index = 0
	t.State = Cover
	t.SkillId = 6

}

func (t *TextBlockMenu) Update() {
	if t.triggered {
		t.ui.Update()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if t.State == Reading {
			t.index++
			if t.index == len(t.text) {
				t.UnTrigger()
				t.State = Completed
				t.index = 0
			} else {
				t.textArea.SetText(t.text[t.index])
			}
		}
		if t.State == Cover {
			t.State = Reading
		}

	}

}

func (t *TextBlockMenu) ReturnSkillID() int {
	return t.SkillId
}

func (t *TextBlockMenu) Reset() {
	t.State = Cover
	t.index = 0
	t.textArea.SetText(t.text[t.index])
}

func (t *TextBlockMenu) Draw(screen *ebiten.Image) {
	if t.triggered {
		if t.State == Cover {
			opts := ebiten.DrawImageOptions{}
			opts.GeoM.Scale(3, 3)
			opts.GeoM.Translate(550, 100)
			screen.DrawImage(t.cover, &opts)
		} else {
			t.ui.Draw(screen)
		}
	}
}

func (t *TextBlockMenu) Trigger() {
	t.triggered = true
}

func (t *TextBlockMenu) UnTrigger() {
	t.triggered = false
}
