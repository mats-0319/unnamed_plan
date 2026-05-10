package model

type GameScore struct {
	Score      int    // 分数，每个游戏可能有多个评价结果的维度，将其加权计算为一个分数
	Result     string // 详细结果，包含每个参与计算得分的维度
	Player     string // user name，游客成绩该字段为空
	PlayerName string // user nickname(at that time)，游客成绩显示'游客+[随机字符]'
}

type FlipGameScore struct {
	Model
	GameScore
}
