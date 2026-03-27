package flip

import "image/color"

// screen = border + body
// body = score + button + card
// offset是基于body的，算总值的时候要加上border

const (
	border = 20

	scoreWidth  = CardBoardWidth - 2*scoreMargin
	scoreHeight = 60
	scoreMargin = 10

	buttonWidth  = CardBoardWidth - 2*buttonMargin
	buttonHeight = 40
	buttonMargin = 10

	cardWidthScale  = 3 // 4:3
	cardHeightScale = 4
	cardMarginScale = 1
	cardScaleBase   = 20

	cardWidth  = cardWidthScale * cardScaleBase
	cardHeight = cardHeightScale * cardScaleBase
	cardMargin = cardMarginScale * cardScaleBase
)

const (
	ScoreBoardOffsetWidth  = scoreMargin
	ScoreBoardOffsetHeight = scoreMargin
	ScoreBoardWidth        = scoreWidth
	ScoreBoardHeight       = scoreHeight

	ButtonBoardOffsetWidth  = buttonMargin
	ButtonBoardOffsetHeight = buttonMargin + ScoreBoardHeight + 2*scoreMargin
	ButtonBoardWidth        = buttonWidth
	ButtonBoardHeight       = buttonHeight

	CardBoardOffsetWidth  = 0
	CardBoardOffsetHeight = ScoreBoardHeight + 2*scoreMargin + ButtonBoardHeight + 2*buttonMargin
	CardBoardWidth        = 4*cardWidth + 5*cardMargin
	CardBoardHeight       = 4*cardHeight + 5*cardMargin

	BodyWidth  = CardBoardWidth
	BodyHeight = ScoreBoardHeight + ButtonBoardOffsetHeight + CardBoardHeight

	ScreenWidth  = BodyWidth + 2*border
	ScreenHeight = BodyHeight + 2*border
)

var (
	backgroundColorLight = color.RGBA{R: 0xf0, G: 0xef, B: 0xe2, A: 0xff}
	backgroundColorDark  = color.RGBA{R: 0xd2, G: 0xd1, B: 0xb1, A: 0xff}
)
