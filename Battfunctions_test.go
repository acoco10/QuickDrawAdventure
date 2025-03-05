package main

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/battle"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"github.com/acoco10/QuickDrawAdventure/gameScenes"
	"log"
	"math/rand"
	"testing"
)

func MakeTestBattle() *battle.Battle {

	stats := map[string]int{
		"health":    4,
		"accuracy":  2,
		"drawSpeed": 3,
		"anger":     0,
		"fear":      0,
	}

	elyse := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "female1", 3)
	george := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1", 3)

	b := battle.NewBattle(&elyse, &george)

	return b
}

func TestDisEq(t *testing.T) {
	var x1, y1, x2, y2 float64
	x1, y1 = 54, 82
	x2, y2 = 345, 97

	answer := gameScenes.DistanceEq(x1, y1, x2, y2)

	if answer != 291.38634147811 {
		println(answer)
	}
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

		TurnInitiative: battle.Enemy,
	}

	stats := map[string]int{
		"health":    4,
		"accuracy":  2,
		"drawSpeed": 3,
		"anger":     0,
		"fear":      0,
	}

	elyse := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1", 3)
	george := battleStats.NewCharacter("george", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1", 3)

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

func TestLoadSkillsJSON(t *testing.T) {
	dialogueSkills, err := battleStats.LoadSkillsFromPath("battleData/equipDialogueSkills.json")
	if err != nil {
		log.Fatal(err)
	}

	for Name, skill := range dialogueSkills {
		println(Name, skill.Target)
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

	elyse := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1", 3)

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

	elyse := battleStats.NewCharacter("elyse", stats, map[string]battleStats.Skill{}, map[string]battleStats.Skill{}, "anger", "male1", 3)

	elyse.UpdateStat(battleStats.Anger, 1)

	if elyse.DisplayStat(battleStats.Anger) != 1 {
		t.Fatalf(`method ResetStatusStats() failed to update correctly %d`, elyse.Stats[battleStats.Anger])
	}

	elyse.ResetStatusStats()

	if elyse.DisplayStat(battleStats.Anger) != 0 {
		t.Fatalf(`method ResetStatusStats() failed to update correctly %d`, elyse.Stats[battleStats.Anger])
	}

}

func TestLoadCharacters(t *testing.T) {
	characters, err := battleStats.LoadCharacters()
	if err != nil {
		log.Fatal("error loading characters.json error:", err)
	}

	elyse := characters[battleStats.Elyse]
	enemy := characters[battleStats.Antonio]

	println("Testing Elyse Dialogue Skills")
	println("-------------------------------------")
	for _, skill := range elyse.CombatSkills {
		println(skill.SkillName)
	}

	println("Testing Elyse Combat Skills")
	println("-------------------------------------")
	for _, skill := range elyse.DialogueSkills {
		println(skill.SkillName)
	}

	println("Testing chosen enemy Dialogue Skills")
	println("-------------------------------------")
	for _, skill := range enemy.CombatSkills {
		println(skill.SkillName)
	}

	println("Testing chosen enemy Combat Skills")
	println("-------------------------------------")
	for _, skill := range enemy.DialogueSkills {
		println(skill.SkillName)
	}

}

func TestEnemyChooseSkill(t *testing.T) {
	characters, err := battleStats.LoadCharacters()
	if err != nil {
		log.Fatal("error loading characters.json error:", err)
	}

	elyse := characters[battleStats.Elyse]
	enemy := characters[battleStats.Antonio]

	testBattle := battle.NewBattle(&elyse, &enemy)

	println("Testing Enemy Choose Dialogue Skills")
	println("-------------------------------------")
	for i := 0; i < 10; i++ {
		_, err := battle.EnemyChooseSkill(*testBattle, enemy.DialogueSkills)
		if err != nil {
			log.Fatal("error when choosing skill")
		}
	}

	println("Testing Enemy Choose Combat Skills")
	println("-------------------------------------")

	for i := 0; i < 10; i++ {
		_, err := battle.EnemyChooseSkill(*testBattle, enemy.CombatSkills)
		if err != nil {
			log.Fatal("error when choosing skill")
		}

	}

	testBattle.Tension = 0

	var drawsChosen = make([]float64, 15)
	var shotInBackChosen float64
	var reloadChosen float64
	index := 1
	for i := 0; i < 1500; i++ {
		if i > 100*index {
			index++
			testBattle.Tension++
		}
		skill, err := battle.EnemyChooseSkill(*testBattle, enemy.DialogueSkills)
		if err != nil {
			log.Fatal("error when choosing skill")
		}
		if skill.SkillName == "draw" {
			drawsChosen[index-1]++
		}
		if skill.SkillName == "shotInTheBack" {
			shotInBackChosen++
		}
	}
	testBattle.BattlePhase = battle.Shooting
	testBattle.CharacterBattleData[battle.Enemy].Ammo = 0
	for i := 0; i < 100; i++ {
		skill, err := battle.EnemyChooseSkill(*testBattle, enemy.CombatSkills)
		if err != nil {
			log.Fatal("error when choosing skill")
		}
		if skill.SkillName == "reload" {
			reloadChosen++
		}

	}

	println("Testing Enemy Choose Dialogue Skills Medium Tension")
	println("-------------------------------------")
	for i, draw := range drawsChosen {
		print("tension =", i)
		fmt.Printf("draw skill chosen %f percent of the time \n", draw/100)
	}

	fmt.Printf("shotInBack chosen %f percent of the time \n", shotInBackChosen/1000)

	println("Testing Enemy reload 0 ammo")
	println("-------------------------------------")
	fmt.Printf("reload chosen %f percent of the time \n", reloadChosen/100)

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
