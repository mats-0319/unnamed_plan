package api

type GameName int8

const (
	GameName_Flip GameName = 1
)

const URI_UploadGameScore = "/game-score/upload"

type UploadGameScoreReq struct {
	GameName GameName `json:"game_name"`
	Score    int      `json:"score"`
	Result   string   `json:"result"` // json str
	Player   string   `json:"player"` // 已登录则为空，否则：'游客+[随机字符]'
}

type UploadGameScoreRes struct{}

const URI_ListGameScore = "/game-score/list"

type GameScore struct {
	Score      int    `json:"score"`
	Result     string `json:"result"`
	PlayerName string `json:"player_name"` // user nickname
	Timestamp  int64  `json:"timestamp"`
}

type ListGameScoreReq struct {
	GameName GameName   `json:"game_name"`
	Page     Pagination `json:"page"`
}

type ListGameScoreRes struct {
	Count  int64        `json:"count"`
	Scores []*GameScore `json:"scores"`
}
