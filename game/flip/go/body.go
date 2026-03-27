package flip

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mats0319/unnamed_plan/game/flip/assets/font"
)

type Body struct {
	score  *ScoreBoard
	button *ButtonBoard
	card   *CardBoard

	textFace *text.GoTextFace

	scoreImage  *ebiten.Image
	buttonImage *ebiten.Image
	cardImage   *ebiten.Image
}

func NewBody() (*Body, error) {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(font.FontRoboto))
	if err != nil {
		return nil, err
	}

	cardBoard, err := NewCardBoard()
	if err != nil {
		return nil, err
	}

	return &Body{
		score:       NewScoreBoard(),
		button:      NewButtonBoard(),
		card:        cardBoard,
		textFace:    &text.GoTextFace{Source: source, Size: 20},
		scoreImage:  ebiten.NewImage(ScoreBoardWidth, ScoreBoardHeight),
		buttonImage: ebiten.NewImage(ButtonBoardWidth, ButtonBoardHeight),
		cardImage:   ebiten.NewImage(CardBoardWidth, CardBoardHeight),
	}, nil
}

func (b *Body) Update(state GameState, input *Input) int {
	b.button.Update(state)
	frontCount, addStep := b.card.Update(input)
	b.score.Update(state, addStep)

	return frontCount
}

func (b *Body) Draw(bodyBoard *ebiten.Image, state GameState) {
	b.score.Draw(b.scoreImage, b.textFace)
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(ScoreBoardOffsetWidth, ScoreBoardOffsetHeight)
	bodyBoard.DrawImage(b.scoreImage, options)

	b.button.Draw(b.buttonImage, b.textFace)
	options = &ebiten.DrawImageOptions{}
	options.GeoM.Translate(ButtonBoardOffsetWidth, ButtonBoardOffsetHeight)
	bodyBoard.DrawImage(b.buttonImage, options)

	b.card.Draw(b.cardImage, state)
	options = &ebiten.DrawImageOptions{}
	options.GeoM.Translate(CardBoardOffsetWidth, CardBoardOffsetHeight)
	bodyBoard.DrawImage(b.cardImage, options)
}

func drawCenterText(image *ebiten.Image, content string, textFace *text.GoTextFace, cx float64, cy float64, color color.Color) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(cx, cy)
	options.ColorScale.ScaleWithColor(color)
	options.PrimaryAlign = text.AlignCenter
	options.SecondaryAlign = text.AlignCenter
	text.Draw(image, content, textFace, options)
}
