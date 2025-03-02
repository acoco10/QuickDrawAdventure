package battleStats

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"log"
)

type CharacterName uint8

const (
	Elyse CharacterName = iota
	Wolf
	CowardlyCowboy
	Sheriff
	Antonio
	None
)

type CharacterJSON struct {
	Name           string         `json:"name"`
	Stats          map[string]int `json:"stats"`
	CombatSkills   []string       `json:"combat_skills"`
	DialogueSkills []string       `json:"dialogue_skills"`
	DialogueSlots  int            `json:"dialogue_slots"`
	SoundFxType    string         `json:"soundFxType"`
	Weakness       string         `json:"weakness"`
}

type CharactersJSON struct {
	Characters []CharacterJSON `json:"characters"`
}

func LoadSingleCharacter(charName string) (CharacterData, error) {
	var char CharacterData
	contents, err := assets.Battle.ReadFile("battleData/characters.json")

	if err != nil {
		log.Fatal(err)
	}

	var charactersJSON CharactersJSON

	err = json.Unmarshal(contents, &charactersJSON)

	if err != nil {
		log.Fatal(contents, err)
	}

	for _, characterJSON := range charactersJSON.Characters {
		if characterJSON.Name == charName {
			char, err = LoadCharacter(characterJSON)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return char, nil
}

func LoadCharacter(characterJSON CharacterJSON) (CharacterData, error) {

	// getting the string keys from all skills to confirm if character json contains valid skills
	combatSkills, dialogueSkills, equipSkills, err := LoadSkills()

	if err != nil {
		log.Fatal("Error loading skills", err)
	}

	if len(equipSkills) == 0 {
		return CharacterData{}, errors.New("no dialogue skills found")
	}
	if len(characterJSON.CombatSkills) == 0 {
		return CharacterData{}, errors.New("no combat skills found")
	}

	dialogueSkillKeys := make([]string, len(dialogueSkills))

	equipSkillKeys := make([]string, len(equipSkills))

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

	i = 0

	for skillName := range equipSkills {
		equipSkillKeys[i] = skillName
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

		if StringInSlice(skill, dialogueSkillKeys) {
			characterDialogueSkills[skill] = dialogueSkills[skill]
		} else if StringInSlice(skill, equipSkillKeys) {
			characterDialogueSkills[skill] = equipSkills[skill]
		} else {
			log.Fatalf("character contains invalid dialogue skill: %s", skill)
		}

	}

	character := NewCharacter(characterJSON.Name, characterJSON.Stats, characterCombatSkills, characterDialogueSkills, characterJSON.Weakness, characterJSON.SoundFxType, characterJSON.DialogueSlots)

	return character, nil

}

func LoadCharacters() (map[CharacterName]CharacterData, error) {

	contents, err := assets.Battle.ReadFile("battleData/characters.json")

	if err != nil {
		log.Fatal(err)
	}

	var charactersJSON CharactersJSON

	err = json.Unmarshal(contents, &charactersJSON)

	if err != nil {
		log.Fatal(contents, err)
	}

	characters := make(map[CharacterName]CharacterData)

	for _, characterJSON := range charactersJSON.Characters {
		fmt.Printf("loading character %s\n", characterJSON.Name)
		char, err := LoadCharacter(characterJSON)
		if err != nil {
			log.Fatal(char, err)
		}
		if char.Name == "elyse" {
			characters[Elyse] = char
		}
		if char.Name == "wolf" {
			characters[Wolf] = char
		}
		if char.Name == "sheriff" {
			characters[Sheriff] = char
		}
		if char.Name == "antonio" {
			characters[Antonio] = char
		}

	}

	return characters, nil
}
