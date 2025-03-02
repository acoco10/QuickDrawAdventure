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
	Damage           []int
	Ammo             int
}

func (c *CharacterTurnData) UpdateAmmo() {

	if c.SkillUsed.SkillName == "reload" {
		c.Ammo = 6
	}

	for _, effect := range c.SkillUsed.Effects {
		c.Ammo -= effect.NShots
	}
}
