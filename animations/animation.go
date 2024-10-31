package animations

type Animation struct {
	FirstF       int
	LastF        int
	Step         int //how many indices do we move per frame
	SpeedinTPS   float32
	frameCounter float32
	frame        int
}

func (a *Animation) Update() {
	//this code iteration assumes each animation loops
	a.frameCounter -= 1.0 // no need to worry about time as ebiten has a locked frame rate
	if a.frameCounter < 0 {
		a.frameCounter = a.SpeedinTPS
		a.frame += a.Step

		if a.frame > a.LastF {
			//loop back to beginning of animation
			a.frame = a.FirstF
		}
	}
}
func (a *Animation) Frame() int {
	return a.frame
}

func (a *Animation) ResetFrame() {
	a.frame = a.FirstF
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
