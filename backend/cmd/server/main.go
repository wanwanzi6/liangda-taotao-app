package main

import (
	"liangda-taotao/config"
	"liangda-taotao/internal/handler"
	"liangda-taotao/internal/middleware"
	"liangda-taotao/internal/model"
	"liangda-taotao/internal/repository"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// initDB 负责初始化数据库连接
func initDB() *gorm.DB {
	dsn := config.GetDSN(config.AppDBConfig)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ 数据库无响应: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&model.User{}, model.Product{}, model.Category{})
	if err != nil {
		log.Fatalf("❌ 自动迁移失败: %v", err)
	}

	log.Println("✅ 数据库初始化完成")
	return db
}

// seedData 负责填充初始分类
func seedData(db *gorm.DB) {
	categories := []string{"代步工具", "数码电子", "美妆护理", "教材资料", "运动户外",
		"生活电器", "技能服务", "零食饮品"}
	for _, name := range categories {
		db.FirstOrCreate(&model.Category{}, model.Category{Name: name})
	}
	log.Println("🌱 种子数据检查/填充完成")
}

func main() {
	// 1. 初始化数据库
	db := initDB()
	seedData(db)

	// 2. 初始化 Repository 层
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 3. 初始化 Handler 层
	categoryHandler := handler.NewCategoryHandler(categoryRepo)
	productHandler := handler.NewProductHandler(productRepo)
	userHandler := handler.NewUserHandler(userRepo)

	// 4. 启动 Gin 路由
	r := gin.Default()

	r.Static("/uploads", "./uploads")

	//设置路由组
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		// 微信登录
		v1.POST("/login", userHandler.WeChatLogin)
		// 获取全部分类
		v1.GET("/categories", categoryHandler.GetAll)
		// 获取商品列表
		v1.GET("/products", productHandler.GetList)
		// 查看商品详情
		v1.GET("/products/:id", productHandler.GetDetail)
		// 上传商品图片
		v1.POST("/upload", productHandler.Upload)

		// 需要登录的接口
		auth := v1.Group("")
		auth.Use(middleware.AuthMiddleware())
		{
			// 获取当前用户信息
			auth.GET("/user", userHandler.GetUserInfo)
			// 更新用户信息
			auth.PUT("/user", userHandler.UpdateUserInfo)
			// 发布商品（需要登录）
			auth.POST("/products", productHandler.Create)
			// 删除商品（需要登录）
			auth.DELETE("/products/:id", productHandler.Delete)
		}
	}

	// 5. 启动并在 8080 端口监听
	log.Println("✨ 量大淘淘后端服务已在 8080 端口启动...")
	r.Run(":8080")

}
