package main

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/dialogueData"
	"math/rand"
	"testing"
)

func MakeTestBattle() *battle.Battle {
	testTurn1 := battle.Turn{
		PlayerMessage:  []string{"foo", "bar", "baz"},
		EnemyMessage:   []string{"foo", "bar", "baz"},
		TurnInitiative: battle.Enemy,
	}

	stats := map[string]int{
		"health":    4,
		"accuracy":  2,
		"drawSpeed": 3,
		"anger":     0,
		"fear":      0,
	}

	elyse := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "female1")
	george := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1")

	b := battle.NewBattle(&elyse, &george)

	b.Turns[0] = &testTurn1

	return b
}

func TestShoot(t *testing.T) {

	damageRange := []int{0, 2}
	var results float64

	for i := 0; i < 10000; i++ {
		damage := battle.Shoot(10, 0, 0, 70, damageRange)
		results += float64(damage)
	}

}

func TestBattleState(t *testing.T) {
	testTurn1 := battle.Turn{
		PlayerMessage:  []string{"foo", "bar", "baz"},
		EnemyMessage:   []string{"foo", "bar", "baz"},
		TurnInitiative: battle.Enemy,
	}

	stats := map[string]int{
		"health":    4,
		"accuracy":  2,
		"drawSpeed": 3,
		"anger":     0,
		"fear":      0,
	}

	elyse := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1")
	george := battleStats.NewCharacter("george", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1")

	b := battle.NewBattle(&elyse, &george)

	b.UpdateInitiative(battle.Enemy)

	b.Turns[0] = &testTurn1

	b.UpdateState()

	if b.State != battle.EnemyTurn {
		t.Fatalf("test1: incorrect state, state = %d", b.State)
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
	combatSkills, _, _ := battleStats.LoadSkills()
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
	chars, _ := battleStats.LoadCharacters()
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

	elyse := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1")

	if elyse.DisplayStat(battleStats.Health) != 4 {
		t.Fatalf(`method displayCharHealth did not work`)
	}

	elyse.UpdateCharAccuracy(-1)

	if elyse.DisplayStat(battleStats.Accuracy) != 1 {
		t.Fatalf(`method updateCharAccuracy did not work`)
	}

	elyse.UpdateCharHealth(-1)

	elyse.UpdateCharHealth(-0)

}

func Test_use_stat_Buff(t *testing.T) {
	stats := map[string]int{
		"health":    4,
		"accuracy":  2,
		"drawspeed": 3,
		"anger":     0,
		"fear":      0,
	}

	elyse := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1")

	elyse.UpdateStat(battleStats.Anger, 1)

	if elyse.DisplayStat(battleStats.Anger) != 1 {
		t.Fatalf(`method ResetStatusStats() failed to update correctly %d`, elyse.Stats[battleStats.Anger])
	}

	elyse.ResetStatusStats()

	if elyse.DisplayStat(battleStats.Anger) != 0 {
		t.Fatalf(`method ResetStatusStats() failed to update correctly %d`, elyse.Stats[battleStats.Anger])
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

func TestDialogueData(t *testing.T) {
	dialogue := dialogueData.GetResponse("bethAnne", 1)
	println(dialogue)
}
