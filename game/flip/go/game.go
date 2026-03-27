package flip

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GameState int8

const (
	GameState_Prepared GameState = 1 // 已就绪
	GameState_Playing            = 2 // 进行中
	GameState_End                = 3 // 已结束
)

type Game struct {
	state      GameState
	frontCount int

	input     *Input
	body      *Body
	bodyBoard *ebiten.Image
}

func NewGame() (*Game, error) {
	body, err := NewBody()
	if err != nil {
		return nil, err
	}

	return &Game{
		state:     GameState_Prepared,
		input:     &Input{},
		body:      body,
		bodyBoard: ebiten.NewImage(BodyWidth, BodyHeight),
	}, nil
}

func (g *Game) Update() error {
	g.input.Update(g.state)
	if g.state == GameState_Prepared && g.input.clickOn == ClickOn_Button {
		g.state = GameState_Playing
	} else if g.state == GameState_End && g.input.clickOn == ClickOn_Button {
		newBody, err := NewBody()
		if err != nil {
			return err
		}

		g.body = newBody // reset game
		g.state = GameState_Prepared
	}

	g.frontCount = g.body.Update(g.state, g.input)
	if g.state == GameState_Playing && g.frontCount >= 16 {
		g.state = GameState_End

		sendFlipScoreToWeb(g.body.score.duration.Milliseconds(), g.body.score.steps)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColorDark)

	g.body.Draw(g.bodyBoard, g.state)

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(border, border)
	screen.DrawImage(g.bodyBoard, options)
}

func (g *Game) Layout(_ int, _ int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
