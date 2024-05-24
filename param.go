/*
@author: sk
@date: 2023/2/5
*/
package main

type Param struct {
	// 固定有的
	EventType int
	ActionEnd bool // 当前行为是否结束
	EventEnd  bool // 是否提前结束 在ActionGroup中使用
	// 经常使用的
	Player *Player
	// 非固定使用
	CardPile         *CardPile
	PlayPhase        *PlayPhase
	Card             *Card
	Invalid          bool       // 各种效果  动作 是否无效
	CardFilter       CardFilter // 选择 单张 卡的 过滤器
	SrcCard, TarCard *Card      // 攻击 双方  若TarCard 为nil 就是 直接攻击 对方
	HurtValue        int        // 对玩家造成对伤害值   正值  对对方伤害  负值 伤害自己
	CardPileType     int
}
