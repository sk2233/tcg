/*
@author: sk
@date: 2023/2/5
*/
package main

import (
	"GameBase2/config"
	"GameBase2/utils"
	"fmt"
	"math/rand"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//===================ActionGroup=======================

type ActionGroup struct {
	handles []IHandle
	index   int
	param   *Param
	current IAction
}

func (a *ActionGroup) Action() {
	if a.param.ActionEnd {
		a.index++
		a.prepare()
	} else {
		a.current.Action()
	}
}

func (a *ActionGroup) prepare() {
	if a.index < len(a.handles) && !a.param.EventEnd {
		a.current = a.handles[a.index].CreateAction(a.param)
		a.param.ActionEnd = false
	} else {
		ActionManager.PopAction()
	}
}

func (a *ActionGroup) SetTop(top bool) {
	InvokeTop(a.current, top)
}

func NewActionGroup(param *Param, handles ...IHandle) *ActionGroup {
	sort.Slice(handles, func(i, j int) bool { // 先 排序
		return utils.InvokeOrder(handles[i]) <= utils.InvokeOrder(handles[j])
	})
	res := &ActionGroup{handles: handles, index: 0, param: param}
	res.prepare()
	return res
}

//==================StepAction=================

type StepAction struct {
	Param *Param
	index int
	Step  []func()
}

func NewStepAction(param *Param) *StepAction {
	return &StepAction{Param: param, index: 0}
}

func (s *StepAction) Next() {
	s.index++
}

func (s *StepAction) End() {
	s.index = len(s.Step)
}

func (s *StepAction) Action() {
	if s.index < len(s.Step) {
		s.Step[s.index]()
	} else {
		s.Param.ActionEnd = true
	}
}

//=======================OnceAction=======================

type OnceAction struct {
	Param  *Param
	action func()
}

func (o *OnceAction) Action() {
	o.action()
	o.Param.ActionEnd = true
}

func NewOnceAction(param *Param, action func()) *OnceAction {
	return &OnceAction{Param: param, action: action}
}

//==================================

type SimpleFuncAction struct {
	*OnceAction
	action func(param *Param)
}

func (a *SimpleFuncAction) mainAction() {
	a.action(a.Param)
}

func NewSimpleFuncAction(param *Param, action func(*Param)) *SimpleFuncAction {
	res := &SimpleFuncAction{action: action}
	res.OnceAction = NewOnceAction(param, res.mainAction)
	return res
}

//=======================SelectAction========================

type SelectAction struct {
	param   *Param
	btns    *ButtonGroup
	actions []func(*Param)
}

func (s *SelectAction) SetTop(top bool) {
	s.btns.Enable = top
}

func NewSelectAction(param *Param, names []string, actions ...func(*Param)) *SelectAction {
	btns := NewButtonGroup(0.5+1i, 4)
	btns.SetPos((1280-330)/2 + 330 + (720i - 110i))
	btns.Replace(names...)
	return &SelectAction{param: param, actions: actions, btns: btns}
}

func (s *SelectAction) Action() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	index := s.btns.ClickBtn(utils.GetCursorPos())
	if index < 0 {
		return
	}
	s.actions[index](s.param)
	s.btns.Die = true
	s.param.ActionEnd = true
}

func EmptyFunc(*Param) {

}

//===================CardPileAction======================

func CreateCardPileAction(param *Param) IAction {
	res := &CardPileAction{param: param}
	res.cardShow = NewCardShow(param.CardPile.GetAllCards())
	res.actionBtns = NewButtonGroup(0.5+1i, 4)
	res.actionBtns.Enable = false
	return res
}

type CardPileAction struct {
	param       *Param
	cardShow    *CardShow
	actionBtns  *ButtonGroup
	nameHandles []INameHandle
	selectCard  *Card
}

func (c *CardPileAction) SetTop(top bool) {
	c.cardShow.Enable = top
	c.actionBtns.Enable = top
}

func (c *CardPileAction) Action() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	pos := utils.GetCursorPos()
	if c.actionBtns.Enable { // 处理按钮 选项
		c.actionBtns.Enable = false
		index := c.actionBtns.ClickBtn(pos)
		if index < 0 {
			return
		} // 这里CardPileType 与 事件是对应的
		ActionManager.PushAction(&Param{EventType: c.param.CardPile.Type, Player: c.param.Player, Card: c.selectCard},
			c.nameHandles[index])
		c.actionEnd()
		return
	}
	if !c.cardShow.CollisionPoint(pos) { // 处理关闭事件
		c.actionEnd() // 点击到外面关闭
		return
	}
	cardUI := c.cardShow.ClickCard(pos) // 处理 卡片 点击
	if cardUI != nil {
		Info.SetCard(cardUI.Card)
		c.nameHandles = cardUI.Card.GetNameHandles(&Param{EventType: cardUI.Card.Place, Player: c.param.Player,
			Card: cardUI.Card})
		if len(c.nameHandles) > 0 {
			c.selectCard = cardUI.Card
			btns := make([]string, 0)
			for i := 0; i < len(c.nameHandles); i++ {
				btns = append(btns, c.nameHandles[i].GetName())
			}
			c.actionBtns.Replace(btns...)
			c.actionBtns.Enable = true
			c.actionBtns.SetPos(cardUI.Pos + 75/2.0)
		} else {
			Tip.AddTip("该卡当前无法使用")
		}
		return
	}
}

func (c *CardPileAction) actionEnd() {
	c.param.ActionEnd = true
	c.cardShow.Die = true
	c.actionBtns.Die = true
}

//===================PickCardAction======================

func CreatePickCardAction(param *Param) IAction { // Pick 选择的 是 4大牌堆的   Select 选的 是 手牌 或 场上的牌
	res := &PickCardAction{param: param}
	res.cardShow = NewCardShow(param.CardPile.GetAllCards())
	res.okBtn = NewButtonGroup(0.5+0.5i, 0)
	res.okBtn.SetPos(1143 + 610i + CardSize/2)
	res.okBtn.Replace("Ok")
	return res
}

type PickCardAction struct {
	param      *Param
	cardShow   *CardShow
	okBtn      *ButtonGroup
	selectCard *CardUI
}

func (c *PickCardAction) SetTop(top bool) {
	c.cardShow.Enable = top
	c.okBtn.Enable = top
}

func (c *PickCardAction) Action() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	pos := utils.GetCursorPos()
	cardUI := c.cardShow.ClickCard(pos) // 处理 卡片 点击
	if cardUI != nil {
		Info.SetCard(cardUI.Card)
		if !c.param.CardFilter(cardUI.Card) {
			return
		}
		if c.selectCard != nil {
			c.selectCard.UnSelect()
			if c.selectCard == cardUI {
				c.selectCard = nil
				return
			}
		}
		c.selectCard = cardUI
		c.selectCard.Select()
		return
	}
	if c.okBtn.ClickBtn(pos) == 0 { // 确定选择
		c.param.ActionEnd = true
		c.cardShow.Die = true
		c.okBtn.Die = true
		if c.selectCard != nil {
			c.param.Card = c.selectCard.Card
		}
		return
	}
}

//=======================SelectCardAction============================

func CreateSelectCardAction(param *Param) IAction {
	res := &SelectCardAction{param: param}
	res.okBtn = NewButtonGroup(0.5+0.5i, 0)
	res.okBtn.SetPos(1143 + 610i + CardSize/2)
	res.okBtn.Replace("Ok")
	return res
}

type SelectCardAction struct { // 一张一张 选  选手卡  或场上的  卡
	param  *Param
	okBtn  *ButtonGroup
	cardUI *CardUI
}

func (s *SelectCardAction) SetTop(top bool) {
	s.okBtn.Enable = top
}

func (s *SelectCardAction) Action() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	pos := utils.GetCursorPos()
	cardUI := s.param.Player.ClickCard(pos) // 选择 卡牌
	if cardUI != nil {
		if !s.param.CardFilter(cardUI.Card) {
			return
		}
		if s.cardUI != nil {
			s.cardUI.UnSelect()
			if s.cardUI == cardUI {
				s.cardUI = nil
				return
			}
		}
		s.cardUI = cardUI
		s.cardUI.Select()
		return
	}
	if s.okBtn.ClickBtn(pos) == 0 { // 点击确定
		if s.cardUI != nil {
			s.cardUI.UnSelect()
			s.param.Card = s.cardUI.Card
		}
		s.param.ActionEnd = true
		s.okBtn.Die = true
	}
}

//=======================DestroyCardAction=====================

type DestroyCardAction struct {
	*StepAction
	selectParam  *Param
	filter       CardFilter
	cardPileType int
}

func (a *DestroyCardAction) selectStep() {
	a.selectParam = &Param{Player: PlayerManager.GetAnother(a.Param.Player), CardFilter: a.filter}
	ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *DestroyCardAction) mainStep() {
	if a.selectParam.Card == nil {
		Tip.AddTip("没有选择要被破坏的卡牌")
		a.End()
		return
	}
	player := PlayerManager.GetAnother(a.Param.Player)
	// 移除 破坏 对象
	player.RemoveCardUI(a.selectParam.Card)
	player.AddCardPile(a.selectParam.Card, a.cardPileType)
	a.Next()
}

func NewMoveCardAction(param *Param, filter CardFilter, cardPileType int) *DestroyCardAction {
	res := &DestroyCardAction{filter: filter, cardPileType: cardPileType}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//======================MonsterAttackAction============================

func CreateMonsterAttackAction(param *Param) IAction {
	res := &MonsterAttackAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.eventStep, res.mainStep, res.endStep}
	return res
}

type MonsterAttackAction struct {
	*StepAction
	beforeAttackParam, afterAttackParam *Param
}

func (a *MonsterAttackAction) eventStep() { // 指定 目标时
	a.Param.SrcCard.AttackNum--
	start := GetCardPos(a.Param.SrcCard, a.Param.Player)
	end := GetCardPos(a.Param.TarCard, PlayerManager.GetAnother(a.Param.Player))
	Tip.AddLine(start, end, config.ColorRed)
	a.beforeAttackParam = &Param{EventType: EventTypeBeforeAttack, Player: a.Param.Player, SrcCard: a.Param.SrcCard,
		TarCard: a.Param.TarCard, Card: a.Param.SrcCard}
	ActionManager.TriggerEvent(a.beforeAttackParam)
	a.Next()
}

func (a *MonsterAttackAction) mainStep() { // 战斗结算
	if a.beforeAttackParam.Invalid {
		Tip.AddTip("战斗取消")
		a.End()
		return
	} // 计算伤害结果
	a.afterAttackParam = &Param{EventType: EventTypeAfterAttack, Player: a.Param.Player, Card: a.Param.SrcCard}
	if a.Param.TarCard == nil { // 直接怼玩家
		a.afterAttackParam.HurtValue = a.Param.SrcCard.GetValue()
	} else { // 怪兽结算
		hurtValue := a.Param.SrcCard.GetValue() - a.Param.TarCard.GetValue()
		if hurtValue > 0 {
			a.afterAttackParam.TarCard = a.Param.TarCard
			if a.Param.TarCard.Defense {
				hurtValue = 0
			}
		} else if hurtValue < 0 {
			if !a.Param.TarCard.Defense {
				a.afterAttackParam.SrcCard = a.Param.SrcCard
			}
		} else {
			a.afterAttackParam.SrcCard = a.Param.SrcCard
			a.afterAttackParam.TarCard = a.Param.TarCard
		}
		a.afterAttackParam.HurtValue = hurtValue
	}
	ActionManager.TriggerEvent(a.afterAttackParam) // 触发事件
	a.Next()
}

func (a *MonsterAttackAction) endStep() { // 判断 是否送墓
	src := a.afterAttackParam.Player
	tar := PlayerManager.GetAnother(src)
	if a.afterAttackParam.HurtValue > 0 {
		Tip.AddTip(fmt.Sprintf("被攻击者受到%d点伤害", a.afterAttackParam.HurtValue))
		tar.ChangeHp(-a.afterAttackParam.HurtValue)
	} else if a.afterAttackParam.HurtValue < 0 {
		Tip.AddTip(fmt.Sprintf("攻击者受到%d点伤害", -a.afterAttackParam.HurtValue))
		src.ChangeHp(a.afterAttackParam.HurtValue)
	}
	if a.afterAttackParam.SrcCard != nil {
		Tip.AddTip(fmt.Sprintf("攻击怪兽[%s]被战斗破坏", a.afterAttackParam.SrcCard.Data.Name))
		src.RemoveCardUI(a.afterAttackParam.SrcCard)
		src.AddCardPile(a.afterAttackParam.SrcCard, CardPileTypeCemetery)
	}
	if a.afterAttackParam.TarCard != nil {
		Tip.AddTip(fmt.Sprintf("被攻击怪兽[%s]被战斗破坏", a.afterAttackParam.TarCard.Data.Name))
		tar.RemoveCardUI(a.afterAttackParam.TarCard)
		tar.AddCardPile(a.afterAttackParam.TarCard, CardPileTypeCemetery)
	}
	a.Next()
}

//=====================NormalSummonAction======================

type NormalSummonAction struct {
	*StepAction
	normalSummonParam, costParam *Param
	allNum, CostNum              int
}

func (a *NormalSummonAction) eventStep() {
	a.Param.PlayPhase.SummonNum--
	a.CostNum = GetCostNum(a.Param.Card.Data.Level)
	a.allNum = a.CostNum
	a.normalSummonParam = &Param{EventType: EventTypeNormalSummon, Player: a.Param.Player, Card: a.Param.Card}
	ActionManager.TriggerEvent(a.normalSummonParam)
	a.Next()
}

func (a *NormalSummonAction) costStep() {
	if a.costParam == nil { // 第一次 进入
		if a.normalSummonParam.Invalid {
			Tip.AddTip("通常召唤被无效化")
			a.End()
			return
		}
	} else { // 处理结果进入
		if a.costParam.Card == nil { // 取消召唤
			Tip.AddTip("没有选择祭品，召唤取消")
			a.End()
			return
		}
		a.Param.Player.RemoveCardUI(a.costParam.Card)
		a.Param.Player.AddCardPile(a.costParam.Card, CardPileTypeCemetery)
		a.CostNum--
	}
	if a.CostNum > 0 {
		Tip.AddTip(fmt.Sprintf("请选择一个怪兽解放(%d/%d)", a.CostNum, a.allNum))
		a.costParam = &Param{Player: a.Param.Player, CardFilter: And(IsMonster, IsField)}
		ActionManager.PushAction(a.costParam, NewWarpHandle(CreateSelectCardAction))
	} else {
		a.Next()
	}
}

func (a *NormalSummonAction) mainStep() {
	cardArea := a.Param.Player.GetEmptyCardArea(true)
	if cardArea == nil {
		Tip.AddTip("没有可以召唤的区域!")
		a.End()
		return
	}
	cardArea.SetCardUI(a.Param.Player.RemoveCardUI(a.Param.Card))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.Param.Card})
	a.Next()
}

func NewNormalSummonAction(param *Param) *NormalSummonAction {
	res := &NormalSummonAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.eventStep, res.costStep, res.mainStep}
	return res
}

//=================MonsterPrepareAction=====================

type MonsterPrepareAction struct {
	*OnceAction
	card *Card
}

func (a *MonsterPrepareAction) mainStep() {
	a.card.HasAdjust = false
	a.card.AttackNum = 1
}

func NewMonsterPrepareAction(card *Card, param *Param) *MonsterPrepareAction {
	res := &MonsterPrepareAction{card: card}
	res.OnceAction = NewOnceAction(param, res.mainStep)
	return res
}

//=====================RetrieveCardAction===================

type RetrieveCardAction struct { // 尝试  从  某一区域检索  一张 卡牌 到手牌
	*OnceAction
	cardPileType int
	filter       CardFilter
}

func (a *RetrieveCardAction) mainStep() {
	cards := a.Param.Player.GetPileCards(a.cardPileType, a.filter)
	if len(cards) > 0 {
		a.Param.Player.RemovePileCard(a.cardPileType, cards[0])
		a.Param.Player.AddHandCards(cards[0])
	}
}

func NewRetrieveCardAction(param *Param, cardPileType int, filter CardFilter) *RetrieveCardAction {
	res := &RetrieveCardAction{cardPileType: cardPileType, filter: filter}
	res.OnceAction = NewOnceAction(param, res.mainStep)
	return res
}

//=======================SummonBlueEye1Action======================

type SummonBlueEye1Action struct {
	*StepAction
	chooseParam *Param
}

func (a *SummonBlueEye1Action) sacrificeStep() {
	Tip.AddTip("请选择一个怪兽解放")
	a.chooseParam = &Param{Player: a.Param.Player, CardFilter: And(IsMonster, IsField)}
	ActionManager.PushAction(a.chooseParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *SummonBlueEye1Action) chooseStep() {
	if a.chooseParam.Card == nil {
		Tip.AddTip("没有选择祭品，特殊召唤取消")
		a.End()
		return
	}
	// 防止 后面选择 这个
	a.Param.Player.RemoveCardUI(a.Param.Card)
	a.Param.Player.AddCardPile(a.Param.Card, CardPileTypeCemetery)
	a.Param.Player.RemoveCardUI(a.chooseParam.Card)
	a.Param.Player.AddCardPile(a.chooseParam.Card, CardPileTypeCemetery)
	Tip.AddTip("请选择一个手牌中的青眼怪兽特殊召唤")
	a.chooseParam = &Param{Player: a.Param.Player, CardFilter: And(IsMonster, IsHand, HasField(FieldBlueEye))}
	ActionManager.PushAction(a.chooseParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *SummonBlueEye1Action) mainStep() {
	if a.chooseParam.Card == nil {
		Tip.AddTip("没有选择被召唤的青眼怪兽，特殊召唤取消")
		a.End()
		return
	}
	cardArea := a.Param.Player.GetEmptyCardArea(true)
	if cardArea == nil {
		Tip.AddTip("没有空格子召唤，特殊召唤取消")
		a.End()
		return
	}
	cardArea.SetCardUI(a.Param.Player.RemoveCardUI(a.chooseParam.Card))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.chooseParam.Card})
	a.Next()
}

func NewSummonBlueEye1Action(param *Param) *SummonBlueEye1Action {
	res := &SummonBlueEye1Action{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.sacrificeStep, res.chooseStep, res.mainStep}
	return res
}

//=======================SummonBlueEye2Action======================

type SummonBlueEye2Action struct {
	*StepAction
	chooseParam *Param
}

func (a *SummonBlueEye2Action) chooseStep() {
	Tip.AddTip("[太古的白石]请选择一个手牌中的青眼怪兽特殊召唤")
	a.chooseParam = &Param{Player: a.Param.Player, CardFilter: And(IsMonster, IsHand, HasField(FieldBlueEye))}
	ActionManager.PushAction(a.chooseParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *SummonBlueEye2Action) mainStep() {
	if a.chooseParam.Card == nil {
		Tip.AddTip("没有选择被召唤的青眼怪兽，特殊召唤取消")
		a.End()
		return
	}
	cardArea := a.Param.Player.GetEmptyCardArea(true)
	if cardArea == nil {
		Tip.AddTip("没有空格子召唤，特殊召唤取消")
		a.End()
		return
	}
	cardArea.SetCardUI(a.Param.Player.RemoveCardUI(a.chooseParam.Card))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.chooseParam.Card})
	a.Next()
}

func NewSummonBlueEye2Action(param *Param) *SummonBlueEye2Action {
	res := &SummonBlueEye2Action{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.chooseStep, res.mainStep}
	return res
}

//====================PickBlueEyeAction======================

type PickBlueEyeAction struct {
	*StepAction
	pickParam *Param
}

func (a *PickBlueEyeAction) pickStep() {
	Tip.AddTip("[太古的白石]请选择墓地的一只青眼怪兽加入手牌")
	// 从墓地 除外  防止 选到他
	a.Param.Player.RemovePileCard(CardPileTypeCemetery, a.Param.Card)
	a.Param.Player.AddCardPile(a.Param.Card, CardPileTypeExcept)
	a.pickParam = &Param{Player: a.Param.Player, CardFilter: And(IsMonster, HasField(FieldBlueEye)),
		CardPile: a.Param.Player.GetPileCard(CardPileTypeCemetery)}
	ActionManager.PushAction(a.pickParam, NewWarpHandle(CreatePickCardAction))
	a.Next()
}

func (a *PickBlueEyeAction) mainStep() {
	if a.pickParam.Card == nil {
		Tip.AddTip("没有选择墓地的青眼怪兽，取消发动效果")
		a.End()
		return
	}
	// 从墓地 回到手牌
	a.Param.Player.RemovePileCard(CardPileTypeCemetery, a.pickParam.Card)
	a.Param.Player.AddHandCards(a.pickParam.Card)
	a.Next()
}

func NewPickBlueEyeAction(param *Param) *PickBlueEyeAction {
	res := &PickBlueEyeAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.pickStep, res.mainStep}
	return res
}

//=======================SummonToEnemyAction=====================

type SummonToEnemyAction struct {
	*StepAction
	costParam *Param
}

func (a *SummonToEnemyAction) costStep() {
	Tip.AddTip("[海龟坏兽]请选择对面的一个怪兽解放")
	a.costParam = &Param{Player: PlayerManager.GetAnother(a.Param.Player), CardFilter: And(IsMonster, IsField)}
	ActionManager.PushAction(a.costParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *SummonToEnemyAction) mainStep() {
	if a.costParam.Card == nil { // 取消召唤
		Tip.AddTip("没有选择祭品，召唤取消")
		a.End()
		return
	}
	player := PlayerManager.GetAnother(a.Param.Player)
	cardArea := player.GetEmptyCardArea(true)
	if cardArea == nil {
		Tip.AddTip("没有可以召唤的区域!")
		a.End()
		return
	}
	// 移除祭品
	player.RemoveCardUI(a.costParam.Card)
	player.AddCardPile(a.costParam.Card, CardPileTypeCemetery)
	// 召唤卡牌
	cardArea.SetCardUI(a.Param.Player.RemoveCardUI(a.Param.Card))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.Param.Card})
	a.Next()
}

func NewSummonToEnemyAction(param *Param) *SummonToEnemyAction {
	res := &SummonToEnemyAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.costStep, res.mainStep}
	return res
}

//=======================ShowSummonAction=====================

type ShowSummonAction struct {
	*StepAction
	selectParam *Param
}

func (a *ShowSummonAction) selectStep() {
	Tip.AddTip("[青眼亚白龙]请选择手牌中的一张青眼白龙展示")
	a.selectParam = &Param{Player: a.Param.Player, CardFilter: And(IsHand, NameEq("青眼白龙"))}
	ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *ShowSummonAction) mainStep() {
	if a.selectParam.Card == nil { // 取消召唤
		Tip.AddTip("没有展示的[青眼白龙]，召唤取消")
		a.End()
		return
	}
	cardArea := a.Param.Player.GetEmptyCardArea(true)
	if cardArea == nil {
		Tip.AddTip("没有可以召唤的区域!")
		a.End()
		return
	}
	// 召唤卡牌
	cardArea.SetCardUI(a.Param.Player.RemoveCardUI(a.Param.Card))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.Param.Card})
	a.Next()
}

func NewShowSummonAction(param *Param) *ShowSummonAction {
	res := &ShowSummonAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//=======================SummonBlueEyeWhiteDragonAction=====================

type SummonBlueEyeWhiteDragonAction struct {
	*StepAction
	selectParam *Param
}

func (a *SummonBlueEyeWhiteDragonAction) selectStep() {
	Tip.AddTip("[白灵龙]请选择手牌中的一张青眼白龙特殊召唤")
	a.selectParam = &Param{Player: a.Param.Player, CardFilter: And(IsHand, NameEq("青眼白龙"))}
	ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *SummonBlueEyeWhiteDragonAction) mainStep() {
	if a.selectParam.Card == nil { // 取消召唤
		Tip.AddTip("没有选择[青眼白龙]，召唤取消")
		a.End()
		return
	}
	// 进入 墓地
	a.Param.Player.RemoveCardUI(a.Param.Card)
	a.Param.Player.AddCardPile(a.Param.Card, CardPileTypeCemetery)
	// 召唤上场
	a.Param.Player.GetEmptyCardArea(true).SetCardUI(a.Param.Player.RemoveCardUI(a.selectParam.Card))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.Param.Card})
	a.Next()
}

func NewSummonBlueEyeWhiteDragonAction(param *Param) *SummonBlueEyeWhiteDragonAction {
	res := &SummonBlueEyeWhiteDragonAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//=======================DrawCard1Action=====================

type DrawCard1Action struct {
	*StepAction
	selectParam *Param
}

func (a *DrawCard1Action) selectStep() {
	a.Param.Player.RemoveCardUI(a.Param.Card)
	a.Param.Player.AddCardPile(a.Param.Card, CardPileTypeCemetery)
	Tip.AddTip("[交易进行]请选择手牌中的一张等级8的怪兽")
	a.selectParam = &Param{Player: a.Param.Player, CardFilter: And(IsHand, IsMonster, LevelEq(8))}
	ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *DrawCard1Action) mainStep() {
	if a.selectParam.Card == nil {
		Tip.AddTip("没有选择等级8的怪兽，效果终止")
		a.End()
		return
	} // 扔卡 摸牌
	a.Param.Player.RemoveCardUI(a.selectParam.Card)
	a.Param.Player.AddCardPile(a.selectParam.Card, CardPileTypeCemetery)
	a.Param.Player.DrawCards(2)
	a.Next()
}

func NewDrawCard1Action(param *Param) *DrawCard1Action {
	res := &DrawCard1Action{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//=======================RetrieveMonsterAction=====================

type RetrieveMonsterAction struct {
	*StepAction
	selectParam *Param
}

func (a *RetrieveMonsterAction) selectStep() {
	a.Param.Player.RemoveCardUI(a.Param.Card)
	a.Param.Player.AddCardPile(a.Param.Card, CardPileTypeCemetery)
	Tip.AddTip("[龙旋律]请选择一张手牌丢弃")
	a.selectParam = &Param{Player: a.Param.Player, CardFilter: IsHand}
	ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *RetrieveMonsterAction) mainStep() {
	if a.selectParam.Card == nil {
		Tip.AddTip("没有选择手牌，效果终止")
		a.End()
		return
	} // 扔卡 摸牌
	a.Param.Player.RemoveCardUI(a.selectParam.Card)
	a.Param.Player.AddCardPile(a.selectParam.Card, CardPileTypeCemetery) // 检索
	cards := a.Param.Player.GetPileCards(CardPileTypeDeck, And(AtkGe(3000), DefLe(2500), RaceEq(RaceDragon)))
	max := utils.Min(2, len(cards))
	for i := 0; i < max; i++ {
		a.Param.Player.RemovePileCard(CardPileTypeDeck, cards[i])
		a.Param.Player.AddHandCards(cards[i])
	}
	a.Next()
}

func NewRetrieveMonsterAction(param *Param) *RetrieveMonsterAction {
	res := &RetrieveMonsterAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//=======================DrawCard2Action=====================

type DrawCard2Action struct {
	*StepAction
	selectParam *Param
}

func (a *DrawCard2Action) selectStep() {
	a.Param.Player.RemoveCardUI(a.Param.Card)
	a.Param.Player.AddCardPile(a.Param.Card, CardPileTypeCemetery)
	Tip.AddTip("[调和的宝牌]请选择手牌中的一张[太古的白石]")
	a.selectParam = &Param{Player: a.Param.Player, CardFilter: And(IsHand, NameEq("太古的白石"))}
	ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *DrawCard2Action) mainStep() {
	if a.selectParam.Card == nil {
		Tip.AddTip("没有选择[太古的白石]，效果终止")
		a.End()
		return
	} // 扔卡 摸牌
	a.Param.Player.RemoveCardUI(a.selectParam.Card)
	a.Param.Player.AddCardPile(a.selectParam.Card, CardPileTypeCemetery)
	a.Param.Player.DrawCards(2)
	a.Next()
}

func NewDrawCard2Action(param *Param) *DrawCard2Action {
	res := &DrawCard2Action{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//=======================ReviveMonsterAction======================

type ReviveMonsterAction struct {
	*StepAction
	pickParam *Param
}

func (a *ReviveMonsterAction) pickStep() {
	a.Param.Player.RemoveCardUI(a.Param.Card)
	a.Param.Player.AddCardPile(a.Param.Card, CardPileTypeCemetery)
	Tip.AddTip("[复活的福音]请选择墓地一只8星龙族怪兽复活")
	a.pickParam = &Param{Player: a.Param.Player, CardFilter: And(IsMonster, LevelEq(8), RaceEq(RaceDragon)),
		CardPile: a.Param.Player.GetPileCard(CardPileTypeCemetery)}
	ActionManager.PushAction(a.pickParam, NewWarpHandle(CreatePickCardAction))
	a.Next()
}

func (a *ReviveMonsterAction) mainStep() {
	if a.pickParam.Card == nil {
		Tip.AddTip("没有选择复活怪兽，效果终止")
		a.End()
		return
	}
	cardArea := a.Param.Player.GetEmptyCardArea(true)
	if cardArea == nil {
		Tip.AddTip("没有空格子召唤，效果取消")
		a.End()
		return
	}
	a.Param.Player.RemovePileCard(CardPileTypeCemetery, a.pickParam.Card)
	cardArea.SetCardUI(NewCardUI(a.pickParam.Card, a.Param.Player.player))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.pickParam.Card})
	a.Next()
}

func NewReviveMonsterAction(param *Param) *ReviveMonsterAction {
	res := &ReviveMonsterAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.pickStep, res.mainStep}
	return res
}

//=======================SummonMaxAction=====================

type SummonMaxAction struct {
	*StepAction
	selectParam *Param
}

func (a *SummonMaxAction) selectStep() {
	Tip.AddTip("请选择场上一只[青眼白龙]")
	a.selectParam = &Param{Player: a.Param.Player, CardFilter: And(IsField, NameEq("青眼白龙"))}
	ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *SummonMaxAction) mainStep() {
	if a.selectParam.Card == nil {
		Tip.AddTip("没有选择[青眼白龙]，终止召唤")
		a.End()
		return
	} // 移除
	a.Param.Player.RemoveCardUI(a.selectParam.Card)
	a.Param.Player.AddCardPile(a.selectParam.Card, CardPileTypeExcept)
	a.Param.Player.RemovePileCard(CardPileTypeExtra, a.Param.Card) // 召唤
	a.Param.Player.GetEmptyCardArea(true).SetCardUI(NewCardUI(a.Param.Card, a.Param.Player.player))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.Param.Card})
	a.Next()
}

func NewSummonMaxAction(param *Param) *SummonMaxAction {
	res := &SummonMaxAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//=====================SummonDoubleAction======================

type SummonDoubleAction struct {
	*StepAction
	selectParam *Param
	costNum     int
}

func (a *SummonDoubleAction) selectStep() {
	if a.selectParam != nil {
		if a.selectParam.Card == nil {
			Tip.AddTip("没有选择[青眼白龙]，召唤取消")
			a.End()
			return
		}
		a.Param.Player.RemoveCardUI(a.selectParam.Card)
		a.Param.Player.AddCardPile(a.selectParam.Card, CardPileTypeCemetery)
		a.costNum--
	}
	if a.costNum > 0 {
		Tip.AddTip(fmt.Sprintf("请选择一个[青眼白龙](%d/2)", a.costNum))
		a.selectParam = &Param{Player: a.Param.Player, CardFilter: And(IsField, HasField(FieldBlueEyeWhiteDragon))}
		ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	} else {
		a.Next()
	}
}

func (a *SummonDoubleAction) mainStep() {
	a.Param.Player.RemovePileCard(CardPileTypeExtra, a.Param.Card) // 召唤
	a.Param.Player.GetEmptyCardArea(true).SetCardUI(NewCardUI(a.Param.Card, a.Param.Player.player))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.Param.Card})
	a.Next()
}

func NewSummonDoubleAction(param *Param) *SummonDoubleAction {
	res := &SummonDoubleAction{costNum: 2}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//=====================SummonChestnutAction======================

type SummonChestnutAction struct {
	*StepAction
	selectParam *Param
	costNum     int
}

func (a *SummonChestnutAction) selectStep() {
	if a.selectParam != nil {
		if a.selectParam.Card == nil {
			Tip.AddTip("没有选择一星怪兽，召唤取消")
			a.End()
			return
		}
		a.Param.Player.RemoveCardUI(a.selectParam.Card)
		a.Param.Player.AddCardPile(a.selectParam.Card, CardPileTypeCemetery)
		a.costNum--
	}
	if a.costNum > 0 {
		Tip.AddTip(fmt.Sprintf("请选择一个一星怪兽(%d/2)", a.costNum))
		a.selectParam = &Param{Player: a.Param.Player, CardFilter: And(IsField, IsMonster, LevelEq(1))}
		ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	} else {
		a.Next()
	}
}

func (a *SummonChestnutAction) mainStep() {
	a.Param.Player.RemovePileCard(CardPileTypeExtra, a.Param.Card) // 召唤
	a.Param.Player.GetEmptyCardArea(true).SetCardUI(NewCardUI(a.Param.Card, a.Param.Player.player))
	ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: a.Param.Player, Card: a.Param.Card})
	a.Next()
}

func NewSummonChestnutAction(param *Param) *SummonChestnutAction {
	res := &SummonChestnutAction{costNum: 2}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}

//====================SummonCommonAction======================

type SummonCommonAction struct {
	*StepAction
	pickParam *Param
}

func (a *SummonCommonAction) pickStep() {
	Tip.AddTip("[苍眼银龙]请选择墓地的一只通常怪兽特殊召唤")
	a.pickParam = &Param{Player: a.Param.Player, CardFilter: IsCommonMonster,
		CardPile: a.Param.Player.GetPileCard(CardPileTypeCemetery)}
	ActionManager.PushAction(a.pickParam, NewWarpHandle(CreatePickCardAction))
	a.Next()
}

func (a *SummonCommonAction) mainStep() {
	if a.pickParam.Card == nil {
		Tip.AddTip("没有选择墓地的通常怪兽，取消发动效果")
		a.End()
		return
	}
	// 从墓地 特殊召唤
	player := a.Param.Player
	player.RemovePileCard(CardPileTypeCemetery, a.pickParam.Card)
	player.GetEmptyCardArea(true).SetCardUI(NewCardUI(a.pickParam.Card, player.player))
	a.Next()
}

func NewSummonCommonAction(param *Param) *SummonCommonAction {
	res := &SummonCommonAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.pickStep, res.mainStep}
	return res
}

//=======================PlayBasketballAction=====================

type PlayBasketballAction struct {
	*StepAction
	selectParam *Param
}

func (a *PlayBasketballAction) selectStep() {
	a.Param.Player.DrawCards(1)
	enemy := PlayerManager.GetAnother(a.Param.Player)
	enemy.DrawCards(1)
	Tip.AddTip("[蔡徐坤]请选择手牌中的一张手牌")
	a.selectParam = &Param{Player: enemy, CardFilter: IsHand}
	ActionManager.PushAction(a.selectParam, NewWarpHandle(CreateSelectCardAction))
	a.Next()
}

func (a *PlayBasketballAction) mainStep() {
	if a.selectParam.Card == nil {
		Tip.AddTip("没有选择手牌，蔡徐坤获胜")
		a.selectParam.Player.SkipPhase(PhaseAttack)
		a.End()
		return
	}
	cards := a.Param.Player.GetHandCards(IsAny)
	playerCard := cards[rand.Intn(len(cards))]
	a.Param.Player.RemoveCardUI(playerCard)
	a.Param.Player.AddCardPile(playerCard, CardPileTypeCemetery)
	enemyCard := a.selectParam.Card
	a.selectParam.Player.RemoveCardUI(enemyCard)
	a.selectParam.Player.AddCardPile(enemyCard, CardPileTypeCemetery)
	if (enemyCard.Data.CardType+1)%3 != playerCard.Data.CardType { // 借助 枚举顺序了
		a.selectParam.Player.SkipPhase(PhaseAttack)
	}
	a.Next()
}

func NewPlayBasketballAction(param *Param) *PlayBasketballAction {
	res := &PlayBasketballAction{}
	res.StepAction = NewStepAction(param)
	res.Step = []func(){res.selectStep, res.mainStep}
	return res
}
