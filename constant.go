/*
@author: sk
@date: 2023/2/4
*/
package main

import (
	"GameBase2/config"
	R "tcg/res"

	"golang.org/x/image/font"
)

const (
	CardPileTypeExcept = iota
	CardPileTypeCemetery
	CardPileTypeExtra
	CardPileTypeDeck
	CardPileTypeNone
)

var (
	Font21 font.Face
	Font24 font.Face
	Font36 font.Face
	Font72 font.Face
)

func InitFont() {
	Font21 = config.FontFactory.CreateFont(R.RAW.IPIX, 36, 21)
	Font24 = config.FontFactory.CreateFont(R.RAW.IPIX, 36, 24)
	Font36 = config.FontFactory.CreateFont(R.RAW.IPIX, 36, 36)
	Font72 = config.FontFactory.CreateFont(R.RAW.IPIX, 36, 72)
}

var (
	Color255_93_0    = Color(255, 93, 0)
	Color63_194_96   = Color(63, 194, 96)
	Color65_147_217  = Color(65, 147, 217)
	Color163_167_166 = Color(163, 167, 166)
	Color0_0_64      = Color(0, 0, 64)
	Color64_0_0      = Color(64, 0, 0)
	Color163_101_70  = Color(163, 101, 70)
	Color61_138_129  = Color(61, 138, 129)
	Color160_72_130  = Color(160, 72, 130)
	Color0_0_0_127   = AlphaColor(0, 0, 0, 127)
	Color71_159_144  = Color(71, 159, 144)
)

var (
	YesOrNo = []string{"Yes", "No"}
)

const (
	CardSize    = 75 + 110i
	TipInterval = 40
	ThinkTime   = 30
	MaxCardNum  = 6
)

const (
	CardTypeMonster = iota
	CardTypeMagic
	CardTypeTrap
	CardTypeNone
)

const (
	MonsterTypeCommon = iota
	MonsterTypeEffect
	MonsterTypeNone
)

const ( // 先仅使用部分
	RaceAngel = iota
	RaceUndead
	RaceSoldier
	RacePlant
	RaceMachine
	RaceDragon
	RaceHeiZi
	RaceNone
)

const ( // 先仅使用部分
	NatureWind = iota
	NatureWater
	NatureFire
	NatureLand
	NatureLight
	NatureGod
	NatureNone
)

const ( // 字段先使用数字表示 用于测试
	FieldBlueEye = 1 << iota
	Field2
	Field3
	Field4
	Field5
	FieldBlueEyeWhiteDragon // 青眼百龙 字段   方便 一些  视为的技能
	FieldAttack             // 电脑 用来 区分 状态的
	FieldBand               // 结束乐队  字段
	FieldNone
)

const (
	MagicTypeCommon = iota
	MagicTypeForever
	MagicTypeEquip
	MagicTypeRush
	MagicTypeNone
)

const (
	TrapTypeCommon = iota
	TrapTypeForever
	TrapTypeNone
)

const (
	PhaseDraw = 1 << iota
	PhasePrepare
	PhasePlay
	PhaseAttack
	PhaseDiscard
	PhaseEnd
	PhaseNone
)

const ( // 必须 与 CardPileXxx 顺序一致
	PlaceExcept = iota
	PlaceCemetery
	PlaceExtra
	PlaceDeck
	PlaceHand // 在手牌中
	PlaceField
	PlaceNone
)

const ( // 前面几个必须与 CardPileXxx 顺序一致 方便 使用
	EventTypeExceptCard   = iota // 选择 除外 中的卡事件
	EventTypeCemeteryCard        // 选择墓地的卡事件
	EventTypeExtraCard           // 选择额外的卡事件
	EventTypeDeckCard            // 选择卡组的卡事件
	// 这两个 顺序与 PlaceXxx 一致 方便使用
	EventTypeHandCard // 选择手牌 的事件
	EventTypeFieldCard
	EventTypeNormalSummon // 通常召唤
	EventTypePreparePhase // 玩家准备
	EventTypeMonsterAttack
	EventTypeBeforeAttack  // 战斗宣言
	EventTypeAfterAttack   // 战斗结算 完毕 打算 执行结果
	EventTypeSummonSuccess // 召唤成功 各种 都算
	EventTypeGoCardPile    // 卡牌 进入 四种牌堆
	EventTypeEndPhase      // 回合结束
	EventTypeAttackPhase   // 攻击阶段
	EventTypeSelectSrc
	EventTypeSelectTar
	EventTypeDiscardCard
	EventTypeNone
)
