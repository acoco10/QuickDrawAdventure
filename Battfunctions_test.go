package main

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/battleStatsDataManagement"
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

	elyse := battleStatsDataManagement.NewCharacter("elyse", stats, map[string]battleStatsDataManagement.Skill{}, map[string]battleStatsDataManagement.Skill{}, "anger")
	george := battleStatsDataManagement.NewCharacter("elyse", stats, map[string]battleStatsDataManagement.Skill{}, map[string]battleStatsDataManagement.Skill{}, "anger")

	b := battle.NewBattle(&elyse, &george)

	b.Turns[0] = &testTurn1

	return b
}

func TestShoot(t *testing.T) {

	damageRange := []int{2, 0}
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

	elyse := battleStatsDataManagement.NewCharacter("elyse", stats, map[string]battleStatsDataManagement.Skill{}, map[string]battleStatsDataManagement.Skill{}, "anger")
	george := battleStatsDataManagement.NewCharacter("elyse", stats, map[string]battleStatsDataManagement.Skill{}, map[string]battleStatsDataManagement.Skill{}, "anger")

	b := battle.NewBattle(&elyse, &george)

	b.Turns[0] = &testTurn1

	b.UpdateState()

	if b.State != battle.EnemyTurn {
		t.Fatalf("test1: incorrect state")
	}

	testTurn2 := battle.Turn{
		PlayerMessage:  []string{"foo", "bar", "baz"},
		EnemyMessage:   []string{},
		TurnInitiative: battle.Enemy,
	}

	b.Turns[0] = &testTurn2
	b.UpdateState()
	if b.State != battle.PlayerTurn {
		t.Fatalf("test2: incorrect state: %d should be %d", b.State, battle.PlayerTurn)
	}

	testTurn3 := battle.Turn{
		PlayerMessage:  []string{"foo"},
		EnemyMessage:   []string{},
		TurnInitiative: battle.Player,
	}

	b.Turns[0] = &testTurn3
	b.UpdateState()
	if b.State != battle.PlayerTurn {
		t.Fatalf("test3: incorrect state")
	}

	testTurn4 := battle.Turn{
		PlayerMessage:  []string{},
		EnemyMessage:   []string{},
		TurnInitiative: battle.Player,
	}

	b.Turns[0] = &testTurn4
	b.UpdateState()
	if b.State != battle.NextTurn {
		t.Fatalf("test4: incorrect state")
	}

	testTurn5 := battle.Turn{
		PlayerMessage:  []string{"foo"},
		EnemyMessage:   []string{"bar"},
		TurnInitiative: battle.Player,
	}
	b.Turns[0] = &testTurn5
	b.UpdateState()
	if b.State != battle.PlayerTurn {
		t.Fatalf("test5: incorrect state")
	}

	testTurn6 := battle.Turn{
		PlayerMessage:  []string{"foo", "bar", "baz"},
		EnemyMessage:   []string{"foo", "bar", "baz"},
		TurnInitiative: battle.Enemy,
	}

	b.Turns[0] = &testTurn6

	b.UpdateState()
	if b.State != battle.EnemyTurn {
		t.Fatalf("test6: incorrect state")
	}

}

func TestLoadBadSkillJSON(t *testing.T) {
	badSkills, err := battleStatsDataManagement.LoadSkillsFromPath("badSkills.json")

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
	combatSkills, _, _ := battleStatsDataManagement.LoadSkills()
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
	chars, _ := battleStatsDataManagement.LoadCharacters()
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

	elyse := battleStatsDataManagement.NewCharacter("elyse", stats, map[string]battleStatsDataManagement.Skill{}, map[string]battleStatsDataManagement.Skill{}, "anger")

	if elyse.DisplayStat(battleStatsDataManagement.Health) != 4 {
		t.Fatalf(`method displayCharHealth did not work`)
	}

	elyse.UpdateCharAccuracy(-1)

	if elyse.DisplayStat(battleStatsDataManagement.Accuracy) != 1 {
		t.Fatalf(`method updateCharAccuracy did not work`)
	}

	elyse.UpdateCharHealth(-1)

	if elyse.DisplayStat(battleStatsDataManagement.Health) != 3 {
		t.Fatalf(`method updateCharhealthdid not work health value:%d expected value 3`, elyse.DisplayStat(battleStatsDataManagement.Health))
	}

	elyse.UpdateCharHealth(-0)
	if elyse.DisplayStat(battleStatsDataManagement.Health) != 4 {
		t.Fatalf(`method updateCharhealth did not work value inserted %d above its maximum: 4 `, elyse.DisplayStat(battleStatsDataManagement.Health))
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

	elyse := battleStatsDataManagement.NewCharacter("elyse", stats, map[string]battleStatsDataManagement.Skill{}, map[string]battleStatsDataManagement.Skill{}, "anger")

	elyse.UpdateStat(battleStatsDataManagement.Anger, 1)

	if elyse.DisplayStat(battleStatsDataManagement.Anger) != 1 {
		t.Fatalf(`method ResetStatusStats() failed to update correctly %d`, elyse.Stats[battleStatsDataManagement.Anger])
	}

	elyse.ResetStatusStats()

	if elyse.DisplayStat(battleStatsDataManagement.Anger) != 0 {
		t.Fatalf(`method ResetStatusStats() failed to update correctly %d`, elyse.Stats[battleStatsDataManagement.Anger])
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
