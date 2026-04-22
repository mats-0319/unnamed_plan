//go:build !js || !wasm

// 本地开发期间使用该文件，包含功能：上传成绩

package flip

import "log"

func sendFlipScoreToWeb(duration int64, steps int) {
	log.Printf("Duration: %d, Steps: %d\n", duration, steps)
}
