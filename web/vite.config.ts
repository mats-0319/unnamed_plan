import { fileURLToPath, URL } from "node:url"
import os from "os"

import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import vueDevTools from "vite-plugin-vue-devtools"

// element-plus auto import
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"
import { ElementPlusResolver } from "unplugin-vue-components/resolvers"

// https://vite.dev/config/
export default defineConfig({
	envPrefix: "Vite_",
	plugins: [
		vue(),
		vueDevTools(),

		AutoImport({ resolvers: [ElementPlusResolver()] }),
		Components({ resolvers: [ElementPlusResolver()] })
	],
	resolve: { alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) } },
	clearScreen: false,
	server: {
		host: getLocalIP(),
		port: 20319,
		open: true
	}
})

export function getLocalIP(): string {
	const networks = os.networkInterfaces()
	for (let key in networks) {
		// @ts-ignore
		for (let ins of networks[key]) {
			if (ins.family === "IPv4" && !ins.internal) {
				return ins.address
			}
		}
	}

	return "127.0.0.1"
}
