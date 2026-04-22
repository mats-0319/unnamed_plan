package flip

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ScoreBoard struct {
	durationStr string
	startTime   time.Time // 计算用
	endTime     time.Time

	steps int
}

func NewScoreBoard() *ScoreBoard {
	return &ScoreBoard{}
}

func (s *ScoreBoard) Update(state GameState, stepOffset int) (gameTime int64) {
	switch state {
	case GameState_Prepared:
		s.durationStr = "00:00.000"
	case GameState_Playing:
		if s.startTime.IsZero() {
			s.startTime = time.Now()
		}

		duration := time.Since(s.startTime)
		s.durationStr = fmt.Sprintf("%02d:%02d.%03d", int(duration/time.Minute)%60, int(duration/time.Second)%60, duration.Milliseconds()%1000)
	case GameState_End:
		if s.endTime.IsZero() { // 只计算一次
			s.endTime = time.Now()
			duration := s.endTime.Sub(s.startTime)
			s.durationStr = fmt.Sprintf("%02d:%02d.%03d", int(duration/time.Minute)%60, int(duration/time.Second)%60, duration.Milliseconds()%1000)

			gameTime = duration.Milliseconds()
		}
	}

	s.steps += stepOffset

	return
}

func (s *ScoreBoard) Draw(scoreBoard *ebiten.Image, textFace *text.GoTextFace) {
	scoreBoard.Fill(backgroundColorLight)

	line1CenterY := ScoreBoardHeight / 4.0
	line2CenterY := ScoreBoardHeight / 4.0 * 3

	durationWidth := ScoreBoardWidth * 6.0 / 10
	durationCenterX := durationWidth / 2

	drawCenterText(scoreBoard, "Duration", textFace, durationCenterX, line1CenterY, backgroundColorDark)
	drawCenterText(scoreBoard, s.durationStr, textFace, durationCenterX, line2CenterY, color.Black)

	stepsWidth := ScoreBoardWidth - durationWidth
	stepsCenterX := durationWidth + stepsWidth/2

	drawCenterText(scoreBoard, "Steps", textFace, stepsCenterX, line1CenterY, backgroundColorDark)
	drawCenterText(scoreBoard, strconv.Itoa(s.steps), textFace, stepsCenterX, line2CenterY, color.Black)
}
