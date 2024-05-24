/*
@author: sk
@date: 2023/2/4
*/
package main

import (
	"GameBase2/config"
	"GameBase2/model"
	"GameBase2/utils"
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func Color(r, g, b float64) *model.Color {
	return &model.Color{R: r / 255, G: g / 255, B: b / 255, A: 1}
}

func AlphaColor(r, g, b, a float64) *model.Color {
	return &model.Color{R: r / 255, G: g / 255, B: b / 255, A: a / 255}
}

// DrawCard 绘制卡面
func DrawCard(screen *ebiten.Image, card *Card) {
	data := card.Data
	switch data.CardType {
	case CardTypeMonster:
		utils.DrawRect(screen, 0, CardSize, 1, Color163_101_70, config.ColorBlack)
		utils.DrawAnchorText(screen, data.Name, complex(37.5, 20), 0.5+0.5i, Font21, colornames.Black)
		utils.DrawAnchorText(screen, fmt.Sprintf("LV:%d,%s族", data.Level, GetRace(card.Race)),
			complex(37.5, 40), 0.5+0.5i, Font21, colornames.Black)
		utils.DrawAnchorText(screen, fmt.Sprintf("ATK:%d", data.Atk), complex(37.5, 70),
			0.5+0.5i, Font21, colornames.Black)
		utils.DrawAnchorText(screen, fmt.Sprintf("DEF:%d", data.Def), complex(37.5, 90),
			0.5+0.5i, Font21, colornames.Black)
	case CardTypeMagic:
		utils.DrawRect(screen, 0, CardSize, 1, Color61_138_129, config.ColorBlack)
		utils.DrawAnchorText(screen, data.Name, complex(37.5, 20), 0.5+0.5i, Font21, colornames.White)
		utils.DrawAnchorText(screen, GetMagicType(data.MagicType), complex(37.5, 90), 0.5+0.5i, Font21,
			colornames.Black)
	case CardTypeTrap:
		utils.DrawRect(screen, 0, CardSize, 1, Color160_72_130, config.ColorBlack)
		utils.DrawAnchorText(screen, data.Name, complex(37.5, 20), 0.5+0.5i, Font21, colornames.White)
		utils.DrawAnchorText(screen, GetTrapType(data.TrapType), complex(37.5, 90), 0.5+0.5i, Font21,
			colornames.Black)
	}
}

var (
	magicTypes   = strings.Split("通常,永续,装备,速攻", ",")
	trapTypes    = strings.Split("通常,永续", ",")
	races        = strings.Split("天使,不死,战士,植物,机械,龙,小黑子", ",")
	natures      = strings.Split("风,水,炎,地,光,神", ",")
	monsterTypes = strings.Split("通常,效果", ",")
	fields       = strings.Split("青眼,字段2,字段3,字段4,字段5,青眼白龙,攻击状态,结束乐队", ",")
)

func GetRace(race int) any {
	return races[race]
}

func GetMagicType(magicType int) string {
	return magicTypes[magicType]
}

func GetNature(nature int) string {
	return natures[nature]
}

func GetFields(field int) string {
	buf := strings.Builder{}
	for i := 0; i < len(fields); i++ {
		if field&(1<<i) > 0 {
			if buf.Len() > 0 {
				buf.WriteRune(',')
			}
			buf.WriteString(fields[i])
		}
	}
	return buf.String()
}

func GetMonsterType(monsterType int) string {
	return monsterTypes[monsterType]
}

func GetTrapType(trapType int) string {
	return trapTypes[trapType]
}

// DrawCardBack 绘制卡背
func DrawCardBack(screen *ebiten.Image) {
	utils.DrawRect(screen, 0, CardSize, 1, config.ColorBlack, Color255_93_0)
	utils.DrawAnchorText(screen, "卡背", CardSize/2, 0.5+0.5i, Font36, colornames.White)
}

func Repeat[T any](value T, count int) []T {
	res := make([]T, 0)
	for i := 0; i < count; i++ {
		res = append(res, value)
	}
	return res
}

func GetCostNum(level int) int {
	if level <= 4 {
		return 0
	}
	if level <= 6 {
		return 1
	}
	return 2
}

func InvokeTop(src any, top0 bool) {
	if top, ok := src.(ITop); ok {
		top.SetTop(top0)
	}
}

func InvokeDrawMark(src any, pos complex128, screen *ebiten.Image) {
	if drawMark, ok := src.(IDrawMark); ok {
		drawMark.DrawMark(pos, screen)
	}
}

func GetCardPos(card *Card, player *Player) complex128 { // 获取 手牌  或 场上 卡牌的中心位置  卡牌为nil 获取玩家位置
	if card == nil {
		return complex((1280-330)/2+330, imag(player.Pos)+110/2)
	}
	return player.GetCardUI(card).Pos + CardSize/2
}
