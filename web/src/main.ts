import {createApp} from 'vue'
import {createPinia} from 'pinia'

import App from './app.vue'
import {router} from './router.ts'

import 'element-plus/dist/index.css'
import "./index.less"

// axios init interceptors
import {initInterceptors} from "@/axios/ts/config_extend.ts";

initInterceptors((): void => {
    router.replace({name: "home", params: {v: "1"}}) // distinguish 'login error' router to 'home' with others
})

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
