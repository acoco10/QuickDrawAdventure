package dialogueData

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/tidwall/gjson"
	"log"
)

type NpcTownDialogue struct {
	Name     string             `json:"name"`
	Dialogue []TownDialogueData `json:"dialogue"`
}

type TownDialogue struct {
	DialogueMap []NpcTownDialogue `json:"characters"`
}

type TownDialogueData struct {
	DialogueID   int    `json:"DialogueID"`
	DialogueText string `json:"DialogueText"`
}

type DialogueTracker struct {
	CharName string
	Index    int
}

func TalkFirst(charName string, storyPointId int) bool {

	data, err := assets.Dialogue.ReadFile("dialogueData/playerTownDialogue.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(data)
	var response string

	query := fmt.Sprintf("storyPoints.#(storyPointId==%d).playerDialogue.#(name==%s).talkFirst", storyPointId, charName)
	results := gjson.Get(jsonString, query)
	response = results.String()
	if response == "true" {
		println("talkfirst function output:", response)
		return true
	}
	println("talkfirst function output: false")
	return false
}

func GetResponse(charName string, dialogueId int) string {

	data, err := assets.Dialogue.ReadFile("dialogueData/townDialogue.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(data)
	var response string

	query := fmt.Sprintf("characters.#(name==%s).dialogues.#(storyPoint == %d).dialogue.#(dialogueId == %d).dialogueText", charName, 1, dialogueId)
	results := gjson.Get(jsonString, query)
	response = results.String()

	return response
}

func GetPlayerResponse(charName string, storyPointId int, dialogueId int) string {
	data, err := assets.Dialogue.ReadFile("dialogueData/playerTownDialogue.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(data)
	var response string

	query := fmt.Sprintf("storyPoints.#(storyPointId==%d).playerDialogue.#(name==%s).dialogue.#(dialogueId == %d).dialogueText", storyPointId, charName, dialogueId)
	results := gjson.Get(jsonString, query)
	response = results.String()

	return response
}
