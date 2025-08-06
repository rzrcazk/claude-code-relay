package router

import (
	"claude-code-relay/controller"
	"claude-code-relay/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetAPIRouter(server *gin.Engine) {

	// 全局限流：每分钟60次请求
	server.Use(middleware.RateLimit(60, time.Minute))

	// 健康检查
	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// API路由组
	api := server.Group("/api/v1")
	{
		// 公开接口
		auth := api.Group("/auth")
		{
			auth.POST("/login", controller.Login)
			auth.POST("/register", controller.Register)
		}

		// 系统状态
		api.GET("/status", controller.GetStatus)

		// 需要认证的接口
		authenticated := api.Group("")
		authenticated.Use(middleware.Auth())
		{
			// 用户相关
			user := authenticated.Group("/user")
			{
				user.GET("/profile", controller.GetProfile)
				user.PUT("/profile", controller.UpdateProfile)
			}

			// 任务相关
			task := authenticated.Group("/tasks")
			{
				task.GET("", controller.GetTasks)
				task.POST("", controller.CreateTask)
				task.GET("/:id", controller.GetTask)
				task.PUT("/:id", controller.UpdateTask)
				task.DELETE("/:id", controller.DeleteTask)
			}

			// 分组相关
			group := authenticated.Group("/groups")
			{
				group.GET("/list", controller.GetGroups)      // 获取分组列表
				group.POST("/create", controller.CreateGroup) // 创建分组
				group.GET("/:id", controller.GetGroup)        // 获取分组详情
				group.PUT("/:id", controller.UpdateGroup)     // 更新分组
				group.DELETE("/:id", controller.DeleteGroup)  // 删除分组
			}

			// 账号管理相关
			account := authenticated.Group("/accounts")
			{
				account.GET("/list", controller.GetAccountList)         // 获取账号列表
				account.POST("/create", controller.CreateAccount)       // 创建账号
				account.GET("/detail/:id", controller.GetAccount)       // 获取账号详情
				account.PUT("/update/:id", controller.UpdateAccount)    // 更新账号
				account.DELETE("/delete/:id", controller.DeleteAccount) // 删除账号
			}

			// Claude OAuth 相关
			oauth := authenticated.Group("/oauth")
			{
				oauth.GET("/generate-auth-url", controller.GetOAuthURL) // 获取OAuth授权URL
				oauth.POST("/exchange-code", controller.ExchangeCode)   // 验证授权码并获取token
			}

			// 管理员接口
			admin := authenticated.Group("/admin")
			admin.Use(middleware.AdminAuth())
			{
				admin.GET("/users", controller.GetUsers)
				admin.GET("/logs", controller.GetApiLogs)
				admin.GET("/dashboard", controller.GetDashboard)
			}
		}
	}

	// 静态文件服务（为前端预留）
	server.Static("/static", "./web/static")
	server.StaticFile("/", "./web/index.html")
}
