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
	"image/color"
	R "tcg/res"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.TIP, createTip)
}

func createTip(data *model.ObjectData) model.IObject {
	res := &tip{msgs: make([]*Msg, 0), lines: make([]*Line, 0)}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	Tip = res
	return res
}

type Msg struct {
	Text   string
	Color  color.Color
	Offset float64 //偏移出  360 消失
}

type Line struct {
	Start, End complex128
	Color      *model.Color
	Time       int
}

// 各种触发操作都可以使用的提示信息  来说明确实发生了(类似日志)
type tip struct {
	*object.PointObject
	msgs  []*Msg
	lines []*Line
}

func (t *tip) Order() int {
	return 33
}

func (t *tip) Draw(screen *ebiten.Image) {
	for _, line := range t.lines {
		utils.DrawLine(screen, line.Start, line.End, 4, line.Color)
	}
	for _, msg := range t.msgs {
		utils.DrawAnchorText(screen, msg.Text, t.Pos-complex(0, msg.Offset), 0.5+0.5i, Font72, msg.Color)
	}
}

func (t *tip) Update() {
	for i := 0; i < len(t.msgs); i++ {
		t.msgs[i].Offset += 2
	}
	if len(t.msgs) > 0 && t.msgs[0].Offset > 360 { // 只有第一个有出界可能性
		t.msgs = t.msgs[1:]
	}
	for i := 0; i < len(t.lines); i++ {
		t.lines[i].Time--
	}
	for i := len(t.lines) - 1; i >= 0; i-- {
		if t.lines[i].Time <= 0 {
			t.lines = t.lines[i+1:]
			break
		}
	}
}

func (t *tip) AddTip(msg string) {
	t.AddColorTip(msg, colornames.Red)
}

func (t *tip) AddColorTip(msg string, clr color.Color) {
	t.msgs = append(t.msgs, &Msg{Text: msg, Color: clr, Offset: 0})
	for i := len(t.msgs) - 2; i >= 0; i-- { // 调整位置
		if t.msgs[i].Offset < t.msgs[i+1].Offset+TipInterval {
			t.msgs[i].Offset = t.msgs[i+1].Offset + TipInterval
		}
	}
	if t.msgs[0].Offset > 360 { // 只有第一个有出界可能性
		t.msgs = t.msgs[1:]
	}
}

func (t *tip) AddLine(start, end complex128, clr *model.Color) {
	t.lines = append(t.lines, &Line{Start: start, End: end, Color: clr, Time: 60})
}
