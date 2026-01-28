<template>
	<top />
	<router-view />
</template>

<script setup lang="ts">
import { useFlagStore } from "@/pinia/flag.ts"
import { onMounted } from "vue"
import { useUserStore } from "@/pinia/user.ts"
import { deepCopy } from "@/ts/util.ts"
import Top from "@/components/top.vue"

let flags = useFlagStore()
let userStore = useUserStore()

onMounted(() => {
	// on re-size screen width
	flags.onScreenWidthChanged(screen.width)

	window.addEventListener("resize", () => {
		flags.onScreenWidthChanged(screen.width)
	})

	// keep 'login info' during refresh
	window.addEventListener("beforeunload", () => {
		sessionStorage.setItem("login_data", JSON.stringify(userStore.user))
	})

	let loginData = sessionStorage.getItem("login_data")
	if (loginData) {
		userStore.user = deepCopy(JSON.parse(loginData))
		sessionStorage.removeItem("login_data")
	}
})
</script>

<style lang="less">
body {
	padding: 0;
	margin: 0;

	overflow-x: hidden;
}
</style>
