# 建设中 🔨

## 社交平台

**基本介绍**

* 用户登录服务
  * 注册、密码加密存储
  * 登录、登录态校验
    * Cookie + Session
    * Session 存储基于 Redis 实现（多实例部署环境）
      * 但是每次请求都要访问 Redis，性能瓶颈问题⚠️
    * 刷新登录状态（利用 Session 有效期机制）
    * 后续换为 JWT（JSON Web Token）机制
  * 保护登录系统
    * 限流
    * 增强登录安全
* 用户关系服务
* 发帖服务
* 支付服务
* 搜索服务
* 即时通讯
* Feed 流

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
  * [cors](https://github.com/gin-contrib/cors) -  Official *cross-origin resource sharing* (CORS) gin's middleware
  * [sessions](https://github.com/gin-contrib/sessions) - Gin middleware for session management
* [dlclark/regexp2](https://github.com/dlclark/regexp2) - full-featured 正则表达式
* [go-gorm/gorm](https://github.com/go-gorm/gorm) - The fantastic ORM library for Golang
  * [go-gorm/mysql](https://github.com/go-gorm/mysql) - GORM mysql driver
* [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - Golang implementation of JSON Web Tokens (JWT)
* 

**Docker**

* [mysql](https://hub.docker.com/_/mysql)
* [redis](https://hub.docker.com/_/redis)

## 技术要点

**[Gin](https://gin-gonic.com/zh-cn/docs/)**

* 定义、注册接口
* 后端处理（校验、处理、返回响应）
* [Middleware](https://github.com/gin-gonic/contrib)
  * AOP-Aspect-Oriented Programming
  * 解决跨域问题

**[GORM](https://gorm.io/zh_CN/)** 

* 模型定义
* 