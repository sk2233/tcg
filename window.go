/*
@author: sk
@date: 2023/2/5
*/
package main

import (
	"GameBase2/config"
	"GameBase2/object"
	"GameBase2/utils"
	R "tcg/res"

	"github.com/hajimehoshi/ebiten/v2"
)

//====================ButtonGroup====================

type ButtonGroup struct {
	*object.PointObject
	btns   []*ButtonUI
	anchor complex128
	gap    float64
	Die    bool
}

func (b *ButtonGroup) Order() int {
	return 22
}

func (b *ButtonGroup) Draw(screen *ebiten.Image) {
	for i := 0; i < len(b.btns); i++ {
		b.btns[i].Draw(screen)
	}
}

func (b *ButtonGroup) IsDie() bool {
	return b.Die
}

func (b *ButtonGroup) SetPos(pos complex128) { // 以改点 为整体的  0.5+1i
	b.Pos = pos
	b.updateBtns()
}

func (b *ButtonGroup) updateBtns() {
	if len(b.btns) <= 0 {
		return
	}
	w := real(b.btns[0].Size)
	h := imag(b.btns[0].Size)
	for i := 1; i < len(b.btns); i++ {
		w += real(b.btns[i].Size) + b.gap
	}
	x := real(b.Pos) - real(b.anchor)*w
	y := imag(b.Pos) - imag(b.anchor)*h
	for i := 0; i < len(b.btns); i++ {
		b.btns[i].Pos = complex(x, y)
		x += real(b.btns[i].Size) + b.gap
	}
}

func (b *ButtonGroup) Replace(btns ...string) {
	b.btns = make([]*ButtonUI, 0)
	for i := 0; i < len(btns); i++ {
		b.btns = append(b.btns, NewButtonUI(btns[i]))
	}
	b.updateBtns()
}

func (b *ButtonGroup) ClickBtn(pos complex128) int {
	for i := 0; i < len(b.btns); i++ {
		if b.btns[i].CollisionPoint(pos) {
			return i
		}
	}
	return -1
}

func NewButtonGroup(anchor complex128, gap float64) *ButtonGroup {
	res := &ButtonGroup{anchor: anchor, gap: gap, btns: make([]*ButtonUI, 0), Die: false}
	res.PointObject = object.NewPointObject()
	utils.AddToLayer(R.LAYER.UI, res) // 自动添加
	return res
}

//================CardShow================

type CardShow struct {
	*object.PointObject
	cardUIs []*CardUI
	Die     bool
	size    complex128
}

func (c *CardShow) CollisionPoint(pos complex128) bool {
	return utils.PointCollision(c, pos)
}

func (c *CardShow) GetMin() complex128 {
	return c.Pos
}

func (c *CardShow) GetMax() complex128 {
	return c.Pos + c.size
}

func (c *CardShow) IsDie() bool {
	return c.Die
}

func NewCardShow(cards []*Card) *CardShow { // 使用 8 * 4的范围显示卡片
	res := &CardShow{cardUIs: make([]*CardUI, 0), Die: false, size: complex(8*75+9*24, 4*110+5*24)}
	res.PointObject = object.NewPointObject()
	res.Pos = (complex(1280-330, 720)-res.size)/2 + 330 // 不要挡住 显示
	for i := 0; i < len(cards); i++ {
		res.cardUIs = append(res.cardUIs, NewCardUI(cards[i], true))
		res.cardUIs[i].Pos = res.Pos + utils.Int2Vector(24+(i%8)*(24+75), 24+(i/8)*(24+110))
	}
	utils.AddToLayer(R.LAYER.UI, res)
	return res
}

func (c *CardShow) Draw(screen *ebiten.Image) {
	utils.DrawRect(screen, c.Pos, c.size, 2, Color71_159_144, config.ColorWhite)
	for i := 0; i < len(c.cardUIs); i++ {
		c.cardUIs[i].Draw(screen)
	}
}

func (c *CardShow) ClickCard(pos complex128) *CardUI {
	for i := 0; i < len(c.cardUIs); i++ {
		if c.cardUIs[i].CollisionPoint(pos) {
			return c.cardUIs[i]
		}
	}
	return nil
}
