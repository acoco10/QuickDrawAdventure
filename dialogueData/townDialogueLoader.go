package dialogueData

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"log"
	"os"
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

func LoadNpcTownDialogue() TownDialogue {
	contents, err := os.ReadFile("dialogueData/townDialogue.json")
	if err != nil {
		log.Fatal(err)
	}

	var townDialogue TownDialogue

	err = json.Unmarshal(contents, &townDialogue)

	if err != nil {
		log.Fatal(contents, err)
	}

	return townDialogue
}

type DialogueTracker struct {
	CharName string
	Index    int
}

func GetResponse(charName string, dialogueId int) string {

	data, err := os.ReadFile("dialogueData/townDialogue.json")
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
	data, err := os.ReadFile("dialogueData/playerTownDialogue.json")
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
