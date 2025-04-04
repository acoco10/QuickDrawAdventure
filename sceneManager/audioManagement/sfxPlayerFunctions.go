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
		a.Play(a.sfxQueue[a.sfxQueueIndex])
		if a.sfxQueueIndex < len(a.sfxQueue)-1 {
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
	log.Printf("Playing audio %q", id)
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

func (a *SFXAudioPlayer) ConfigureAttackResultSoundQueue(damage []int, target string) {
	var soundList []resource.AudioID
	for _, result := range damage {
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
		if result < 0 {
			soundList = append(soundList, NoAmmo)
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
