package battle

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/battleStatsDataManagement"
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
)

type Battle struct {
	EnemyTurn          bool
	Player             *battleStatsDataManagement.Character
	Enemy              *battleStatsDataManagement.Character
	turnInitiative     Initiative
	nextTurnInitiative Initiative
	battlePhase        Phase
	Turns              map[int]*Turn
	Turn               int
	BattleLost         bool
	BattleWon          bool
	PlayerAmmo         int
	EnemyAmmo          int
	WinningProb        int
	State              State
}

type Turn struct {
	Phase                  Phase                           `json:"phase"`
	PlayerSkillUsed        battleStatsDataManagement.Skill `json:"PlayerSkillUsed"`
	EnemySkillUsed         battleStatsDataManagement.Skill `json:"enemySkillUsed"`
	PlayerRoll             bool
	PlayerSecondaryRoll    bool
	EnemyRoll              bool
	EnemySecondaryRoll     bool
	DamageToPlayer         []int
	DamageToEnemy          []int
	EnemyEventTriggered    bool
	PlayerEventTriggered   bool
	WinProbAfterPlayerTurn int
	WinProbAfterEnemyTurn  int
	PlayerStartIndex       int
	EnemyStartIndex        int
	PlayerMessage          []string
	EnemyMessage           []string
	TurnInitiative         Initiative
	PlayerEffectsTriggered bool
	PlayerIndex            int
}

func NewBattle(player *battleStatsDataManagement.Character, enemy *battleStatsDataManagement.Character) *Battle {
	battle := Battle{}

	battle.EnemyTurn = false
	battle.Player = player
	battle.Enemy = enemy
	battle.battlePhase = Dialogue
	battle.Turns = make(map[int]*Turn)
	battle.Turn = 0
	battle.Turns[0] = &Turn{}
	battle.PlayerAmmo = 6
	battle.EnemyAmmo = 6
	battle.WinningProb = 50
	battle.turnInitiative = Player
	battle.State = NotStarted
	return &battle
}

func (b *Battle) GetTurn() *Turn {
	return b.Turns[b.Turn]
}

func (b *Battle) UpdateState() {
	turn := b.GetTurn()
	if turn.TurnInitiative == Player {
		if len(turn.PlayerMessage) > 0 {
			println("player message length:", len(turn.PlayerMessage))
			b.State = PlayerTurn
		}
		if len(turn.PlayerMessage) == 0 {
			b.State = EnemyTurn
		}
		if len(turn.PlayerMessage) == 0 && len(turn.EnemyMessage) == 0 {
			b.State = NextTurn
		}
	}

	if turn.TurnInitiative == Enemy {
		if len(turn.EnemyMessage) > 0 {
			b.State = EnemyTurn
		}
		if len(turn.EnemyMessage) == 0 {
			b.State = PlayerTurn
		}
		if len(turn.PlayerMessage) == 0 && len(turn.EnemyMessage) == 0 {
			b.State = NextTurn
		}
	}
}

func (b *Battle) UpdateWinProbability(winProb int) {
	b.WinningProb = winProb
}

func (b *Battle) UpdatePlayerAmmo() {
	turn := b.GetTurn()
	playerAmmoUsed := 0
	for _, effect := range turn.PlayerSkillUsed.Effects {
		if effect.EffectType == "shot" {
			playerAmmoUsed += effect.NShots
		}
	}
	b.PlayerAmmo -= playerAmmoUsed
}
func (b *Battle) UpdateEnemyAmmo() {
	turn := b.GetTurn()
	enemyAmmoUsed := 0
	for _, effect := range turn.EnemySkillUsed.Effects {
		if effect.EffectType == "shot" {
			b.EnemyAmmo += effect.NShots
		}
	}
	b.EnemyAmmo -= enemyAmmoUsed
}

func (b *Battle) RandTurnInitiative() Initiative {
	i := rand.IntN(2)
	if i == 0 {
		return Player
	}
	return Enemy
}

func (b *Battle) GetPhase() Phase {
	return b.battlePhase
}
func (b *Battle) incrementTurn() {
	b.Turn++
}

func (b *Battle) UpdateBattlePhase() {
	b.battlePhase = Shooting
}

func (b *Battle) UpdateInitiative(initiative Initiative) {
	b.turnInitiative = initiative
}

func CapitalizeWord(word string) string {
	return strings.ToUpper(string(word[0])) + word[1:]
}

func (b *Battle) Buff(usedOn *battleStatsDataManagement.Character, effect battleStatsDataManagement.Effect) {
	stats := usedOn.DisplayStats()
	affectedStat, err := battleStatsDataManagement.StringToStat(effect.Stat)
	if err != nil {
		log.Fatal(err)
	}
	affectedStatValue := stats[affectedStat]

	fmt.Printf("%s %s Stat before buff skils: %d\n", usedOn.Name, effect.Stat, affectedStatValue)

	usedOn.UpdateStat(affectedStat, effect.Amount)

	fmt.Printf("%s %s Stat after buff skil: %d\n", usedOn.Name, effect.Stat, usedOn.DisplayStat(affectedStat))
}

func (b *Battle) DamagePlayer() {
	turn := b.Turns[b.Turn]
	for _, shot := range turn.DamageToPlayer {
		if shot > 0 {
			b.Player.UpdateCharHealth(-shot)
		}
	}
}

func (b *Battle) DamageEnemy() {
	turn := b.Turns[b.Turn]
	for _, shot := range turn.DamageToEnemy {
		if shot > 0 {
			b.Enemy.UpdateCharHealth(-shot)
		}
	}
}

func (b *Battle) GenerateTurn(playerSkill battleStatsDataManagement.Skill) {

	if b.turnInitiative != b.nextTurnInitiative && b.battlePhase != Dialogue { //no idea how but we were getting into this loop during the shooting phase sometimes

		fmt.Printf("battle.go:81: changing turn initiative from:%d to:%d", b.turnInitiative, b.nextTurnInitiative)

		b.UpdateInitiative(b.nextTurnInitiative)
	}

	b.incrementTurn()
	b.Turns[b.Turn] = &Turn{}

	turn := b.Turns[b.Turn]
	turn.PlayerEventTriggered = false
	turn.PlayerEffectsTriggered = false
	turn.EnemyEventTriggered = false

	enemySkill, err := EnemyChooseSkill(*b, b.Enemy.DialogueSkills)

	if err != nil {
		log.Fatal(err)
	}

	turn.EnemySkillUsed = enemySkill
	enemyRoll := Roll(enemySkill.Effects[0].SuccessPer)
	turn.EnemyRoll = enemyRoll
	enemySecondaryRoll := false

	if len(enemySkill.Effects) > 1 {
		enemySecondaryRoll = Roll(enemySkill.Effects[1].SuccessPer)
	}

	turn.EnemySecondaryRoll = enemySecondaryRoll

	turn.PlayerSkillUsed = playerSkill
	playerRoll := Roll(playerSkill.Effects[0].SuccessPer)
	playerSecondaryRoll := false
	turn.PlayerRoll = playerRoll

	if len(playerSkill.Effects) > 1 {
		playerSecondaryRoll = Roll(playerSkill.Effects[1].SuccessPer)
	}

	turn.PlayerSecondaryRoll = playerSecondaryRoll

	var battleInitiative bool
	var playerSkillDialogue []string
	var enemySkillDialogue []string

	if enemySkill.SkillName == "draw" || playerSkill.SkillName == "draw" {

		battleInitiative = Draw(b.Player.DisplayStats(), b.Enemy.DisplayStats())

		if playerSkill.SkillName == "draw" {
			playerSkillDialogue = append(playerSkillDialogue, b.DrawFunction(b.Player, b.Enemy, battleInitiative)...)
		}

		if enemySkill.SkillName == "draw" {
			enemySkillDialogue = append(enemySkillDialogue, b.DrawFunction(b.Enemy, b.Player, battleInitiative)...)
		}
	}

	//apply buffs from playerBattleSprite based on roll results
	if playerSkill.SkillName != "draw" {
		playerSkillDialogue = append(playerSkillDialogue, b.generateMessageForUsedDialogueSkill(*b.Player, *b.Enemy, playerSkill, playerRoll, playerSecondaryRoll)...)
	}

	turn.PlayerMessage = playerSkillDialogue
	//apply buffs from enemyBattleSprite based on roll results

}

func (b *Battle) DrawFunction(user *battleStatsDataManagement.Character, opponent *battleStatsDataManagement.Character, battleInitiative bool) []string {

	drawSkillDialogue := b.generateDrawMessage(user, opponent, battleInitiative)

	if battleInitiative {
		b.nextTurnInitiative = Player
	} else {
		b.nextTurnInitiative = Enemy
	}
	return drawSkillDialogue
}

func (b *Battle) EnactEffects(skill battleStatsDataManagement.Skill, user *battleStatsDataManagement.Character, opponent *battleStatsDataManagement.Character, roll bool, secondaryRoll bool) {
	SkillEffectOne := skill.Effects[0]
	var skillEffectTwo battleStatsDataManagement.Effect

	if len(skill.Effects) > 1 {
		skillEffectTwo = skill.Effects[1]
	}

	if roll {
		if SkillEffectOne.On == "enemy" {
			b.Buff(opponent, SkillEffectOne)
		}
		if SkillEffectOne.On == "self" {
			b.Buff(user, SkillEffectOne)
		}
	}

	if secondaryRoll {
		if skillEffectTwo.On == "enemy" {
			b.Buff(opponent, SkillEffectOne)
		}

		if skillEffectTwo.On == "unsuccessfulSelf" && !roll {
			b.Buff(user, skillEffectTwo)
		}
		if skillEffectTwo.On == "successfulSelf" && roll {
			b.Buff(user, skillEffectTwo)
		}

	}
}

func (b *Battle) generateDrawMessage(turnTaker *battleStatsDataManagement.Character, opponent *battleStatsDataManagement.Character, battleInitiative bool) (drawMessage []string) {

	if b.Player.Name == turnTaker.Name {
		if battleInitiative {
			drawMessage = append(drawMessage, "elyse Reaches for her gun!")
			drawMessage = append(drawMessage, fmt.Sprintf("She's  a second faster!"))
			return drawMessage

		} else {
			drawMessage = append(drawMessage, "elyse reaches for her gun!")
			drawMessage = append(drawMessage, fmt.Sprintf("%s draws a second faster!", opponent.Name))
			return drawMessage

		}
	}

	if battleInitiative {

		drawMessage = append(drawMessage, fmt.Sprintf("%s reaches for their gun!", turnTaker.Name))
		drawMessage = append(drawMessage, fmt.Sprintf("elyse is a second faster!"))
		return drawMessage

	} else {

		drawMessage = append(drawMessage, fmt.Sprintf("%s reaches for their gun", turnTaker.Name))
		drawMessage = append(drawMessage, fmt.Sprintf("He's a second faster!"))
		return drawMessage
	}
}

func (b *Battle) generateMessageForUsedDialogueSkill(turnTaker battleStatsDataManagement.Character, opponent battleStatsDataManagement.Character, skill battleStatsDataManagement.Skill, roll bool, secondaryRoll bool) (message []string) {

	fmt.Printf("entering dialogue loop")

	output := fmt.Sprintf("%s uses %s", turnTaker.Name, skill.SkillName)

	message = append(message, output)

	dialogue := GetSkillDialogue(turnTaker.Name, skill.SkillName, roll)
	response := GetResponse(opponent.Name, skill.SkillName, roll)

	message = append(message, dialogue)
	message = append(message, response)

	skillName := CapitalizeWord(skill.SkillName)

	if !roll {
		message = append(message, fmt.Sprintf("%s is ineffective", skillName))
		if secondaryRoll && skill.Effects[0].On == "unsuccessfulSelf" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", turnTaker.Name, skill.Effects[1].Stat, skill.Effects[1].Amount))
		}
	}

	if roll {
		message = append(message, fmt.Sprintf("%s is effective", skillName))

		if skill.Effects[0].On == "self" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", turnTaker.Name, skill.Effects[0].Stat, skill.Effects[0].Amount))
		}

		if skill.Effects[0].On == "enemy" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", opponent.Name, skill.Effects[0].Stat, skill.Effects[0].Amount))
		}

		if secondaryRoll && skill.Effects[0].On == "successfulSelf" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", turnTaker.Name, skill.Effects[1].Stat, skill.Effects[1].Amount))
		}
	}
	return message
}

func (b *Battle) GenerateMessageForUsedCombatSkill(name string, skillName string, damage []int) (message []string) {

	output1 := fmt.Sprintf("%s uses %s", name, skillName)
	message = append(message, output1)

	for _, shot := range damage {

		output2 := fmt.Sprintf("%s takes a shot!", name)
		output3 := "It's a complete miss!"

		if shot > 0 {
			output3 = fmt.Sprintf("It does %d damage", shot)
		}
		if shot < 0 {
			output3 = "No Ammo!"
		}

		message = append(message, output2)
		message = append(message, output3)

	}

	return message
}

func (b *Battle) TakeCombatTurn(playerSkill battleStatsDataManagement.Skill) []string {
	fmt.Printf("battle.go:369: entering combat loop\n")
	if b.turnInitiative != b.nextTurnInitiative {
		fmt.Printf("battle.go:371: changing turn initiative from:%d to:%d\n", b.turnInitiative, b.nextTurnInitiative)
		b.UpdateInitiative(b.nextTurnInitiative)
	}

	var message []string

	b.incrementTurn()
	b.Turns[b.Turn] = &Turn{}

	if playerSkill.SkillName == "focused_shot" {
		b.nextTurnInitiative = Enemy
	}

	turn := b.Turns[b.Turn]
	enemySkill, err := EnemyChooseSkill(*b, b.Enemy.CombatSkills)

	turn.PlayerEventTriggered = false
	turn.EnemyEventTriggered = false
	if err != nil {
		log.Fatal(err)
	}

	turn.EnemySkillUsed = enemySkill

	if enemySkill.SkillName == "focused_shot" {
		b.nextTurnInitiative = Player
	}

	eOneTurnBuffAmount := 0
	eTurnAmmo := b.EnemyAmmo
	var eAffectedStat battleStatsDataManagement.Stat
	for _, effect := range enemySkill.Effects {

		if effect.EffectType == "buff" {
			if effect.On == "self" {

				eAffectedStat, err = battleStatsDataManagement.StringToStat(effect.Stat)
				if err != nil {
					log.Fatal(err)
				}

				eOneTurnBuffAmount = effect.Amount
				b.Buff(b.Enemy, effect)
			}
			if effect.On == "enemy" {
				b.Buff(b.Player, effect)
			}
		}

		if effect.EffectType == "shot" {
			for shot := 0; shot < effect.NShots; shot++ {
				if eTurnAmmo > 0 {
					eTurnAmmo--
					damage := Shoot(b.Enemy.DisplayStat(battleStatsDataManagement.Accuracy), b.Enemy.DisplayStat(battleStatsDataManagement.Fear), b.Enemy.DisplayStat(battleStatsDataManagement.Anger), effect.SuccessPer, effect.DamageRange)
					turn.DamageToPlayer = append(turn.DamageToPlayer, damage)
				} else {
					turn.DamageToPlayer = append(turn.DamageToPlayer, -1)
				}
			}
		}
		if effect.EffectType == "reload" {
			b.EnemyAmmo = 6
		}
	}

	if eOneTurnBuffAmount > 0 {
		fmt.Printf("battle line 350: Resetting %s one turn buff from: %d ", b.Enemy.Name, b.Enemy.DisplayStat(eAffectedStat))
		b.Enemy.UpdateStat(eAffectedStat, -eOneTurnBuffAmount)
		fmt.Printf("to: %d\n", b.Enemy.DisplayStat(eAffectedStat))
	}

	pTurnAmmo := b.PlayerAmmo
	turn.PlayerSkillUsed = playerSkill
	pOneTurnBuffAmount := 0
	var pAffectedStat battleStatsDataManagement.Stat
	for _, effect := range playerSkill.Effects {

		if effect.EffectType == "buff" {
			if effect.On == "self" {
				pAffectedStat, err = battleStatsDataManagement.StringToStat(effect.Stat)
				b.Buff(b.Player, effect)
				pOneTurnBuffAmount = effect.Amount
			}
			if effect.On == "enemy" {
				b.Buff(b.Enemy, effect)
			}
		}

		if effect.EffectType == "shot" {
			for i := 0; i < effect.NShots; i++ {
				if pTurnAmmo > 0 {
					pTurnAmmo--
					damage := Shoot(b.Player.DisplayStat(battleStatsDataManagement.Accuracy), b.Player.DisplayStat(battleStatsDataManagement.Fear), b.Player.DisplayStat(battleStatsDataManagement.Anger), effect.SuccessPer, effect.DamageRange)
					turn.DamageToEnemy = append(turn.DamageToEnemy, damage)
				} else {
					turn.DamageToEnemy = append(turn.DamageToEnemy, -1)
				}
			}
		}
		if effect.EffectType == "reload" {
			b.PlayerAmmo = 6
		}
	}

	if pOneTurnBuffAmount > 0 {
		fmt.Printf("battle line 379:Resetting %s one turn buff from: %d ", b.Player.Name, b.Player.DisplayStat(pAffectedStat))
		b.Player.UpdateStat(pAffectedStat, -pOneTurnBuffAmount)
		fmt.Printf("to: %d\n", b.Player.DisplayStat(pAffectedStat))
	}

	enemyMessage := b.GenerateMessageForUsedCombatSkill(b.Enemy.Name, enemySkill.SkillName, turn.DamageToPlayer)

	playerMessage := b.GenerateMessageForUsedCombatSkill(b.Player.Name, playerSkill.SkillName, turn.DamageToEnemy)

	if b.turnInitiative == Enemy {
		b.Turns[b.Turn].EnemyStartIndex = 0
		b.Turns[b.Turn].PlayerStartIndex = len(enemyMessage) - 1
		message = append(message, enemyMessage...)

		if b.Player.DisplayStat(battleStatsDataManagement.Health) <= 0 {
			message = append(message, enemyMessage...)
			message = append(message, fmt.Sprintf("Oh no! %s bit the dust!", b.Player.Name))
		}

		message = append(message, playerMessage...)

		if b.Enemy.DisplayStat(battleStatsDataManagement.Health) <= 0 {
			message = append(message, fmt.Sprintf("You win! %s bit the dust!", b.Enemy.Name))
		}
		return message
	}
	b.Turns[b.Turn].PlayerStartIndex = 0
	b.Turns[b.Turn].EnemyStartIndex = len(playerMessage) - 1
	message = append(message, playerMessage...)

	if b.Enemy.DisplayStat(battleStatsDataManagement.Health) <= 0 {
		message = append(message, fmt.Sprintf("you win! %s bit the dust!", b.Enemy.Name))
		return message
	}
	message = append(message, enemyMessage...)

	if b.Player.DisplayStat(battleStatsDataManagement.Health) <= 0 {
		message = append(message, fmt.Sprintf("Oh no! %s bit the dust!", b.Player.Name))
	}
	return message
}
