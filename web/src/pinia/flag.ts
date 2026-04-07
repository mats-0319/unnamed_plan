import { ref } from "vue"
import { defineStore } from "pinia"

export const useFlagStore = defineStore("flag", () => {
    const wildScreenFlag = ref<boolean>(true)

    function onScreenWidthChanged(width: number): void {
        wildScreenFlag.value = width > 1280
    }

    return { wildScreenFlag, onScreenWidthChanged }
})
