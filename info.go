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
	"fmt"
	R "tcg/res"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.INFO, createInfo)
}

func createInfo(data *model.ObjectData) model.IObject {
	res := &info{card: ebiten.NewImage(330, 480)}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	Info = res
	return res
}

// 一边的卡片信息提示板
type info struct {
	*object.PointObject
	card *ebiten.Image
	desc string // 4边距
}

func (i *info) Init() {
	// 默认绘制卡背
	utils.StrokeRect(i.card, 2+2i, 330+480i-4-4i, 2, Color255_93_0)
	utils.DrawAnchorText(i.card, "卡背", (330+480i)/2, 0.5+0.5i, Font72, colornames.White)
	i.desc = "卡牌描述"
}

func (i *info) SetCard(card *Card) {
	data := card.Data
	_, i.desc = utils.WarpText(data.Desc, 330-8*2, Font36)
	switch data.CardType {
	case CardTypeMonster:
		utils.DrawRect(i.card, 0, 330+480i, 2, Color163_101_70, config.ColorBlack)
		utils.DrawAnchorText(i.card, data.Name, complex(165, 40), 0.5+0.5i, Font72, colornames.Black)
		utils.DrawAnchorText(i.card, fmt.Sprintf("LV:%d,%s怪兽", data.Level, GetMonsterType(data.MonsterType)),
			complex(165, 100), 0.5+0.5i, Font72, colornames.Black)
		utils.DrawAnchorText(i.card, fmt.Sprintf("%s属性,%s族", GetNature(data.Nature),
			GetRace(card.Race)), complex(165, 160), 0.5+0.5i, Font72, colornames.Black)
		utils.DrawAnchorText(i.card, "字段:"+GetFields(data.Fields), complex(165, 220), 0.5+0.5i, Font72, colornames.Black)
		utils.DrawAnchorText(i.card, fmt.Sprintf("ATK:%d", data.Atk), complex(165, 380),
			0.5+0.5i, Font72, colornames.Black)
		utils.DrawAnchorText(i.card, fmt.Sprintf("DEF:%d", data.Def), complex(165, 440),
			0.5+0.5i, Font72, colornames.Black)
	case CardTypeMagic:
		utils.DrawRect(i.card, 0, 330+480i, 2, Color61_138_129, config.ColorBlack)
		utils.DrawAnchorText(i.card, data.Name, complex(165, 40), 0.5+0.5i, Font72, colornames.White)
		utils.DrawAnchorText(i.card, "字段:"+GetFields(data.Fields), complex(165, 220), 0.5+0.5i, Font72, colornames.Black)
		utils.DrawAnchorText(i.card, GetMagicType(data.MagicType), complex(165, 440), 0.5+0.5i, Font72,
			colornames.Black)
	case CardTypeTrap:
		utils.DrawRect(i.card, 0, 330+480i, 2, Color160_72_130, config.ColorBlack)
		utils.DrawAnchorText(i.card, data.Name, complex(165, 40), 0.5+0.5i, Font72, colornames.White)
		utils.DrawAnchorText(i.card, "字段:"+GetFields(data.Fields), complex(165, 220), 0.5+0.5i, Font72, colornames.Black)
		utils.DrawAnchorText(i.card, GetTrapType(data.TrapType), complex(165, 440), 0.5+0.5i, Font72,
			colornames.Black)
	}
}

func (i *info) Draw(screen *ebiten.Image) {
	utils.DrawImage(screen, i.card, i.Pos)
	utils.StrokeRect(screen, 480i+2+2i, 330+240i-4-4i, 2, config.ColorWhite)
	utils.DrawAnchorText(screen, i.desc, 480i+8+8i, 0, Font36, colornames.White)
}
