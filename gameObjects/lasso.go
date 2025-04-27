package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/camera"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type Lasso struct {
	chargeMax                    float32
	charge                       float32
	Triggered                    bool
	x1, y1, x2, y2, Length, maxX float32
	direction                    Direction
	timer                        int
}

func (l *Lasso) SetPosition(char Character) {
	l.timer = 20
	l.x1 = float32(char.X*4) + float32(char.SpriteSheet.SpriteWidth)*3.8
	l.y1 = float32(char.Y*4) + float32(char.SpriteSheet.SpriteHeight)*2
	l.y2 = l.y1 - 4
	l.x2 = l.x1
	l.direction = char.Direction

}

func (l *Lasso) SetLength() {
	l.Length = l.charge * 100
	if l.direction == Left {
		l.maxX = l.x2 - l.Length
	}
	if l.direction == Right {
		l.maxX = l.x1 + l.Length
	}
}

func (l *Lasso) Update() Signal {
	if l.Triggered {
		if l.x2 <= l.maxX && l.direction == Right {
			l.x2 += 50
			l.y2 -= .5
		} else if l.x2 >= l.maxX && l.direction == Left {
			l.x2 += 50
			l.y2 -= .5
		} else {
			l.timer--
			if l.timer == 0 {
				l.Triggered = false
				l.x1 = l.x2
				l.timer = 20
				charge := l.charge
				l.charge = 0
				return Signal{value: charge}
			}
		}
	} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if l.charge < l.chargeMax {
			l.charge += 0.1
		}
	} else if inpututil.IsKeyJustReleased(ebiten.KeySpace) {
		l.SetLength()
		l.Trigger()
	}
	return Signal{}
}

func (l *Lasso) Trigger() float32 {
	l.Triggered = true
	return l.charge
}

func (l *Lasso) Draw(screen *ebiten.Image, cam camera.Camera) {
	if l.Triggered {
		x1 := l.x1 + float32(cam.X*4)
		x2 := l.x2 + float32(cam.X*4)
		y1 := l.y1 + float32(cam.Y*4)
		y2 := l.y2 + float32(cam.Y*4)
		vector.StrokeLine(screen, x1, y1, x2, y2, 4, color.RGBA{150, 100, 100, 255}, false)

	}
}
