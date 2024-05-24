/*
@author: sk
@date: 2023/2/5
*/
package main

import (
	"GameBase2/model"

	"github.com/hajimehoshi/ebiten/v2"
)

type CardFilter func(card *Card) bool

type HandleFactory func(player *Player, card *Card) IHandle

type CanSelectFactory func(player *Player, card *Card) ICanSelect

type ICanSelect interface {
	CanSelect(param *Param) bool
}

type IHandle interface {
	CanHandle(param *Param) bool
	CreateAction(param *Param) IAction
}

type INameHandle interface {
	model.IName
	IHandle
}

type IAction interface {
	Action()
}

type IPhase interface {
	model.IInit
	IAction
}

type ITop interface {
	SetTop(top bool)
}

type IDrawMark interface {
	DrawMark(pos complex128, screen *ebiten.Image)
}
