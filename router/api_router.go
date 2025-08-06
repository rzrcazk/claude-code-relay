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
