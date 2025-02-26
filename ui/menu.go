package ui

import (
	"github.com/acoco10/QuickDrawAdventure/assetManagement"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"log"
)

type Menu struct {
	MenuContainer    *widget.Container
	Buttons          []*widget.Button
	ButtonVisibility bool
}

func (m *Menu) DisableButtonVisibility() {
	for _, b := range m.Buttons {
		b.GetWidget().Visibility = widget.Visibility_Hide
	}
}

func (m *Menu) DisableButtons() {
	for _, b := range m.Buttons {
		b.GetWidget().Disabled = true
	}
}

func (m *Menu) EnableButtonVisibility() {
	for _, b := range m.Buttons {
		b.GetWidget().Visibility = widget.Visibility_Show
		b.GetWidget().Disabled = false
	}
}

func LoadButtonImage() *widget.ButtonImage {
	/*hoverImg, _, err := ebitenutil.NewImageFromFile("buttonIdle.png")
	if err != nil {
	}*/

	//click imp prints a line through our button to indicate which option was selected
	clickImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/buttonClicked.png")

	if err != nil {
		log.Fatal(err)
	}

	//alpha can be changed for debugging purposes
	hover := image.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 100})
	idle := image.NewNineSliceColor(color.RGBA{R: 0, G: 100, B: 120, A: 0})
	pressed := image.NewNineSlice(clickImg, [3]int{54, 132 - 54 - 6, 6}, [3]int{2, 20, 2})
	disabled := image.NewNineSlice(clickImg, [3]int{54, 132 - 54 - 6, 6}, [3]int{2, 20, 2})

	return &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
}

func LoadSlotButtonImage() *widget.ButtonImage {
	/*hoverImg, _, err := ebitenutil.NewImageFromFile("buttonIdle.png")
	if err != nil {
	}*/

	//click imp prints a line through our button to indicate which option was selected
	clickImg, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/buttonClicked.png")

	if err != nil {
		log.Fatal(err)
	}

	//alpha can be changed for debugging purposes
	hover := image.NewNineSliceColor(color.RGBA{R: 100, G: 100, B: 120, A: 100})
	idle := image.NewNineSliceColor(color.RGBA{R: 0, G: 100, B: 120, A: 0})
	pressed := image.NewNineSlice(clickImg, [3]int{54, 132 - 54 - 6, 6}, [3]int{2, 20, 2})
	disabled := image.NewNineSlice(clickImg, [3]int{54, 132 - 54 - 6, 6}, [3]int{2, 20, 2})

	return &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: disabled,
	}
}

func SkillBoxContainer(headerText string) *widget.Container {
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/menuBackground.png")
	if err != nil {
	}

	nineSliceImage := image.NewNineSlice(img, [3]int{12, 600 - 24, 12}, [3]int{12, 200 - 24, 12})
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

func SkillsContainer() *widget.Container {
	skillContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(
				widget.Insets{Right: 0, Left: 0, Top: 0, Bottom: 20}),
			widget.RowLayoutOpts.Spacing(5),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 50)),
	)

	return skillContainer
}
