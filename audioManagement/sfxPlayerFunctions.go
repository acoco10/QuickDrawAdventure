package audioManagement

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
)

type SFXAudioPlayer struct {
	loader        *resource.Loader
	countDown     int
	sfxQueue      []resource.AudioID
	sfxQueueIndex int
	sfxPlayer     *audio.Player
}

func (a *SFXAudioPlayer) Update() {

	if a.countDown > 0 {
		a.countDown--
	}

	if a.countDown == 1 && a.sfxQueue != nil {
		if a.sfxQueueIndex <= len(a.sfxQueue)-1 {
			a.Play(a.sfxQueue[a.sfxQueueIndex])
			a.countDown = 12
			a.sfxQueueIndex++
		} else {
			a.countDown = 0
			a.sfxQueueIndex = 0
			a.sfxQueue = nil
		}
	}
}

func (a *SFXAudioPlayer) Play(id resource.AudioID) {
	sfxPlayer := a.loader.LoadWAV(id).Player
	a.sfxPlayer = sfxPlayer
	err := sfxPlayer.Rewind()
	if err != nil {
		log.Printf("%q Rewind: %s", id, err)
	}
	sfxPlayer.Play()
}

func (a *SFXAudioPlayer) QueueSound(id resource.AudioID) {
	a.sfxQueue[a.sfxQueueIndex] = id
}

func (a *SFXAudioPlayer) ConfigureAttackResultSoundQueue(damage []int, target string, attacker string) {
	log.Printf("Formatting attack sound queue")
	var soundList []resource.AudioID
	for _, result := range damage {
		if attacker == "wolf" {
			soundList = append(soundList, WolfBite)
			if result > 0 {
				soundList = append(soundList, ElyseGrunt)

			}
		} else {
			if result > 0 {
				soundList = append(soundList, GunShot)
				soundList = append(soundList, BulletHit)
				if target == "Player" {
					soundList = append(soundList, ElyseGrunt)
				} else {
					soundList = append(soundList, EnemyGrunt)
				}
				soundList = append(soundList, GunCock)
			}
			if result == 0 {
				soundList = append(soundList, GunShot)
				soundList = append(soundList, BulletMiss)
				soundList = append(soundList, GunCock)
			}
			if result == -1 {
				soundList = append(soundList, NoAmmo)
			}
		}
	}
	a.countDown = 12
	a.sfxQueue = soundList
}

func (a *SFXAudioPlayer) ConfigureSoundQueue(soundList []resource.AudioID) {
	a.sfxQueue = soundList
	a.countDown = 2
}

func NewAudioPlayer() *SFXAudioPlayer {
	return &SFXAudioPlayer{
		loader:        GameSFXLoader(),
		sfxQueueIndex: 0,
	}
}

func (a *SFXAudioPlayer) Stop() {
	a.sfxPlayer.Pause()
}
