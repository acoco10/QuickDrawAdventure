package battleStats

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

type CharacterName uint8

const (
	Elyse CharacterName = iota
	Wolf
	CowardlyCowboy
	Sheriff
	GunSlinger
	None
)

type CharacterJSON struct {
	Name           string         `json:"name"`
	Stats          map[string]int `json:"stats"`
	CombatSkills   []string       `json:"combat_skills"`
	DialogueSkills []string       `json:"dialogue_skills"`
	Weakness       string         `json:"weakness"`
}

type CharactersJSON struct {
	Characters []CharacterJSON `json:"characters"`
}

func LoadCharacter(characterJSON CharacterJSON) (Character, error) {

	// getting the string keys from all skills to confirm if character json contains valid skills
	combatSkills, dialogueSkills, err := LoadSkills()

	if err != nil {
		log.Fatal("Error loading skills", err)
	}

	if len(dialogueSkills) == 0 {
		return Character{}, errors.New("no dialogue skills found")
	}
	if len(characterJSON.CombatSkills) == 0 {
		return Character{}, errors.New("no combat skills found")
	}

	dialogueSkillKeys := make([]string, len(dialogueSkills))

	combatSkillKeys := make([]string, len(combatSkills))

	i := 0

	for skillName := range dialogueSkills {
		dialogueSkillKeys[i] = skillName
		i++
	}

	i = 0

	for skillName := range combatSkills {
		combatSkillKeys[i] = skillName
		i++
	}

	characterCombatSkills := make(map[string]Skill)

	for _, skill := range characterJSON.CombatSkills {

		// if skill is invalid pass an error and print the invalid skill

		if !StringInSlice(skill, combatSkillKeys) {
			log.Fatalf("character contains invalid skill: %s %s", skill, combatSkillKeys)
		}

		characterCombatSkills[skill] = combatSkills[skill]
	}

	characterDialogueSkills := make(map[string]Skill)

	for _, skill := range characterJSON.DialogueSkills {

		if !StringInSlice(skill, dialogueSkillKeys) {
			log.Fatalf("character contains invalid skill: %s", skill)
		}

		characterDialogueSkills[skill] = dialogueSkills[skill]
	}

	character := NewCharacter(characterJSON.Name, characterJSON.Stats, characterCombatSkills, characterDialogueSkills, characterJSON.Weakness)

	return character, nil

}

func LoadCharacters() (map[CharacterName]Character, error) {

	contents, err := os.ReadFile("battleStats/data/characters.json")

	if err != nil {
		log.Fatal(err)
	}

	var charactersJSON CharactersJSON

	err = json.Unmarshal(contents, &charactersJSON)

	if err != nil {
		log.Fatal(contents, err)
	}

	characters := make(map[CharacterName]Character)

	for _, characterJSON := range charactersJSON.Characters {
		fmt.Printf("loading character %s\n", characterJSON.Name)
		char, err := LoadCharacter(characterJSON)
		if err != nil {
			log.Fatal(char, err)
		}
		if char.Name == "Elyse" {
			characters[Elyse] = char
		}
		if char.Name == "Wolf" {
			characters[Wolf] = char
		}
		if char.Name == "Sheriff" {
			characters[Sheriff] = char
		}
		if char.Name == "Gunslinger" {
			characters[GunSlinger] = char
		}

	}

	return characters, nil
}
