package battle

import (
	"fmt"
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
	EnemyTurn          bool
	Player             *battleStats.Character
	Enemy              *battleStats.Character
	turnInitiative     Initiative
	nextTurnInitiative Initiative
	BattlePhase        Phase
	Turns              map[int]*Turn
	Turn               int
	BattleLost         bool
	BattleWon          bool
	PlayerAmmo         int
	EnemyAmmo          int
	WinningProb        int
	State              State
	Tension            int
}

type Turn struct {
	Phase                  Phase
	PlayerSkillUsed        battleStats.Skill
	EnemySkillUsed         battleStats.Skill
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
	PlayerMessage          []string
	EnemyMessage           []string
	TurnInitiative         Initiative
	PlayerEffectsTriggered bool
	EnemyEffectsTriggered  bool
	PlayerIndex            int
	EnemyIndex             int
	PlayerTurnCompleted    bool
	EnemyTurnCompleted     bool
	EnemyWeakness          bool
	PlayerWeakness         bool
}

func NewBattle(player *battleStats.Character, enemy *battleStats.Character) *Battle {
	battle := Battle{}
	battle.EnemyTurn = false
	battle.Player = player
	battle.Enemy = enemy
	battle.BattlePhase = Dialogue
	battle.Turns = make(map[int]*Turn)
	battle.Turn = 0
	battle.Turns[0] = &Turn{}
	battle.PlayerAmmo = 6
	battle.EnemyAmmo = 6
	battle.WinningProb = 50
	battle.turnInitiative = Player
	battle.State = NotStarted

	if len(enemy.DialogueSkills) == 0 {
		battle.BattlePhase = Shooting
	}

	return &battle
}

func (b *Battle) GetTurn() *Turn {
	return b.Turns[b.Turn]
}

func (b *Battle) UpdateState() {
	turn := b.GetTurn()
	if b.turnInitiative == Player && len(turn.PlayerMessage) > 0 {
		if !turn.PlayerTurnCompleted {
			b.State = PlayerTurn
		} else if !turn.EnemyTurnCompleted {
			b.State = EnemyTurn
		} else {
			b.State = NextTurn
		}
	}
	if b.turnInitiative == Enemy && len(turn.PlayerMessage) > 0 {
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

func (b *Battle) UpdatePlayerAmmo() {
	turn := b.GetTurn()
	playerAmmoUsed := 0
	for _, effect := range turn.PlayerSkillUsed.Effects {
		if effect.EffectType == "shot" {
			playerAmmoUsed += effect.NShots
		}
	}
	b.PlayerAmmo -= playerAmmoUsed
	if b.PlayerAmmo < 0 {
		b.PlayerAmmo = 0
	}
}
func (b *Battle) UpdateEnemyAmmo() {
	turn := b.GetTurn()
	enemyAmmoUsed := 0
	for _, effect := range turn.EnemySkillUsed.Effects {
		if effect.EffectType == "shot" {
			enemyAmmoUsed += effect.NShots
		}
	}
	b.EnemyAmmo -= enemyAmmoUsed
	if b.EnemyAmmo < 0 {
		b.EnemyAmmo = 0
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
}

func (b *Battle) UpdateInitiative(initiative Initiative) {
	b.turnInitiative = initiative
}

func CapitalizeWord(word string) string {
	return strings.ToUpper(string(word[0])) + word[1:]
}

func (b *Battle) Buff(usedOn *battleStats.Character, effect battleStats.Effect) {
	weakness, err := battleStats.StringToStat(effect.Stat)
	stats := usedOn.DisplayStats()
	affectedStat, err := battleStats.StringToStat(effect.Stat)
	if err != nil {
		log.Fatal(err)
	}
	affectedStatValue := stats[affectedStat]

	fmt.Printf("%s %s Stat before buff skils: %d\n", usedOn.Name, effect.Stat, affectedStatValue)
	if weakness == usedOn.Weakness {
		println("turn:", b.Turn, "weakness triggered")
		usedOn.UpdateStat(affectedStat, effect.Amount*2)
	}
	if weakness != usedOn.Weakness {
		usedOn.UpdateStat(affectedStat, effect.Amount)
	}

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

func (b *Battle) GenerateTurn(playerSkill battleStats.Skill) {

	if b.turnInitiative != b.nextTurnInitiative && b.BattlePhase != Dialogue { //no idea how but we were getting into this loop during the shooting phase sometimes

		fmt.Printf("battle.go:81: changing turn initiative from:%d to:%d", b.turnInitiative, b.nextTurnInitiative)

		b.UpdateInitiative(b.nextTurnInitiative)
	}

	b.incrementTurn()
	println("turn incremented", "turn:", b.Turn)
	b.Turns[b.Turn] = &Turn{}

	turn := b.Turns[b.Turn]

	turn.PlayerEventTriggered = false
	turn.PlayerEffectsTriggered = false
	turn.EnemyEventTriggered = false
	turn.PlayerTurnCompleted = false
	turn.EnemyTurnCompleted = false
	turn.EnemyWeakness = false
	turn.PlayerWeakness = false

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
	var playerStatAffectedBySkill battleStats.Stat
	if playerSkill.SkillName != "draw" {
		playerStatAffectedBySkill, err = battleStats.StringToStat(playerSkill.Effects[0].Stat)
		if err != nil {
			log.Fatal(err)
		}
	}
	if playerRoll && playerStatAffectedBySkill == b.Enemy.Weakness {
		turn.EnemyWeakness = true
	}
	var battleInitiative bool
	var playerSkillDialogue []string
	var enemySkillDialogue []string

	if enemySkill.SkillName == "draw" || playerSkill.SkillName == "draw" {

		battleInitiative = ReadyDraw(b.Player.DisplayStats(), b.Enemy.DisplayStats())
		if battleInitiative {
			b.nextTurnInitiative = Player
		} else {
			b.nextTurnInitiative = Enemy
		}

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

	if enemySkill.SkillName != "draw" {
		enemySkillDialogue = append(enemySkillDialogue, b.generateMessageForUsedDialogueSkill(*b.Enemy, *b.Player, enemySkill, enemyRoll, enemySecondaryRoll)...)
	}

	if turn.TurnInitiative == Player {
		switch turn.PlayerSkillUsed.SkillName {
		case "draw":
			turn.PlayerMessage = playerSkillDialogue
			turn.EnemyMessage = []string{}
		default:
			turn.PlayerMessage = playerSkillDialogue
			turn.EnemyMessage = enemySkillDialogue
		}
	}

	if turn.TurnInitiative == Enemy {
		switch turn.EnemySkillUsed.SkillName {
		case "draw":
			turn.EnemyMessage = enemySkillDialogue
			turn.PlayerMessage = []string{}
		default:
			turn.PlayerMessage = playerSkillDialogue
			turn.EnemyMessage = enemySkillDialogue
		}
	}

	//apply buffs from enemyBattleSprite based on roll results

}

func (b *Battle) DrawFunction(user *battleStats.Character, opponent *battleStats.Character, battleInitiative bool) []string {

	drawSkillDialogue := b.generateDrawMessage(user, opponent, battleInitiative)

	if battleInitiative {
		b.nextTurnInitiative = Player
	} else {
		b.nextTurnInitiative = Enemy
	}
	return drawSkillDialogue
}

func (b *Battle) EnactEffects(skill battleStats.Skill, user *battleStats.Character, opponent *battleStats.Character, roll bool, secondaryRoll bool) {
	b.Tension = b.Tension + skill.Tension
	SkillEffectOne := skill.Effects[0]
	var skillEffectTwo battleStats.Effect

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

func (b *Battle) generateDrawMessage(turnTaker *battleStats.Character, opponent *battleStats.Character, battleInitiative bool) (drawMessage []string) {

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

func (b *Battle) generateMessageForUsedDialogueSkill(turnTaker battleStats.Character, opponent battleStats.Character, skill battleStats.Skill, roll bool, secondaryRoll bool) (message []string) {

	output := fmt.Sprintf("%s uses %s", turnTaker.Name, skill.SkillName)

	message = append(message, output)

	dialogue := GetSkillDialogue(turnTaker.Name, skill.SkillName, roll)
	response := GetResponse(opponent.Name, skill.SkillName, roll)

	message = append(message, dialogue)
	message = append(message, response)

	skillName := CapitalizeWord(skill.SkillName)
	name := CapitalizeWord(turnTaker.Name)

	if !roll {
		message = append(message, fmt.Sprintf("%s is ineffective", skillName))
		if secondaryRoll && skill.Effects[0].On == "unsuccessfulSelf" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", name, skill.Effects[1].Stat, skill.Effects[1].Amount))
		}
	}

	if roll {
		message = append(message, fmt.Sprintf("%s is effective", skillName))

		if skill.Effects[0].On == "self" {
			message = append(message, fmt.Sprintf("%s's %s increased by %d", name, skill.Effects[0].Stat, skill.Effects[0].Amount))
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

func (b *Battle) TakeCombatTurn(playerSkill battleStats.Skill) {
	fmt.Printf("battle.go:369: entering combat loop\n")
	if b.turnInitiative != b.nextTurnInitiative {
		fmt.Printf("battle.go:371: changing turn initiative from:%d to:%d\n", b.turnInitiative, b.nextTurnInitiative)
		b.UpdateInitiative(b.nextTurnInitiative)
	}
	b.incrementTurn()
	b.Turns[b.Turn] = &Turn{}
	turn := b.Turns[b.Turn]
	turn.PlayerEventTriggered = false
	turn.PlayerEffectsTriggered = false
	turn.EnemyEventTriggered = false
	turn.PlayerTurnCompleted = false
	turn.EnemyTurnCompleted = false
	turn.EnemyWeakness = false
	turn.PlayerWeakness = false

	if playerSkill.SkillName == "focusedShot" {
		b.nextTurnInitiative = Enemy
	}

	enemySkill, err := EnemyChooseSkill(*b, b.Enemy.CombatSkills)

	if err != nil {
		log.Fatal(err)
	}

	turn.EnemySkillUsed = enemySkill
	turn.PlayerSkillUsed = playerSkill

	if enemySkill.SkillName == "focusedShot" {
		b.nextTurnInitiative = Player
	}

	eOneTurnBuffAmount := 0
	eTurnAmmo := b.EnemyAmmo
	var eAffectedStat battleStats.Stat
	for _, effect := range enemySkill.Effects {

		if effect.EffectType == "buff" {
			if effect.On == "self" {

				eAffectedStat, err = battleStats.StringToStat(effect.Stat)
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
					damage := Shoot(b.Enemy.DisplayStat(battleStats.Accuracy), b.Enemy.DisplayStat(battleStats.Fear), b.Enemy.DisplayStat(battleStats.Anger), effect.SuccessPer, effect.DamageRange)
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

	var pAffectedStat battleStats.Stat
	for _, effect := range playerSkill.Effects {

		if effect.EffectType == "buff" {
			if effect.On == "self" {
				pAffectedStat, err = battleStats.StringToStat(effect.Stat)
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
					damage := Shoot(b.Player.DisplayStat(battleStats.Accuracy), b.Player.DisplayStat(battleStats.Fear), b.Player.DisplayStat(battleStats.Anger), effect.SuccessPer, effect.DamageRange)
					turn.DamageToEnemy = append(turn.DamageToEnemy, damage)
				} else {
					println("out of ammo")
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
	turn.PlayerMessage = playerMessage
	turn.EnemyMessage = enemyMessage

}
