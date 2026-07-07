import { ref } from "vue"
import { defineStore } from "pinia"

export const useFlagStore = defineStore("flag", () => {
    const routeParam = ref<string>("")
    const wildScreenFlag = ref<boolean>(true)
    const loading = ref<boolean>(false)

    function setRouteParam(value: string): void {
        routeParam.value = value

        setTimeout(() => {
            routeParam.value = ""
        }, 1000)
    }

    function onScreenWidthChanged(width: number): void {
        wildScreenFlag.value = width > 1280
    }

    function setLoading(flag: boolean): void {
        loading.value = flag
    }

    return { routeParam, setRouteParam, wildScreenFlag, onScreenWidthChanged, loading, setLoading }
})
