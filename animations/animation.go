package animations

type AnimationObj interface {
	Update()
	Reset()
	Frame() int
	Stop() int
	Coord() (float64, float64)
}

type MovementAnimation struct {
	x, y       float64
	StartingX  float64
	StartingY  float64
	XUpdate    float64
	YUpdate    float64
	SpeedInTPS int
	Counter    int
	frame      int
	LastFrame  int
}

func (ma *MovementAnimation) Update() {
	ma.Counter--
	if ma.Counter < 0 {
		ma.frame++
		ma.Counter = ma.SpeedInTPS
		ma.x += ma.XUpdate
		ma.y += ma.YUpdate
	}
}

func (ma *MovementAnimation) Reset() {
	ma.x = ma.StartingX
	ma.y = ma.StartingY
	ma.frame = 0
}

func (ma *MovementAnimation) Frame() int {
	return ma.frame
}

func (ma *MovementAnimation) Stop() int {
	return ma.LastFrame
}
func (ma *MovementAnimation) Coord() (float64, float64) {
	return ma.x, ma.y
}

func NewMAnimation(startingX float64, startingY float64, XUp float64, YUp float64, speed int, lastF int) *MovementAnimation {
	return &MovementAnimation{
		startingX,
		startingY,
		startingX,
		startingY,
		XUp,
		YUp,
		speed,
		speed,
		0,
		lastF,
	}
}

type Animation struct {
	FirstF       int
	LastF        int
	Step         int
	SpeedInTPS   float32
	frameCounter float32
	frame        int
}

func (a *Animation) Update() {
	//this code iteration assumes each animation loops
	a.frameCounter -= 1.0 // no need to worry about time as ebiten has a locked frame rate
	if a.frameCounter < 0 {
		a.frameCounter = a.SpeedInTPS
		a.frame += a.Step

		if a.frame > a.LastF {
			a.frame = a.FirstF
		}
	}
}

func (a *Animation) Frame() int {
	return a.frame
}

func (a *Animation) Reset() {
	a.frame = a.FirstF
}

func (a *Animation) Stop() int {
	return a.LastF
}

func (a *Animation) Coord() (float64, float64) {
	return 0, 0
}

func NewAnimation(firstF int, lastF int, step int, speedinTPS float32) *Animation {
	return &Animation{
		firstF,
		lastF,
		step,
		speedinTPS,
		speedinTPS,
		firstF,
	}
}
