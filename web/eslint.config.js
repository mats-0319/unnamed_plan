import stylistic from "@stylistic/eslint-plugin" // eslint将格式化的内容提取到这个包维护了
import tsEslint from "typescript-eslint" // 比原生的defineConfig更智能
import vueEslintParser from "vue-eslint-parser" // 解析vue代码
import pluginVue from "eslint-plugin-vue" // 用于特殊条件（例如本文件中的vue html部分缩进控制）

export default tsEslint.config(
    { ignores: ["node_modules/**", "dist/**", "public/**","format_result.html"] }, // 全局忽略
	// 没有引入任何其他配置，保证eslint不会执行任何非预期行为
    {
        files: ["**/*.{js,ts,vue}"],
        plugins: {
            "@stylistic": stylistic,
            "vue": pluginVue
        },
        languageOptions: {
            parser: vueEslintParser, // 解析vue html
            parserOptions: {
                parser: tsEslint.parser, // 解析vue script(ts)
                ecmaVersion: "latest",
                sourceType: "module",
                extraFileExtensions: [".vue"]
            }
        },
        rules: {
            "@stylistic/indent": ["warn", 4], // 缩进
            "vue/html-indent": ["warn", 2],
            "@stylistic/max-len": ["warn", { code: 120, ignoreComments: true, ignoreUrls: true }], // 单行代码长度
            "@stylistic/semi": ["warn", "never"] // 分号
        }
    }
)
