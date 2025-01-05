package main

import (
	"fmt"
	"github.com/acoco10/qdabattlesystem/battle"
	"github.com/acoco10/qdabattlesystem/dataManagement"
	"math/rand"
	"testing"
)

func TestShoot(t *testing.T) {

	damageRange := []int{2, 0}
	var results float64

	for i := 0; i < 10000; i++ {
		damage := battle.Shoot(10, 0, 0, 70, damageRange)
		results += float64(damage)
	}

}

func TestLoadBadSkillJSON(t *testing.T) {
	badSkills, err := dataManagement.LoadSkillsFromPath("badSkills.json")

	if err != nil {
		t.Fatalf(`LoadSkillsFromPath("combatSkills.json") did not load due to %s`, err)
	}

	nBadSkillsLoaded := len(badSkills)
	if nBadSkillsLoaded > 0 {
		t.Fatalf(`LoadSkillsFromPath() loaded %d skills from badSkills.json that it should not have`, nBadSkillsLoaded)
	}
}

func TestRoll(t *testing.T) {
	successPer := 80
	i := 0
	success := 0
	fail := 0
	for i < 10000 {
		if battle.Roll(successPer) {
			success++
		} else {
			fail++
		}
		i++
	}
	fmt.Printf("successful rolls:%d, failed rolls: %d", success, fail)
}

func TestLoadGoodSkillsJSON(t *testing.T) {
	combatSkills, _, _ := dataManagement.LoadSkills()
	dSkillsLength := len(combatSkills)
	dialogueSkillNames := make([]string, dSkillsLength)
	for _, skill := range combatSkills {
		i := skill.Index
		dialogueSkillNames[i] = skill.SkillName
	}

	for _, dialogueSkillName := range dialogueSkillNames {
		println(dialogueSkillName)
	}
}

func TestHowRandWorks(t *testing.T) {
	randNums := make([]int, 10)
	for i := 0; i < 10; i++ {
		randInt := rand.Intn(2)
		randNums = append(randNums, randInt)
	}

	for _, randNum := range randNums {
		println(randNum)
	}

}

func TestDialogueLoading(t *testing.T) {
	charName := "elyse"
	skillType := "insult"
	println(len(battle.GetSkillDialogue(charName, skillType, false)))
}

func TestLoadCharacter(t *testing.T) {
	chars, _ := dataManagement.LoadCharacters()
	for skill := range chars[0].DialogueSkills {
		println(skill)
	}
	for skill := range chars[0].DialogueSkills {
		println(skill)
	}
}

func TestCharacterMethods(t *testing.T) {

	stats := map[string]int{
		"health":    4,
		"accuracy":  2,
		"drawSpeed": 3,
		"anger":     0,
		"fear":      0,
	}

	elyse := dataManagement.NewCharacter("elyse", stats, map[string]dataManagement.Skill{}, map[string]dataManagement.Skill{})

	if elyse.DisplayStat(dataManagement.Health) != 4 {
		t.Fatalf(`method displayCharHealth did not work`)
	}

	elyse.UpdateCharAccuracy(-1)

	if elyse.DisplayStat(dataManagement.Accuracy) != 1 {
		t.Fatalf(`method updateCharAccuracy did not work`)
	}

	elyse.UpdateCharHealth(-1)

	if elyse.DisplayStat(dataManagement.Health) != 3 {
		t.Fatalf(`method updateCharhealthdid not work health value:%d expected value 3`, elyse.DisplayStat(dataManagement.Health))
	}

	elyse.UpdateCharHealth(-0)
	if elyse.DisplayStat(dataManagement.Health) != 4 {
		t.Fatalf(`method updateCharhealth did not work value inserted %d above its maximum: 4 `, elyse.DisplayStat(dataManagement.Health))
	}

}

func Test_use_stat_Buff(t *testing.T) {
	stats := map[string]int{
		"health":    4,
		"accuracy":  2,
		"drawspeed": 3,
		"anger":     0,
		"fear":      0,
	}

	elyse := dataManagement.NewCharacter("elyse", stats, map[string]dataManagement.Skill{}, map[string]dataManagement.Skill{})

	elyse.UpdateStat(dataManagement.Anger, 1)

	if elyse.DisplayStat(dataManagement.Anger) != 1 {
		t.Fatalf(`method ResetStatusStats() failed to update correctly %d`, elyse.Stats[dataManagement.Anger])
	}

	elyse.ResetStatusStats()

	if elyse.DisplayStat(dataManagement.Anger) != 0 {
		t.Fatalf(`method ResetStatusStats() failed to update correctly %d`, elyse.Stats[dataManagement.Anger])
	}

}

/*func Test_map_stringPrinter(t *testing.T) {
	testString1 := "the brown dog jumped over the blue fence rapidly"
	funcResult := textWrapper(testString1, 75)
	expectedResult := "the brown dog jumped over the blue fence rapidly"
	if funcResult != expectedResult {
		t.Fatalf("UpdateLines function incorrect output expected:%s, got:%s", expectedResult, funcResult)
	}
	testString2 := "y heard o me"
	funcResult = textWrapper(testString2, 10)
	expectedResult = "y heard o"
	if funcResult != expectedResult {
		t.Fatalf("UpdateLines function incorrect output expected:%s, got:%s", expectedResult, funcResult)
	}
}*/

/*func TestLineWrapper(t *testing.T) {

	texP := textPrinter{
		TextInput:         []string{"the brown dog jumped over the blue fence rapidly"},
		Counter:           0,
		CounterOn:         false,
		Checkpoint:        0,
		stringPosition:    1,
		charactersPerLine: 75,
		NextMessage:       true,
	}
	funcResult := texP.MessageWrapperToLines(texP.TextInput)
	for _, line := range funcResult {
		t.Logf("Result:%s", line)
	}

}
*/
