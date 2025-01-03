package animations

type CyclicAnimation struct {
	Animation *Animation
	cycles    int
}

func (ca *CyclicAnimation) GetCycles() int {
	return ca.cycles
}

func NewCyclicAnimation(firstF int, lastF int, step int, speedinTPS float32, cycles int) *CyclicAnimation {
	cyclicAnimation := &CyclicAnimation{
		NewAnimation(firstF, lastF, step, speedinTPS),
		cycles,
	}
	return cyclicAnimation
}
