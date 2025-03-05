package battle

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/tidwall/gjson"
	"math"
	"math/rand/v2"
)

func Shoot(aC int, a int, f int, ChanceToHit int, damageRange []int) (damage int) {

	if (aC*2)+ChanceToHit-a-f > rand.IntN(100) {
		if len(damageRange) > 1 {
			return damageRange[0] + rand.IntN(damageRange[1])
		} else if len(damageRange) > 0 {
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
		if battle.Tension >= battle.CharacterBattleData[Enemy].Stats[battleStats.TensionThreshold] {
			drawChance = int(math.Pow(float64(battle.Tension), 1.9))
			if drawChance > 100 {
				drawChance = 95
			}
			println("drawChance =", drawChance)
		}

		if rand.IntN(101) <= drawChance {
			println("Choosing skill for enemy:draw", "drawChance =", drawChance)
			return enemySkills["draw"], nil
		}
	}

	if len(enemySkills) == 1 {
		return enemySkills["bite"], nil
	}

	if battle.CharacterBattleData[Enemy].Ammo == 0 {
		return enemySkills["reload"], nil
	}

	skill = SkillRandomizer(enemySkills)

	fmt.Printf("Choosing skill for enemy:%s\n", skill.SkillName)

	return skill, nil
}

func SkillRandomizer(skills map[string]battleStats.Skill) (skill battleStats.Skill) {

	var skillKeys []string

	for key, skillop := range skills {
		if skillop.SkillName != "reload" && skillop.SkillName != "draw" {
			skillKeys = append(skillKeys, key)
		}
	}

	randSkillInt := rand.IntN(len(skillKeys))

	return skills[skillKeys[randSkillInt]]
}

func ReadyDraw(userStats map[battleStats.Stat]int, oppStats map[battleStats.Stat]int) bool {

	initiative := true

	userDS := float64(userStats[battleStats.DrawSpeed]) - float64(userStats[battleStats.Anger])*5 - float64(userStats[battleStats.Fear])*5
	opponentDS := float64(oppStats[battleStats.DrawSpeed]) - float64(oppStats[battleStats.Anger])*5 - float64(oppStats[battleStats.Fear])*5

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

func GetSkillDialogue(char battleStats.CharacterData, skillName string, status bool, data string) string {
	charName := char.Name
	var response string
	if status {
		if char.Name == "elyse" {
			selInt := rand.IntN(3)
			if skillName == "insult" {
				query := fmt.Sprintf("insults.#(id = %d).insult", selInt)
				results := gjson.Get(char.DialogueData, query)
				response = results.String()
			}
			if skillName == "brag" {
				query := fmt.Sprintf("brags.#(id = %d).brag", selInt)
				results := gjson.Get(char.DialogueData, query)
				response = results.String()
			}
		} else {
			query := fmt.Sprintf("characters.#(name==%s).skills.#(type == %s).successful", charName, skillName)
			results := gjson.Get(data, query)
			response = results.String()
		}
	} else {
		query := fmt.Sprintf("characters.#(name==%s).skills.#(type == %s).unsuccessful", charName, skillName)
		results := gjson.Get(data, query)
		response = results.String()
	}
	return response
}

func GetResponse(char battleStats.CharacterData, opponent battleStats.CharacterData, skillName string, status bool, data string) string {
	var response string
	oppCharName := opponent.Name
	if status {
		query := fmt.Sprintf("characters.#(name==%s).responses.#(type == %s).successful", oppCharName, skillName)
		results := gjson.Get(data, query)
		response = results.String()
	} else {
		query := fmt.Sprintf("characters.#(name==%s).responses.#(type == %s).unsuccessful", oppCharName, skillName)
		results := gjson.Get(data, query)
		response = results.String()
	}
	return response
}
