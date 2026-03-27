package flip

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ButtonBoard struct {
	buttonText string

	buttonImage *ebiten.Image
}

func NewButtonBoard() *ButtonBoard {
	return &ButtonBoard{buttonImage: ebiten.NewImage(buttonWidth, buttonHeight)}
}

func (b *ButtonBoard) Update(state GameState) {
	switch state {
	case GameState_Prepared:
		b.buttonText = "Start Game"
	case GameState_End:
		b.buttonText = "Restart Game"
	default:
		b.buttonText = ""
	}
}

func (b *ButtonBoard) Draw(buttonBoard *ebiten.Image, textFace *text.GoTextFace) {
	buttonBoard.Fill(backgroundColorLight)

	if len(b.buttonText) > 0 {
		drawCenterText(buttonBoard, b.buttonText, textFace, ButtonBoardWidth/2.0, ButtonBoardHeight/2.0, color.Black)
	}
}
