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
	R "tcg/res"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.CARD_AREA, createCardArea)
}

func createCardArea(data *model.ObjectData) model.IObject {
	res := &CardArea{Monster: utils.GetBool(data.Properties, R.PROP.MONSTER, false)}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	return res
}

// CardArea 放置卡牌的  怪兽区/魔陷区
type CardArea struct {
	*object.PointObject
	cardUI  *CardUI // 用来绘制 显示
	Monster bool
}

func (c *CardArea) GetMin() complex128 {
	return c.Pos
}

func (c *CardArea) GetMax() complex128 {
	return c.Pos + CardSize
}

func (c *CardArea) Draw(screen *ebiten.Image) {
	if c.Monster {
		utils.StrokeRect(screen, c.Pos, CardSize, 2, Color255_93_0)
		utils.DrawAnchorText(screen, "怪兽区", c.Pos+CardSize/2, 0.5+0.5i, Font36, colornames.White)
	} else {
		utils.StrokeRect(screen, c.Pos, CardSize, 2, Color63_194_96)
		utils.DrawAnchorText(screen, "魔陷区", c.Pos+CardSize/2, 0.5+0.5i, Font36, colornames.White)
	}
	if c.cardUI != nil {
		c.cardUI.Draw(screen)
	}
}

func (c *CardArea) ClickCard(pos complex128) *CardUI {
	if c.cardUI == nil || !c.cardUI.CollisionPoint(pos) {
		return nil
	}
	return c.cardUI
}

func (c *CardArea) GetHandles(param *Param) []IHandle {
	if c.cardUI == nil {
		return make([]IHandle, 0)
	}
	return c.cardUI.Card.GetHandles(param)
}

func (c *CardArea) SetCardUI(cardUI *CardUI) {
	cardUI.Card.Place = PlaceField
	cardUI.Pos = c.Pos
	c.cardUI = cardUI
	cardUI.UpdateImg()
}

func (c *CardArea) GetCardUI() *CardUI {
	return c.cardUI
}

func (c *CardArea) ClearCardUI() {
	c.cardUI = nil
}
