/*
@author: sk
@date: 2023/2/4
*/
package main

import "GameBase2/utils"

// GetYugiMutoCards Dark♂游戏
func GetYugiMutoCards() []*CardData {
	res := make([]*CardData, 0)
	res = append(res, &CardData{
		Name: "神圣",
	})
	return res
}

// GetBZhanCards B站鬼畜卡组
func GetBZhanCards() ([]*CardData, []*CardData) {
	deck := make([]*CardData, 0)
	deck = append(deck, Repeat(&CardData{
		Name:        "东雪莲",
		Desc:        "被攻击时可以变为防守状态取消改攻击，并可以从额外特殊召唤[孙笑川]",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         1500,
		Def:         2500,
		Level:       4,
		Race:        RaceUndead,
		Nature:      NatureFire,
		Fields:      FieldAttack,
		Actions:     []string{"东雪莲-免疫攻击"},
	}, 3)...)
	deck = append(deck, Repeat(&CardData{
		Name:        "嘉然",
		Desc:        "1,可以放弃攻击，并舍弃一张手牌，把一只敌方怪兽变成[狗]衍生物，2，不能成为[狗]衍生物的攻击目标",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         1500,
		Def:         1000,
		Level:       4,
		Race:        RaceUndead,
		Nature:      NatureFire,
		Fields:      Field5,
		Actions:     []string{"嘉然-嘉然小姐的狗"},
		CanSelects:  []string{"嘉然-被攻击"},
	}, 3)...)
	deck = append(deck, Repeat(&CardData{
		Name:        "摇摆阳",
		Desc:        "1，准备阶段可以从手牌特殊召唤一只[摇摆伴舞]或将一只[摇摆伴舞]加入手牌，场上最多存在2只，2，攻击阶段攻击力上升相邻[摇摆伴舞]的数量X1000",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         2000,
		Def:         1000,
		Level:       4,
		Race:        RaceUndead,
		Nature:      NatureFire,
		Fields:      Field5,
		Actions:     []string{"摇摆阳-召唤伴舞", "摇摆阳-摇摆攻击"},
	}, 3)...)
	deck = append(deck, Repeat(&CardData{
		Name:        "摇摆伴舞",
		Desc:        "大白板",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeCommon,
		Atk:         1000,
		Def:         500,
		Level:       2,
		Race:        RaceUndead,
		Nature:      NatureFire,
		Fields:      Field5,
	}, 5)...)
	deck = append(deck, Repeat(&CardData{
		Name:      "喜多的心动魔法",
		Desc:      "每次成功召唤怪兽，自己恢复500HP",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeForever,
		Fields:    FieldBand,
		Actions:   []string{"喜多的心动魔法-回血"},
	}, 3)...)
	deck = append(deck, Repeat(&CardData{
		Name:      "你就是歌姬吧",
		Desc:      "支付500HP，把对手的血量降到500，本回合结束再恢复到原来的体力值",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeCommon,
		Fields:    Field5,
		Actions:   []string{"你就是歌姬吧-降低血量"},
	}, 3)...)
	deck = append(deck, &CardData{
		Name:      "九转大肠",
		Desc:      "1，每次回合结束减少一点标记，标记为0时你立即获得胜利，2，每次回合结束移动到对面到区域，若没有空位随机炸一张卡",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeForever,
		Fields:    Field5,
		Actions:   []string{"九转大肠-减少标记"},
	})
	deck = append(deck, Repeat(&CardData{
		Name:      "十七张牌你能秒我",
		Desc:      "1，支付4000HP，将手牌补充到17张，2，本回合弃牌阶段每弃置一张牌恢复200HP，且弃置的牌置入卡组",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeCommon,
		Fields:    Field5,
		Actions:   []string{"十七张牌你能秒我-发动效果"},
	}, 3)...)
	deck = append(deck, Repeat(&CardData{
		Name:     "你这是违法行为",
		Desc:     "当敌人一回合召唤3只8星以上怪兽时发动，破坏其所以怪兽",
		CardType: CardTypeTrap,
		TrapType: TrapTypeCommon,
		Actions:  []string{"你这是违法行为-全部破坏", "覆盖"}, // 覆盖行为 是陷阱延时必须的
	}, 3)...)
	deck = append(deck, Repeat(&CardData{
		Name: "蔡徐坤",
		Desc: "1，回合开始时可以修改一名怪兽的种族为黑子，2，可以放弃攻击改为进行篮球邀约(各摸一张牌进行拼点(怪兽>魔法,魔法>陷阱,陷阱>怪兽" +
			"))若你没输对方跳过下个攻击阶段，3，黑子怪兽攻击你时攻击力减半",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         1500,
		Def:         2000,
		Level:       4,
		Race:        RaceUndead,
		Nature:      NatureFire,
		Fields:      Field5,
		Actions:     []string{"蔡徐坤-变小黑子", "蔡徐坤-篮球邀约", "蔡徐坤-打小黑子"},
	}, 3)...)
	utils.RandomArr(deck)
	for i := 0; i < len(deck); i++ {
		if deck[i].CardType == CardTypeMonster {
			deck[i].Actions = append(deck[i].Actions, "怪兽准备")
			deck[i].CanSelects = append(deck[i].CanSelects, "攻击源")
		}
	}
	extra := make([]*CardData, 0) // 额外暂时仅用来展示没有效果
	extra = append(extra, Repeat(&CardData{
		Name:        "孙笑川",
		Desc:        "我和我的猫，还有你妈，都很想念你，不过，那是骗你的啦,其实我没有猫，你也没有妈",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeCommon,
		Atk:         3100,
		Def:         2000,
		Level:       8,
		Race:        RaceUndead,
		Nature:      NatureWind,
		Fields:      Field5,
	}, 3)...)
	extra = append(extra, Repeat(&CardData{
		Name:        "Pop Cat",
		Desc:        "Test",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         3000,
		Def:         2000,
		Level:       8,
		Race:        RaceDragon,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | Field5,
	}, 2)...)
	utils.RandomArr(extra)
	for i := 0; i < len(extra); i++ {
		if extra[i].CardType == CardTypeMonster {
			extra[i].Actions = append(extra[i].Actions, "怪兽准备")
			extra[i].CanSelects = append(extra[i].CanSelects, "攻击源")
		}
	}
	return deck, extra
}

// GetSetoKaibaCards 海马卡组
func GetSetoKaibaCards() ([]*CardData, []*CardData) {
	deck := make([]*CardData, 0)
	deck = append(deck, Repeat(&CardData{
		Name:        "青眼白龙",
		Desc:        "大白板",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeCommon,
		Atk:         3000,
		Def:         2500,
		Level:       8,
		Race:        RaceDragon,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | FieldBlueEyeWhiteDragon,
	}, 3)...)
	deck = append(deck, &CardData{
		Name:        "青眼的贤士",
		Desc:        "1，召唤成功从卡组检索白石，2，从手牌丢弃，解放我方一只怪兽，特殊召唤一只青眼怪兽",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         0,
		Def:         1500,
		Level:       1,
		Race:        RaceSoldier,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | Field5,
		Actions:     []string{"青眼的贤士-检索白石", "青眼的贤士-召唤青眼怪兽"},
	})
	deck = append(deck, Repeat(&CardData{
		Name:        "太古的白石",
		Desc:        "1,进入墓地后的回合结束阶段特招青眼怪兽，2，将墓地的白石除外，将一只墓地的青眼怪兽加入手牌",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         600,
		Def:         500,
		Level:       1,
		Race:        RaceDragon,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | Field5,
		Actions:     []string{"太古的白石-召唤青眼怪兽", "太古的白石-拣回青眼怪兽"},
	}, 3)...)
	deck = append(deck, &CardData{
		Name:        "海龟坏兽",
		Desc:        "1，解放对手一只怪兽，把该卡攻击表示特殊召唤给对方",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         2200,
		Def:         3000,
		Level:       8,
		Race:        RacePlant,
		Nature:      NatureWater,
		Fields:      Field5,
		Actions:     []string{"海龟坏兽-送给对面"},
	})
	deck = append(deck, Repeat(&CardData{
		Name:        "青眼亚白龙",
		Desc:        "1，手牌中有青眼白龙可以展示青眼白龙来特殊召唤，2，可以当做青眼白龙使用，且看作通常怪兽，3，破坏对方一只怪兽，当前回合不能攻击",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeCommon,
		Atk:         3000,
		Def:         2500,
		Level:       8,
		Race:        RaceDragon,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | FieldBlueEyeWhiteDragon, // 当青眼白龙 使用 检索的是 这个字段
		Actions:     []string{"青眼亚白龙-展示召唤", "青眼亚白龙-破坏怪兽"},
	}, 3)...)
	deck = append(deck, Repeat(&CardData{
		Name:        "白灵龙",
		Desc:        "1，可以当做青眼白龙使用，2，召唤成功除外对手一张魔法陷阱卡，3，当对方场上有怪兽时，解放该怪兽从手牌特招青眼白龙",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         2500,
		Def:         2000,
		Level:       8,
		Race:        RaceDragon,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | FieldBlueEyeWhiteDragon,
		Actions:     []string{"白灵龙-除外陷阱", "白灵龙-特招青眼白龙"},
	}, 3)...)
	deck = append(deck, Repeat(&CardData{
		Name:        "增值的G",
		Desc:        "从手牌舍弃发动，对方每次召唤怪兽，我方抽一张牌", // 对方不会特招
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         500,
		Def:         200,
		Level:       2,
		Race:        RacePlant,
		Nature:      NatureLand,
		Fields:      FieldBlueEye | Field5,
		Actions:     []string{"增值的G-发动效果"},
	}, 2)...)
	deck = append(deck, &CardData{ // TODO 暂时  没有 需求 当白板使用
		Name:        "灰流丽",
		Desc:        "从手牌舍弃，阻止对手，检索，特招，把手牌丢入墓地",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         0,
		Def:         1800,
		Level:       3,
		Race:        RaceUndead,
		Nature:      NatureWind,
		Fields:      FieldBlueEye | Field5,
	})
	deck = append(deck, &CardData{ // 小修改 ，对面 没有 卡牌效果  效果都混到来一起  攻击也是效果 也会受牵连
		Name:     "无限泡影",
		Desc:     "我方场上没有卡时可以直接发动，移除对方一只怪兽，若是覆盖发动额外移除对方魔法陷阱卡",
		CardType: CardTypeTrap,
		TrapType: TrapTypeCommon,
		Actions:  []string{"无限泡影-直接发动", "无限泡影-覆盖发动"},
	})
	deck = append(deck, &CardData{
		Name:      "羽毛扫",
		Desc:      "破坏对手所有魔法陷阱卡",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeCommon,
		Actions:   []string{"羽毛扫-破坏魔法陷阱"},
	})
	deck = append(deck, &CardData{
		Name:      "交易进行",
		Desc:      "把手牌等级8的怪兽送入墓地抽两张牌",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeCommon,
		Actions:   []string{"交易进行-抽牌"},
	})
	deck = append(deck, Repeat(&CardData{
		Name:      "龙旋律",
		Desc:      "舍弃一张手牌，从卡组检索两张攻击力3000以上，守备力2500以下的龙族怪兽",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeCommon,
		Actions:   []string{"龙旋律-检索怪兽"},
	}, 2)...)
	deck = append(deck, Repeat(&CardData{
		Name:      "调和的宝牌",
		Desc:      "舍弃一张白石，从牌组抽两张牌",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeCommon,
		Actions:   []string{"调和的宝牌-抽牌"},
	}, 2)...)
	deck = append(deck, Repeat(&CardData{
		Name:      "龙之灵庙",
		Desc:      "1，可把卡组的一只龙族怪兽送入墓地，如果是通常怪兽，可以再送入墓地一张龙族怪兽，2，一回合只能发动一次",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeCommon,
		Actions:   []string{"龙之灵庙-堆墓"},
	}, 2)...)
	deck = append(deck, Repeat(&CardData{
		Name:      "复活的福音",
		Desc:      "1，复活我方墓地一只等级8的龙族怪兽，2，在墓地时可以代替一次战斗破坏或效果破坏",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeCommon,
		Actions:   []string{"复活的福音-复活怪兽", "复活的福音-抵挡破坏"},
	}, 3)...)
	deck = append(deck, &CardData{ // TODO 没有 需求 当 放在场上的 白板使用
		Name:      "技能抽取",
		Desc:      "支付1000生命值，场上所有怪兽效果无效化",
		CardType:  CardTypeMagic,
		MagicType: MagicTypeForever,
		Actions:   []string{"技能抽取-放置"},
	})
	for i := 0; i < len(deck); i++ {
		switch deck[i].CardType {
		case CardTypeMonster:
			deck[i].Actions = append(deck[i].Actions, "通常召唤", "切换状态", "怪兽准备")
			deck[i].CanSelects = append(deck[i].CanSelects, "攻击源")
		case CardTypeTrap:
			deck[i].Actions = append(deck[i].Actions, "覆盖")
		}
	}
	extra := make([]*CardData, 0)
	extra = append(extra, &CardData{
		Name:        "青眼混沌MAX龙",
		Desc:        "1，移除场上一张青眼白龙特殊召唤，2，不受效果影响，3，攻击守备怪兽两倍穿防",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         4000,
		Def:         0,
		Level:       8,
		Race:        RaceDragon,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | Field5,
		Actions:     []string{"青眼混沌MAX龙-召唤", "青眼混沌MAX龙-2倍穿防"},
	})
	extra = append(extra, &CardData{
		Name:        "青眼双爆裂龙",
		Desc:        "1，场上2只青眼白龙特殊召唤，2，不会战斗破坏，3，可以攻击两次",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         3000,
		Def:         2500,
		Level:       10,
		Race:        RaceDragon,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | Field5,
		Actions:     []string{"青眼双爆裂龙-召唤", "青眼双爆裂龙-免战破", "青眼双爆裂龙-攻击2次"},
	})
	extra = append(extra, &CardData{
		Name:        "苍眼银龙",
		Desc:        "1，等级之和为9的至少两只怪兽(需要至少一张青眼怪兽)特殊召唤，2，特超成功，使场上龙族怪兽获取效果抗性一回合，3，我方准备阶段可以特招墓地一只通常怪兽",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         2500,
		Def:         3000,
		Level:       9,
		Race:        RaceDragon,
		Nature:      NatureLight,
		Fields:      FieldBlueEye | Field5,
		Actions:     []string{"苍眼银龙-召唤", "苍眼银龙-召唤通常怪兽"},
	})
	extra = append(extra, &CardData{
		Name:        "联结栗子球",
		Desc:        "1，需要解放两只等级为1的怪兽特殊召唤，2，对方怪兽攻击时可以解放该牌使其攻击力变为0",
		CardType:    CardTypeMonster,
		MonsterType: MonsterTypeEffect,
		Atk:         300,
		Def:         0,
		Level:       1,
		Race:        RaceMachine,
		Nature:      NatureWind,
		Fields:      FieldBlueEye | Field5,
		Actions:     []string{"联结栗子球-召唤", "联结栗子球-减攻击力"},
	})
	for i := 0; i < len(extra); i++ {
		extra[i].Actions = append(extra[i].Actions, "切换状态", "怪兽准备")
		extra[i].CanSelects = append(extra[i].CanSelects, "攻击源")
	}
	utils.RandomArr(deck)
	return deck, extra
}
