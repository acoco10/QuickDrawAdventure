package dialogueData

import (
	"fmt"
	"github.com/tidwall/gjson"
)

type DialogueTracker struct {
	CharName string
	Index    int
}

func TalkFirst(charName string, storyPointId int, data string) bool {
	query := fmt.Sprintf("#(name==%s).dialogues.#(storyPoint == %d).talkFirst", charName, storyPointId)
	results := gjson.Get(data, query)
	response := results.String()
	if response == "true" {
		println("player should talk first")
		return true
	}
	return false
}

func GetResponse(charName string, dialogueId int, data string) string {

	query := fmt.Sprintf("#(name==%s).dialogues.#(storyPoint == %d).dialogue.#(dialogueId == %d).dialogueText", charName, 1, dialogueId)
	results := gjson.Get(data, query)
	response := results.String()

	return response
}

func GetPlayerResponse(charName string, storyPointId int, dialogueId int, data string) string {

	query := fmt.Sprintf("storyPoints.#(storyPointId==%d).playerDialogue.#(name==%s).dialogue.#(dialogueId == %d).dialogueText", storyPointId, charName, dialogueId)
	results := gjson.Get(data, query)
	response := results.String()

	return response
}
