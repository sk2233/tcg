/*
@author: sk
@date: 2023/2/5
*/
package main

import (
	"GameBase2/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//================DrawPhase==================

type DrawPhase struct {
	player *Player
	time   int
}

func (d *DrawPhase) Init() {
	Tip.AddTip("抽卡阶段")
	d.time = ThinkTime
}

func NewDrawPhase(player *Player) *DrawPhase {
	return &DrawPhase{player: player}
}

func (d *DrawPhase) Action() {
	if d.time > 0 {
		d.time--
	} else {
		d.player.DrawCards(1)
		d.player.NextPhase()
	}
}

//================PreparePhase==================

type PreparePhase struct {
	player *Player
	time   int
}

func (d *PreparePhase) Init() {
	Tip.AddTip("准备阶段")
	d.time = ThinkTime
}

func NewPreparePhase(player *Player) *PreparePhase {
	return &PreparePhase{player: player}
}

func (d *PreparePhase) Action() {
	if d.time > 0 {
		d.time--
	} else {
		d.player.NextPhase()
		ActionManager.TriggerEvent(&Param{EventType: EventTypePreparePhase, Player: d.player})
	}
}

//================PlayPhase==================

type PlayPhase struct {
	player    *Player
	SummonNum int
	// 电脑使用
	time int
	// 玩家使用
	nextBtn     *ButtonGroup
	actionBtns  *ButtonGroup
	nameHandles []INameHandle
	selectCard  *Card
}

func (d *PlayPhase) SetTop(top bool) {
	if d.player.player {
		d.nextBtn.Enable = top
	}
}

func (d *PlayPhase) Init() {
	Tip.AddTip("主要阶段")
	d.SummonNum = 1 // 仅允许普通召唤一次
	if d.player.player {
		d.nextBtn = NewButtonGroup(0.5+0.5i, 0)
		d.nextBtn.SetPos(1143 + 610i + CardSize/2)
		d.nextBtn.Replace("Next")
		d.actionBtns = NewButtonGroup(0.5+1i, 4)
		d.actionBtns.Enable = false
	} else {
		d.time = ThinkTime
	}
}

func NewPlayPhase(player *Player) *PlayPhase {
	return &PlayPhase{player: player}
}

func (d *PlayPhase) Action() {
	if d.player.player {
		d.playerPlay()
	} else {
		d.enemyPlay()
	}
}

func (d *PlayPhase) playerPlay() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	pos := utils.GetCursorPos()
	if d.actionBtns.Enable { // 行为按钮启用时仅判断行为按钮
		d.actionBtns.Enable = false // 不管选择 什么 都消失  什么都没点也消失
		index := d.actionBtns.ClickBtn(pos)
		if index < 0 {
			return
		}
		ActionManager.PushAction(&Param{EventType: d.selectCard.Place, PlayPhase: d,
			Player: d.player, Card: d.selectCard}, d.nameHandles[index])
		return
	}
	if d.nextBtn.ClickBtn(pos) == 0 { // 结束判断
		d.playerEnd()
		return
	}
	card := d.player.ClickCard(pos) // 点击手卡或 桌子上的 卡  都是 卡不过状态 不同 涉及的行为不同 行为是卡自己判断的
	if card != nil {
		Info.SetCard(card.Card) // 实际内部也可以通过 自身的 card.Card.Place 判断
		d.selectCard = card.Card
		d.nameHandles = card.Card.GetNameHandles(&Param{EventType: card.Card.Place, PlayPhase: d, Player: d.player,
			Card: card.Card})
		if len(d.nameHandles) > 0 {
			btns := make([]string, 0)
			for i := 0; i < len(d.nameHandles); i++ {
				btns = append(btns, d.nameHandles[i].GetName())
			}
			d.actionBtns.Replace(btns...)
			d.actionBtns.Enable = true
			d.actionBtns.SetPos(card.Pos + 75/2.0)
		} else {
			Tip.AddTip("该卡当前无法使用")
		}
		return
	}
	cardPile := d.player.ClickCardPile(pos) // 额外 墓地  除外  卡组 选择  暂时 仅处理 额外与墓地
	if cardPile != nil && (cardPile.Type == CardPileTypeExtra || cardPile.Type == CardPileTypeCemetery) {
		ActionManager.PushAction(&Param{CardPile: cardPile, Player: d.player}, NewWarpHandle(CreateCardPileAction))
		return
	}
	card = PlayerManager.Enemy.ClickCard(pos) // 敌人 卡牌 只能查看 显示的卡牌
	if card != nil {
		if card.Card.Place == PlaceField && !card.Card.Flip {
			Info.SetCard(card.Card)
		}
	}
}

func (d *PlayPhase) enemyPlay() {
	if d.time > 0 {
		d.time--
		return
	}
	if d.player.computerPlay(d) {
		d.time = ThinkTime
		return
	}
	d.player.NextPhase()
}

func (d *PlayPhase) playerEnd() {
	d.nextBtn.Die = true
	d.actionBtns.Die = true
	d.player.NextPhase()
}

//================AttackPhase==================

type AttackPhase struct {
	player *Player
	// 电脑使用
	time int
	// 玩家使用
	nextBtn      *ButtonGroup
	selectCardUI *CardUI
}

func (d *AttackPhase) SetTop(top bool) {
	if d.player.player {
		d.nextBtn.Enable = top
	}
}

func (d *AttackPhase) Init() {
	Tip.AddTip("攻击阶段")
	ActionManager.TriggerEvent(&Param{EventType: EventTypeAttackPhase, Player: d.player})
	if d.player.player {
		d.nextBtn = NewButtonGroup(0.5+0.5i, 0)
		d.nextBtn.SetPos(1143 + 610i + CardSize/2)
		d.nextBtn.Replace("Next")
	} else {
		d.time = ThinkTime
	}
}

func NewAttackPhase(player *Player) *AttackPhase {
	return &AttackPhase{player: player}
}

func (d *AttackPhase) Action() {
	if d.player.player {
		d.playerAttack()
	} else {
		d.enemyAttack()
	}
}

func (d *AttackPhase) playerAttack() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	pos := utils.GetCursorPos()
	if d.nextBtn.ClickBtn(pos) == 0 { // 结束判断
		d.playerEnd()
		return
	}
	cardUI := d.player.ClickCard(pos) // 选择 自己的 怪兽  其他都是查看
	if cardUI != nil {
		Info.SetCard(cardUI.Card) // 实际内部也可以通过 自身的 cardUI.Card.Place 判断
		// 必须在场上的 怪兽
		if cardUI.Card.Data.CardType == CardTypeMonster && cardUI.Card.Place == PlaceField {
			if !cardUI.Card.CanSelect(&Param{EventType: EventTypeSelectSrc}) {
				Tip.AddTip("选择的怪兽无法攻击")
				return
			}
			if d.selectCardUI != nil {
				d.selectCardUI.UnSelect()
				if d.selectCardUI == cardUI {
					d.selectCardUI = nil
					return
				}
			}
			d.selectCardUI = cardUI
			d.selectCardUI.Select()
		}
		return
	}
	cardUI = PlayerManager.Enemy.ClickCard(pos) // 尝试攻击 敌人的怪兽
	if cardUI != nil {
		if cardUI.Card.Place == PlaceField && !cardUI.Card.Flip {
			Info.SetCard(cardUI.Card)
			if d.selectCardUI == nil || cardUI.Card.Data.CardType != CardTypeMonster {
				return // 必须 选择了 自己的 怪兽 且  选择 别人的 卡牌 也是怪兽
			}
			if !cardUI.Card.CanSelect(&Param{EventType: EventTypeSelectTar, Card: d.selectCardUI.Card}) {
				Tip.AddTip("不能选择目标怪兽作为攻击对象")
				return
			}
			ActionManager.PushAction(&Param{EventType: EventTypeMonsterAttack, Player: d.player,
				SrcCard: d.selectCardUI.Card, TarCard: cardUI.Card}, NewWarpHandle(CreateMonsterAttackAction))
			d.selectCardUI.UnSelect()
			d.selectCardUI = nil
		}
	}
	if d.selectCardUI != nil && PlayerManager.Enemy.Click(pos) { // 直接攻击 敌方
		if len(PlayerManager.Enemy.GetFieldCards(IsMonster)) > 0 {
			Tip.AddTip("对方场上还存在怪兽没有解除，不能直接攻击主战者")
			return
		}
		ActionManager.PushAction(&Param{EventType: EventTypeMonsterAttack, Player: d.player,
			SrcCard: d.selectCardUI.Card}, NewWarpHandle(CreateMonsterAttackAction))
		d.selectCardUI.UnSelect()
		d.selectCardUI = nil
	}
}

func (d *AttackPhase) enemyAttack() {
	if d.time > 0 {
		d.time--
		return
	}
	if d.player.computerAttack() {
		d.time = ThinkTime
		return
	}
	d.player.NextPhase()
}

func (d *AttackPhase) playerEnd() {
	d.nextBtn.Die = true
	if d.selectCardUI != nil {
		d.selectCardUI.UnSelect()
	}
	d.player.NextPhase()
}

//================DiscardPhase==================

type DiscardPhase struct {
	player *Player
	// 玩家使用
	okBtn        *ButtonGroup
	selectCardUI *CardUI
	// 电脑特有
	time int
}

func (d *DiscardPhase) SetTop(top bool) {
	if d.player.player {
		d.okBtn.Enable = top
	}
}

func (d *DiscardPhase) Init() {
	Tip.AddTip("弃牌阶段")
	if len(d.player.GetHandCards(IsAny)) <= MaxCardNum {
		d.player.NextPhase() // 无需弃牌
		return
	}
	if d.player.player {
		d.okBtn = NewButtonGroup(0.5+0.5i, 0)
		d.okBtn.SetPos(1143 + 610i + CardSize/2)
		d.okBtn.Replace("Ok")
	} else {
		d.time = ThinkTime
	}
}

func NewDiscardPhase(player *Player) *DiscardPhase {
	return &DiscardPhase{player: player}
}

func (d *DiscardPhase) Action() {
	if d.player.player {
		d.playerDiscard()
	} else {
		d.enemyDiscard()
	}
}

func (d *DiscardPhase) playerDiscard() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	pos := utils.GetCursorPos()
	cardUI := d.player.ClickCard(pos) // 选择 自己的 手牌  其他都是查看
	if cardUI != nil {
		Info.SetCard(cardUI.Card)           // 实际内部也可以通过 自身的 cardUI.Card.Place 判断
		if cardUI.Card.Place != PlaceHand { // 手牌判断
			return
		}
		if d.selectCardUI != nil {
			d.selectCardUI.UnSelect()
			if d.selectCardUI == cardUI {
				d.selectCardUI = nil
				return
			}
		}
		d.selectCardUI = cardUI
		d.selectCardUI.Select()
		return
	}
	if d.okBtn.ClickBtn(pos) == 0 && d.selectCardUI != nil { // 确定 弃牌
		d.player.RemoveCardUI(d.selectCardUI.Card)
		d.player.AddCardPile(d.selectCardUI.Card, CardPileTypeCemetery)
		if len(d.player.GetHandCards(IsAny)) <= MaxCardNum { // 不用 弃牌了
			d.okBtn.Die = true
			d.player.NextPhase()
		}
	}
}

func (d *DiscardPhase) enemyDiscard() {
	if d.time > 0 {
		d.time--
		return
	}
	cards := d.player.GetHandCards(IsAny)
	if len(cards) > MaxCardNum {
		card := cards[0]
		d.player.RemoveCardUI(card)
		d.player.AddCardPile(card, CardPileTypeCemetery)
		ActionManager.TriggerEvent(&Param{EventType: EventTypeDiscardCard, Card: card, Player: d.player})
		d.time = ThinkTime
		return
	}
	d.player.NextPhase()
}

//================EndPhase==================

type EndPhase struct {
	player *Player
	time   int
}

func (d *EndPhase) Init() {
	d.player.Reset() // 这里清除状态 若要添加  请在 回合结束提加
	Tip.AddTip("结束阶段")
	d.time = ThinkTime
}

func NewEndPhase(player *Player) *EndPhase {
	return &EndPhase{player: player}
}

func (d *EndPhase) Action() {
	if d.time > 0 {
		d.time--
	} else {
		ActionManager.TriggerEvent(&Param{EventType: EventTypeEndPhase, Player: d.player})
		d.player.NextPhase()
	}
}
