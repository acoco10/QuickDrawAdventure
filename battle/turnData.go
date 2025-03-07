package battle

import "github.com/acoco10/QuickDrawAdventure/battleStats"

type CharacterTurnData struct {
	SkillUsed        battleStats.Skill
	WeaknessTargeted bool
	Index            int
	Message          []string
	EffectsTriggered bool
	Roll             bool
	SecondaryRoll    bool
	DamageOutput     []int
	EventTriggered   bool
	Completed        bool
	ComeBackEquipped bool
	Stunned          bool
}

type CharacterBattleData struct {
	*battleStats.CharacterData
	Ammo      int
	DrawBonus bool
	*CharacterTurnData
}

func (cb *CharacterBattleData) UpdateAmmo() {

	if cb.SkillUsed.SkillName == "reload" {
		cb.Ammo = 6
	}

	for _, effect := range cb.SkillUsed.Effects {
		cb.Ammo -= effect.NShots
		if cb.Ammo < 0 {
			cb.Ammo = 0
		}
	}
}
