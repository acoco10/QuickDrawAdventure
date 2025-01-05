package gameScenes

import (
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/spritesheet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
	"math/rand/v2"
	"strconv"
)

type GraphicalEffectManager struct {
	PlayerEffects GraphicalEffectSequencer
	EnemyEffects  GraphicalEffectSequencer
	GameEffects   GraphicalEffectSequencer
}
type EffectType uint8

const (
	NoEffect EffectType = iota
	DrawEffect
	DamageEffect
	StareEffect
	EnemyStareEffect
)

type GraphicalEffectSequencer struct {
	effects     map[EffectType]GraphicEffects
	EffectQueue []GraphicEffects
	state       EffectState
	counter     int
	configured  bool
}

func (e *GraphicalEffectSequencer) TriggerEffectQueue() {
	if len(e.EffectQueue) > 0 {
		e.EffectQueue[0].Trigger()
		e.state = Triggered
	}
}

func (e *GraphicalEffectSequencer) Update() {
	if e.state == Triggered {
		effect := e.EffectQueue[0]
		effect.Update()
		if effect.CheckState() == NotTriggered {
			e.EffectQueue = e.EffectQueue[1:]
			if len(e.EffectQueue) > 0 {
				e.EffectQueue[0].Trigger()
			} else {
				e.state = NotTriggered
			}
		}
	}
}

func (e *GraphicalEffectSequencer) ProcessPlayerTurnData(turn *battle.Turn) {
	e.EffectQueue = make([]GraphicEffects, 0)
	e.configured = true
	if turn.PlayerSkillUsed.SkillName == "stare down" {
		e.EffectQueue = append(e.EffectQueue, e.effects[StareEffect])
	}
	if turn.PlayerSkillUsed.SkillName == "draw" {
		if turn.EnemySkillUsed.SkillName != "draw" {
			e.EffectQueue = append(e.EffectQueue, e.effects[DrawEffect])
		} else if turn.PlayerStartIndex == 1 {
			e.EffectQueue = append(e.EffectQueue, e.effects[DrawEffect])
		}
	}

	for _, result := range turn.DamageToEnemy {
		if result > 0 {

			println("appending damage effect to playerBattleSprite effects\n", "result:", result, "\n")
			face, err := LoadFont(15)
			if err != nil {
				log.Fatal(err)
			}

			oneDamageImg, _, err := ebitenutil.NewImageFromFile("assets/images/effectAssets/1damage.png")
			if err != nil {
				log.Printf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
			}

			damage := strconv.FormatInt(int64(result), 10)
			dopts := text.DrawOptions{}
			dopts.GeoM.Translate(16, 10)
			text.Draw(oneDamageImg, damage, face, &dopts)
			damageEffect := NewStaticEffect(oneDamageImg, 620+float64(rand.IntN(35)), 170+float64(rand.IntN(80)))
			e.EffectQueue = append(e.EffectQueue, damageEffect)
			println("graphicalEffectManager.go:76 length of effect Queue =", len(e.EffectQueue), "\n")
		}
		if result == 0 {
			missImage, _, err := ebitenutil.NewImageFromFile("assets/images/miss.png")
			if err != nil {
				log.Printf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
			}

			missEffect := NewStaticEffect(missImage, 520+float64(rand.IntN(80)), 200+float64(rand.IntN(80)))
			e.EffectQueue = append(e.EffectQueue, missEffect)
		}
	}
}

func (e *GraphicalEffectSequencer) ProcessEnemyTurnData(turn *battle.Turn) {
	e.EffectQueue = make([]GraphicEffects, 0)
	e.configured = true
	if turn.EnemySkillUsed.SkillName == "draw" {
		if turn.PlayerSkillUsed.SkillName != "draw" {
			e.EffectQueue = append(e.EffectQueue, e.effects[DrawEffect])
		} else if turn.EnemyStartIndex == 1 {
			e.EffectQueue = append(e.EffectQueue, e.effects[DrawEffect])
		}
	}
	if turn.EnemySkillUsed.SkillName == "stare down" {
		e.EffectQueue = append(e.EffectQueue, e.effects[EnemyStareEffect])
	}

	for _, result := range turn.DamageToPlayer {
		if result > 0 {
			println("appending damage effect to enemyBattleSprite effects\n", "result:", result)
			face, err := LoadFont(15)
			if err != nil {
				log.Fatal(err)
			}
			oneDamageImg, _, err := ebitenutil.NewImageFromFile("assets/images/effectAssets/1damage.png")
			if err != nil {
				log.Printf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
			}

			damage := strconv.FormatInt(int64(result), 10)
			dopts := text.DrawOptions{}
			dopts.GeoM.Translate(16, 10)
			text.Draw(oneDamageImg, damage, face, &dopts)
			damageEffect := NewStaticEffect(oneDamageImg, 740+float64(rand.IntN(20)), 430+float64(rand.IntN(35)))
			e.EffectQueue = append(e.EffectQueue, damageEffect)
		}
		if result == 0 {
			missImage, _, err := ebitenutil.NewImageFromFile("assets/images/miss.png")
			if err != nil {
				log.Printf("ebitenutil.NewImageFromFile file not found due to: %s\n", err)
			}

			missEffect := NewStaticEffect(missImage, 750+float64(rand.IntN(15)), 350+float64(rand.IntN(35)))
			e.EffectQueue = append(e.EffectQueue, missEffect)
		}
		println("graphicalEffectManager.go:76 length of enemyBattleSprite effect Queue =", len(e.EffectQueue), "\n")

	}

}

func (e *GraphicalEffectSequencer) loadCharacterEffects() {

	drawImg, _, err := ebitenutil.NewImageFromFile("assets/images/effectAssets/DrawAffect.png")
	if err != nil {
		log.Fatalf("ebitenutil.NewImageFromFile file not found%s\n", err)
	}

	drawEffectSpriteSheet := spritesheet.NewSpritesheet(1, 4, 199, 125)
	drawEffect := NewEffect(drawImg, drawEffectSpriteSheet, 450, 200, 3, 0, 1, 12, 4)

	enemyStaredownimg, _, err := ebitenutil.NewImageFromFile("assets/images/sheriffStaredownAnimationSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}

	sheriffStareSpriteSheet := spritesheet.NewSpritesheet(7, 1, 320, 180)
	enemyStareEffect := NewEffect(enemyStaredownimg, sheriffStareSpriteSheet, 0, 0, 6, 0, 1, 30, 5)

	staredownimg, _, err := ebitenutil.NewImageFromFile("assets/images/staredownAnimationSpriteSheet.png")
	if err != nil {
		log.Fatal(err)
	}
	stareSpriteSheet := spritesheet.NewSpritesheet(4, 1, 320, 180)
	stareEffect := NewEffect(staredownimg, stareSpriteSheet, 0, 0, 3, 0, 1, 30, 5)

	effects := map[EffectType]GraphicEffects{
		DrawEffect:       drawEffect,
		StareEffect:      stareEffect,
		EnemyStareEffect: enemyStareEffect,
	}
	e.effects = effects
}

func (e *GraphicalEffectSequencer) Draw(screen *ebiten.Image) {
	if e.state == Triggered {
		effect := e.EffectQueue[0]
		effect.Draw(screen)
	}
}

func NewGraphicalEffectManager() GraphicalEffectManager {
	gef := GraphicalEffectManager{
		PlayerEffects: GraphicalEffectSequencer{},
		EnemyEffects:  GraphicalEffectSequencer{}}

	gef.PlayerEffects.loadCharacterEffects()
	gef.EnemyEffects.loadCharacterEffects()
	gef.PlayerEffects.counter = 10
	gef.PlayerEffects.state = NotTriggered
	gef.EnemyEffects.state = NotTriggered

	return gef
}
