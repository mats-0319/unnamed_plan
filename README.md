# 未命名计划

[//]: # ([![Go Reference]&#40;https://pkg.go.dev/badge/github.com/mats9693/unnamed_plan.svg&#41;]&#40;https://pkg.go.dev/github.com/mats9693/unnamed_plan&#41;)

[//]: # ([![Go Report Card]&#40;https://goreportcard.com/badge/github.com/mats9693/unnamed_plan&#41;]&#40;https://goreportcard.com/report/github.com/mats9693/unnamed_plan&#41;)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)

一个web项目，用来学习web开发和部署，项目会部署到公网，如果[公网地址](https://115.190.167.134)不可用，多半是我买的云服务器到期了 ^_^ （2026.9.25）

## 使用

如果你想在本地运行该项目，需要以下环境：（版本号为经过测试的推荐版本）

- go 1.25
- node 24 (corepack-pnpm)
- postgresql 16

本地运行：`server/.run/`文件夹包含服务端程序构建配置（jetbrains ide可用），修改配置文件即可启动服务端程序。

云端部署：执行`scripts/build.sh`，会在项目根目录下生成`build`文件夹，将该文件夹上传到云服务器，执行其中的shell脚本即可启动服务端程序。

## 设计图

> 感谢：[绘图工具](https://excalidraw.com)，[压缩工具](https://tinypng.com)

todo

## 项目介绍(draft)

目录：

- build：生成的部署用内容
- doc：文档
    - deploy：部署文档，介绍把程序发布到公网的主要步骤（含nginx反向代理配置）
- scripts：脚本
    - build.sh：在本地使用，生成可以部署到云服务器的内容（包括几个服务端程序及启动脚本、web页面）
    - restart_server.sh：在云服务器使用，启动所有服务端程序
- server：服务端代码
- web：ui代码

主要技术、框架与工具：

- go 1.25
- gorm
    - ORM: Object Relation Mapping，对象关系映射，将数据库的行、列甚至数据库本身，映射成编程语言的对象或字段。
      使用ORM，可以通过形如`db.create(user)`的方式操作数据库，而不需要编写形如`insert into user values (...)`的sql
    - DAO: Data Access Object，数据访问对象，将数据库操作包装在一起，与业务代码分离。
      主要应用场景有：需要接多个数据库（pg、mysql、sqlite）、数据操作复杂（例如复杂查询，sql写出来上KB的）
- gocts：自研工具，可以根据go定义的接口结构，生成对应的ts结构（class）和axios client代码。

- vue3 (html+ts+less)
- node 24 (corepack-pnpm)
- vite
- axios
- vue-router
- pinia
- element-plus
- eslint+prettier

- nginx
- postgresql
- shell
