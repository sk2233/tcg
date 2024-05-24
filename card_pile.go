/*
@author: sk
@date: 2023/2/4
*/
package main

import (
	"GameBase2/config"
	"GameBase2/factory"
	"GameBase2/model"
	"GameBase2/object"
	"GameBase2/utils"
	"strconv"
	"strings"
	R "tcg/res"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.CARD_PILE, createCardPile)
}

func createCardPile(data *model.ObjectData) model.IObject {
	res := &CardPile{cards: make([]*Card, 0), cardImg: ebiten.NewImage(75, 110),
		Type: utils.GetInt(data.Properties, R.PROP.TYPE, CardPileTypeNone)}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	return res
}

// CardPile 除外/墓地/额外/卡组 区域  正面的是 最后一张在上面  反面的 是第一张在上面 方便使用
type CardPile struct {
	*object.PointObject
	cards   []*Card
	cardImg *ebiten.Image // 新加入卡牌时进行变化
	Type    int
}

func (c *CardPile) GetMin() complex128 {
	return c.Pos
}

func (c *CardPile) GetMax() complex128 {
	return c.Pos + CardSize
}

func (c *CardPile) CollisionPoint(pos complex128) bool {
	return utils.PointCollision(c, pos)
}

func (c *CardPile) Init() {
	DrawCardBack(c.cardImg) // 默认绘制为卡背
}

func (c *CardPile) Draw(screen *ebiten.Image) {
	if len(c.cards) > 0 {
		utils.DrawImage(screen, c.cardImg, c.Pos)
		utils.DrawAnchorText(screen, strconv.Itoa(len(c.cards)), c.Pos+CardSize/2, 0.5+0.5i, Font72, colornames.Red)
	} else {
		utils.StrokeRect(screen, c.Pos, CardSize, 2, c.getColor())
		utils.DrawAnchorText(screen, c.getName(), c.Pos+CardSize/2, 0.5+0.5i, Font36, colornames.White)
	}
}

func (c *CardPile) isFlip() bool {
	return c.Type == CardPileTypeDeck || c.Type == CardPileTypeExtra
}

var (
	colors = []*model.Color{Color65_147_217, Color163_167_166, Color163_167_166, Color255_93_0}
	names  = strings.Split("除外区,墓地,额外,卡组", ",")
)

func (c *CardPile) getColor() *model.Color {
	return colors[c.Type]
}

func (c *CardPile) getName() string {
	return names[c.Type]
}

func (c *CardPile) InitCards(player *Player, datas []*CardData) {
	for i := 0; i < len(datas); i++ {
		c.cards = append(c.cards, NewCard(datas[i], player, c.Type)) // 要求 type与place 必须对应
	}
	c.updateCard()
}

func (c *CardPile) updateCard() { // 更新顶端显示
	if c.isFlip() {
		return
	}
	if len(c.cards) > 0 {
		DrawCard(c.cardImg, c.cards[len(c.cards)-1])
	} else {
		DrawCardBack(c.cardImg)
	}
}

func (c *CardPile) GetCards(num int) []*Card {
	if len(c.cards) < num {
		return nil
	}
	res := c.cards[:num]
	c.cards = c.cards[num:]
	c.updateCard()
	return res
}

func (c *CardPile) GetAllCards() []*Card {
	return c.cards
}

func (c *CardPile) GetHandles(param *Param) []IHandle {
	res := make([]IHandle, 0)
	for i := 0; i < len(c.cards); i++ {
		res = append(res, c.cards[i].GetHandles(param)...)
	}
	return res
}

func (c *CardPile) AddCard(card *Card) {
	c.cards = append(c.cards, card)
	c.updateCard()
}

func (c *CardPile) RemoveCard(card *Card) {
	for i := 0; i < len(c.cards); i++ {
		if c.cards[i] == card {
			c.cards = append(c.cards[:i], c.cards[i+1:]...)
			c.updateCard()
			return
		}
	}
}
