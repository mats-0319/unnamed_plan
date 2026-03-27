package flip

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mats0319/unnamed_plan/game/flip/assets"
)

type CardBoard struct {
	cards     [16]*card
	images    [8]*ebiten.Image
	backImage *ebiten.Image

	lastClick []int // last click card index
}

type card struct {
	number      int
	isFrontSide bool

	isFlipping bool
	flipAngle  float64
}

func NewCardBoard() (*CardBoard, error) {
	res := &CardBoard{}

	// init cards
	cardInitFlag := [16]bool{}
	for i := 0; i < 8; i++ {
		for range 2 {
			cardIndex := rand.IntN(16)
			for cardInitFlag[cardIndex] { // 如果随机索引已被占用，则向后找到下一个可用索引
				cardIndex = (cardIndex + 1) % 16
			}

			if res.cards[cardIndex] == nil {
				res.cards[cardIndex] = &card{}
			}
			res.cards[cardIndex].number = i
			cardInitFlag[cardIndex] = true
		}
	}

	// init card images
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.AssetsFS, "back.png")
	if err != nil {
		return nil, err
	}
	res.backImage = img

	for i := range res.images {
		img, _, err := ebitenutil.NewImageFromFileSystem(assets.AssetsFS, fmt.Sprintf("front_%d.png", i))
		if err != nil {
			return nil, err
		}
		res.images[i] = img
	}

	return res, nil
}

func (c *CardBoard) Update(input *Input) (frontCount int, addStep bool) {
	defer func() {
		for _, v := range c.cards {
			if v.isFrontSide {
				frontCount++
			}
		}
	}()

	// 翻转动画参数计算
	for _, v := range c.cards {
		if v.isFlipping {
			v.flipAngle += math.Pi / 12

			if v.flipAngle >= math.Pi {
				v.isFlipping = false
				v.flipAngle = 0
			}
		}
	}

	// 获取点击的卡牌索引（或者没有点击到卡牌上）
	index := input.clickCardIndex
	if input.clickOn != ClickOn_Card || c.cards[index].isFrontSide { // 未选中卡牌或选中已经是正面向上的卡牌
		return
	}

	{
		if len(c.lastClick) > 1 { // 已翻转过两张卡牌，已经点击了第三张但尚未执行翻转
			// 它们的序号不同，则在下一次翻转卡牌前（即此刻）将它们重置为反面向上（无动画）
			if c.cards[c.lastClick[0]].number != c.cards[c.lastClick[1]].number {
				c.cards[c.lastClick[0]].isFrontSide = false
				c.cards[c.lastClick[1]].isFrontSide = false
			}

			c.lastClick = make([]int, 0, 2)
		}

		c.lastClick = append(c.lastClick, index) // 记录翻转卡牌索引

		if len(c.lastClick) > 1 { // 此时刚翻转两张卡牌
			addStep = true // 步数统计+1
		}
	}

	c.cards[index].isFrontSide = true
	c.cards[index].isFlipping = true

	return
}

func (c *CardBoard) Draw(cardBoard *ebiten.Image, state GameState) {
	cardBoard.Fill(backgroundColorLight)

	if state == GameState_Prepared {
		return
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			options := &ebiten.DrawImageOptions{}
			x := j*cardWidth + (j+1)*cardMargin
			y := i*cardHeight + (i+1)*cardMargin

			cardItem := c.cards[i*4+j]
			if cardItem.isFlipping {
				scaleX := math.Cos(cardItem.flipAngle)
				options.GeoM.Translate(-cardWidth/2, -cardHeight/2) // 将卡牌中心点移动到画布原点(0,0)
				options.GeoM.Scale(scaleX, 1)
				if scaleX < 0 {
					options.GeoM.Scale(-1, 1)
				}
				options.GeoM.Translate(float64(x+cardWidth/2), float64(y+cardHeight/2)) // 将卡牌移动到目标位置
				if scaleX > 0 {
					cardBoard.DrawImage(c.backImage, options)
				} else {
					cardBoard.DrawImage(c.images[cardItem.number], options)
				}
			} else {
				options.GeoM.Translate(float64(x), float64(y))

				if cardItem.isFrontSide {
					cardBoard.DrawImage(c.images[cardItem.number], options)
				} else {
					cardBoard.DrawImage(c.backImage, options)
				}
			}
		}
	}
}
