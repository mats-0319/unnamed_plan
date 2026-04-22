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
	Cards [16]*CardItem

	BackImage   *ebiten.Image
	FrontImages [8]*ebiten.Image

	Selected []int // selected card index, waiting for further handle
}

type CardItem struct {
	Number      int
	IsFrontSide bool

	IsFlipping bool
	FlipAngle  float64
}

func NewCardBoard() (*CardBoard, error) {
	res := &CardBoard{}
	res.reset() // init cards

	// init images
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.AssetsFS, "back.png")
	if err != nil {
		return nil, err
	}
	res.BackImage = img

	for i := range res.FrontImages {
		img, _, err := ebitenutil.NewImageFromFileSystem(assets.AssetsFS, fmt.Sprintf("front_%d.png", i))
		if err != nil {
			return nil, err
		}
		res.FrontImages[i] = img
	}

	return res, nil
}

func (c *CardBoard) reset() {
	cardInitFlag := [16]bool{}
	for i := 0; i < 8; i++ { // i: [0,7]
		for range 2 { // 每个i用两次
			cardIndex := rand.IntN(16)
			for cardInitFlag[cardIndex] { // 如果随机索引已被占用，则向后找到下一个可用索引
				cardIndex = (cardIndex + 1) % 16
			}

			if c.Cards[cardIndex] == nil {
				c.Cards[cardIndex] = &CardItem{}
			}
			c.Cards[cardIndex].Number = i
			cardInitFlag[cardIndex] = true
		}
	}
}

func (c *CardBoard) Update(input *Input) (frontCount int, stepOffset int) {
	defer func() {
		for _, v := range c.Cards {
			if v.IsFrontSide {
				frontCount++
			}
		}
	}()

	// 翻转动画参数计算
	for _, v := range c.Cards {
		if v.IsFlipping {
			v.FlipAngle += math.Pi / 12 // 12帧完成翻转

			if v.FlipAngle >= math.Pi {
				v.IsFlipping = false
				v.FlipAngle = 0
			}
		}
	}

	// 获取点击的卡牌索引，没有点击卡牌或点击已经是正面向上的卡牌时不进行后续处理
	index := input.clickCardIndex
	if input.clickOn != ClickOn_Card || c.Cards[index].IsFrontSide { // 未选中卡牌或选中已经是正面向上的卡牌
		return
	}

	{
		if len(c.Selected) > 1 { // 已翻转过两张卡牌，已经点击了第三张但尚未执行翻转
			// 它们的序号不同，则在下一次翻转卡牌前（即此刻）将它们重置为反面向上（无动画）；
			// 序号相同则保持正面向上（无动作）。
			if c.Cards[c.Selected[0]].Number != c.Cards[c.Selected[1]].Number {
				c.Cards[c.Selected[0]].IsFrontSide = false
				c.Cards[c.Selected[1]].IsFrontSide = false
			}

			c.Selected = make([]int, 0, 2) // 移除处理过的卡牌
		}

		c.Selected = append(c.Selected, index) // 执行翻转

		if len(c.Selected) > 1 { // 此时刚翻转两张卡牌
			stepOffset = 1 // 步数统计+1
		}
	}

	c.Cards[index].IsFrontSide = true
	c.Cards[index].IsFlipping = true

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

			cardItem := c.Cards[i*4+j]
			if cardItem.IsFlipping { // 如果卡牌正在翻转
				scaleX := math.Cos(cardItem.FlipAngle) // scaleX: 1 -> 0 -> -1

				options.GeoM.Translate(-cardWidth/2, -cardHeight/2) // 将卡牌中心点移动到画布原点(0,0)
				options.GeoM.Scale(scaleX, 1)
				if scaleX < 0 { // 翻转到后半段
					// 不能和下面的if合并，因为缩放是基于原点的，如果简单合并会导致展示效果未定义
					options.GeoM.Scale(-1, 1) // 图片水平镜像
				}
				options.GeoM.Translate(float64(x+cardWidth/2), float64(y+cardHeight/2)) // 将卡牌移动到目标位置
				if scaleX > 0 {
					cardBoard.DrawImage(c.BackImage, options)
				} else {
					cardBoard.DrawImage(c.FrontImages[cardItem.Number], options)
				}
			} else {
				options.GeoM.Translate(float64(x), float64(y))

				if cardItem.IsFrontSide {
					cardBoard.DrawImage(c.FrontImages[cardItem.Number], options)
				} else {
					cardBoard.DrawImage(c.BackImage, options)
				}
			}
		}
	}
}
