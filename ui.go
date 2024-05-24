/*
@author: sk
@date: 2023/2/4
*/
package main

import (
	"GameBase2/config"
	"GameBase2/utils"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
)

//=============CardUI================

type CardUI struct { // 手牌 等显示出的UI  非实例  显示需要借助  object
	Card    *Card
	Pos     complex128
	player  bool // 是否是玩家的卡  展示手牌时需要特殊处理
	cardImg *ebiten.Image
	check   bool
}

func (u *CardUI) GetMin() complex128 {
	return u.Pos
}

func (u *CardUI) GetMax() complex128 {
	return u.Pos + CardSize
}

func (u *CardUI) CollisionPoint(pos complex128) bool {
	return utils.PointCollision(u, pos)
}

func NewCardUI(card *Card, player bool) *CardUI {
	res := &CardUI{Card: card, Pos: 0, player: player, cardImg: ebiten.NewImage(75, 110)}
	res.UpdateImg()
	return res
}

func (u *CardUI) Draw(screen *ebiten.Image) {
	if u.Card.Place == PlaceField && u.Card.Defense {
		utils.DrawAngleImage(screen, u.cardImg, u.Pos, CardSize/2, math.Pi/2)
	} else {
		utils.DrawImage(screen, u.cardImg, u.Pos)
	}
	if u.Card.Place == PlaceField {
		u.Card.DrawMark(u.Pos, screen)
	}
	if u.Card.Data.CardType == CardTypeMonster && u.Card.Place == PlaceField {
		utils.DrawAnchorText(screen, fmt.Sprintf("%d/%d", u.Card.Atk, u.Card.Def),
			u.Pos+CardSize-75/2, 0.5+0, Font24, colornames.White)
	}
	if u.check {
		utils.StrokeCircle(screen, u.Pos+CardSize/2, 25, 12, 2, config.ColorGreen)
	}
}

func (u *CardUI) UpdateImg() { // 改变状态后 记得更新图片
	if u.Card.Place == PlaceHand && !u.player {
		DrawCardBack(u.cardImg) // 暂时只可能在 这两个地方
	} else if u.Card.Place == PlaceField && u.Card.Flip {
		DrawCardBack(u.cardImg)
	} else {
		DrawCard(u.cardImg, u.Card)
	}
}

func (u *CardUI) UnSelect() {
	u.check = false
}

func (u *CardUI) Select() {
	u.check = true
}

//===========ButtonUI===========

type ButtonUI struct {
	Pos, Size complex128 // 4 边距
	show      string
}

func (b *ButtonUI) GetMin() complex128 {
	return b.Pos
}

func (b *ButtonUI) GetMax() complex128 {
	return b.Pos + b.Size
}

func (b *ButtonUI) CollisionPoint(pos complex128) bool {
	return utils.PointCollision(b, pos)
}

func NewButtonUI(show string) *ButtonUI {
	res := &ButtonUI{show: show}
	bound := text.BoundString(Font36, show)
	res.Size = utils.Int2Vector(bound.Dx()+8, bound.Dy()+8)
	return res
}

func (b *ButtonUI) Draw(screen *ebiten.Image) {
	utils.DrawRect(screen, b.Pos, b.Size, 1, config.ColorBlue, config.ColorWhite)
	utils.DrawAnchorText(screen, b.show, b.Pos+b.Size/2, 0.5+0.5i, Font36, colornames.White)
}
