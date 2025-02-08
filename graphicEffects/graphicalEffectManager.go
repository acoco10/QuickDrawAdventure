package graphicEffects

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EffectType uint8

const (
	NoEffect EffectType = iota
	DrawEffect
	StareEffect
	WeaknessEffect
	FearEffect
	HitSplatEffect
	SuccessfulEffect
	UnsuccessfulEffect
	MissEffect
	TookDamageEffect
	UnsuccesfulStareEffect
	FillerEffect
	MuzzleEffect
)

type GraphicalEffectManager struct {
	PlayerEffects *GraphicalEffectSequencer
	EnemyEffects  *GraphicalEffectSequencer
	GameEffects   *GraphicalEffectSequencer
}

type GraphicalEffectSequencer struct {
	effects     map[EffectType]GraphicEffect
	EffectQueue []GraphicEffect
	effectIndex int
	state       EffectState
	Counter     int
	Configured  bool
}

func (e *GraphicalEffectSequencer) TriggerEffectQueue() {
	println("triggering effect queue len =", len(e.EffectQueue))
	if len(e.EffectQueue) > 0 {
		e.EffectQueue[0].Trigger()
		e.state = Triggered
	}
}

func (e *GraphicalEffectSequencer) GetState() EffectState {
	return e.state
}

func (e *GraphicalEffectSequencer) Update() {
	if e.Counter > 0 {
		e.Counter--
	}
	if e.state == Triggered && e.Counter < 1 {
		effect := e.EffectQueue[0]
		effect.Update()
		if effect.CheckState() == NotTriggered {
			e.effectIndex++
			e.EffectQueue = e.EffectQueue[1:]
			if len(e.EffectQueue) > 0 && e.EffectQueue[0] != nil {
				e.EffectQueue[0].Trigger()
				e.Counter = 1
			} else {
				e.effectIndex = 0
				e.Counter = 0
				e.state = NotTriggered
			}
		}
	}
}

func (e *GraphicalEffectSequencer) Draw(screen *ebiten.Image) {
	if e.state == Triggered {
		if len(e.EffectQueue) > 0 && e.EffectQueue[0] != nil {
			effect := e.EffectQueue[0]
			effect.Draw(screen)
		}

	}
}

func NewGraphicalEffectManager() *GraphicalEffectManager {
	gef := GraphicalEffectManager{
		PlayerEffects: &GraphicalEffectSequencer{},
		EnemyEffects:  &GraphicalEffectSequencer{},
		GameEffects:   &GraphicalEffectSequencer{},
	}
	gef.PlayerEffects.state = NotTriggered
	gef.EnemyEffects.state = NotTriggered
	gef.GameEffects.state = NotTriggered

	return &gef
}
