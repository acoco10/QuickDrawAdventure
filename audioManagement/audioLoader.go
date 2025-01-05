package audioManagement

import (
	"bytes"
	"embed"
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"io"
	"log"
)

//go:embed sounds/combatSFX
var combatSFX embed.FS

//go:embed sounds/music
var music embed.FS

//go:embed sounds/menuSFX

var menuSFX embed.FS

const (
	audioNone resource.AudioID = iota
	GunCock
	GunShot
	BulletHit
	BulletMiss
	ElyseGrunt
	EnemyGrunt
	PistolUnHolster
	DrawButton
	NoAmmo
	Reload
	IntroMusic
	DialogueMusic
	BattleMusic
	BlowingSmoke
	TextOutput
	StareDownEffect
)

var audioContext = audio.NewContext(44100)

func GameSFXLoader() *resource.Loader {

	gunShot, err := combatSFX.ReadFile("sounds/combatSFX/gunShot.wav")
	if err != nil {
		log.Fatal("cannot load gunShot.wav")
	}

	miss, err := combatSFX.ReadFile("sounds/combatSFX/whizz.wav")
	if err != nil {
		log.Fatal("cannot load whizz.wav")
	}

	hit, err := combatSFX.ReadFile("sounds/combatSFX/splat.wav")
	if err != nil {
		log.Fatal("cannot load hit.wav")
	}

	gunCock, err := combatSFX.ReadFile("sounds/combatSFX/cockingRevolver.wav")
	if err != nil {
		log.Fatal("cannot load cockingRevolver.wav")
	}
	maleEnemyGrunt, err := combatSFX.ReadFile("sounds/combatSFX/maleEnemyGrunt1.wav")
	if err != nil {
		log.Fatal("cannot load maleEnemyGrunt.wav")
	}

	elyseGrunt, err := combatSFX.ReadFile("sounds/combatSFX/elyseGrunt1.wav")
	if err != nil {
		log.Fatal("cannot load elyseGrunt1.wav")
	}

	pistolUnHolster, err := combatSFX.ReadFile("sounds/combatSFX/pistolUnholster.wav")
	if err != nil {
		log.Fatal("cannot load pistolunholster.wav")
	}

	drawButton, err := combatSFX.ReadFile("sounds/combatSFX/drawButtonClick.wav")
	if err != nil {
		log.Fatal("cannot load drawButtonClick.wav")
	}

	noAmmo, err := combatSFX.ReadFile("sounds/combatSFX/noAmmo.wav")
	if err != nil {
		log.Fatal("cannot load noAmmo.wav")
	}

	reload, err := combatSFX.ReadFile("sounds/combatSFX/reload.wav")
	if err != nil {
		log.Fatal("cannot load reload.wav")
	}

	victoryBlow, err := combatSFX.ReadFile("sounds/combatSFX/blowingSmoke.wav")
	if err != nil {
		log.Fatal("cannot load blowing.wav")
	}

	textSFX, err := menuSFX.ReadFile("sounds/menuSFX/rpg-text-speech-sound-131477.wav")
	if err != nil {
		log.Fatal("cannot load textSFX")
	}

	stareDownEffect, err := combatSFX.ReadFile("sounds/combatSFX/stareDownSoundEffect.wav")
	if err != nil {
		log.Fatal("cannot load stare down SFX")
	}

	var SoundData = map[string][]byte{
		"sounds/combatSFX/gunShot.wav":                    gunShot,
		"sounds/combatSFX/whizz.wav":                      miss,
		"sounds/combatSFX/splat.wav":                      hit,
		"sounds/combatSFX/cockingRevolver.wav":            gunCock,
		"sounds/combatSFX/maleEnemyGrunt1.wav":            maleEnemyGrunt,
		"sounds/combatSFX/elyseGrunt1.wav":                elyseGrunt,
		"sounds/combatSFX/pistolUnholster.wav":            pistolUnHolster,
		"sounds/combatSFX/drawButtonClick.wav":            drawButton,
		"sounds/combatSFX/noAmmo.wav":                     noAmmo,
		"sounds/combatSFX/reload.wav":                     reload,
		"sounds/combatSFX/blowingSmoke.wav":               victoryBlow,
		"sounds/menuSFX/rpg-text-speech-sound-131477.wav": textSFX,
		"sounds/combatSFX/stareDownSoundEffect.wav":       stareDownEffect,
	}

	l := resource.NewLoader(audioContext)

	l.OpenAssetFunc = func(path string) io.ReadCloser {
		return io.NopCloser(bytes.NewReader(SoundData[path]))
	}

	l.AudioRegistry.Assign(map[resource.AudioID]resource.AudioInfo{
		GunShot:         {Path: "sounds/combatSFX/gunShot.wav", Volume: 0.2},
		BulletHit:       {Path: "sounds/combatSFX/splat.wav", Volume: 0.2},
		BulletMiss:      {Path: "sounds/combatSFX/whizz.wav", Volume: 0.2},
		GunCock:         {Path: "sounds/combatSFX/cockingRevolver.wav", Volume: 0.2},
		ElyseGrunt:      {Path: "sounds/combatSFX/elyseGrunt1.wav", Volume: 0.2},
		EnemyGrunt:      {Path: "sounds/combatSFX/maleEnemyGrunt1.wav", Volume: 0.2},
		PistolUnHolster: {Path: "sounds/combatSFX/pistolUnholster.wav", Volume: 0.5},
		DrawButton:      {Path: "sounds/combatSFX/drawButtonClick.wav", Volume: -0.3},
		NoAmmo:          {Path: "sounds/combatSFX/noAmmo.wav", Volume: 0.2},
		Reload:          {Path: "sounds/combatSFX/reload.wav", Volume: -0.3},
		BlowingSmoke:    {Path: "sounds/combatSFX/blowingSmoke.wav", Volume: 0.2},
		TextOutput:      {Path: "sounds/menuSFX/rpg-text-speech-sound-131477.wav", Volume: 0.1},
		StareDownEffect: {Path: "sounds/combatSFX/stareDownSoundEffect.wav", Volume: 0.1},
	})

	return l

}

func MusicLoader() *resource.Loader {
	dialogueMusic, err := music.ReadFile("sounds/music/dialogueMusic.wav")
	if err != nil {
		log.Fatal("cannot load dialogue music")
	}

	battleMusic, err := music.ReadFile("sounds/music/battleMusic.wav")
	if err != nil {
		log.Fatal("cannot load battle music")
	}

	introMusic, err := music.ReadFile("sounds/music/dirty-outlaws-dust-and-echoes-244224.wav")
	if err != nil {
		log.Fatal("cannot load intro music")
	}

	var SoundData = map[string][]byte{
		"sounds/music/dirty-outlaws-dust-and-echoes-244224.wav": introMusic,
		"sounds/music/dialogueMusic.wav":                        dialogueMusic,
		"sounds/music/battleMusic.wav":                          battleMusic,
	}

	l := resource.NewLoader(audioContext)

	l.OpenAssetFunc = func(path string) io.ReadCloser {
		return io.NopCloser(bytes.NewReader(SoundData[path]))
	}

	l.AudioRegistry.Assign(map[resource.AudioID]resource.AudioInfo{
		IntroMusic:    {Path: "sounds/music/dirty-outlaws-dust-and-echoes-244224.wav", Volume: 0.2},
		DialogueMusic: {Path: "sounds/music/dialogueMusic.wav", Volume: 0.2},
		BattleMusic:   {Path: "sounds/music/battleMusic.wav", Volume: 0.2},
	})

	return l
}
