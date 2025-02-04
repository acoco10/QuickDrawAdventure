package battle

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/tidwall/gjson"
	"log"
	"math"
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

func EnemyChooseSkill(battle Battle, enemySkills map[string]battleStats.Skill) (skill battleStats.Skill, err error) {
	if battle.BattlePhase == Dialogue {
		drawChance := 0
		if battle.Tension > battle.Enemy.Stats[battleStats.TensionThreshold] {
			drawChance = int(math.Pow(float64(battle.Tension), 1.9))
		}

		if rand.IntN(100)+1 < drawChance {

			return enemySkills["draw"], nil

		}
	}
	if len(enemySkills) == 1 {
		return enemySkills["bite"], nil
	}

	if battle.EnemyAmmo == 0 {
		return enemySkills["reload"], nil
	}

	randSkillInt := rand.IntN(len(enemySkills) - 1) //reload/draw are the last index and randn is exclusive

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
	fmt.Printf("Choosing skill for enemy:%s\n", skill.SkillName)
	if skill.SkillName == "reload" {
		skill = enemySkills["focusedShot"]
	}

	return skill, nil
}

func Draw(userStats map[battleStats.Stat]int, oppStats map[battleStats.Stat]int) bool {

	initiative := true

	userDS := (float64(userStats[battleStats.DrawSpeed]) - float64(userStats[battleStats.Anger])/2 - float64(userStats[battleStats.Fear])/2) * 10
	opponentDS := (float64(oppStats[battleStats.DrawSpeed]) - float64(oppStats[battleStats.Anger])/2 - float64(oppStats[battleStats.Fear])/2) * 10

	fmt.Printf("enemyDS: %d, enemy fear: %d, enemy anger: %d", oppStats[battleStats.DrawSpeed], oppStats[battleStats.Anger], oppStats[battleStats.Fear])

	if float64(rand.IntN(101)) > 50+userDS-opponentDS {
		initiative = false
	}
	return initiative
}

func DrawProb(userStats map[battleStats.Stat]int, oppStats map[battleStats.Stat]int) int {

	userDS := userStats[battleStats.DrawSpeed] - userStats[battleStats.Anger]*5 - userStats[battleStats.Fear]*5
	opponentDS := oppStats[battleStats.DrawSpeed] - oppStats[battleStats.Anger]*5 - oppStats[battleStats.Fear]*5

	fmt.Printf("enemyDS: %d, enemy fear: %d, enemy anger: %d\n", oppStats[battleStats.DrawSpeed], oppStats[battleStats.Anger], oppStats[battleStats.Fear])
	fmt.Printf("playerDS: %d, player fear: %d, player anger: %d\n", userStats[battleStats.DrawSpeed], userStats[battleStats.Anger], userStats[battleStats.Fear])

	println("player draw speed(user) =", userDS)
	println("enemy draw  speed(opponent) =", opponentDS)

	return 50 + userDS - opponentDS
}

func GetSkillDialogue(charName string, skillName string, status bool) string {
	data, err := os.ReadFile("battle/battleDialogue.json")
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
	data, err := os.ReadFile("battle/battleDialogue.json")
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
