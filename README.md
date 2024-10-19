## 1. 基本介绍
> go-admin 是基于 Gin 和 Gorm 开发的管理系统，主要具备 JWT + Casbin 鉴权、动态路由、动态菜单以及文件上传下载等功能。该项目旨在**学习和实践** Gin 的使用，源自于某开源项目的改进与扩展。

## 2. 使用说明
```
- golang版本 >= v1.22
- IDE推荐：Goland
```

### 2.1 初始化项目
```bash
# 下载项目
git clone https://github.com/town-coding/go-admin
# 进入项目文件夹
cd go-admin
# 使用 go mod 并安装go依赖包
go mod tidy
# 修改config.yaml文件，配置数据库类型，配置数据信息
# 初始化 数据库
go run main.go init
# 启动项目
go run main.go
```
## 3. 技术选型
- 后端：使用 Gin 构建基础的 RESTful 风格 API。
- 数据库：项目集成 PostgreSQL 和 MySQL 作为数据存储，通过 Gorm 实现基本的数据库操作。
- 日志：集成 Zap 实现高效的日志记录。
- 缓存：使用 Redis 存储当前活跃用户的 JWT 令牌，并实现多点登录限制。
- 其他中间件：
  - 使用 Viper 动态读取系统配置，结合 Pflag 匹配命令行参数，并通过 Cobra 实现命令行操作。
  - 集成 Cron 用于定时任务管理。

## 4. 主要功能

- 权限管理：基于 JWT 和 Casbin 实现的权限管理系统。
- 多点登录：利用 JWT 令牌和 Redis 缓存实现多点登录功能。
- 文件上传下载：支持基于阿里云和本地的文件上传操作。
- 用户管理：系统管理员可分配用户角色及角色权限。
- 角色管理：创建主要的权限控制对象，为角色分配不同的 API 权限和菜单权限。
- 菜单管理：实现用户动态菜单配置，支持不同角色的菜单显示。
