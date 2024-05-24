/*
@author: sk
@date: 2023/2/4
*/
package main

import (
	"GameBase2/model"
	"GameBase2/utils"
	R "tcg/res"
)

//=================playerManager==================

type playerManager struct {
	Player, Enemy *Player
	player        bool
}

func NewPlayerManager() *playerManager {
	return &playerManager{}
}

func (p *playerManager) Init() {
	objLayer := utils.GetObjectLayer(R.LAYER.BASE)
	p.Player = objLayer.GetObject(R.OBJECT.PLAYER).(*Player)
	p.Enemy = objLayer.GetObject(R.OBJECT.ENEMY).(*Player)
	player := utils.RandomItem(p.Player, p.Enemy)
	p.player = player.player
	player.SkipPhase(PhaseAttack | PhaseDraw)
	ActionManager.SetCurrent(player)
}

func (p *playerManager) NextPlayer(player bool) {
	if p.player != player {
		Round.AddRound()
	}
	if player { // 交替
		ActionManager.SetCurrent(p.Enemy)
	} else {
		ActionManager.SetCurrent(p.Player)
	}
}

func (p *playerManager) GetAnother(player *Player) *Player {
	if player == p.Player {
		return p.Enemy
	}
	return p.Player
}

//======================actionManager=======================

type actionManager struct {
	actions *model.Stack[*ActionGroup]
	current *Player
	GameEnd bool
}

func NewActionManager() *actionManager {
	return &actionManager{actions: model.NewStack[*ActionGroup]()}
}

func (a *actionManager) Update() {
	if a.GameEnd { // 暂时这样  胜利应该使用阶段的
		return
	}
	if a.actions.IsEmpty() {
		a.current.Action()
	} else {
		a.actions.Peek().Action()
	}
}

func (a *actionManager) SetCurrent(player *Player) {
	a.current = player
	player.Prepare()
}

func (a *actionManager) PushAction(param *Param, handle ...IHandle) { // 一般主动 触发 的只有 一个 被动触发的可能有多个
	if a.actions.IsEmpty() {
		a.current.SetTop(false)
	} else {
		a.actions.Peek().SetTop(false)
	}
	a.actions.Push(NewActionGroup(param, handle...)) // 最终都按多个处理
}

func (a *actionManager) PopAction() {
	a.actions.Pop()
	if a.actions.IsEmpty() {
		a.current.SetTop(true)
	} else {
		a.actions.Peek().SetTop(true)
	}
}

func (a *actionManager) TriggerEvent(param *Param) { // 暂时 不管 效果 优先 级 有需要  排序一下即可
	handles := PlayerManager.Player.GetHandles(param)
	handles = append(handles, PlayerManager.Enemy.GetHandles(param)...)
	if len(handles) > 0 { // 大于0 才有 意义
		a.PushAction(param, handles...)
	}
}
