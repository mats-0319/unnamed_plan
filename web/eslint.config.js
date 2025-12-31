import { defineConfig, globalIgnores } from "eslint/config"
import vueTs from "@vue/eslint-config-typescript"
import vuePrettier from "@vue/eslint-config-prettier"
import vueParser from "vue-eslint-parser"
import tsParser from "@typescript-eslint/parser"

export default defineConfig([
	globalIgnores(["node_modules/", "dist/"]),
	vueTs(),
	vuePrettier,
	{
		files: ["**/*.{js,ts,vue}"],
		languageOptions: { parser: vueParser, parserOptions: { parser: tsParser, sourceType: "esnext" } },
		rules: {
			"@typescript-eslint/ban-ts-comment": 0,
			"@typescript-eslint/no-array-constructor": 0,
			"prefer-const": 0,
			"@typescript-eslint/no-unused-vars": 0,
			"@typescript-eslint/no-explicit-any": 0
		}
	}
])
