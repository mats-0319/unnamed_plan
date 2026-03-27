//go:build js && wasm
// +build js,wasm

package flip

import (
	"encoding/json"
	"log"
	"syscall/js"
)

// Flip game 分数计算规则：每一步计2秒，用时（毫秒，包含步数折算时间）作为cost，使用100_0000 - cost即为最终得分，负分记为0

func sendFlipScoreToWeb(duration int64, steps int) {
	result := &struct { // flip game result
		Duration int64 `json:"duration"`
		Steps    int   `json:"steps"`
	}{
		Duration: duration,
		Steps:    steps,
	}

	// marshal game result into str
	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Println("Send Flip Score Failed", err)
		return
	}

	// calc score
	score := 100_0000 - int(duration) - steps*2_000
	if score < 0 {
		score = 0
	}

	// 构造数据结构：游戏成绩
	jsData := js.ValueOf(map[string]any{
		"game_name": 1, //flip=1
		"score":     score,
		"result":    string(resultBytes),
	})

	// postMessage
	js.Global().Get("parent").Call("postMessage", jsData, "*")
}
