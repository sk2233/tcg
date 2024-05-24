/*
@author: sk
@date: 2023/2/11
*/
package main

var (
	canSelectFactories = make(map[string]CanSelectFactory)
)

func CreateCanSelect(name string, player *Player, card *Card) ICanSelect {
	return canSelectFactories[name](player, card)
}

func init() {
	canSelectFactories["攻击源"] = createAttackSrcCanSelect
	canSelectFactories["嘉然-被攻击"] = createDianaAttackTar
}

func createDianaAttackTar(player *Player, card *Card) ICanSelect {
	return &DianaAttackTar{}
}

type DianaAttackTar struct {
}

func (d *DianaAttackTar) CanSelect(param *Param) bool {
	if param.EventType != EventTypeSelectTar {
		return true
	}
	return param.Card.Data.Name != "狗"
}

//=================AttackSrcCanSelect====================

func createAttackSrcCanSelect(player *Player, card *Card) ICanSelect {
	return &AttackSrcCanSelect{card: card}
}

type AttackSrcCanSelect struct {
	card *Card
}

func (a *AttackSrcCanSelect) CanSelect(param *Param) bool {
	if param.EventType != EventTypeSelectSrc {
		return true
	}
	return a.card.AttackNum > 0 && !a.card.Defense
}
