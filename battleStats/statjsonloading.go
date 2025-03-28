package battleStats

import (
	"encoding/json"
	"github.com/acoco10/QuickDrawAdventure/assets"
)

type Effect struct {
	EffectType  string `json:"type"`
	Stat        string `json:"stat"`
	Amount      int    `json:"amount"`
	Duration    int    `json:"duration"`
	DamageRange []int  `json:"damage"`
	SuccessPer  int    `json:"successPer"`
	NShots      int    `json:"nShots"`
	On          string `json:"on"`
}

type SkillJson struct {
	Skills []Skill `json:"skills"`
}

type Skill struct {
	Index     int      `json:"index"`
	SkillName string   `json:"name"`
	Text      string   `json:"text"`
	Type      string   `json:"type"`
	Target    string   `json:"target"`
	Tension   int      `json:"tension"`
	Effects   []Effect `json:"effects"`
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func LoadSkillsFromPath(fileName string) (map[string]Skill, error) {

	skillMap := make(map[string]Skill)

	contents, err := assets.Battle.ReadFile(fileName)

	if err != nil {

		return skillMap, err
	}

	var skillsJSON SkillJson

	err = json.Unmarshal(contents, &skillsJSON)

	if err != nil {
		return skillMap, err
	}

	for _, skill := range skillsJSON.Skills {
		skillMap[skill.SkillName] = skill
	}

	return skillMap, nil
}

func LoadSkills() (combatSkills map[string]Skill, dialogueSkills map[string]Skill, equipSkills map[string]Skill, err error) {
	combatSkills, err = LoadSkillsFromPath("battleData/combatSkills.json")
	if err != nil {
		return combatSkills, dialogueSkills, equipSkills, err
	}
	dialogueSkills, err = LoadSkillsFromPath("battleData/dialogueSkills.json")
	if err != nil {
		return combatSkills, dialogueSkills, equipSkills, err
	}
	equipabSkills, err := LoadSkillsFromPath("battleData/equipDialogueSkills.json")
	if err != nil {
		return combatSkills, dialogueSkills, equipSkills, err
	}

	return combatSkills, dialogueSkills, equipabSkills, nil

}

type ResponseText struct {
	Successful   string `json:"successful"`
	Unsuccessful string `json:"unsuccessful"`
}

type DialogueSkills struct {
	Insult     ResponseText `json:"insult"`
	Brag       ResponseText `json:"brag"`
	Intimidate ResponseText `json:"intimidate"`
}

type DialogueResponses struct {
	Insult     ResponseText `json:"insult"`
	Brag       ResponseText `json:"brag"`
	Intimidate ResponseText `json:"intimidate"`
}

type Response struct {
	Skills    DialogueSkills    `json:"skills"`
	Responses DialogueResponses `json:"responses"`
}

type DialogueJSON struct {
	DialogueMessages []CharacterDialogueData `json:"characters"`
}

type CharacterDialogueData struct {
	CharacterName    string     `json:"name"`
	DialogueMessages []Response `json:"response"`
}
