package flip

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ButtonBoard struct {
	ButtonText string
}

func NewButtonBoard() *ButtonBoard {
	return &ButtonBoard{}
}

func (b *ButtonBoard) Update(state GameState) {
	switch state {
	case GameState_Prepared:
		b.ButtonText = "Start Game"
	case GameState_End:
		b.ButtonText = "Restart Game"
	default:
		b.ButtonText = ""
	}
}

func (b *ButtonBoard) Draw(buttonBoard *ebiten.Image, textFace *text.GoTextFace) {
	buttonBoard.Fill(backgroundColorLight)

	if len(b.ButtonText) > 0 {
		drawCenterText(buttonBoard, b.ButtonText, textFace, ButtonBoardWidth/2.0, ButtonBoardHeight/2.0, color.Black)
	}
}
