package battle

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/assets"
	"github.com/acoco10/QuickDrawAdventure/battleStats"
	"log"
	"math/rand/v2"
	"strings"
)

type Phase uint8

const (
	Dialogue Phase = iota
	Shooting
)

type Initiative uint8

const (
	Player Initiative = iota
	Enemy
)

type State uint8

const (
	PlayerTurn State = iota
	EnemyTurn
	NotStarted
	NextTurn
	Over
)

type Battle struct {
	EnemyTurn           bool
	CharacterBattleData map[Initiative]*CharacterBattleData
	turnInitiative      Initiative
	nextTurnInitiative  Initiative
	BattlePhase         Phase
	Turns               map[int]*Turn
	Turn                int
	BattleLost          bool
	BattleWon           bool
	WinningProb         int
	State               State
	Tension             int
	DialogueText        string
}

type Turn struct {
	CharacterData          []*CharacterTurnData
	Phase                  Phase
	WinProbAfterPlayerTurn int
	WinProbAfterEnemyTurn  int
	TurnInitiative         Initiative
	PlayerIndex            int
	EnemyIndex             int
	PlayerTurnCompleted    bool
	EnemyTurnCompleted     bool
}

func NewBattle(player *battleStats.CharacterData, enemy *battleStats.CharacterData) *Battle {
	enemyBattleData := CharacterBattleData{
		CharacterData:     enemy,
		Ammo:              6,
		DrawBonus:         true,
		CharacterTurnData: &CharacterTurnData{},
	}

	playerBattleData := CharacterBattleData{
		CharacterData:     player,
		Ammo:              6,
		DrawBonus:         false,
		CharacterTurnData: &CharacterTurnData{},
	}

	charData := make(map[Initiative]*CharacterBattleData)

	charData[Enemy] = &enemyBattleData
	charData[Player] = &playerBattleData

	battle := Battle{}
	battle.EnemyTurn = false
	battle.CharacterBattleData = charData
	battle.BattlePhase = Dialogue
	battle.Turns = make(map[int]*Turn)
	battle.Turn = 0
	battle.Turns[0] = &Turn{}
	battle.WinningProb = 50
	battle.turnInitiative = Player
	battle.State = NotStarted

	if len(enemy.DialogueSkills) == 0 {
		battle.BattlePhase = Shooting
	}

	data, err := assets.Dialogue.ReadFile("dialogueData/battleDialogue.json")
	if err != nil {
		log.Fatal(err)
	}

	battle.DialogueText = string(data)

	return &battle
}

func (b *Battle) GetTurn() *Turn {
	return b.Turns[b.Turn]
}

func (b *Battle) UpdateState() {
	turn := b.GetTurn()
	if b.turnInitiative == Player && len(b.CharacterBattleData[Player].Message) > 0 {
		if !turn.PlayerTurnCompleted {
			b.State = PlayerTurn
		} else if !turn.EnemyTurnCompleted {
			b.State = EnemyTurn
		} else {
			b.State = NextTurn
		}
	}
	if b.turnInitiative == Enemy && len(b.CharacterBattleData[Enemy].Message) > 0 {
		if !turn.EnemyTurnCompleted {
			b.State = EnemyTurn
		} else if !turn.PlayerTurnCompleted {
			b.State = PlayerTurn
		} else {
			b.State = NextTurn
		}
	}
	if b.BattleWon || b.BattleLost {
		b.State = Over
	}
}

func (b *Battle) UpdateWinProbability(winProb int) {
	b.WinningProb = winProb
}

func (b *Battle) UpdateAmmo() {
	for _, char := range b.CharacterBattleData {
		char.UpdateAmmo()
	}
}

func (b *Battle) RandTurnInitiative() Initiative {
	i := rand.IntN(2)
	if i == 0 {
		return Player
	}
	return Enemy
}

func (b *Battle) GetPhase() Phase {
	return b.BattlePhase
}
func (b *Battle) incrementTurn() {
	b.Turn++
}

func (b *Battle) UpdateBattlePhase() {
	b.BattlePhase = Shooting
	b.State = NextTurn
}

func (b *Battle) UpdateInitiative(initiative Initiative) {
	b.turnInitiative = initiative
}

func CapitalizeWord(word string) string {
	return strings.ToUpper(string(word[0])) + word[1:]
}

func (b *Battle) Buff(usedOn *battleStats.CharacterData, effect battleStats.Effect, target string) {

	stats := usedOn.DisplayStats()
	affectedStat, err := battleStats.StringToStat(effect.Stat)
	if err != nil {
		log.Fatal(err)
	}
	affectedStatValue := stats[affectedStat]

	fmt.Printf("%s %s Stat before buff skils: %d\n", usedOn.Name, effect.Stat, affectedStatValue)
	if target == usedOn.Weakness {
		println("turn:", b.Turn, "weakness triggered")
		usedOn.UpdateStat(affectedStat, effect.Amount*2)
	}
	if target != usedOn.Weakness {
		usedOn.UpdateStat(affectedStat, effect.Amount)
	}

	fmt.Printf("%s %s Stat after buff skil: %d\n", usedOn.Name, effect.Stat, usedOn.DisplayStat(affectedStat))
}

func (b *Battle) DamageCharacter(charDamaged Initiative, attacker Initiative) {
	for _, dmg := range b.CharacterBattleData[attacker].DamageOutput {
		b.CharacterBattleData[charDamaged].UpdateCharHealth(-dmg)
	}
}

func (b *Battle) UpdateChar(char *CharacterBattleData, enemy *CharacterBattleData, skillUsed battleStats.Skill) {
	char.EffectsTriggered = false
	char.EventTriggered = false
	char.Completed = false
	char.WeaknessTargeted = false
	char.SkillUsed = skillUsed
	char.Roll = Roll(skillUsed.Effects[0].SuccessPer)
	if len(skillUsed.Effects) > 1 {
		char.SecondaryRoll = Roll(skillUsed.Effects[1].SuccessPer)
	}
	CheckWeakness(char, enemy)
	if skillUsed.SkillName != "draw" {
		char.Message = b.generateMessageForUsedDialogueSkill(*char, *enemy)
	}
	if skillUsed.SkillName == "draw" {
		drawMSG := b.DrawFunction(char, enemy)
		char.Message = drawMSG
	}

}

func CheckWeakness(char *CharacterBattleData, enemy *CharacterBattleData) {
	if char.SkillUsed.SkillName != "draw" {
		if char.SkillUsed.Target == enemy.Weakness {
			enemy.WeaknessTargeted = true
		}
	}
}

func DrawInitiative(player *CharacterBattleData, enemy *CharacterBattleData) Initiative {
	battleInitiative := ReadyDraw(player.DisplayStats(), enemy.DisplayStats())
	if battleInitiative {
		return Player
	} else {
		return Enemy
	}
}

func (b *Battle) GenerateTurn(playerSkill battleStats.Skill) {
	if b.turnInitiative != b.nextTurnInitiative && b.BattlePhase != Dialogue { //no idea how but we were getting into this loop during the shooting phase sometimes
		fmt.Printf("battle.go:81: changing turn initiative from:%d to:%d", b.turnInitiative, b.nextTurnInitiative)
		b.UpdateInitiative(b.nextTurnInitiative)
	}

	b.incrementTurn()

	println("turn incremented", "turn:", b.Turn)

	b.Turns[b.Turn] = &Turn{}

	for _, char := range b.CharacterBattleData {
		println("updating character:", char.Name)
		if char.Name == "elyse" {
			b.UpdateChar(char, b.CharacterBattleData[Enemy], playerSkill)
		} else {
			skillUsed, err := EnemyChooseSkill(*b, char.DialogueSkills)
			if err != nil {
				log.Fatal(err)
			}

			println("enemy chose:", skillUsed.SkillName)
			b.UpdateChar(char, b.CharacterBattleData[Player], skillUsed)
		}
	}
}

func (b *Battle) DrawFunction(user *CharacterBattleData, opponent *CharacterBattleData) []string {
	dInitiative := DrawInitiative(b.CharacterBattleData[Player], b.CharacterBattleData[Enemy])
	drawSkillDialogue := b.GenDrawMessage(*user, *opponent, dInitiative)
	b.nextTurnInitiative = dInitiative
	return drawSkillDialogue
}

func (b *Battle) SetDrawBonus() {
	if b.nextTurnInitiative == Player {
		b.CharacterBattleData[Player].DrawBonus = true
	} else {
		b.CharacterBattleData[Enemy].DrawBonus = true
	}
}

func (b *Battle) EnactEffects(skill battleStats.Skill, user *battleStats.CharacterData, opponent *battleStats.CharacterData, roll bool, secondaryRoll bool) {
	b.Tension = b.Tension + skill.Tension
	SkillEffectOne := skill.Effects[0]
	var skillEffectTwo battleStats.Effect

	if len(skill.Effects) > 1 {
		skillEffectTwo = skill.Effects[1]
	}

	if roll {
		if SkillEffectOne.On == "enemy" {
			b.Buff(opponent, SkillEffectOne, skill.Target)
		}
		if SkillEffectOne.On == "self" {
			b.Buff(user, SkillEffectOne, skill.Target)
		}
	}

	if secondaryRoll {
		if skillEffectTwo.On == "enemy" {
			b.Buff(opponent, SkillEffectOne, skill.Target)
		}

		if skillEffectTwo.On == "unsuccessfulSelf" && !roll {
			b.Buff(user, skillEffectTwo, skill.Target)
		}
		if skillEffectTwo.On == "successfulSelf" && roll {
			b.Buff(user, skillEffectTwo, skill.Target)
		}

	}
}

func (b *Battle) GenDrawMessage(turnTaker CharacterBattleData, opponent CharacterBattleData, battleInitiative Initiative) (drawMessage []string) {
	if turnTaker.Name == "elyse" {
		if battleInitiative == Player {
			drawMessage = append(drawMessage, "elyse Reaches for her gun!")
			drawMessage = append(drawMessage, fmt.Sprintf("She's  a second faster!"))
			return drawMessage

		} else {
			drawMessage = append(drawMessage, "elyse reaches for her gun!")
			drawMessage = append(drawMessage, fmt.Sprintf("%s draws a second faster!", opponent.Name))
			return drawMessage
		}
	}
	if battleInitiative == Enemy {
		drawMessage = append(drawMessage, fmt.Sprintf("%s reaches for their gun", turnTaker.Name))
		drawMessage = append(drawMessage, fmt.Sprintf("He's a second faster!"))
		return drawMessage
	} else {
		drawMessage = append(drawMessage, fmt.Sprintf("%s reaches for their gun!", turnTaker.Name))
		drawMessage = append(drawMessage, fmt.Sprintf("Elyse is a second faster!"))
		return drawMessage
	}
}

func (b *Battle) generateMessageForUsedDialogueSkill(turnTaker CharacterBattleData, opponent CharacterBattleData) (message []string) {

	name := CapitalizeWord(turnTaker.Name)
	skillName := CapitalizeWord(turnTaker.SkillUsed.SkillName)
	skill := turnTaker.SkillUsed
	output := fmt.Sprintf("%s uses %s", name, skillName)

	message = append(message, output)
	dialogue := turnTaker.SkillUsed.Text
	message = append(message, dialogue)

	if !turnTaker.Roll {
		message = append(message, fmt.Sprintf("%s is ineffective", skillName))
		if turnTaker.SecondaryRoll && skill.Effects[0].On == "unsuccessfulSelf" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", name, skill.Effects[1].Stat, skill.Effects[1].Amount))
		}
	}

	if turnTaker.Roll {
		message = append(message, fmt.Sprintf("%s is effective", skillName))

		if skill.Effects[0].On == "self" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", name, skill.Effects[0].Stat, skill.Effects[0].Amount))
		}

		if skill.Effects[0].On == "enemy" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", opponent.Name, skill.Effects[0].Stat, skill.Effects[0].Amount))
		}

		if turnTaker.SecondaryRoll && skill.Effects[0].On == "successfulSelf" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", turnTaker.Name, skill.Effects[1].Stat, skill.Effects[1].Amount))
		}
	}
	return message
}

func (b *Battle) GenerateMessageForUsedCombatSkill(name string, skillName string, damage []int) (message []string) {

	output1 := fmt.Sprintf("%s uses %s", name, skillName)
	message = append(message, output1)

	if skillName != "reload" {
		totalDmg := 0
		for _, shot := range damage {
			if shot > 0 {
				totalDmg += shot
			}
		}

		output2 := fmt.Sprintf("It does %d total damage!,", totalDmg)
		message = append(message, output2)
	}
	return message
}

func (b *Battle) UpdateCharCombat(char *CharacterBattleData, enemy *CharacterBattleData, skillUsed battleStats.Skill) {
	char.EffectsTriggered = false
	char.EventTriggered = false
	char.Completed = false
	char.WeaknessTargeted = false
	char.SkillUsed = skillUsed
	char.Roll = Roll(skillUsed.Effects[0].SuccessPer)
	if len(skillUsed.Effects) > 1 {
		char.SecondaryRoll = Roll(skillUsed.Effects[1].SuccessPer)
	}
	CheckWeakness(char, enemy)

	currentAmmo := char.Ammo
	var effectedStat battleStats.Stat
	var oneTurnBuffAmount int
	for _, effect := range char.SkillUsed.Effects {

		if effect.EffectType == "buff" {
			if effect.On == "self" {

				effStat, err := battleStats.StringToStat(effect.Stat)

				if err != nil {
					log.Fatal(err)
				}
				effectedStat = effStat

				oneTurnBuffAmount = effect.Amount
				b.Buff(char.CharacterData, effect, char.SkillUsed.Target)
			}
			if effect.On == "enemy" {
				b.Buff(enemy.CharacterData, effect, char.SkillUsed.Target)
			}
		}

		if effect.EffectType == "shot" {
			for shot := 0; shot < effect.NShots; shot++ {
				if currentAmmo > 0 {
					currentAmmo--
					damage := Shoot(char.DisplayStat(battleStats.Accuracy), char.DisplayStat(battleStats.Fear), char.DisplayStat(battleStats.Anger), effect.SuccessPer, effect.DamageRange)
					enemy.DamageTaken = append(enemy.DamageTaken, damage)
				} else {
					enemy.DamageTaken = append(enemy.DamageTaken, -1)
				}
			}
		}
		if effect.EffectType == "reload" {
			char.Ammo = 6
		}
	}

	if oneTurnBuffAmount > 0 {
		fmt.Printf("battle line 350: Resetting %s one turn buff from: %d ", char.Name, char.DisplayStat(effectedStat))
		char.UpdateStat(effectedStat, -oneTurnBuffAmount)
		fmt.Printf("to: %d\n buff amount was: %d\n", char.DisplayStat(effectedStat), oneTurnBuffAmount)
	}
	message := b.GenerateMessageForUsedCombatSkill(char.Name, char.SkillUsed.SkillName, enemy.DamageTaken)
	char.Message = message

}

func (b *Battle) TakeCombatTurn(playerSkill battleStats.Skill) {
	fmt.Printf("battle.go:369: entering combat loop\n")
	if b.turnInitiative != b.nextTurnInitiative {
		fmt.Printf("battle.go:371: changing turn initiative from:%d to:%d\n", b.turnInitiative, b.nextTurnInitiative)
		b.UpdateInitiative(b.nextTurnInitiative)
	}
	b.incrementTurn()
	b.Turns[b.Turn] = &Turn{}

	for _, char := range b.CharacterBattleData {
		println("updating character:", char.Name)
		if char.Name == "elyse" {
			b.UpdateCharCombat(char, b.CharacterBattleData[Player], playerSkill)
		} else {
			skillUsed, err := EnemyChooseSkill(*b, char.CombatSkills)
			if err != nil {
				log.Fatal(err)
			}
			b.UpdateCharCombat(char, b.CharacterBattleData[Enemy], skillUsed)
		}
	}
}
