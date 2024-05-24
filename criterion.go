/*
@author: sk
@date: 2023/2/11
*/
package main

func And(filters ...CardFilter) CardFilter {
	return func(card *Card) bool {
		for i := 0; i < len(filters); i++ {
			if !filters[i](card) {
				return false
			}
		}
		return true
	}
}

func Or(filters ...CardFilter) CardFilter {
	return func(card *Card) bool {
		for i := 0; i < len(filters); i++ {
			if filters[i](card) {
				return true
			}
		}
		return false
	}
}

func Not(filter CardFilter) CardFilter {
	return func(card *Card) bool {
		return !filter(card)
	}
}

func IsMonster(card *Card) bool {
	return card.Data.CardType == CardTypeMonster
}

func IsCommonMonster(card *Card) bool {
	return IsMonster(card) && card.Data.MonsterType == MonsterTypeCommon
}

func CanAttackSrc(card *Card) bool {
	return card.CanSelect(&Param{EventType: EventTypeSelectSrc})
}

func IsField(card *Card) bool {
	return card.Place == PlaceField
}

func IsHand(card *Card) bool {
	return card.Place == PlaceHand
}

func IsAny(card *Card) bool {
	return true
}

func NameEq(name string) CardFilter {
	return func(card *Card) bool {
		return card.Data.Name == name
	}
}

func LevelEq(level int) CardFilter {
	return func(card *Card) bool {
		return card.Data.Level == level
	}
}

func RaceEq(race int) CardFilter {
	return func(card *Card) bool {
		return card.Race == race
	}
}

func AtkGe(atk int) CardFilter {
	return func(card *Card) bool {
		return card.Data.Atk >= atk
	}
}

func DefLe(def int) CardFilter {
	return func(card *Card) bool {
		return card.Data.Def <= def
	}
}

func HasField(field int) CardFilter {
	return func(card *Card) bool {
		return card.Data.Fields&field > 0
	}
}

func Except(exceptCard *Card) CardFilter {
	return func(card *Card) bool {
		return exceptCard != card
	}
}
