<template>
  <div class="game color-bg-0">
    <div class="g-left">
      <elevated_button class="gl-item" @click="routerLink('gFlip')">Flip</elevated_button>
    </div>

    <el-divider direction="vertical" />

    <div class="g-right">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import Elevated_button from "@/components/elevated_button.vue"
import { routerLink } from "@/ts/util.ts"
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
