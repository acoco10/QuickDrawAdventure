package gameScenes

import events "github.com/acoco10/QuickDrawAdventure/eventSystem"

type DialogueStateManager struct {
	*events.EventHub
	characterStoryProgression map[string]*CharacterState
}

type CharacterState struct {
	StoryPoint int
}

func NewDialogueStateManager(hub *events.EventHub) {
	DSM := DialogueStateManager{}
	DSM.EventHub = hub

}

func (d *DialogueStateManager) UpdateDialogueTracker(name string) {
	d.characterStoryProgression[name].StoryPoint++
}

func (d *DialogueStateManager) init() {
	d.EventHub = events.NewEventHub()
	d.Subscribe(events.DialogueEvent{}, func(e events.Event) {
		ev := e.(events.DialogueEvent)
		for _, char := range ev.Characters {
			d.UpdateDialogueTracker(char)
		}
	})
}
