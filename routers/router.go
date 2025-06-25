package routers

import (
	"promanage/backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 项目相关
	project := r.Group("/api/projects")
	{
		project.GET("", handlers.GetProjects)
		project.POST("", handlers.CreateProject)
		project.PUT(":id", handlers.UpdateProject)
		project.DELETE(":id", handlers.DeleteProject)

		// 任务相关
		project.GET(":id/tasks", handlers.GetTasks)
		project.POST(":id/tasks", handlers.CreateTask)
		project.PUT(":id/tasks/:tid", handlers.UpdateTask)
		project.DELETE(":id/tasks/:tid", handlers.DeleteTask)

		// 需求相关
		project.GET(":id/requirements", handlers.GetRequirements)
		project.POST(":id/requirements", handlers.CreateRequirement)
		project.PUT(":id/requirements/:rid", handlers.UpdateRequirementStatus)
	}

	r.POST("/api/projects/sort", handlers.SortProjects)
	r.POST("/api/ai/chat", handlers.AIChatHandler)
	// 操作密钥校验接口
	r.POST("/api/check-manager-key", handlers.CheckManagerKeyHandler)
}
