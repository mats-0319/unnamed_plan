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
	Number    int
	Status    CardStatus
	FlipAngle float64 // 仅在 选中状态-执行翻转动画 期间生效
}

type CardStatus int8

const (
	CardStatus_Default  CardStatus = 0 // 默认状态，反面向上
	CardStatus_Flipping            = 1 // 选中状态，执行翻转动画
	CardStatus_Selected            = 2 // 选中状态，正面向上
	CardStatus_Matched             = 3 // 配对成功，正面向上
)

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
	c.Cards = [16]*CardItem{}

	cardInitFlag := [16]bool{}
	for i := 0; i < 8; i++ { // i: [0,7]
		for range 2 { // 每个i用两次
			cardIndex := rand.IntN(16)
			for cardInitFlag[cardIndex] { // 如果随机索引已被占用，则向后找到下一个可用索引
				cardIndex = (cardIndex + 1) % 16
			}

			if c.Cards[cardIndex] == nil { // 蹭数据循环来初始化数组的每个元素
				c.Cards[cardIndex] = &CardItem{}
			}
			c.Cards[cardIndex].Number = i
			cardInitFlag[cardIndex] = true
		}
	}
}

func (c *CardBoard) Update(input *Input) (matchedCount int, stepOffset int) {
	defer func() {
		for _, v := range c.Cards {
			if v.Status == CardStatus_Matched {
				matchedCount++
			}
		}
	}()

	// 翻转动画
	for _, card := range c.Cards {
		if card.Status == CardStatus_Flipping {
			card.FlipAngle += math.Pi / 12 // 0.2秒完成翻转

			if card.FlipAngle >= math.Pi {
				card.Status = CardStatus_Selected
				card.FlipAngle = 0
			}
		}
	}

	// 获取点击的卡牌索引，只有满足以下条件之一的，才会进入后续处理：
	// 1. 选中默认状态的卡牌
	// 2. 选中已选择状态的卡牌、且此前已选择两张卡牌、且卡牌已完成翻转动画
	// 因为我希望能在点击第三张卡牌时处理前两张（而不是翻转动画结束之后就处理），所以这里要写的复杂一点
	index := input.clickCardIndex
	if input.clickOn != ClickOn_Card ||
		!(c.Cards[index].Status == CardStatus_Default ||
			(c.Cards[index].Status == CardStatus_Selected &&
				len(c.Selected) == 2 && c.Cards[c.Selected[1]].Status == CardStatus_Selected)) {
		return
	}

	{
		if len(c.Selected) > 1 { // 已翻转过两张卡牌，正准备翻转第三张卡牌
			if c.Cards[c.Selected[0]].Number == c.Cards[c.Selected[1]].Number {
				c.Cards[c.Selected[0]].Status = CardStatus_Matched
				c.Cards[c.Selected[1]].Status = CardStatus_Matched

				if c.Cards[index].Status == CardStatus_Matched {
					c.Selected = make([]int, 0, 2)
					return
				}
			}

			c.Cards[c.Selected[0]].Status = CardStatus_Default
			c.Cards[c.Selected[1]].Status = CardStatus_Default

			c.Selected = make([]int, 0, 2) // 移除处理过的卡牌
		}

		c.Cards[index].Status = CardStatus_Flipping
		c.Selected = append(c.Selected, index) // 执行翻转

		if len(c.Selected) > 1 { // 此时刚翻转两张卡牌
			stepOffset = 1 // 步数统计+1
		}
	}

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
			x := j*cardWidth + (j+1)*cardMargin  // offset x
			y := i*cardHeight + (i+1)*cardMargin // offset y

			cardItem := c.Cards[i*4+j]
			switch cardItem.Status {
			case CardStatus_Flipping: // 如果卡牌正在翻转
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
			case CardStatus_Default:
				options.GeoM.Translate(float64(x), float64(y))
				cardBoard.DrawImage(c.BackImage, options)
			case CardStatus_Selected, CardStatus_Matched:
				options.GeoM.Translate(float64(x), float64(y))
				cardBoard.DrawImage(c.FrontImages[cardItem.Number], options)
			}
		}
	}
}
