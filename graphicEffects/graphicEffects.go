package graphicEffects

import (
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
)

type EffectState uint8

const (
	Triggered EffectState = iota
	NotTriggered
)

type GraphicEffectType uint8

const (
	Static GraphicEffectType = iota
	Animated
)

type GraphicEffect interface {
	Draw(screen *ebiten.Image, depth int)
	Update()
	AccessImage() *ebiten.Image

	Trigger()
	CheckState() EffectState
	UnTrigger()
	Type() GraphicEffectType
	SetCoord(x float64, y float64)
	SetDepth(depth int)
	SetDuration(frames int)
}

type StaticEffect struct {
	name       EffectType
	img        *ebiten.Image
	state      EffectState
	duration   int
	indefinite bool
	x, y       float64
	counter    int
	scale      float64
	depth      int
}

func (se *StaticEffect) CheckState() EffectState {
	return se.state
}

func (se *StaticEffect) Trigger() {
	se.state = Triggered
}
func (se *StaticEffect) UnTrigger() {
	se.state = NotTriggered
}
func (se *StaticEffect) SetCoord(x float64, y float64) {
	se.x = x
	se.y = y
}

func (se *StaticEffect) Draw(screen *ebiten.Image, depth int) {
	if se.state == Triggered && se.depth == depth {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Reset()
		opts.GeoM.Translate(se.x, se.y)
		opts.GeoM.Scale(se.scale, se.scale)
		screen.DrawImage(se.img, opts)
	}

}

func (se *StaticEffect) SetDepth(depth int) {
	se.depth = depth

}

func (se *StaticEffect) Update() {
	if se.state == Triggered {
		if !se.indefinite {
			se.counter--
			if se.counter < 0 {
				se.state = NotTriggered
				se.counter = se.duration
			}
		}
	}
}

func (se *StaticEffect) Frame() int {
	return se.counter
}
func (se *StaticEffect) SetDuration(frames int) {
	se.duration = frames
}

func (se *StaticEffect) AccessImage() *ebiten.Image {
	return se.img
}
func (se *StaticEffect) Type() GraphicEffectType {
	return Static
}

func NewStaticEffect(img *ebiten.Image, x, y float64, duration int, scale float64, effectType EffectType) *StaticEffect {
	se := StaticEffect{
		img:        img,
		x:          x,
		y:          y,
		counter:    50,
		duration:   duration,
		state:      NotTriggered,
		scale:      scale,
		indefinite: false,
		name:       effectType,
	}
	if duration == 0 {
		se.indefinite = true
	}
	return &se
}

type AnimatedEffect struct {
	img          *ebiten.Image
	spriteSheet  *spritesheet.SpriteSheet
	x, y         float64
	firstFrame   int
	LastFrame    int
	frame        int
	step         int
	speed        int
	FrameCounter int
	cycleCounter int
	cycles       int
	state        EffectState
	scale        float64
	effectType   GraphicEffectType
	visible      bool
	depth        int
}

func (e *AnimatedEffect) SetDuration(frames int) {
	//filler func as animated effects have set duration based on animation/frame speed
	//but needed for static effect
}

func (e *AnimatedEffect) UnTrigger() {
	e.state = NotTriggered
}
func (e *AnimatedEffect) SetCoord(x float64, y float64) {
	e.x = x
	e.y = y
}
func (e *AnimatedEffect) SetDepth(depth int) {
	e.depth = depth
}

func (e *AnimatedEffect) Draw(screen *ebiten.Image, depth int) {
	if (e.state == Triggered || e.visible) && e.depth == depth {
		img := e.img.SubImage(e.spriteSheet.Rect(e.frame)).(*ebiten.Image)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(e.scale, e.scale)
		opts.GeoM.Translate(e.x, e.y)
		screen.DrawImage(img, opts)
	}
}
func (e *AnimatedEffect) MakeVisible() {
	e.visible = true
}
func (e *AnimatedEffect) MakeNoyVisible() {
	e.visible = false
}

func (e *AnimatedEffect) Frame() int {
	return e.frame
}
func (e *AnimatedEffect) Reset() {
	e.frame = 0
	e.FrameCounter = 0
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
		e.FrameCounter -= 1.0
		if e.FrameCounter < 0 {
			e.FrameCounter = e.speed
			e.frame += e.step

			if e.frame == e.LastFrame {
				e.frame = e.firstFrame
				e.state = NotTriggered
			}
		}
	}
}

func (e *AnimatedEffect) Image() *ebiten.Image {
	return e.img
}

func (e *AnimatedEffect) Type() GraphicEffectType {
	return Animated
}

func NewEffect(img *ebiten.Image, sheet *spritesheet.SpriteSheet, x float64, y float64, lastF int, firstF int, step int, speed int, scale float64) *AnimatedEffect {
	effect := &AnimatedEffect{
		spriteSheet:  sheet,
		img:          img,
		x:            x,
		y:            y,
		firstFrame:   firstF,
		LastFrame:    lastF,
		frame:        firstF,
		step:         step,
		speed:        speed,
		FrameCounter: speed,
		cycleCounter: 1,
		cycles:       1,
		state:        NotTriggered,
		scale:        scale,
		effectType:   Animated,
	}
	return effect
}
