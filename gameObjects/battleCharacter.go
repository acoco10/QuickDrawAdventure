package gameObjects

import (
	"github.com/acoco10/QuickDrawAdventure/battleStats"
)

type BattleCharacter struct {
	batteSprite *BattleSprite
	charStats   *battleStats.CharacterData
	charType    CharType
}
