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
	config.ObjectFactory.RegisterPointFactory(R.CLASS.ROUND, createRound)
}

func createRound(data *model.ObjectData) model.IObject {
	res := &round{num: 1}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	Round = res
	return res
}

type round struct {
	*object.PointObject
	num int
}

func (r *round) Draw(screen *ebiten.Image) {
	utils.DrawAnchorText(screen, fmt.Sprintf("第%d轮", r.num), r.Pos, 0, Font36, colornames.White)
}

func (r *round) AddRound() {
	r.num++
}
