import { defineStore } from "pinia"
import { GameName, GameScore, ListGameScoreRes, UploadGameScoreRes } from "@/axios/ts/game.go.ts"
import { gameAxios } from "@/axios/ts/game.http.ts"
import { log } from "@/ts/log.ts"

export const useGameScoreStore = defineStore("game_score", () => {
    function uploadGameScore(game_name: GameName, score: number, result: string, player: string, cb: () => void): void {
        gameAxios.uploadGameScore(game_name, score, result, player).then(({}: { data: UploadGameScoreRes }) => {
            cb()

            log.success("upload game score")
        })
    }

    function listGameScore(
        game_name: GameName,
        pageSize: number,
        pageNum: number,
        cb: (count: number, scores: Array<GameScore>) => void,
    ): void {
        gameAxios.
            listGameScore(game_name, { size: pageSize, num: pageNum }).
            then(({ data }: { data: ListGameScoreRes }) => {
                cb(data.count, data.scores)

                log.success("list game score")
            })
    }

    return { uploadGameScore, listGameScore }
})
