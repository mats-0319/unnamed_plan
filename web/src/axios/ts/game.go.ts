// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.4

import { Pagination } from "./common.go"

export class UploadGameScoreReq {
	game_name: GameName = GameName.placeholder
	score: number = 0
	result: string = "" // json str
	player: string = "" // 已登录则为空，否则：'游客+[随机字符]'
}

export class UploadGameScoreRes {}

export class GameScore {
	score: number = 0
	result: string = ""
	player_name: string = "" // user nickname
	timestamp: number = 0
}

export class ListGameScoreReq {
	game_name: GameName = GameName.placeholder
	page: Pagination = new Pagination()
}

export class ListGameScoreRes {
	count: number = 0
	scores: Array<GameScore> = new Array<GameScore>()
}

export enum GameName {
	placeholder = -1,
	Flip = 1
}
