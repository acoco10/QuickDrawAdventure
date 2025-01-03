package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
)

type EffectState uint8

const (
	Triggered EffectState = iota
	NotTriggered
)

type GraphicEffects interface {
	Draw(screen *ebiten.Image)
	Update()
	AccessImage() *ebiten.Image
	Trigger()
	CheckState() EffectState
}

type StaticEffect struct {
	img      *ebiten.Image
	state    EffectState
	duration int
	x, y     float64
	counter  int
}

func (se *StaticEffect) CheckState() EffectState {
	return se.state
}

func (se *StaticEffect) Trigger() {
	se.state = Triggered
}

func (se *StaticEffect) Draw(screen *ebiten.Image) {

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Reset()
	opts.GeoM.Translate(se.x, se.y)
	screen.DrawImage(se.img, opts)
}

func (se *StaticEffect) Update() {
	if se.state == Triggered {
		se.counter--
		if se.counter < 0 {
			se.state = NotTriggered
			se.counter = se.duration
		}
	}
}

func (se *StaticEffect) Frame() int {
	return se.counter
}

func (se *StaticEffect) AccessImage() *ebiten.Image {
	return se.img
}

func NewStaticEffect(img *ebiten.Image, x, y float64) *StaticEffect {
	return &StaticEffect{
		img:      img,
		x:        x,
		y:        y,
		counter:  50,
		duration: 50,
		state:    NotTriggered,
	}
}

type AnimatedEffect struct {
	img          *ebiten.Image
	spriteSheet  *spritesheet.SpriteSheet
	x, y         float64
	firstFrame   int
	lastFrame    int
	frame        int
	step         int
	speed        int
	frameCounter int
	cycleCounter int
	cycles       int
	state        EffectState
}

func (e *AnimatedEffect) Draw(screen *ebiten.Image) {
	img := e.img.SubImage(e.spriteSheet.Rect(e.frame)).(*ebiten.Image)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(4, 4)
	opts.GeoM.Translate(e.x, e.y)
	screen.DrawImage(img, opts)
}

func (e *AnimatedEffect) CheckState() EffectState {
	return e.state
}

func (e *AnimatedEffect) AccessImage() *ebiten.Image {
	return e.img
}

func (e *AnimatedEffect) Trigger() {
	e.state = Triggered
}

func (e *AnimatedEffect) Update() {
	if e.state == Triggered {
		e.frameCounter -= 1.0
		if e.frameCounter < 0 {
			e.frameCounter = e.speed
			e.frame += e.step

			if e.frame == e.lastFrame {
				e.frame = e.lastFrame

				e.state = NotTriggered
			}
		}
	}
}

func (e *AnimatedEffect) Image() *ebiten.Image {
	return e.img
}

func NewEffect(img *ebiten.Image, sheet *spritesheet.SpriteSheet, x float64, y float64, lastF int, firstF int, step int, speed int) *AnimatedEffect {
	effect := &AnimatedEffect{
		spriteSheet:  sheet,
		img:          img,
		x:            x,
		y:            y,
		firstFrame:   firstF,
		lastFrame:    lastF,
		frame:        firstF,
		step:         step,
		speed:        speed,
		frameCounter: speed,
		cycleCounter: 1,
		cycles:       1,
		state:        NotTriggered,
	}
	return effect
}
