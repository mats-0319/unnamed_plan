<template>
  <div class="game color-bg-0">
    <div class="g-left">
      <elevated-button class="gl-item" @click="routerLink('gFlip')">Flip</elevated-button>
    </div>

    <el-divider direction="vertical" />

    <div class="g-right">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import ElevatedButton from "@/components/elevated_button.vue"
import { randomVisitorName, routerLink } from "@/ts/util.ts"
import { onMounted, onUnmounted } from "vue"
import { GameName } from "@/axios/ts/game.go.ts"
import { useUserStore } from "@/pinia/user.ts"
import { useGameScoreStore } from "@/pinia/game_score.ts"

let userStore = useUserStore()
let gameScoreStore = useGameScoreStore()

onMounted(() => {
    window.addEventListener("message", handleMessage)
})

onUnmounted(() => {
    window.removeEventListener("message", handleMessage)
})

function handleMessage(event: any) {
    console.log("> Node: test post-message event. ", event)

    const trustedOrigins = Array<string>(getBaseUrl())
    if (!trustedOrigins.includes(event.origin)) {
        console.log("> PostMessage - Invalid Event Origin: ", event)
        return
    }

    const { game_name, score, result } = event.data

    switch (game_name) {
        case GameName.Flip:
            let player = userStore.isLogin() ? "" : randomVisitorName()

            gameScoreStore.uploadGameScore(game_name, score, result, player)

            break
        default:
            console.log("> PostMessage - Invalid Game Name: ", game_name)
            break
    }
}

function getBaseUrl(): string {
    let url = import.meta.env.Vite_axios_base_url
    let localIP = window.location.hostname

    return import.meta.env.DEV ? url.replace("127.0.0.1", localIP) : url
}
</script>

<style lang="less" scoped>
.game {
	display: flex;
	height: calc(100vh - 6.25rem);

	.g-left {
		width: calc(20vw - 1px);
		margin: 6rem 2vw;
		font-size: 1.4rem;

		.gl-item {
			margin: 2rem auto;
		}
	}

	.g-right {
		width: 60vw;
		margin: 1rem 8vw;
	}

	.el-divider--vertical {
		width: 0;
		height: 100%;
		margin: 0;
		border-left: 1px solid darkgray;
	}
}
</style>
