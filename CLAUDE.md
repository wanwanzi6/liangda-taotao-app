# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# 量大淘淘 (Liangda Taotao) 项目文档

## 项目简介
量大淘淘是一个专为量大校友开发的校园二手交易市场小程序。
- **后端架构**: Go (Gin) + GORM (MySQL) + JWT 认证
- **后端端口**: `8080`
- **API 前缀**: `/api/v1`
- **前端架构**: 微信小程序 (原生开发)

## 当前开发进度

### 后端 (Backend) - `已实现`

**分层架构**: Handler/Repository/Model

**数据库模型** (`internal/model/`):
- `User`: 微信 OpenID、昵称、头像、校内认证状态 (IsVerify)
- `Product`: 商品信息、价格 (Decimal)、状态管理、软删除
- `Category`: 8 大校园分类 (代步工具、数码电子、美妆护理、教材资料、运动户外、生活电器、技能服务、零食饮品)

**公开接口** (`/api/v1/`):
- `POST /login` - 微信登录凭证校验并生成 JWT Token
- `GET /categories` - 获取商品分类列表
- `GET /products` - 分页及按分类筛选商品列表
- `GET /products/:id` - 商品详情
- `POST /upload` - 商品图片上传至本地服务器

**鉴权接口** (需 JWT Token):
- `POST /refresh` - 刷新 Token
- `POST /logout` - 退出登录
- `GET /user` - 获取当前用户信息
- `PUT /user` - 更新用户信息
- `GET /user/products` - 获取我的发布列表
- `POST /products` - 发布商品
- `PUT /products/:id` - 更新商品
- `DELETE /products/:id` - 下架商品
- `POST /products/:id/refresh` - 一键擦亮

### 前端 (Frontend) - `已实现`

**页面结构** (`pages/`):
- `index/` - 首页（搜索框、分类滑动切换、瀑布流商品展示）
- `detail/` - 商品详情页
- `post/` - 发布商品页面
- `profile/` - 个人中心（我的发布、擦亮、下架商品）
- `logs/` - 调试日志页

**核心逻辑** (`utils/` 和根目录):
- `app.js` - 全局登录态校验、`requestWithAuth` 鉴权请求封装
- `utils/auth.js` - Token 和用户信息的本地同步存储

## 待办事项 (TODO)
1. **校内认证**: 实现 `User.IsVerify` 认证逻辑（学号/校内邮箱验证）
2. **交互功能**: 留言、收藏、即时通讯功能
3. **性能优化**: Redis 缓存分类数据、优化大图加载

## 开发命令

### 后端
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 前端
使用 **微信开发者工具** 打开 `frontend` 目录即可调试。

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

### 后端结构
```
backend/
├── cmd/server/main.go      # 入口，初始化 DB/Repository/Handler，启动 Gin
├── config/config.go       # 配置管理，从 .env 读取
├── internal/
│   ├── handler/           # HTTP 请求处理
│   ├── repository/       # 数据库操作
│   ├── model/             # GORM 模型
│   └── middleware/       # JWT 鉴权中间件
└── uploads/               # 上传的图片存储
```

### 前端结构
```
frontend/
├── app.js                 # 全局逻辑、登录态校验
├── app.json               # 页面配置、tabBar
├── utils/auth.js          # Token/用户信息本地存储
└── pages/
    ├── index/             # 首页
    ├── detail/            # 详情页
    ├── post/              # 发布页
    ├── profile/           # 个人中心（我的发布）
    └── logs/              # 日志页
```
