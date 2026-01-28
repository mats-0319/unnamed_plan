package main

import (
	"github.com/mats0319/unnamed_plan/server/cmd/model"
	"github.com/mats0319/unnamed_plan/server/internal/utils"
)

var defaultUsers = []*model.User{
	newUser("mats0319", "Mario", "", true),
	newUser("admin", "", "5SSFNNEJUENPCCKP", true),
	newUser("user", "", "", false),
}

func newUser(userName string, nickname string, totpKey string, isAdmin bool) *model.User {
	salt := utils.GenerateRandomBytes[string](10)
	pwd := "123456"
	pwd = utils.HmacSHA256[string](pwd)
	pwd = utils.HmacSHA256(pwd, salt)

	if len(nickname) < 1 {
		nickname = userName
	}

	return &model.User{
		UserName: userName,
		Nickname: nickname,
		Password: pwd,
		Salt:     salt,
		TotpKey:  totpKey,
		IsAdmin:  isAdmin,
	}
}

var testNotes = []*model.Note{
	model.NewNote(1001, "Mario", false, noteTitle_1, noteContent_1),
	model.NewNote(1001, "Mario", false, noteTitle_2, noteContent_2),
}

const (
	noteTitle_1 = "对上暗号的微妙感"
	noteTitle_2 = "当你认为身边的人都是傻子的时候，你应当想一想，是不是在周边人眼中你才是那个傻子，而他们形成了某种默契在哄着你玩。"

	noteContent_1 = "一次坐公交的时候，听到喇叭播放：“尊老爱幼是中华民族的传统美德......”，最开始没当回事，" +
		"还以为是司机的任务指标，一趟车要放几遍之类的；可没一会儿，一个拎着编织袋子的老人站到我的座位前面，" +
		"我一下子想通了：司机播放广播是在他看到有需要帮助的乘客上车的时候，放广播也不是（或者说不全是）因为指标，" +
		"也是想帮助有需要的乘客。\n" +
		"想到这我就站起来让座了，老人家拦了一下，我只说我快要下车了。当时我的脑子里有一种特殊的兴奋感，" +
		"那是一种我和司机师傅对上了暗号的微妙感受。经此一事，我好像明白了文艺青年之间那种“我念一首冷门诗的上句，" +
		"旁边就有人能接上下句，并且还能讨论两句”微妙的感觉了。"
	noteContent_2 = "小学3、4年级的时候，网络游戏天龙八部正火，而我家里当时连电脑都没有，我又想显得合群，就谎称自己也在玩，" +
		"希望能加入大家的讨论，班级里的带头大哥也接受了。然而因为问到我游戏相关的问题时，我只能含糊其词，几次之后，" +
		"他们虽然没有拒绝我的参与，但是我还是能感受到，我和他们之间已经隔了一层可悲的厚障壁了。\n" +
		"随着家里买了电脑、扯了网线，我真玩上游戏了，又去找带头大哥，和大家一起玩了几天，他们开始逐渐接纳我。\n" +
		"到了小学6年级前后吧，发生了一件事。某一天，有一个同学甲来找我们聊游戏，但是言语之间透露出他似乎对游戏并不熟悉，" +
		"带头大哥让我们一人回去想一个游戏相关的问题，下一次甲来找我们玩的时候问他。下一次甲来的时候，大家都把问题藏在聊天里，" +
		"结果发现甲有很多常识性的错误，例如一个中午开放的限时副本他说他早上上学之前就打完了、" +
		"高等级才能携带的宠物他等级不够也说他有等等，我们大概得出一个结论：甲是不玩游戏的。正当我情绪高涨、" +
		"想要站在道德的制高点上对甲指指点点的时候，带头大哥只是点点头，说了一句：“嗯，快上课了，回教室吧”就走了，" +
		"给我整的情绪都不连贯了。又过了一段时间，大家都默契的不愿意和甲玩，我忽然意识到，此时此刻恰如彼时彼刻，" +
		"当年的我还小心翼翼地伪装自己也在玩游戏的时候，其他人是不是也揣着明白装糊涂的在糊弄我这个傻子玩呢？"
)
