/*
@author: sk
@date: 2023/2/4
*/
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Card struct { // 卡牌效果 对象 主要处理卡牌的效果逻辑
	Data       *CardData
	handles    map[string]IHandle
	canSelects map[string]ICanSelect
	Race       int
	Place      int
	Flip       bool
	Defense    bool
	AttackNum  int  // 攻击 次数 默认一次
	HasAdjust  bool // 是否已经调整过了
	CanUse     bool // 陷阱卡专用  开始为false [设置] 效果 中 在回合开始是会设置为true
	Atk, Def   int  // 最终 有效的  临时 效果 直接 修改 这个 可能会出问题
}

func (c *Card) GetNameHandles(param *Param) []INameHandle {
	res := make([]INameHandle, 0)
	for _, handle := range c.handles {
		if nameHandle, ok := handle.(INameHandle); ok && nameHandle.CanHandle(param) {
			res = append(res, nameHandle)
		}
	}
	return res
}

func (c *Card) GetHandles(param *Param) []IHandle {
	res := make([]IHandle, 0)
	for _, handle := range c.handles {
		if handle.CanHandle(param) {
			res = append(res, handle)
		}
	}
	return res
}

func (c *Card) GetValue() int {
	if c.Defense {
		return c.Def
	}
	return c.Atk
}

// CanSelect 会在 一个大前提下 调用  默认直接 返回 true 即可  这里仅用于处理 特殊情况
func (c *Card) CanSelect(param *Param) bool {
	for _, canSelect := range c.canSelects {
		if !canSelect.CanSelect(param) {
			return false
		}
	}
	return true
}

func (c *Card) DrawMark(pos complex128, screen *ebiten.Image) {
	for _, handle := range c.handles {
		InvokeDrawMark(handle, pos, screen)
	}
}

func NewCard(data *CardData, player *Player, place int) *Card {
	res := &Card{Data: data, Place: place, Flip: false, Defense: false, handles: make(map[string]IHandle),
		canSelects: make(map[string]ICanSelect), Atk: data.Atk, Def: data.Def, Race: data.Race}
	for i := 0; i < len(data.Actions); i++ {
		res.handles[data.Actions[i]] = CreateHandle(data.Actions[i], player, res)
	}
	for i := 0; i < len(data.CanSelects); i++ {
		res.canSelects[data.CanSelects[i]] = CreateCanSelect(data.CanSelects[i], player, res)
	}
	return res
}
