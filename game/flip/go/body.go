package flip

import (
	"bytes"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/mats0319/unnamed_plan/game/flip/assets/font"
)

type Body struct {
	Score  *ScoreBoard
	Button *ButtonBoard
	Card   *CardBoard

	TextFace *text.GoTextFace

	ScoreImage  *ebiten.Image
	ButtonImage *ebiten.Image
	CardImage   *ebiten.Image
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
		Score:       NewScoreBoard(),
		Button:      NewButtonBoard(),
		Card:        cardBoard,
		TextFace:    &text.GoTextFace{Source: source, Size: 20},
		ScoreImage:  ebiten.NewImage(ScoreBoardWidth, ScoreBoardHeight),
		ButtonImage: ebiten.NewImage(ButtonBoardWidth, ButtonBoardHeight),
		CardImage:   ebiten.NewImage(CardBoardWidth, CardBoardHeight),
	}, nil
}

func (b *Body) reset() {
	b.Score = NewScoreBoard()
	b.Button = NewButtonBoard()
	b.Card.reset()
}

func (b *Body) Update(state GameState, input *Input) int64 {
	b.Button.Update(state)
	frontCount, stepOffset := b.Card.Update(input)

	if state == GameState_Playing && frontCount >= 16 {
		state = GameState_End
	}

	return b.Score.Update(state, stepOffset)
}

func (b *Body) Draw(bodyBoard *ebiten.Image, state GameState) {
	b.Score.Draw(b.ScoreImage, b.TextFace)
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(ScoreBoardOffsetWidth, ScoreBoardOffsetHeight)
	bodyBoard.DrawImage(b.ScoreImage, options)

	b.Button.Draw(b.ButtonImage, b.TextFace)
	options = &ebiten.DrawImageOptions{}
	options.GeoM.Translate(ButtonBoardOffsetWidth, ButtonBoardOffsetHeight)
	bodyBoard.DrawImage(b.ButtonImage, options)

	b.Card.Draw(b.CardImage, state)
	options = &ebiten.DrawImageOptions{}
	options.GeoM.Translate(CardBoardOffsetWidth, CardBoardOffsetHeight)
	bodyBoard.DrawImage(b.CardImage, options)
}

func drawCenterText(image *ebiten.Image, content string, textFace *text.GoTextFace, cx float64, cy float64, color color.Color) {
	options := &text.DrawOptions{}
	options.GeoM.Translate(cx, cy)
	options.ColorScale.ScaleWithColor(color)
	options.PrimaryAlign = text.AlignCenter
	options.SecondaryAlign = text.AlignCenter
	text.Draw(image, content, textFace, options)
}
