package battleStats

import (
	"fmt"
)

type Character struct {
	Name           string
	Stats          map[Stat]int
	baselineStats  map[Stat]int
	CombatSkills   map[string]Skill
	DialogueSkills map[string]Skill
	Weakness       Stat
}

type StatusEffect struct {
	Duration     int
	AffectedStat Stat
	Amount       int
}

type Stat uint8

const (
	Accuracy Stat = iota
	Health
	Anger
	Fear
	DrawSpeed
	TensionThreshold
)

func (pc *Character) UpdateCharHealth(change int) {
	pc.Stats[Health] += change
	if pc.Stats[Health] > pc.baselineStats[Health] {
		pc.Stats[Health] = pc.baselineStats[Health]
	}
}

func (pc *Character) UpdateCharAnger(change int) {

	pc.Stats[Anger] += change

	if pc.Stats[Anger] < 0 {
		pc.Stats[Anger] = 0
	}
}

func (pc *Character) UpdateCharFear(change int) {
	pc.Stats[Fear] += change

	if pc.Stats[Fear] < 0 {
		pc.Stats[Fear] = 0
	}

}

func (pc *Character) UpdateCharAccuracy(change int) {
	pc.Stats[Accuracy] += change
	if pc.Stats[Accuracy] < 0 {
		pc.Stats[Accuracy] = 0
	}
}

func (pc *Character) UpdateCharDrawSpeed(change int) {
	pc.Stats[DrawSpeed] += change
	if pc.Stats[DrawSpeed] < 0 {
		pc.Stats[DrawSpeed] = 0
	}
}

func (pc *Character) DisplayStat(stat Stat) int {
	return pc.Stats[stat]
}

func (pc *Character) DisplayStats() map[Stat]int {
	return pc.Stats
}
func (pc *Character) UpdateStat(stat Stat, amt int) {
	if stat == Fear {
		pc.UpdateCharFear(amt)
	}
	if stat == Accuracy {
		pc.UpdateCharAccuracy(amt)
	}
	if stat == Anger {
		pc.UpdateCharAnger(amt)
	}
	if stat == DrawSpeed {
		pc.UpdateCharDrawSpeed(amt)
	}
}

func (pc *Character) ResetStatusStats() {
	pc.Stats[Anger] = pc.baselineStats[Anger]
	pc.Stats[Fear] = pc.baselineStats[Fear]
	pc.Stats[Accuracy] = pc.baselineStats[Accuracy]
	pc.Stats[DrawSpeed] = pc.baselineStats[DrawSpeed]
}

func (pc *Character) ResetHealth() {
	pc.Stats[Health] = pc.baselineStats[Health]
}

func NewCharacter(name string, stats map[string]int, combatSkills map[string]Skill, dialogueSkills map[string]Skill, Weakness string) Character {
	var weakness Stat
	charStats := map[Stat]int{}
	for key, stat := range stats {
		if key == "health" {
			charStats[Health] = stat
		}
		if key == "accuracy" {
			charStats[Accuracy] = stat
		}
		if key == "anger" {
			charStats[Anger] = stat
		}
		if key == "fear" {
			charStats[Fear] = stat
		}
		if key == "drawSpeed" {
			charStats[DrawSpeed] = stat
		}
		if key == "tensionThreshold" {
			charStats[TensionThreshold] = stat
		}
	}
	if Weakness == "fear" {
		weakness = Fear
	}
	if Weakness == "anger" {
		weakness = Anger
	}

	return Character{
		Name:           name,
		Stats:          charStats,
		baselineStats:  charStats,
		CombatSkills:   combatSkills,
		DialogueSkills: dialogueSkills,
		Weakness:       weakness,
	}
}

func StringToStat(s string) (Stat, error) {
	if s == "health" {
		return Health, nil
	}
	if s == "accuracy" {
		return Accuracy, nil
	}
	if s == "anger" {
		return Anger, nil
	}
	if s == "fear" {
		return Fear, nil
	}
	if s == "DrawSpeed" {
		return DrawSpeed, nil
	}
	if s == "TensionThreshold" {
		return TensionThreshold, nil
	}
	return 0, fmt.Errorf("not a valid stat: %s", s)
}

func StatToString(s Stat) (string, error) {
	if s == Health {
		return "health", nil
	}
	if s == Accuracy {
		return "accuracy", nil
	}
	if s == Anger {
		return "anger", nil
	}
	if s == Fear {
		return "fear", nil
	}
	if s == DrawSpeed {
		return "draw speed", nil
	}
	if s == TensionThreshold {
		return "tension threshold", nil
	}
	return "", fmt.Errorf("not a valid stat")
}
