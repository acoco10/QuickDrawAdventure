package battle

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/dataManagement"
	"github.com/tidwall/gjson"
	"log"
	"math/rand/v2"
	"os"
)

func Shoot(aC int, a int, f int, ChanceToHit int, damageRange []int) (damage int) {

	if (aC*2)+ChanceToHit-a-f > rand.IntN(100) {
		if len(damageRange) > 1 {
			return damageRange[0] + rand.IntN(damageRange[1])
		} else {
			return damageRange[0]
		}
	}
	return 0
}

func Roll(successPer int) bool {
	return rand.IntN(100) < successPer
}

func EnemyChooseSkill(battle Battle, enemySkills map[string]dataManagement.Skill) (skill dataManagement.Skill, err error) {

	if battle.battlePhase == Dialogue && battle.WinningProb < 40 {
		return enemySkills["draw"], nil
	}

	if battle.EnemyAmmo == 0 {
		println("enemy chose skill:reload")
		return enemySkills["reload"], nil
	}

	randSkillInt := rand.IntN(len(enemySkills)) //reload is the last index and randn is exclusive

	skillIndexes := make([]int, 0)
	for _, skillOption := range enemySkills {
		skillIndexes = append(skillIndexes, skillOption.Index)
	}

	i := 0
	chosenIndex := 0
	for _, skillIndex := range skillIndexes {
		if i == randSkillInt {
			chosenIndex = skillIndex
		}
		i++
	}

	for _, eSkill := range enemySkills {
		if eSkill.Index == chosenIndex {
			skill = eSkill
		}
	}

	return skill, nil
}

func Draw(userStats map[dataManagement.Stat]int, oppStats map[dataManagement.Stat]int) bool {

	initiative := true

	userDS := userStats[dataManagement.DrawSpeed] - userStats[dataManagement.Anger]*10 - userStats[dataManagement.Fear]*10
	opponentDS := oppStats[dataManagement.DrawSpeed] - oppStats[dataManagement.Anger]*10 - oppStats[dataManagement.Fear]*10

	if rand.IntN(101) > 50+userDS-opponentDS {
		initiative = false
	}
	return initiative
}

func DrawProb(userStats map[dataManagement.Stat]int, oppStats map[dataManagement.Stat]int) int {
	userDS := userStats[dataManagement.DrawSpeed] - userStats[dataManagement.Anger]*10 - userStats[dataManagement.Fear]*10
	opponentDS := oppStats[dataManagement.DrawSpeed] - oppStats[dataManagement.Anger]*10 - oppStats[dataManagement.Fear]*10
	return 50 + userDS - opponentDS
}

func GetSkillDialogue(charName string, skillName string, status bool) string {
	data, err := os.ReadFile("battle/dialogue.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(data)
	var response string
	if status {
		query := fmt.Sprintf("characters.#(name==%s).skills.#(type == %s).successful", charName, skillName)
		results := gjson.Get(jsonString, query)
		response = results.String()
	} else {
		query := fmt.Sprintf("characters.#(name==%s).skills.#(type == %s).unsuccessful", charName, skillName)
		results := gjson.Get(jsonString, query)
		response = results.String()
	}
	return response
}

func GetResponse(charName string, skillName string, status bool) string {
	data, err := os.ReadFile("battle/dialogue.json")
	if err != nil {
		log.Fatal(err)
	}
	jsonString := string(data)
	var response string
	if status {
		query := fmt.Sprintf("characters.#(name==%s).responses.#(type == %s).successful", charName, skillName)
		results := gjson.Get(jsonString, query)
		response = results.String()
	} else {
		query := fmt.Sprintf("characters.#(name==%s).responses.#(type == %s).unsuccessful", charName, skillName)
		results := gjson.Get(jsonString, query)
		response = results.String()
	}
	return response
}
