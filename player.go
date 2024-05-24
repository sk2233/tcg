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
	"reflect"
	R "tcg/res"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func init() {
	config.ObjectFactory.RegisterPointFactory(R.CLASS.PLAYER, createPlayer)
}

func createPlayer(data *model.ObjectData) model.IObject {
	res := &Player{player: data.Name == R.OBJECT.PLAYER, Hp: 8000, cardPiles: make([]*CardPile, 4),
		cardAreas: make([]*CardArea, 0)}
	res.PointObject = object.NewPointObject()
	factory.FillPointObject(data, res.PointObject)
	res.phases = []IPhase{NewDrawPhase(res), NewPreparePhase(res), NewPlayPhase(res), NewAttackPhase(res),
		NewDiscardPhase(res), NewEndPhase(res)}
	return res
}

type Player struct {
	*object.PointObject
	player        bool // 是否为玩家
	Hp            int
	cardPiles     []*CardPile // 4个 按CardPileType 常量访问
	cardAreas     []*CardArea // 怪兽区  魔 陷 区
	skipPhaseMask int
	phases        []IPhase
	phaseIndex    int
	cardUIs       []*CardUI
}

func (p *Player) GetMin() complex128 {
	return p.Pos
}

func (p *Player) GetMax() complex128 {
	return p.Pos + CardSize // 与一张 卡牌 一样大
}

func (p *Player) Init() {
	objLayer := utils.GetObjectLayer(R.LAYER.BASE)
	cardPiles := objLayer.GetObjectsByType(reflect.TypeOf(&CardPile{}))
	for i := 0; i < len(cardPiles); i++ {
		cardPile := cardPiles[i].(*CardPile)
		if cardPile.GetBool(R.PROP.PLAYER, false) == p.player {
			p.cardPiles[cardPile.Type] = cardPile
		}
	}
	cardAreas := objLayer.GetObjectsByType(reflect.TypeOf(&CardArea{}))
	for i := 0; i < len(cardAreas); i++ {
		cardArea := cardAreas[i].(*CardArea)
		if cardArea.GetBool(R.PROP.PLAYER, false) != p.player {
			continue
		}
		p.cardAreas = append(p.cardAreas, cardArea)
	}
	deck, extra := GetSetoKaibaCards() // 初始化卡片  不会触发事件
	if !p.player {
		deck, extra = GetBZhanCards()
	}
	p.cardPiles[CardPileTypeDeck].InitCards(p, deck)
	p.cardPiles[CardPileTypeExtra].InitCards(p, extra)
	p.DrawCards(5)
}

func (p *Player) Order() int {
	return -2233
}

func (p *Player) Draw(screen *ebiten.Image) {
	var name string
	if p.player { // 绘制区域背景
		utils.FillRect(screen, 330+360i, 950+360i, Color0_0_64)
		name = "玩家"
	} else {
		utils.FillRect(screen, 330, 950+360i, Color64_0_0)
		name = "敌人"
	}
	utils.DrawAnchorText(screen, fmt.Sprintf("%s\nHP:%d", name, p.Hp), p.Pos, 0, Font36, colornames.White)
	for i := len(p.cardUIs) - 1; i >= 0; i-- {
		p.cardUIs[i].Draw(screen)
	}
}

func (p *Player) SkipPhase(phase int) {
	p.skipPhaseMask |= phase
}

func (p *Player) Prepare() {
	p.phaseIndex = 0
	p.adjustPhaseIndex()
}

func (p *Player) Action() {
	p.phases[p.phaseIndex].Action()
}

func (p *Player) adjustPhaseIndex() {
	for p.phaseIndex < len(p.phases) && (p.skipPhaseMask&(1<<p.phaseIndex)) > 0 {
		p.phaseIndex++
	}
	if p.phaseIndex < len(p.phases) {
		p.phases[p.phaseIndex].Init()
	} else {
		PlayerManager.NextPlayer(p.player)
	}
}

func (p *Player) NextPhase() {
	p.phaseIndex++
	p.adjustPhaseIndex()
}

func (p *Player) DrawCards(num int) {
	cards := p.cardPiles[CardPileTypeDeck].GetCards(num)
	if cards == nil {
		StackRoom.PushLayer(NewWinUILayer(!p.player, "对手无卡可抽"))
		return
	}
	p.AddHandCards(cards...)
}

func (p *Player) AddHandCards(cards ...*Card) {
	for i := 0; i < len(cards); i++ {
		cards[i].Place = PlaceHand
		p.cardUIs = append(p.cardUIs, NewCardUI(cards[i], p.player))
	}
	p.tidyCard()
}

func (p *Player) tidyCard() {
	// 518 ~ 1093
	l := float64(len(p.cardUIs))
	w := utils.Min(l*75, 1093-518)
	x := 518 + (1093-518-w)/2
	y := imag(p.Pos)
	offset := (w - 75) / (l - 1)
	for i := 0; i < len(p.cardUIs); i++ {
		p.cardUIs[i].Pos = complex(x, y)
		x += offset
	}
}

func (p *Player) Reset() {
	p.skipPhaseMask = 0
}

func (p *Player) ClickCard(pos complex128) *CardUI {
	for i := 0; i < len(p.cardUIs); i++ {
		if p.cardUIs[i].CollisionPoint(pos) {
			return p.cardUIs[i]
		}
	}
	for i := 0; i < len(p.cardAreas); i++ {
		card := p.cardAreas[i].ClickCard(pos)
		if card != nil {
			return card
		}
	}
	return nil
}

func (p *Player) ClickCardPile(pos complex128) *CardPile {
	for i := 0; i < len(p.cardPiles); i++ {
		if p.cardPiles[i].CollisionPoint(pos) {
			return p.cardPiles[i]
		}
	}
	return nil
}

func (p *Player) GetHandles(param *Param) []IHandle { // 暂时 不考虑 人物技能
	res := make([]IHandle, 0)
	for i := 0; i < len(p.cardUIs); i++ { // 手牌 中的 卡
		res = append(res, p.cardUIs[i].Card.GetHandles(param)...)
	}
	for i := 0; i < len(p.cardPiles); i++ { // 卡组  墓地  额外   除外  牌堆效果
		res = append(res, p.cardPiles[i].GetHandles(param)...)
	}
	for i := 0; i < len(p.cardAreas); i++ { // 怪兽效果 魔陷 效果
		res = append(res, p.cardAreas[i].GetHandles(param)...)
	}
	return res
}

func (p *Player) GetEmptyCardArea(monster bool) *CardArea {
	for i := 0; i < len(p.cardAreas); i++ {
		if p.cardAreas[i].Monster == monster && p.cardAreas[i].GetCardUI() == nil {
			return p.cardAreas[i]
		}
	}
	return nil
}

func (p *Player) GetCardAreas(monster bool) []*CardArea {
	res := make([]*CardArea, 0)
	for i := 0; i < len(p.cardAreas); i++ {
		if p.cardAreas[i].Monster == monster {
			res = append(res, p.cardAreas[i])
		}
	}
	return res
}

func (p *Player) GetFieldCards(filter CardFilter) []*Card { // 获取场上 符合 条件的牌
	res := make([]*Card, 0)
	for i := 0; i < len(p.cardAreas); i++ {
		cardUI := p.cardAreas[i].GetCardUI()
		if cardUI != nil && filter(cardUI.Card) {
			res = append(res, cardUI.Card)
		}
	}
	return res
}

func (p *Player) GetHandCards(filter CardFilter) []*Card { // 获取手牌 中 符合的牌
	res := make([]*Card, 0)
	for i := 0; i < len(p.cardUIs); i++ {
		if filter(p.cardUIs[i].Card) {
			res = append(res, p.cardUIs[i].Card)
		}
	}
	return res
}

func (p *Player) RemoveCardUI(card *Card) *CardUI { // 移除  手牌  场上的 牌 并 返回
	for i := 0; i < len(p.cardUIs); i++ {
		if p.cardUIs[i].Card == card {
			res := p.cardUIs[i]
			p.cardUIs = append(p.cardUIs[:i], p.cardUIs[i+1:]...)
			p.tidyCard()
			return res
		}
	}
	for i := 0; i < len(p.cardAreas); i++ {
		cardUI := p.cardAreas[i].GetCardUI()
		if cardUI != nil && cardUI.Card == card {
			p.cardAreas[i].ClearCardUI()
			return cardUI
		}
	}
	return nil
}

func (p *Player) GetCardUI(card *Card) *CardUI { // 获取 手上 或场上的牌
	if card == nil {
		return nil
	}
	for i := 0; i < len(p.cardUIs); i++ {
		if p.cardUIs[i].Card == card {
			return p.cardUIs[i]
		}
	}
	for i := 0; i < len(p.cardAreas); i++ {
		cardUI := p.cardAreas[i].GetCardUI()
		if cardUI != nil && cardUI.Card == card {
			return cardUI
		}
	}
	return nil
}

func (p *Player) AddCardPile(card *Card, cardPileType int) { // 移动到 墓地 除外  牌堆  额外
	card.Place = cardPileType // 是一一对应的
	card.Defense = false      // 重制状态
	p.cardPiles[cardPileType].AddCard(card)
	ActionManager.TriggerEvent(&Param{EventType: EventTypeGoCardPile, Player: p, Card: card, CardPileType: cardPileType})
}

func (p *Player) SetTop(top bool) {
	InvokeTop(p.phases[p.phaseIndex], top)
}

// 电脑专用  是否 还可以行动
func (p *Player) computerPlay(playPhase *PlayPhase) bool {
	for i := 0; i < len(p.cardUIs); i++ {
		card := p.cardUIs[i].Card
		switch card.Data.CardType {
		case CardTypeMonster:
			if playPhase.SummonNum > 0 {
				cardArea := p.GetEmptyCardArea(true)
				if cardArea != nil {
					Info.SetCard(card)
					cardArea.SetCardUI(p.RemoveCardUI(card))
					ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: p, Card: card})
					playPhase.SummonNum--
					return true
				}
			}
		case CardTypeMagic: // 普通魔法卡 直接生效
			if card.Data.MagicType == MagicTypeCommon {
				param := &Param{EventType: EventTypeHandCard, Player: p}
				handles := card.GetHandles(param)
				if len(handles) == 0 {
					continue
				}
				Info.SetCard(card)
				p.RemoveCardUI(card)
				p.AddCardPile(card, CardPileTypeCemetery)
				ActionManager.PushAction(param, handles...)
				return true
			} // 其他魔法卡
			cardArea := p.GetEmptyCardArea(false)
			if cardArea != nil {
				Info.SetCard(card)
				cardArea.SetCardUI(p.RemoveCardUI(card))
				return true
			}
		case CardTypeTrap:
			card.Flip = true
			card.CanUse = false
			cardArea := p.GetEmptyCardArea(false)
			if cardArea != nil {
				cardArea.SetCardUI(p.RemoveCardUI(card))
				return true
			}
		}
	}
	// 调整怪兽状态
	playerMonsters := PlayerManager.Player.GetFieldCards(IsMonster)
	playerQueue := model.NewPriorityQueue(func(c1, c2 *Card) bool {
		return c1.GetValue() > c2.GetValue()
	})
	for i := 0; i < len(playerMonsters); i++ {
		playerQueue.Add(playerMonsters[i])
	}
	monsters := p.GetFieldCards(IsMonster)
	queue := model.NewPriorityQueue(func(c1, c2 *Card) bool {
		return c1.GetValue() > c2.GetValue()
	})
	for i := 0; i < len(monsters); i++ {
		queue.Add(monsters[i])
		monsters[i].Defense = false // 先全部调整为 进攻 防御 后面 再调整
	}
	has := false
	for !playerQueue.IsEmpty() && !queue.IsEmpty() {
		if queue.Peek().GetValue() >= playerQueue.Poll().GetValue() {
			queue.Poll()
		} else {
			has = true
		}
	}
	if has { // 存在打不死的 且自己还有 就防御
		for !queue.IsEmpty() {
			queue.Poll().Defense = true
		}
	} // 特殊 攻击字段怪兽处理
	cards := p.GetFieldCards(And(IsMonster, HasField(FieldAttack)))
	for i := 0; i < len(cards); i++ {
		cards[i].Defense = false
	}
	return false
}

func (p *Player) computerAttack() bool {
	monsters := p.GetFieldCards(And(IsMonster, CanAttackSrc))
	if len(monsters) <= 0 {
		return false
	}
	maxMonster := monsters[0]
	for i := 1; i < len(monsters); i++ {
		if monsters[i].GetValue() > maxMonster.GetValue() {
			maxMonster = monsters[i]
		}
	}
	// 尝试 进攻
	playerMonsters := PlayerManager.Player.GetFieldCards(IsMonster)
	if len(playerMonsters) <= 0 { // 直接攻击
		ActionManager.PushAction(&Param{EventType: EventTypeMonsterAttack, Player: p,
			SrcCard: maxMonster}, NewWarpHandle(CreateMonsterAttackAction))
		return true
	}
	playerQueue := model.NewPriorityQueue(func(c1, c2 *Card) bool {
		return c1.GetValue() > c2.GetValue()
	})
	for i := 0; i < len(playerMonsters); i++ {
		playerQueue.Add(playerMonsters[i])
	}
	for !playerQueue.IsEmpty() {
		target := playerQueue.Poll()
		if maxMonster.GetValue() >= target.GetValue() && target.CanSelect(&Param{EventType: EventTypeSelectTar, Card: maxMonster}) {
			ActionManager.PushAction(&Param{EventType: EventTypeMonsterAttack, Player: p,
				SrcCard: maxMonster, TarCard: target}, NewWarpHandle(CreateMonsterAttackAction))
			return true
		}
	}
	return false // 当前最大的 什么 也打不了 只能直接过了
}

func (p *Player) ChangeHp(value int) {
	p.Hp += value
	if p.Hp <= 0 {
		StackRoom.PushLayer(NewWinUILayer(!p.player, "对手生命值降到0"))
		return
	}
}

func (p *Player) Click(pos complex128) bool {
	return utils.PointCollision(p, pos)
}

func (p *Player) GetPileCards(cardPileType int, filter CardFilter) []*Card {
	cards := p.cardPiles[cardPileType].GetAllCards()
	res := make([]*Card, 0)
	for i := 0; i < len(cards); i++ {
		if filter(cards[i]) {
			res = append(res, cards[i])
		}
	}
	return res
}

func (p *Player) RemovePileCard(cardPileType int, card *Card) {
	p.cardPiles[cardPileType].RemoveCard(card)
}

func (p *Player) GetPileCard(cardPileType int) *CardPile {
	return p.cardPiles[cardPileType]
}
