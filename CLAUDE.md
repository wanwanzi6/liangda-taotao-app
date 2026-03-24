# 量大淘淘 (Liangda Taotao) 项目文档

## 项目简介
量大淘淘是一个专为量大校友开发的校园二手交易市场小程序。
- **后端架构**: Go (Gin) + GORM (MySQL) + JWT 认证
- **前端架构**: 微信小程序 (原生开发)

## 当前开发进度

### 后端 (Backend) - `已实现`
- **基础架构**: 采用 Handler/Repository/Model 分层架构。
- **数据库设计**:
  - `User`: 支持微信 OpenID、昵称、头像及校内认证状态。
  - `Product`: 商品信息、价格（Decimal）、状态管理及软删除。
  - `Category`: 预设 8 大常用校园分类（代步、数码、美妆等）。
- **核心接口**:
  - `POST /login`: 微信登录凭证校验并生成 JWT Token。
  - `GET /categories`: 获取商品分类列表。
  - `GET /products`: 支持分页及按分类筛选商品列表。
  - `POST /products`: 鉴权发布商品。
  - `POST /upload`: 商品图片上传至本地服务器。

### 前端 (Frontend) - `已实现`
- **核心页面**:
  - `pages/index`: 沉浸式搜索框、分类滑动切换、瀑布流商品展示。
  - `pages/detail`: 商品详情展示，包含价格、描述及占位图逻辑。
  - `pages/post`: 完整的发布表单，支持图片选择、上传预览及分类选择。
- **功能逻辑**:
  - `utils/auth.js`: 封装本地同步存储 Token 和用户信息。
  - `app.js`: 全局登录态校验及 `requestWithAuth` 带鉴权请求封装。

## 待办事项 (TODO)
1. **校内认证**: 实现 `User` 模型中的 `IsVerify` 认证逻辑（如学号/校内邮箱验证）。
2. **交互功能**: 实现详情页中的“留言”、“收藏”以及“聊一聊”即时通讯功能。
3. **商品管理**: 前端增加“我的发布”页面，调用后端的 `DELETE` 接口下架商品。
4. **性能优化**: 引入 Redis 缓存分类数据，优化大图加载速度。

## 开发规范
- **后端端口**: `8080`。
- **API 前缀**: `/api/v1`。
- **代码规范**: 后端使用 GORM 自动迁移，前端样式遵循”闲鱼黄”视觉风格。

## 开发命令

### 后端 (Backend)
```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 运行开发服务器（自动迁移数据库）
go run cmd/server/main.go

# 或直接运行已编译的二进制
./server.exe
```

### 前端 (Frontend)
使用 **微信开发者工具** 打开 `frontend` 目录即可开发调试。

### 环境变量 (backend/.env)
```
DB_USER=root
DB_PASSWORD=your_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=liang_da_tao_tao
WECHAT_APPID=your_appid
WECHAT_SECRET=your_secret
JWT_SECRET=your_jwt_secret
```

## 架构概览

### 后端分层
- `cmd/server/main.go` - 入口，初始化 DB、Repository、Handler，启动 Gin 服务
- `config/` - 配置管理，从 .env 读取数据库、微信、JWT 配置
- `internal/handler/` - HTTP 请求处理，解析参数、调用 repository、返回响应
- `internal/repository/` - 数据库操作，GORM 查询封装
- `internal/model/` - GORM 模型定义 (User, Product, Category)
- `internal/middleware/` - 中间件 (JWT 鉴权)

### 前端结构
- `app.js` - 全局逻辑，登录态校验、requestWithAuth 封装
- `utils/auth.js` - Token 和用户信息的本地存储
- `pages/index/` - 首页（搜索、分类切换、商品列表）
- `pages/detail/` - 商品详情页
- `pages/post/` - 发布商品页面
- `pages/profile/` - 个人中心页面