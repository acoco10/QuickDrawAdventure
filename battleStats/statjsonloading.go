package battleStats

import (
	"encoding/json"
	"log"
)

type Effect struct {
	EffectType  string `json:"type"`
	Stat        string `json:"stat"`
	Amount      int    `json:"amount"`
	Duration    int    `json:"duration"`
	InitiativeL bool   `json:"initiativeL"`
	DamageRange []int  `json:"damage"`
	SuccessPer  int    `json:"successPer"`
	NShots      int    `json:"nShots"`
	On          string `json:"on"`
}

type SkillJson struct {
	Skills []Skill `json:"skills"`
}

type Skill struct {
	SkillName string   `json:"name"`
	Index     int      `json:"index"`
	Effects   []Effect `json:"effects"`
	Tension   int      `json:"tension"`
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func EffectsTest(effect Effect, skillName string) bool {

	if effect.EffectType == "" {
		log.Printf(`Skill:%s has no parameter for EffectType`, skillName)
		return false
	}

	//checking for valid attack type parameters
	if effect.EffectType == "attack" {
		if effect.DamageRange == nil {
			log.Printf(`Skill:%s is an attack but has no damage range`, skillName)
			return false
		}
		if effect.SuccessPer == 0 {
			log.Printf(`Skill:%s is an attack but has no success percentage`, skillName)
			return false
		}
		if effect.NShots == 0 {
			log.Printf(`Skill:%s is an attack but does not have number of shots`, skillName)
			return false
		}
	}

	//checking for valid buff type parameters
	if effect.EffectType == "buff" {
		if effect.Amount == 0 {
			log.Printf(`Skill:%s is a buff but does not include an effect amount`, skillName)
			return false
		}
	}

	return true
}

func LoadSkillsFromPath(fileName string) (map[string]Skill, error) {

	skillMap := make(map[string]Skill)

	contents, err := battleStatsData.ReadFile(fileName)

	if err != nil {

		return skillMap, err
	}

	var skillsJSON SkillJson

	err = json.Unmarshal(contents, &skillsJSON)

	if err != nil {
		return skillMap, err
	}

	for _, skill := range skillsJSON.Skills {
		for _, effect := range skill.Effects {
			if EffectsTest(effect, skill.SkillName) {
				continue
			}

		}
		if skill.SkillName == "" {
			log.Printf(`Skill:%s has no parameter for skillname`, skill.SkillName)
			continue
		}

		skillMap[skill.SkillName] = skill
	}

	return skillMap, nil
}

func LoadSkills() (combatSkills map[string]Skill, dialogueSkills map[string]Skill, err error) {
	combatSkills, err = LoadSkillsFromPath("data/combatSkills.json")
	if err != nil {
		return combatSkills, dialogueSkills, err
	}
	dialogueSkills, err = LoadSkillsFromPath("data/dialogueSkills.json")
	if err != nil {
		return combatSkills, dialogueSkills, err
	}

	return combatSkills, dialogueSkills, nil

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
