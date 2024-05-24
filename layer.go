/*
@author: sk
@date: 2023/2/5
*/
package main

import (
	"GameBase2/layer"
	"GameBase2/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

type WinUILayer struct {
	*layer.UILayer
	win, reason string
}

func NewWinUILayer(player bool, reason string) *WinUILayer {
	res := &WinUILayer{reason: reason}
	res.UILayer = layer.NewUILayer()
	if player {
		res.win = "玩家获得胜利"
	} else {
		res.win = "敌人获得胜利"
	}
	ActionManager.GameEnd = true
	return res
}

func (u *WinUILayer) Draw(screen *ebiten.Image) {
	utils.FillRect(screen, 300i, 1280+120i, Color0_0_0_127)
	utils.DrawAnchorText(screen, u.win, complex(640, 340), 0.5+0.5i, Font72, colornames.White)
	utils.DrawAnchorText(screen, u.reason, complex(640, 380), 0.5+0.5i, Font36, colornames.White)
}
