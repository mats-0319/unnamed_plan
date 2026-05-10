package flip

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// GameState 状态转换：
// 1. 初始化游戏，进入 就绪 状态
// 2. 点击按钮进入 进行中 状态
// 3. 完成一局游戏，自动进入 结束 状态
// 4. 点击按钮重新初始化游戏，进入 就绪 状态，完成一个循环
type GameState int8

const (
	GameState_Prepared GameState = 1 // 已就绪
	GameState_Playing            = 2 // 进行中
	GameState_End                = 3 // 已结束
)

type Game struct {
	State GameState

	Input     *Input
	Body      *Body
	BodyBoard *ebiten.Image
}

func NewGame() (*Game, error) {
	body, err := NewBody()
	if err != nil {
		return nil, err
	}

	return &Game{
		State:     GameState_Prepared,
		Input:     &Input{},
		Body:      body,
		BodyBoard: ebiten.NewImage(BodyWidth, BodyHeight),
	}, nil
}

func (g *Game) Update() error {
	g.Input.Update(g.State)

	if g.State == GameState_Prepared && g.Input.clickOn == ClickOn_Button {
		g.State = GameState_Playing
	} else if g.State == GameState_End && g.Input.clickOn == ClickOn_Button {
		g.Body.reset()
		g.State = GameState_Prepared
	}

	gameTime := g.Body.Update(g.State, g.Input)

	if g.State == GameState_Playing && gameTime > 0 {
		g.State = GameState_End

		sendFlipScoreToWeb(gameTime, g.Body.Score.steps)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColorDark)

	g.Body.Draw(g.BodyBoard, g.State)

	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(border, border)
	screen.DrawImage(g.BodyBoard, options)
}

func (g *Game) Layout(_ int, _ int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
