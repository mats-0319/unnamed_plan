# web

使用工具：

- node: 24 (corepack-pnpm:10)
- vue3 (html ts:5.9 less)
- vite
- vue-router
- pinia
- element plus
- axios

node在16.13开始，提供`corepack`工具，这是一个**包管理器的管理器**

```cmd
npm  install --global corepack@latest
corepack enable pnpm
corepack use pnpm@latest // 升级pnpm版本
```

pnpm常用功能：

- `pnpm  audit`：检查已安装的包是否存在已知的安全问题
- `pnpm  outdated`：检查依赖是否有更新
- `pnpm why/ls [package name]`：列举一个依赖项在项目中的依赖关系，未指定依赖项时列举全部

## 开发小技巧

- TS7016:找不到指定模块的声明文件
  下载对应的声明文件，例如`import CryptoJs from "crypto-js"`出现该错误，则安装`@types/crypto-js`为开发依赖即可
