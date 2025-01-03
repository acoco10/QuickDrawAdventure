package audioManagement

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
)

type SongPlayerStatus uint8

const (
	Playing SongPlayerStatus = iota
	Paused
	NotPlaying
)

type Song struct {
	songID resource.AudioID
	length int
}

type DJ struct {
	mixCounter    int
	mixPlayer     *audio.Player
	currentSongID resource.AudioID
	mixSongID     resource.AudioID
	status        SongPlayerStatus
	loader        *resource.Loader
	Volume        float64
	player        *audio.Player
	mixVolume     float64
}

func (a *DJ) play() {
	if a.currentSongID != 0 {
		if a.player.IsPlaying() == false {
			a.player = a.loader.LoadWAV(a.currentSongID).Player
			a.player.SetVolume(a.Volume)
			err := a.player.Rewind()
			if err != nil {
				log.Printf("%q Rewind: %s", a.currentSongID, err)
			}
			a.player.Play()
		}
	}
}

func (a *DJ) QueueSong(songID resource.AudioID) {
	log.Printf("Queuing song %q", songID)
	a.currentSongID = songID
}

func (a *DJ) Update() {
	if a.currentSongID != 0 {
		a.play()
	}

	if a.mixCounter > 0 {
		a.mixCounter--
	}

	if a.mixCounter == 1 {
		if a.mixVolume <= 0.5 {
			a.mixVolume += 0.1
			a.Volume -= 0.1
			a.mixPlayer.SetVolume(a.mixVolume)
			a.mixCounter = 50
		} else {
			println("switching main player to mix player\n")
			a.Stop()
			a.Volume = 0.5
			a.player = a.mixPlayer
			a.currentSongID = a.mixSongID
		}
	}

}

func (a *DJ) Stop() {
	a.player.Pause()
}

func (a *DJ) Mix(songID resource.AudioID) {
	a.mixCounter = 20
	a.mixSongID = songID
	a.mixPlayer = a.loader.LoadWAV(songID).Player
	a.mixVolume = 0
	a.mixPlayer.SetVolume(a.mixVolume)
	a.mixPlayer.Play()
}

func NewSongPlayer(firsSong resource.AudioID) *DJ {
	s := DJ{
		loader:        MusicLoader(),
		Volume:        0.5,
		currentSongID: firsSong,
	}
	s.player = s.loader.LoadWAV(firsSong).Player

	return &s
}
