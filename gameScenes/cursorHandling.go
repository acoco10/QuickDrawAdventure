package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"io/fs"
	"log"
)

type BattleMenuCursorUpdater struct {
	currentPosition image.Point
	systemPosition  image.Point
	statusX         int
	statusY         int
	cursorImages    map[string]*ebiten.Image
	counter         int
	maxY            int
	minY            int
	maxX            int
	minX            int
	countdown       int
	pressed         bool
	Event           EventName
}

func CreateCursorUpdater(resWidth int, resHeight int) *BattleMenuCursorUpdater {
	cu := BattleMenuCursorUpdater{}
	X, Y := ebiten.CursorPosition()
	X1, Y1 := int(0.66*float64(resWidth)), int(0.75*float64(resHeight))
	cu.statusX = X1
	cu.statusY = Y1
	cu.systemPosition = image.Point{X1, Y1}
	cu.currentPosition = image.Point{X, Y}
	cu.cursorImages = make(map[string]*ebiten.Image)
	cu.cursorImages[input.CURSOR_DEFAULT] = loadNormalCursorImage()
	cu.cursorImages["buttonHover"] = loadHoverCursorImage()
	cu.cursorImages["buttonPressed"] = loadPressedCursorImage()
	cu.cursorImages["statusBar"] = loadNormalCursorImage()
	cu.countdown = 0
	cu.maxY = Y1
	cu.minY = Y1
	cu.minX = X1
	return &cu
}

func (cu *BattleMenuCursorUpdater) MoveToLockedSpecificPosition(x, y int) {
	cu.minX = x
	cu.minY = y
	cu.maxY = y
}

func (cu *BattleMenuCursorUpdater) MoveCursorToSkillMenu() {
	cu.minX = 138
	cu.minY = 564
	cu.maxY = 564 + 35 + 35 + 35
}

func (cu *BattleMenuCursorUpdater) MoveCursorToStatusBar() {
	cu.minX = cu.statusX
	cu.minY = cu.statusY
	cu.maxY = cu.statusY
}

const (
	KeepPressedOff EventName = iota
	noEvent
)

func (cu *BattleMenuCursorUpdater) changeEvent(name EventName, timer int) {
	cu.Event = name
	cu.countdown = timer
}

func (cu *BattleMenuCursorUpdater) TriggerEvent(name EventName) {

	if name == KeepPressedOff {
		cu.PressedOff()
		cu.MoveCursorToStatusBar()
	}

}

func (cu *BattleMenuCursorUpdater) keepPressed(timer int) {
	cu.pressed = true
	cu.changeEvent(KeepPressedOff, timer)

}

func (cu *BattleMenuCursorUpdater) PressedOff() {
	cu.pressed = false
	cu.countdown = 0
}

func (cu *BattleMenuCursorUpdater) MoveCursorToCombatMenu() {
	cu.minX = 138
	cu.minY = 564
	cu.maxY = 564 + 35 + 35 + 35
}

// Called every Update call from Ebiten
// Note that before this is called the current cursor shape is reset to DEFAULT every cycle

func (cu *BattleMenuCursorUpdater) Update() {

	X, Y := ebiten.CursorPosition()

	diffX := cu.systemPosition.X - X
	diffY := cu.systemPosition.Y - Y

	cu.currentPosition.X -= diffX
	cu.currentPosition.Y -= diffY

	cu.currentPosition.X = cu.minX

	/*if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		cu.currentPosition.X -= 10
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		cu.currentPosition.X += 10
	}*/
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		cu.currentPosition.Y -= 35
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		cu.currentPosition.Y += 35
	}
	/*if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		cu.cuCounterOn = true
	}*/

	if cu.Event != noEvent && cu.countdown == 1 {
		cu.TriggerEvent(cu.Event)
		cu.MoveCursorToStatusBar()
	}

	if cu.currentPosition.Y <= cu.minY {
		cu.currentPosition.Y = cu.minY
	}

	if cu.currentPosition.Y >= cu.maxY {
		cu.currentPosition.Y = cu.maxY
	}

	cu.systemPosition = image.Point{X, Y}

}
func (cu *BattleMenuCursorUpdater) Draw(screen *ebiten.Image) {
}
func (cu *BattleMenuCursorUpdater) AfterDraw(screen *ebiten.Image) {
}

// MouseButtonPressed returns whether mouse button b is currently pressed.
func (cu *BattleMenuCursorUpdater) MouseButtonPressed(b ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(b) || ebiten.IsKeyPressed(ebiten.KeyEnter) || cu.pressed
}

// MouseButtonJustPressed returns whether mouse button b has just been pressed.
// It only returns true during the first frame that the button is pressed.
func (cu *BattleMenuCursorUpdater) MouseButtonJustPressed(b ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(b) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// CursorPosition returns the current cursor position.
// If you define a CursorPosition that doesn't align with a system cursor you will need to
// set the CursorDrawMode to Custom. This is because ebiten doesn't have a way to set the
// cursor location manually
func (cu *BattleMenuCursorUpdater) CursorPosition() (int, int) {
	return cu.currentPosition.X, cu.currentPosition.Y
}

// GetCursorImage Returns the image to use as the cursor
// EbitenUI by default will look for the following cursors:
//
//	"EWResize"
//	"NSResize"
//	"Default"
func (cu *BattleMenuCursorUpdater) GetCursorImage(name string) *ebiten.Image {
	return cu.cursorImages[name]
}

// GetCursorOffset Returns how far from the CursorPosition to offset the cursor image.
// This is best used with cursors such as resizing.
func (cu *BattleMenuCursorUpdater) GetCursorOffset(name string) image.Point {
	return image.Point{}
}

// Layout implements gameScenes.

func loadNormalCursorImage() *ebiten.Image {
	f, err := assets.ImagesDir.Open("images/menuAssets/buttonCursor.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f fs.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 0, 32, 25)))
	//(64, 0, 87, 16)
}

func loadHoverCursorImage() *ebiten.Image {
	f, err := assets.ImagesDir.Open("images/menuAssets/buttonCursor.png")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 0, 32, 25)))
}

func loadPressedCursorImage() *ebiten.Image {
	f, err := assets.ImagesDir.Open("images/menuAssets/buttonCursor.png")
	if err != nil {
		return nil
	}
	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(32, 0, 64, 25)))
}
