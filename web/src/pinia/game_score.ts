import { defineStore } from "pinia"
import { GameName, GameScore, ListGameScoreRes, UploadGameScoreRes } from "@/axios/ts/game.go.ts"
import { gameAxios } from "@/axios/ts/game.http.ts"
import { log } from "@/ts/log.ts"
import { ref } from "vue"

export const useGameScoreStore = defineStore("game_score", () => {
    let count = ref<number>(0)
    let scores = ref<Array<GameScore>>(new Array<GameScore>())

    function uploadGameScore(game_name: GameName, score: number, result: string, player: string): void {
        gameAxios.uploadGameScore(game_name, score, result, player).then(({}: { data: UploadGameScoreRes }) => {
            listGameScore(game_name, 10, 1)

            log.success("upload game score")
        })
    }

    function listGameScore(game_name: GameName, pageSize: number, pageNum: number): void {
        gameAxios.
            listGameScore(game_name, { size: pageSize, num: pageNum }).
            then(({ data }: { data: ListGameScoreRes }) => {
                count.value = data.count
                scores.value = data.scores

                log.success("list game score")
            })
    }

    return { count, scores, uploadGameScore, listGameScore }
})
