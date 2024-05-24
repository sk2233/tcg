/*
@author: sk
@date: 2023/2/5
*/
package main

import (
	"GameBase2/config"
	"GameBase2/utils"
	"math/rand"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

var (
	handleFactories = make(map[string]HandleFactory)
)

func CreateHandle(name string, player *Player, card *Card) IHandle {
	return handleFactories[name](player, card)
}

func init() {
	// 通用效果
	handleFactories["通常召唤"] = createNormalSummon
	handleFactories["切换状态"] = createSwitchState
	handleFactories["怪兽准备"] = createMonsterPrepare
	handleFactories["覆盖"] = createSetTrap
	// 个别效果
	handleFactories["青眼的贤士-检索白石"] = createRetrieveStone
	handleFactories["青眼的贤士-召唤青眼怪兽"] = createSummonBlueEye1
	handleFactories["太古的白石-召唤青眼怪兽"] = createSummonBlueEye2
	handleFactories["太古的白石-拣回青眼怪兽"] = createPickBlueEye
	handleFactories["海龟坏兽-送给对面"] = createSummonToEnemy
	handleFactories["青眼亚白龙-展示召唤"] = createShowSummon
	handleFactories["青眼亚白龙-破坏怪兽"] = createDestroyMonster
	handleFactories["白灵龙-除外陷阱"] = creatExceptTrap
	handleFactories["白灵龙-特招青眼白龙"] = createSummonBlueEyeWhiteDragon
	handleFactories["增值的G-发动效果"] = createSummonDraw
	handleFactories["无限泡影-直接发动"] = createUseDirect
	handleFactories["无限泡影-覆盖发动"] = createUseSet
	handleFactories["羽毛扫-破坏魔法陷阱"] = createDestroyMagicTrap
	handleFactories["交易进行-抽牌"] = createDrawCard1
	handleFactories["龙旋律-检索怪兽"] = createRetrieveMonster
	handleFactories["调和的宝牌-抽牌"] = createDrawCard2
	handleFactories["龙之灵庙-堆墓"] = createMoveCemetery
	handleFactories["复活的福音-复活怪兽"] = createReviveMonster
	handleFactories["复活的福音-抵挡破坏"] = createStopDestroy
	handleFactories["技能抽取-放置"] = createSetField
	handleFactories["青眼混沌MAX龙-召唤"] = createSummonMax
	handleFactories["青眼混沌MAX龙-2倍穿防"] = createHurtDouble
	handleFactories["青眼双爆裂龙-召唤"] = createSummonDouble
	handleFactories["青眼双爆裂龙-免战破"] = createNoAttackDestroy
	handleFactories["青眼双爆裂龙-攻击2次"] = createAttackDouble
	handleFactories["苍眼银龙-召唤"] = createSummonDragon
	handleFactories["苍眼银龙-召唤通常怪兽"] = createSummonCommon
	handleFactories["联结栗子球-召唤"] = createSummonChestnut
	handleFactories["联结栗子球-减攻击力"] = createReduceAttack
	// 敌人的卡牌效果  (为了方便，做出自动发动)
	handleFactories["东雪莲-免疫攻击"] = createImmuneAttack
	handleFactories["嘉然-嘉然小姐的狗"] = createDianaDog
	handleFactories["摇摆阳-召唤伴舞"] = createSummonDancer
	handleFactories["摇摆阳-摇摆攻击"] = createRockAttack
	handleFactories["喜多的心动魔法-回血"] = createSummonRecover
	handleFactories["你就是歌姬吧-降低血量"] = createReduceHp
	handleFactories["九转大肠-减少标记"] = createReduceMark
	handleFactories["十七张牌你能秒我-发动效果"] = createDraw17Card
	handleFactories["你这是违法行为-全部破坏"] = createDestroyAll
	handleFactories["蔡徐坤-变小黑子"] = createMakeHeiZi
	handleFactories["蔡徐坤-篮球邀约"] = createPlayBasketball
	handleFactories["蔡徐坤-打小黑子"] = createAttackHeiZi
}

//====================================

func createAttackHeiZi(player *Player, card *Card) IHandle {
	return &AttackHeiZiHandle{card: card}
}

type AttackHeiZiHandle struct {
	player *Player
	card   *Card
	target *Card
}

func (a *AttackHeiZiHandle) CanHandle(param *Param) bool {
	if param.EventType == EventTypeAfterAttack && a.player != param.Player && a.target != nil {
		a.target.Atk = a.target.Data.Atk
		a.target = nil
	}
	return param.EventType == EventTypeBeforeAttack && a.player != param.Player && param.TarCard == a.card &&
		param.SrcCard.Race == RaceHeiZi
}

func (a *AttackHeiZiHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, a.AttackHeiZi)
}

func (a *AttackHeiZiHandle) AttackHeiZi(param *Param) {
	Tip.AddTip("[蔡徐坤]效果发动，对方攻击力减半")
	a.target = param.SrcCard
	a.target.Atk /= 2
}

//================PlayBasketballHandle==============

func createPlayBasketball(player *Player, card *Card) IHandle {
	return &PlayBasketballHandle{player: player, card: card}
}

type PlayBasketballHandle struct {
	player *Player
	card   *Card
}

func (p *PlayBasketballHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeAttackPhase && param.Player == p.player && p.card.Place == PlaceField &&
		p.card.AttackNum > 0
}

func (p *PlayBasketballHandle) CreateAction(param *Param) IAction {
	p.card.AttackNum = 0
	return NewPlayBasketballAction(param)
}

//===============MakeHeiZiHandle===================

func createMakeHeiZi(player *Player, card *Card) IHandle {
	return &MakeHeiZiHandle{card: card, player: player}
}

type MakeHeiZiHandle struct {
	card   *Card
	player *Player
}

func (m *MakeHeiZiHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypePreparePhase && param.Player == m.player && m.card.Place == PlaceField
}

func (m *MakeHeiZiHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, m.MakeHeiZi)
}

func (m *MakeHeiZiHandle) MakeHeiZi(param *Param) {
	cards := PlayerManager.GetAnother(m.player).GetFieldCards(And(IsMonster, Not(RaceEq(RaceHeiZi))))
	if len(cards) > 0 {
		Tip.AddTip("[蔡徐坤]发动能力，把一只怪兽修改为黑子种族")
		cards[0].Race = RaceHeiZi
	}
}

//===================DestroyAllHandle================

func createDestroyAll(player *Player, card *Card) IHandle {
	return &DestroyAllHandle{player: player, card: card}
}

type DestroyAllHandle struct {
	player *Player
	card   *Card
	count  int
}

func (d *DestroyAllHandle) CanHandle(param *Param) bool {
	if d.card.Place != PlaceField || !d.card.CanUse {
		return false
	}
	if param.EventType == EventTypePreparePhase && param.Player != d.player {
		d.count = 0
		return false
	}
	if param.EventType == EventTypeSummonSuccess && param.Player != d.player && param.Card.Data.Level >= 8 {
		d.count++
	}
	return d.count >= 3
}

func (d *DestroyAllHandle) CreateAction(param *Param) IAction {
	d.count = 0
	return NewSimpleFuncAction(param, d.DestroyAll)
}

func (d *DestroyAllHandle) DestroyAll(param *Param) {
	Tip.AddTip("[你这是违法行为]效果发动,破坏对手所以怪兽")
	Info.SetCard(d.card)
	cards := param.Player.GetFieldCards(IsMonster)
	for i := 0; i < len(cards); i++ {
		param.Player.RemoveCardUI(cards[i])
		param.Player.AddCardPile(cards[i], CardPileTypeCemetery)
	}
	d.player.RemoveCardUI(d.card)
	d.player.AddCardPile(d.card, CardPileTypeCemetery)
}

//==================Draw17CardHandle==================

func createDraw17Card(player *Player, card *Card) IHandle {
	return &Draw17CardHandle{player: player}
}

type Draw17CardHandle struct {
	player *Player
	enable bool
}

func (d *Draw17CardHandle) CanHandle(param *Param) bool {
	if d.enable {
		if param.EventType == EventTypeDiscardCard {
			Tip.AddTip("[十七张牌你能秒我]效果发动")
			d.player.ChangeHp(200)
			d.player.RemovePileCard(CardPileTypeCemetery, param.Card)
			d.player.AddCardPile(param.Card, CardPileTypeDeck)
		} else if param.EventType == EventTypeEndPhase {
			d.enable = false
		}
	}
	return param.EventType == EventTypeHandCard && d.player.Hp > 4000 &&
		len(d.player.GetHandCards(IsAny))+len(d.player.GetPileCards(CardPileTypeDeck, IsAny)) > 17
}

func (d *Draw17CardHandle) CreateAction(param *Param) IAction {
	d.enable = true
	return NewSimpleFuncAction(param, d.Draw17Card)
}

func (d *Draw17CardHandle) Draw17Card(param *Param) {
	Tip.AddTip("[十七张牌你能秒我]效果发动")
	d.player.ChangeHp(-4000)
	d.player.DrawCards(17 - len(d.player.GetHandCards(IsAny)))
}

//=================ReduceMarkHandle======================

func createReduceMark(player *Player, card *Card) IHandle {
	return &ReduceMarkHandle{player: player, mark: 9, current: player, card: card}
}

type ReduceMarkHandle struct {
	player  *Player
	current *Player
	card    *Card
	mark    int
}

func (r *ReduceMarkHandle) DrawMark(pos complex128, screen *ebiten.Image) {
	utils.FillCircle(screen, pos, 8, 12, config.ColorRed)
	utils.DrawAnchorText(screen, strconv.Itoa(r.mark), pos, 0.5+0.5i, Font24, colornames.Black)
}

func (r *ReduceMarkHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeEndPhase && r.card.Place == PlaceField
}

func (r *ReduceMarkHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, r.ReduceMark)
}

func (r *ReduceMarkHandle) ReduceMark(param *Param) {
	Tip.AddTip("[九转大肠]效果发动")
	r.mark--
	if r.mark <= 0 {
		StackRoom.PushLayer(NewWinUILayer(r.player.player, "因[九转大肠]的效果胜利"))
		return
	}
	target := PlayerManager.GetAnother(r.current)
	cardArea := target.GetEmptyCardArea(false)
	if cardArea == nil { // 没有区域就炸出一张
		cards := target.GetFieldCards(Not(IsMonster))
		card := cards[rand.Intn(len(cards))]
		target.RemoveCardUI(card)
		target.AddCardPile(card, CardPileTypeCemetery)
		cardArea = target.GetEmptyCardArea(false)
	} // 转到 对方区域
	cardUI := r.current.RemoveCardUI(r.card)
	cardUI.player = target.player
	cardArea.SetCardUI(cardUI)
	r.current = target
}

//===================ReduceHpHandle====================

func createReduceHp(player *Player, card *Card) IHandle {
	return &ReduceHpHandle{player: player}
}

type ReduceHpHandle struct {
	player *Player
	target *Player
	oldHp  int
}

func (r *ReduceHpHandle) CanHandle(param *Param) bool {
	if param.EventType == EventTypeEndPhase && r.target != nil {
		r.target.Hp = r.oldHp
		r.target = nil
	}
	return param.EventType == EventTypeHandCard && r.player.Hp > 500
}

func (r *ReduceHpHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, r.ReduceHp)
}

func (r *ReduceHpHandle) ReduceHp(param *Param) {
	Tip.AddTip("[你就是歌姬吧]效果发动，降低敌人HP到500")
	r.player.Hp -= 500
	r.target = PlayerManager.GetAnother(r.player)
	r.oldHp = r.target.Hp
	r.target.Hp = 500
}

//====================SummonRecoverHandle========================

func createSummonRecover(player *Player, card *Card) IHandle {
	return &SummonRecoverHandle{player: player, card: card}
}

type SummonRecoverHandle struct {
	player *Player
	card   *Card
}

func (s *SummonRecoverHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeSummonSuccess && s.card.Place == PlaceField && s.player == param.Player
}

func (s *SummonRecoverHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, s.SummonRecover)
}

func (s *SummonRecoverHandle) SummonRecover(param *Param) {
	Tip.AddTip("[喜多的心动魔法]效果发动，恢复500HP")
	s.player.ChangeHp(500)
}

//==================RockAttackHandle=========================

func createRockAttack(player *Player, card *Card) IHandle {
	return &RockAttackHandle{card: card, player: player}
}

type RockAttackHandle struct {
	card   *Card
	player *Player
}

func (r *RockAttackHandle) CanHandle(param *Param) bool {
	if param.EventType == EventTypeAfterAttack && param.Card == r.card {
		r.card.Atk = r.card.Data.Atk
		return false
	}
	return param.EventType == EventTypeBeforeAttack && param.Card == r.card
}

func (r *RockAttackHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, r.RockAttack)
}

func (r *RockAttackHandle) RockAttack(param *Param) {
	Tip.AddTip("[摇摆阳]发动效果，攻击力上升")
	cardAreas := r.player.GetCardAreas(true)
	for i := 0; i < len(cardAreas); i++ {
		if cardAreas[i].cardUI != nil && cardAreas[i].cardUI.Card == r.card {
			if i > 0 && cardAreas[i-1].cardUI != nil && cardAreas[i-1].cardUI.Card.Data.Name == "摇摆伴舞" {
				r.card.Atk += 1000
			}
			if i < len(cardAreas)-1 && cardAreas[i+1].cardUI != nil && cardAreas[i+1].cardUI.Card.Data.Name == "摇摆伴舞" {
				r.card.Atk += 1000
			}
			return
		}
	}
}

//====================SummonDancerHandle==================

func createSummonDancer(player *Player, card *Card) IHandle {
	return &SummonDancerHandle{player: player, card: card}
}

type SummonDancerHandle struct {
	player *Player
	card   *Card
}

func (s *SummonDancerHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypePreparePhase && s.player == param.Player && s.card.Place == PlaceField &&
		((len(s.player.GetHandCards(NameEq("摇摆伴舞"))) > 0 && s.player.GetEmptyCardArea(true) != nil &&
			len(s.player.GetFieldCards(NameEq("摇摆伴舞"))) < 2) ||
			len(s.player.GetPileCards(CardPileTypeDeck, NameEq("摇摆伴舞"))) > 0)
}

func (s *SummonDancerHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, s.SummonDancer)
}

func (s *SummonDancerHandle) SummonDancer(param *Param) {
	cards := s.player.GetHandCards(NameEq("摇摆伴舞"))
	cardArea := s.player.GetEmptyCardArea(true)
	if len(cards) > 0 && cardArea != nil && len(s.player.GetFieldCards(NameEq("摇摆伴舞"))) < 2 {
		Tip.AddTip("[摇摆阳]发动效果，从手牌特殊召唤[摇摆伴舞]")
		cardArea.SetCardUI(s.player.RemoveCardUI(cards[0]))
		return
	}
	cards = s.player.GetPileCards(CardPileTypeDeck, NameEq("摇摆伴舞"))
	if len(cards) > 0 {
		Tip.AddTip("[摇摆阳]发动效果，从手牌检索一张[摇摆伴舞]")
		s.player.RemovePileCard(CardPileTypeDeck, cards[0])
		s.player.AddHandCards(cards[0])
	}
}

//=================DianaDogHandle================

func createDianaDog(player *Player, card *Card) IHandle {
	return &DianaDogHandle{player: player, card: card}
}

type DianaDogHandle struct {
	player *Player
	card   *Card
}

func (d *DianaDogHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeAttackPhase && param.Player == d.player && d.card.Place == PlaceField &&
		d.card.AttackNum > 0 && len(d.player.GetHandCards(IsAny)) > 0 &&
		len(PlayerManager.GetAnother(d.player).GetFieldCards(IsMonster)) > 0
}

func (d *DianaDogHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, d.DianaDog)
}

func (d *DianaDogHandle) DianaDog(param *Param) {
	Tip.AddTip("[嘉然]发动效果，把敌方一只怪兽变成[狗]衍生物")
	d.card.AttackNum = 0
	cards := d.player.GetHandCards(IsAny)
	d.player.RemoveCardUI(cards[0])
	d.player.AddCardPile(cards[0], CardPileTypeCemetery)
	enemy := PlayerManager.GetAnother(d.player)
	// 移除怪兽
	cards = enemy.GetFieldCards(IsMonster)
	enemy.RemoveCardUI(cards[0])
	card := NewCard(d.CreateDog(), enemy, PlaceField)
	enemy.GetEmptyCardArea(true).SetCardUI(NewCardUI(card, enemy.player))
	Info.SetCard(card)
}

func (d *DianaDogHandle) CreateDog() *CardData {
	return &CardData{
		Name:        "狗",
		Desc:        "衍生物",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeCommon,
		Atk:         500,
		Def:         500,
		Level:       1,
		Race:        RaceUndead,
		Nature:      NatureFire,
		Fields:      Field5,
	}
}

//===============ImmuneAttackHandle===================

func createImmuneAttack(player *Player, card *Card) IHandle {
	return &ImmuneAttackHandle{card: card, player: player}
}

type ImmuneAttackHandle struct {
	card   *Card
	player *Player
}

func (r *ImmuneAttackHandle) CanHandle(param *Param) bool { // 自己非防御  被攻击 且 对方 比自己大
	return param.EventType == EventTypeBeforeAttack && !r.card.Defense && param.TarCard == r.card &&
		param.SrcCard.GetValue() > r.card.GetValue()
}

func (r *ImmuneAttackHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("[东雪莲]发动效果，无效该攻击")
	return NewSimpleFuncAction(param, r.ImmuneAttack)
}

func (r *ImmuneAttackHandle) ImmuneAttack(param *Param) {
	param.Invalid = true
	r.card.Defense = true
	cardArea := r.player.GetEmptyCardArea(true)
	cards := r.player.GetPileCards(CardPileTypeExtra, NameEq("孙笑川"))
	if cardArea == nil || len(cards) == 0 {
		return
	}
	Tip.AddTip("[东雪莲]发动效果，特殊召唤[孙笑川]")
	Info.SetCard(cards[0])
	r.player.RemovePileCard(CardPileTypeExtra, cards[0])
	cardArea.SetCardUI(NewCardUI(cards[0], r.player.player))
}

//=================ReduceAttackHandle==================

func createReduceAttack(player *Player, card *Card) IHandle {
	return &ReduceAttackHandle{card: card, player: player}
}

type ReduceAttackHandle struct {
	card   *Card
	player *Player
}

func (r *ReduceAttackHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeBeforeAttack && param.Player != r.player && r.card.Place == PlaceField // 对方攻击 尝试发动
}

func (r *ReduceAttackHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("[联结栗子球]是否发动效果把攻击怪兽攻击力变为0")
	return NewSelectAction(param, YesOrNo, r.ReduceAttack, EmptyFunc)
}

func (r *ReduceAttackHandle) ReduceAttack(param *Param) {
	param.SrcCard.Atk = 0 //
	r.player.RemoveCardUI(r.card)
	r.player.AddCardPile(r.card, CardPileTypeCemetery)
	if param.TarCard == r.card { // 若攻击 自己 取消攻击
		param.Invalid = true
	}
}

//========================SummonCommonHandle============================

func createSummonCommon(player *Player, card *Card) IHandle {
	return &SummonCommonHandle{player: player, card: card}
}

type SummonCommonHandle struct {
	player *Player
	card   *Card
}

func (s *SummonCommonHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypePreparePhase && param.Player == s.player && s.card.Place == PlaceField &&
		len(s.player.GetPileCards(CardPileTypeCemetery, IsCommonMonster)) > 0 &&
		s.player.GetEmptyCardArea(true) != nil
}

func (s *SummonCommonHandle) CreateAction(param *Param) IAction {
	return NewSummonCommonAction(param)
}

//======================AttackDoubleHandle========================

func createAttackDouble(player *Player, card *Card) IHandle {
	return &AttackDoubleHandle{player: player, card: card}
}

type AttackDoubleHandle struct {
	player *Player
	card   *Card
}

func (m *AttackDoubleHandle) Order() int {
	return 2233
}

func (m *AttackDoubleHandle) CanHandle(param *Param) bool { // 自己的准备 阶段 在场上的 怪兽
	return param.EventType == EventTypePreparePhase && m.player == param.Player && m.card.Place == PlaceField
}

func (m *AttackDoubleHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, m.AttackDouble)
}

func (m *AttackDoubleHandle) AttackDouble(param *Param) {
	m.card.AttackNum = 2
}

//===================NoAttackDestroyHandle================

func createNoAttackDestroy(player *Player, card *Card) IHandle {
	return &NoAttackDestroyHandle{card: card}
}

type NoAttackDestroyHandle struct {
	card *Card
}

func (h *NoAttackDestroyHandle) CanHandle(param *Param) bool { // 攻击结算后   自己不会被战斗破坏
	return param.EventType == EventTypeAfterAttack && (param.TarCard == h.card || param.SrcCard == h.card)
}

func (h *NoAttackDestroyHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("[青眼双爆裂龙]不会被战斗破坏")
	return NewSimpleFuncAction(param, h.NoAttackDestroy)
}

func (h *NoAttackDestroyHandle) NoAttackDestroy(param *Param) {
	if param.SrcCard == h.card {
		param.SrcCard = nil
	} else {
		param.TarCard = nil
	}
}

//======================HurtDoubleHandle============================

func createHurtDouble(player *Player, card *Card) IHandle {
	return &HurtDoubleHandle{player: player, card: card}
}

type HurtDoubleHandle struct {
	card   *Card
	player *Player
}

func (h *HurtDoubleHandle) CanHandle(param *Param) bool { // 自己攻击结算后  被攻击方处于防御状态
	return param.EventType == EventTypeAfterAttack && h.player == param.Player && h.card == param.Card &&
		param.SrcCard == nil && param.TarCard.Defense // 自己不能死
}

func (h *HurtDoubleHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("[青眼混沌MAX龙]对防御怪兽进行2倍穿防")
	return NewSimpleFuncAction(param, HurtDouble)
}

func HurtDouble(param *Param) {
	param.HurtValue = (param.Card.GetValue() - param.TarCard.GetValue()) * 2
}

//======================SummonChestnutHandle=========================

func createSummonChestnut(player *Player, card *Card) IHandle {
	return &SummonChestnutHandle{}
}

type SummonChestnutHandle struct {
}

func (s *SummonChestnutHandle) GetName() string {
	return "特殊召唤"
}

func (s *SummonChestnutHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeExtraCard && len(param.Player.GetFieldCards(LevelEq(1))) > 1
}

func (s *SummonChestnutHandle) CreateAction(param *Param) IAction {
	return NewSummonChestnutAction(param)
}

//==========================SummonDragonHandle==============================

func createSummonDragon(player *Player, card *Card) IHandle {
	return &SummonDragonHandle{}
}

type SummonDragonHandle struct {
}

func (s *SummonDragonHandle) GetName() string {
	return "特殊召唤"
}

func (s *SummonDragonHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeExtraCard && len(CalculatePlan(param.Player.GetFieldCards(IsMonster))) > 0
}

func CalculatePlan(cards []*Card) [][]*Card {
	return Next(cards, 0, make([]*Card, 0), 0)
}

func Next(cards []*Card, index int, temp []*Card, level int) [][]*Card {
	if level < 9 && index < len(cards) {
		res := Next(cards, index+1, temp, level)                                                              // 不要
		res = append(res, Next(cards, index+1, append(temp, cards[index]), level+cards[index].Data.Level)...) // 要的情况
		return res
	} else if level == 9 && len(temp) >= 2 {
		for i := 0; i < len(temp); i++ {
			if temp[i].Data.Fields&FieldBlueEye > 0 {
				return [][]*Card{temp}
			}
		}
	}
	return [][]*Card{}
}

func (s *SummonDragonHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("请选择一种组合召唤方式")
	res := CalculatePlan(param.Player.GetFieldCards(IsMonster))
	names0 := make([]string, 0)
	actions := make([]func(*Param), 0)
	for i := 0; i < len(res); i++ {
		temp := res[i]
		buf := strings.Builder{}
		for j := 0; j < len(temp); j++ {
			if j > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(temp[j].Data.Name)
		}
		names0 = append(names0, buf.String())
		actions = append(actions, func(param *Param) {
			player := param.Player // 献祭
			for j := 0; j < len(temp); j++ {
				player.RemoveCardUI(temp[j])
				player.AddCardPile(temp[j], CardPileTypeCemetery)
			}
			player.RemovePileCard(CardPileTypeExtra, param.Card) // 召唤
			player.GetEmptyCardArea(true).SetCardUI(NewCardUI(param.Card, player.player))
			ActionManager.TriggerEvent(&Param{EventType: EventTypeSummonSuccess, Player: player, Card: param.Card})
		})
	}
	names0 = append(names0, "取消")
	actions = append(actions, EmptyFunc)
	return NewSelectAction(param, names0, actions...)
}

//=====================SummonDoubleHandle==============================

func createSummonDouble(player *Player, card *Card) IHandle {
	return &SummonDoubleHandle{}
}

type SummonDoubleHandle struct {
}

func (s *SummonDoubleHandle) GetName() string {
	return "特殊召唤"
}

func (s *SummonDoubleHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeExtraCard && len(param.Player.GetFieldCards(HasField(FieldBlueEyeWhiteDragon))) > 1
}

func (s *SummonDoubleHandle) CreateAction(param *Param) IAction {
	return NewSummonDoubleAction(param)
}

//===================SummonMaxHandle=======================

func createSummonMax(player *Player, card *Card) IHandle {
	return &SummonMaxHandle{}
}

type SummonMaxHandle struct {
}

func (s *SummonMaxHandle) GetName() string {
	return "特殊召唤"
}

func (s *SummonMaxHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeExtraCard && len(param.Player.GetFieldCards(NameEq("青眼白龙"))) > 0
}

func (s *SummonMaxHandle) CreateAction(param *Param) IAction {
	return NewSummonMaxAction(param)
}

//================SetFieldHandle===================

func createSetField(player *Player, card *Card) IHandle {
	return &SetFieldHandle{}
}

type SetFieldHandle struct {
}

func (s *SetFieldHandle) GetName() string {
	return "发动效果"
}

func (s *SetFieldHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeHandCard && param.Player.Hp > 1000 &&
		param.Player.GetEmptyCardArea(false) != nil
}

func (s *SetFieldHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, SetField)
}

func SetField(param *Param) {
	param.Player.ChangeHp(-1000) // 简单 移动到 场上  技能与行动 混在 了一起 暂时  无法禁止
	param.Player.GetEmptyCardArea(false).SetCardUI(param.Player.RemoveCardUI(param.Card))
}

//===================StopDestroyHandle======================

func createStopDestroy(player *Player, card *Card) IHandle {
	return &StopDestroyHandle{card: card, player: player, canUse: false}
}

type StopDestroyHandle struct {
	card   *Card
	player *Player
	canUse bool
}

func (s *StopDestroyHandle) CanHandle(param *Param) bool {
	if param.EventType == EventTypeGoCardPile && param.Card == s.card {
		s.canUse = true // 检测进入 墓地
		return false
	}
	if !s.canUse || param.EventType != EventTypeAfterAttack {
		return false
	}
	if param.Player == s.player { // 自己发动的 攻击
		return param.SrcCard != nil
	} else { // 别人发动的
		return param.TarCard != nil
	}
}

func (s *StopDestroyHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("[复活的福音]是否防止我方怪兽被破坏")
	return NewSelectAction(param, YesOrNo, s.UseEffect, EmptyFunc)
}

func (s *StopDestroyHandle) UseEffect(param *Param) {
	s.canUse = false
	if param.Player == s.player { // 自己发动的 攻击
		param.SrcCard = nil
	} else { // 别人发动的
		param.TarCard = nil
	}
}

//=====================ReviveMonsterHandle=======================

func createReviveMonster(player *Player, card *Card) IHandle {
	return &ReviveMonsterHandle{}
}

type ReviveMonsterHandle struct {
}

func (r *ReviveMonsterHandle) GetName() string {
	return "发动效果"
}

func (r *ReviveMonsterHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeHandCard &&
		len(param.Player.GetPileCards(CardPileTypeCemetery, And(IsMonster, LevelEq(8), RaceEq(RaceDragon)))) > 0
}

func (r *ReviveMonsterHandle) CreateAction(param *Param) IAction {
	return NewReviveMonsterAction(param)
}

//=====================MoveCemeteryHandle========================

func createMoveCemetery(player *Player, card *Card) IHandle {
	return &MoveCemeteryHandle{}
}

type MoveCemeteryHandle struct {
}

func (m *MoveCemeteryHandle) GetName() string {
	return "发动效果"
}

func (m *MoveCemeteryHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeHandCard
}

func (m *MoveCemeteryHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, MoveCemetery)
}

func MoveCemetery(param *Param) {
	param.Player.RemoveCardUI(param.Card) // 先移除自己的卡
	param.Player.AddCardPile(param.Card, CardPileTypeCemetery)
	card := MoveCemeteryMove(param.Player)
	if card != nil && card.Data.MonsterType == MonsterTypeCommon {
		MoveCemeteryMove(param.Player)
	}
}

func MoveCemeteryMove(player *Player) *Card {
	Tip.AddTip("[龙之灵庙]把卡组中一张龙族怪兽送入墓地")
	cards := player.GetPileCards(CardPileTypeDeck, And(IsMonster, RaceEq(RaceDragon)))
	if len(cards) <= 0 {
		return nil
	}
	player.RemovePileCard(CardPileTypeDeck, cards[0])
	player.AddCardPile(cards[0], CardPileTypeCemetery)
	return cards[0]
}

//====================DrawCard2Handle====================

func createDrawCard2(player *Player, card *Card) IHandle {
	return &DrawCard2Handle{}
}

type DrawCard2Handle struct {
}

func (d *DrawCard2Handle) GetName() string {
	return "发动效果"
}

func (d *DrawCard2Handle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeHandCard && len(param.Player.GetHandCards(NameEq("太古的白石"))) > 0
}

func (d *DrawCard2Handle) CreateAction(param *Param) IAction {
	return NewDrawCard2Action(param)
}

//=======================RetrieveMonsterHandle=========================

func createRetrieveMonster(player *Player, card *Card) IHandle {
	return &RetrieveMonsterHandle{}
}

type RetrieveMonsterHandle struct {
}

func (r *RetrieveMonsterHandle) GetName() string {
	return "发动效果"
}

func (r *RetrieveMonsterHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeHandCard && len(param.Player.GetHandCards(IsAny)) > 1
}

func (r *RetrieveMonsterHandle) CreateAction(param *Param) IAction {
	return NewRetrieveMonsterAction(param)
}

//====================DrawCard1Handle====================

func createDrawCard1(player *Player, card *Card) IHandle {
	return &DrawCard1Handle{}
}

type DrawCard1Handle struct {
}

func (d *DrawCard1Handle) GetName() string {
	return "发动效果"
}

func (d *DrawCard1Handle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeHandCard &&
		len(param.Player.GetHandCards(And(IsHand, IsMonster, LevelEq(8)))) > 0
}

func (d *DrawCard1Handle) CreateAction(param *Param) IAction {
	return NewDrawCard1Action(param)
}

//======================DestroyMagicTrapHandle=======================

func createDestroyMagicTrap(player *Player, card *Card) IHandle {
	return &DestroyMagicTrapHandle{}
}

type DestroyMagicTrapHandle struct {
}

func (d *DestroyMagicTrapHandle) GetName() string {
	return "发动效果"
}

func (d *DestroyMagicTrapHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeHandCard
}

func (d *DestroyMagicTrapHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, DestroyMagicTrap)
}

func DestroyMagicTrap(param *Param) {
	param.Player.RemoveCardUI(param.Card) // 先移除自己的卡
	param.Player.AddCardPile(param.Card, CardPileTypeCemetery)
	enemy := PlayerManager.GetAnother(param.Player) // 再移除敌人的 陷阱 魔法卡
	cards := enemy.GetFieldCards(Not(IsMonster))
	for i := 0; i < len(cards); i++ {
		enemy.RemoveCardUI(cards[i])
		enemy.AddCardPile(cards[i], CardPileTypeCemetery)
	}
	Tip.AddTip("[羽毛扫]破坏对方所有魔法陷阱卡")
}

//========================UseSetHandle============================

func createUseSet(player *Player, card *Card) IHandle {
	return &UseSetHandle{card: card, player: player}
}

type UseSetHandle struct {
	card   *Card
	player *Player
}

func (u *UseSetHandle) GetName() string {
	return "发动"
}

func (u *UseSetHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeFieldCard && u.card.CanUse &&
		len(PlayerManager.GetAnother(u.player).GetFieldCards(IsMonster)) > 0
}

func (u *UseSetHandle) CreateAction(param *Param) IAction {
	u.player.RemoveCardUI(u.card) // 先移除自己的卡
	u.player.AddCardPile(u.card, CardPileTypeCemetery)
	Tip.AddTip("[无限泡影]破坏对方所有魔法陷阱卡")
	enemy := PlayerManager.GetAnother(u.player) // 再移除敌人的 陷阱 魔法卡
	cards := enemy.GetFieldCards(Not(IsMonster))
	for i := 0; i < len(cards); i++ {
		enemy.RemoveCardUI(cards[i])
		enemy.AddCardPile(cards[i], CardPileTypeCemetery)
	}
	Tip.AddTip("[无限泡影]请选择对面的一个怪兽破坏")
	return NewMoveCardAction(param, And(IsMonster, IsField), CardPileTypeCemetery)
}

//=================UseDirectHandle====================

func createUseDirect(player *Player, card *Card) IHandle {
	return &UseDirectHandle{player: player, card: card}
}

type UseDirectHandle struct {
	card   *Card
	player *Player
}

func (u *UseDirectHandle) GetName() string {
	return "直接发动"
}

func (u *UseDirectHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeHandCard && len(u.player.GetFieldCards(IsAny)) <= 0 &&
		len(PlayerManager.GetAnother(u.player).GetFieldCards(IsMonster)) > 0
}

func (u *UseDirectHandle) CreateAction(param *Param) IAction {
	u.player.RemoveCardUI(u.card) // 先移除自己的卡
	u.player.AddCardPile(u.card, CardPileTypeCemetery)
	Tip.AddTip("[无限泡影]请选择对面的一个怪兽破坏")
	return NewMoveCardAction(param, And(IsMonster, IsField), CardPileTypeCemetery)
}

//===================SetTrapHandle=====================

func createSetTrap(player *Player, card *Card) IHandle {
	return &SetTrapHandle{player: player, card: card}
}

type SetTrapHandle struct {
	player *Player
	card   *Card
}

func (s *SetTrapHandle) GetName() string {
	return "放置"
}

func (s *SetTrapHandle) CanHandle(param *Param) bool {
	if param.EventType == EventTypePreparePhase && s.player == param.Player {
		s.card.CanUse = true
	}
	return param.EventType == EventTypeHandCard && param.Card.Data.CardType == CardTypeTrap &&
		param.Player.GetEmptyCardArea(false) != nil
}

func (s *SetTrapHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, SetTrapFunc)
}

func SetTrapFunc(param *Param) {
	param.Card.Flip = true
	param.Card.CanUse = false
	param.Player.GetEmptyCardArea(false).SetCardUI(param.Player.RemoveCardUI(param.Card))
}

//==========SummonDrawHandle===========

func createSummonDraw(player *Player, card *Card) IHandle {
	return &SummonDrawHandle{enableDraw: false, player: player}
}

type SummonDrawHandle struct {
	player     *Player
	enableDraw bool
}

func (s *SummonDrawHandle) GetName() string {
	return "发动效果"
}

func (s *SummonDrawHandle) CanHandle(param *Param) bool {
	if s.enableDraw && param.EventType == EventTypeSummonSuccess && s.player != param.Player {
		s.player.DrawCards(1) // 别人召唤 成功 自己 摸牌
		return false
	}
	if s.enableDraw && param.EventType == EventTypePreparePhase && s.player == param.Player {
		s.enableDraw = false // 自己的回合开始 结束效果
		return false
	}
	return param.EventType == EventTypeHandCard // 在手牌中 发动
}

func (s *SummonDrawHandle) CreateAction(param *Param) IAction {
	s.enableDraw = true
	return NewSimpleFuncAction(param, SummonDrawFunc)
}

func SummonDrawFunc(param *Param) {
	param.Player.RemoveCardUI(param.Card)
	param.Player.AddCardPile(param.Card, CardPileTypeCemetery)
}

//===================SummonBlueEyeWhiteDragonHandle====================

func createSummonBlueEyeWhiteDragon(player *Player, card *Card) IHandle {
	return &SummonBlueEyeWhiteDragonHandle{}
}

type SummonBlueEyeWhiteDragonHandle struct {
}

func (d *SummonBlueEyeWhiteDragonHandle) GetName() string {
	return "特招青眼白龙"
}

func (d *SummonBlueEyeWhiteDragonHandle) CanHandle(param *Param) bool { // 必须 还有 攻击的能力
	return param.EventType == EventTypeFieldCard && param.Card.Data.CardType == CardTypeMonster &&
		len(PlayerManager.GetAnother(param.Player).GetFieldCards(IsMonster)) > 0 &&
		len(param.Player.GetHandCards(NameEq("青眼白龙"))) > 0
}

func (d *SummonBlueEyeWhiteDragonHandle) CreateAction(param *Param) IAction {
	return NewSummonBlueEyeWhiteDragonAction(param)
}

//================ExceptTrapHandle================

func creatExceptTrap(player *Player, card *Card) IHandle {
	return &ExceptTrapHandle{card: card}
}

type ExceptTrapHandle struct {
	card *Card
}

func (r *ExceptTrapHandle) CanHandle(param *Param) bool { // 自己召唤成功  且 对方场上 有陷阱 魔法
	return param.EventType == EventTypeSummonSuccess && param.Card == r.card &&
		len(PlayerManager.GetAnother(param.Player).GetFieldCards(Not(IsMonster))) > 0
}

func (r *ExceptTrapHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("[白灵龙]召唤成功，可以除外对方一张魔法或陷阱")
	return NewMoveCardAction(param, And(Not(IsMonster), IsField), CardPileTypeExcept)
}

//===================DestroyMonsterHandle====================

func createDestroyMonster(player *Player, card *Card) IHandle {
	return &DestroyMonsterHandle{}
}

type DestroyMonsterHandle struct {
}

func (d *DestroyMonsterHandle) GetName() string {
	return "破坏怪兽"
}

func (d *DestroyMonsterHandle) CanHandle(param *Param) bool { // 必须 还有 攻击的能力
	return param.EventType == EventTypeFieldCard && param.Card.Data.CardType == CardTypeMonster &&
		param.Card.AttackNum > 0 && len(PlayerManager.GetAnother(param.Player).GetFieldCards(IsMonster)) > 0
}

func (d *DestroyMonsterHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("[青眼亚白龙]请选择对面的一个怪兽破坏")
	param.Card.AttackNum--
	return NewMoveCardAction(param, And(IsMonster, IsField), CardPileTypeCemetery)
}

//===================ShowSummonHandle=====================

func createShowSummon(player *Player, card *Card) IHandle {
	return &ShowSummonHandle{}
}

type ShowSummonHandle struct {
}

func (s *ShowSummonHandle) GetName() string {
	return "展示召唤"
}

func (s *ShowSummonHandle) CanHandle(param *Param) bool { // 展示效果 不能 使用 自己当青眼白龙使用
	return param.EventType == EventTypeHandCard && len(param.Player.GetHandCards(NameEq("青眼白龙"))) > 0 &&
		param.Player.GetEmptyCardArea(true) != nil
}

func (s *ShowSummonHandle) CreateAction(param *Param) IAction {
	return NewShowSummonAction(param)
}

//===================SummonToEnemyHandle=====================

func createSummonToEnemy(player *Player, card *Card) IHandle {
	return &SummonToEnemyHandle{}
}

type SummonToEnemyHandle struct {
}

func (s *SummonToEnemyHandle) GetName() string {
	return "送给对面"
}

func (s *SummonToEnemyHandle) CanHandle(param *Param) bool {
	if param.EventType != EventTypeHandCard {
		return false
	} // 在手牌 中  且  对面 有怪兽 可以当祭品
	return len(PlayerManager.GetAnother(param.Player).GetFieldCards(IsMonster)) > 0
}

func (s *SummonToEnemyHandle) CreateAction(param *Param) IAction {
	return NewSummonToEnemyAction(param)
}

//====================PickBlueEyeHandle===================

func createPickBlueEye(player *Player, card *Card) IHandle {
	return &PickBlueEyeHandle{}
}

type PickBlueEyeHandle struct {
}

func (p *PickBlueEyeHandle) GetName() string {
	return "拣回青眼怪兽"
}

func (p *PickBlueEyeHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeCemeteryCard && param.Card.Place == PlaceCemetery &&
		len(param.Player.GetPileCards(CardPileTypeCemetery, And(IsMonster, HasField(FieldBlueEye), Except(param.Card)))) > 0
}

func (p *PickBlueEyeHandle) CreateAction(param *Param) IAction {
	return NewPickBlueEyeAction(param)
}

//========================SummonBlueEye2Handle======================

func createSummonBlueEye2(player *Player, card *Card) IHandle {
	return &SummonBlueEye2Handle{intoCemetery: false, player: player, card: card}
}

type SummonBlueEye2Handle struct {
	intoCemetery bool
	player       *Player
	card         *Card
}

func (s *SummonBlueEye2Handle) CanHandle(param *Param) bool {
	if param.EventType == EventTypeGoCardPile && param.Player == s.player && param.Card == s.card {
		s.intoCemetery = true // 检测进入 墓地
		return false
	}
	if s.intoCemetery && param.EventType == EventTypeEndPhase && param.Player == s.player &&
		len(param.Player.GetHandCards(And(IsMonster, HasField(FieldBlueEye)))) > 0 &&
		param.Player.GetEmptyCardArea(true) != nil {
		s.intoCemetery = false
		return true
	}
	return false
}

func (s *SummonBlueEye2Handle) CreateAction(param *Param) IAction {
	return NewSummonBlueEye2Action(param)
}

//===================SummonBlueEye1Handle========================

func createSummonBlueEye1(player *Player, card *Card) IHandle {
	return &SummonBlueEye1Handle{}
}

type SummonBlueEye1Handle struct {
}

func (s *SummonBlueEye1Handle) GetName() string {
	return "召唤青眼怪兽"
}

func (s *SummonBlueEye1Handle) CanHandle(param *Param) bool { // 在手牌 中  有怪兽 解放   且手牌 中  有 青眼怪兽
	return param.EventType == EventTypeHandCard && param.Card.Place == PlaceHand &&
		len(param.Player.GetFieldCards(IsMonster)) > 0 && len(param.Player.GetHandCards(And(IsMonster, HasField(FieldBlueEye)))) > 0
}

func (s *SummonBlueEye1Handle) CreateAction(param *Param) IAction {
	return NewSummonBlueEye1Action(param)
}

//================RetrieveStoneHandle================

func createRetrieveStone(player *Player, card *Card) IHandle {
	return &RetrieveStoneHandle{card: card}
}

type RetrieveStoneHandle struct {
	card *Card
}

func (r *RetrieveStoneHandle) CanHandle(param *Param) bool { // 自己召唤成功
	return param.EventType == EventTypeSummonSuccess && param.Card == r.card
}

func (r *RetrieveStoneHandle) CreateAction(param *Param) IAction {
	Tip.AddTip("[青眼的贤士]召唤成功，检索一张[太古的白石]到手牌中")
	return NewRetrieveCardAction(param, CardPileTypeDeck, NameEq("太古的白石"))
}

//===================MonsterPrepareHandle=========================

func createMonsterPrepare(player *Player, card *Card) IHandle {
	return &MonsterPrepareHandle{player: player, card: card}
}

type MonsterPrepareHandle struct {
	player *Player
	card   *Card
}

func (m *MonsterPrepareHandle) CanHandle(param *Param) bool { // 自己的准备 阶段 在场上的 怪兽
	return param.EventType == EventTypePreparePhase && m.player == param.Player && m.card.Place == PlaceField
}

func (m *MonsterPrepareHandle) CreateAction(param *Param) IAction {
	return NewMonsterPrepareAction(m.card, param)
}

//=======================SwitchStateHandle==========================

func createSwitchState(player *Player, card *Card) IHandle {
	return &SwitchStateHandle{}
}

type SwitchStateHandle struct {
}

func (s *SwitchStateHandle) GetName() string {
	return "切换状态"
}

func (s *SwitchStateHandle) CanHandle(param *Param) bool {
	return param.EventType == EventTypeFieldCard && param.Card.Data.CardType == CardTypeMonster && !param.Card.HasAdjust
}

func (s *SwitchStateHandle) CreateAction(param *Param) IAction {
	return NewSimpleFuncAction(param, SwitchStateFunc)
}

func SwitchStateFunc(param *Param) {
	param.Card.Defense = !param.Card.Defense
	param.Card.HasAdjust = true
}

//======================NormalSummonHandle=====================

func createNormalSummon(player *Player, card *Card) IHandle {
	return &NormalSummonHandle{}
}

type NormalSummonHandle struct {
}

func (n *NormalSummonHandle) GetName() string {
	return "通常召唤"
}

func (n *NormalSummonHandle) CanHandle(param *Param) bool { // 选择的是怪兽 卡  且 还有召唤次数
	if param.EventType != EventTypeHandCard || param.PlayPhase.SummonNum <= 0 {
		return false
	}
	costNum := GetCostNum(param.Card.Data.Level)
	if costNum > 0 { // 需要消耗祭品，保证足够  位置肯定够 需要 有
		return len(param.Player.GetFieldCards(IsMonster)) >= costNum
	}
	return param.Player.GetEmptyCardArea(true) != nil // 有位置
}

func (n *NormalSummonHandle) CreateAction(param *Param) IAction {
	return NewNormalSummonAction(param)
}

//===================WarpHandle========================

type WarpHandle struct {
	creator func(param *Param) IAction
}

func NewWarpHandle(creator func(param *Param) IAction) *WarpHandle {
	return &WarpHandle{creator: creator}
}

func (w *WarpHandle) CanHandle(param *Param) bool {
	return true
}

func (w *WarpHandle) CreateAction(param *Param) IAction {
	return w.creator(param)
}
