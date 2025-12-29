import { ref } from "vue"
import { defineStore } from "pinia"

export let useFlagStore = defineStore("flag", () => {
	let wildScreenFlag = ref<boolean>(true)

	function onScreenWidthChanged(width: number): void {
		wildScreenFlag.value = width > 1280
	}

	return { wildScreenFlag, onScreenWidthChanged }
})
