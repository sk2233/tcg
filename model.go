/*
@author: sk
@date: 2023/2/4
*/
package main

type CardData struct { // 存储卡牌的基本数据
	// 共用属性
	Name       string
	Desc       string
	Fields     int // 多个 & 关系  每张卡 都可以有字段
	CardType   int
	Actions    []string // 效果 处理对象 名称
	CanSelects []string // 是否 可选
	// 怪兽卡
	MonsterType  int
	Atk, Def     int
	Level        int // 最高 12星
	Race, Nature int // 枚举
	// 魔法卡
	MagicType int
	// 陷阱卡
	TrapType int
}
