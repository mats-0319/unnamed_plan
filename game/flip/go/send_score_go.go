//go:build !js || !wasm

package flip

import "log"

type ScoreToWeb struct {
	GameName int8   `json:"game_name"` // flip=1
	Score    int    `json:"score"`
	Result   string `json:"result"` // json str
}

// Flip game 分数计算规则：每一步计2秒，用时（毫秒，包含步数折算时间）作为cost，使用100_0000 - cost即为最终得分，负分记为0

func sendFlipScoreToWeb(duration int64, steps int) {
	log.Printf("Duration: %d, Steps: %d\n", duration, steps)
}
