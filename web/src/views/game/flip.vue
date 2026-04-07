<template>
  <div class="game-flip">
    <iframe :src="gameUrl" width="380" height="610" />

    <div class="gf-score">
      <el-collapse accordion>
        <div>Score Count:&nbsp;{{ scoreCount }}</div>
        <el-collapse-item v-for="(item, index) in topScore" :key="index">
          <template #title>{{ index + 1 }}{{ ". " + item.player_name +" "+ item.score }}</template>

          <div>{{ item.result }}</div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue"
import { useGameScoreStore } from "@/pinia/game_score.ts"
import { useUserStore } from "@/pinia/user.ts"
import { randomVisitorName } from "@/ts/util.ts"
import { GameName, GameScore } from "@/axios/ts/game.go.ts"

let userStore = useUserStore()
let gameScoreStore = useGameScoreStore()

let scoreCount = ref<number>(0)
let topScore = ref<Array<GameScore>>(new Array<GameScore>())

let gameUrl = ref<string>(import.meta.env.Vite_axios_flip_game_url)

onMounted(() => {
    window.addEventListener("message", handleMessage)
    listScore()
})

onUnmounted(() => {
    window.removeEventListener("message", handleMessage)
})

function handleMessage(event: any) {
    const { game_name, score, result } = event.data

    if (game_name != GameName.Flip) {
        console.log("invalid game name")
        return
    }

    let player = userStore.isLogin() ? "" : randomVisitorName()

    gameScoreStore.uploadGameScore(game_name, score, result, player, () => {
        listScore()
    })
}

function listScore() {
    gameScoreStore.listGameScore(GameName.Flip, 10, 1, (count: number, scores: Array<GameScore>) => {
        scoreCount.value = count
        topScore.value = scores
    })
}
</script>

<style lang="less" scoped>
.game-flip {
	display: flex;

	.gf-score {
		margin-left: 8vw;
		width: 20vw;

		font-size: 1.6rem;
	}

	iframe {
		transform-origin: left top;
		@media (min-height: 1440px) {
			transform: scale(1.4, 1.4); // 大分辨率屏幕放大一些，1920*1080不用缩放
		}
	}
}
</style>
