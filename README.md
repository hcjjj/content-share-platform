# 建设中 🔨

## 社交平台

**基本介绍**

* 用户登录功能
* 

> 如何启动前端：在 webook-fe 目录下先 `npm install` 后 `npm run dev`

**项目结构**

* 参考 [Kratos](https://go-kratos.dev/)、[go-zero](https://go-zero.dev/) 、[Domain-Driven Design](https://zhuanlan.zhihu.com/p/91525839)
* Service - Repository - DAO (Data Access Object) 三层结构 
  * service：领域服务（domain service），一个业务的完整处理过程
  * repository：领域对象的存储，存储数据的抽象
    * dao：数据库操作
  * domain：领域对象
* handler（和HTTP打交道） → service（主要业务逻辑） → repository（数据存储抽象） → dao（数据库操作）

## 技术栈

**第三方库**

* [gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP web 框架
  * [gin-contrib/cors](https://github.com/gin-contrib/cors) -  gin's middleware for *cross-origin resource sharing*

* [dlclark/regexp2](https://github.com/dlclark/regexp2) - full-featured 正则表达式

* [go-gorm/gorm](https://github.com/go-gorm/gorm) - The fantastic ORM library for Golang
  * [go-gorm/mysql](https://github.com/go-gorm/mysql) - GORM mysql driver

**Docker**

* [mysql](https://hub.docker.com/_/mysql)
* 

## 技术要点

**Gin**

* 定义、注册接口
* 后端处理（校验、处理、返回响应）
* [Middleware](https://github.com/gin-gonic/contrib)
  * AOP-Aspect-Oriented Programming
  * 解决跨域问题

**GORM**

* 模型定义
* 