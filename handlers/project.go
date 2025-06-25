package handlers

import (
	"net/http"
	"promanage/backend/db"
	"promanage/backend/models"

	"github.com/gin-gonic/gin"
)

// 获取所有项目
func GetProjects(c *gin.Context) {
	var projects []models.Project
	db.DB.Order("sort_order asc, id asc").Preload("Tasks").Find(&projects)
	c.JSON(http.StatusOK, gin.H{"data": projects})
}

// 新建项目
func CreateProject(c *gin.Context) {
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var p models.Project
	if v, ok := raw["name"].(string); ok {
		p.Name = v
	}
	if v, ok := raw["description"].(string); ok {
		p.Description = v
	}
	if v, ok := raw["due_date"].(string); ok {
		p.DueDate = parseDate(v)
	}
	if p.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "项目名称不能为空"})
		return
	}
	if err := db.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// 更新项目
func UpdateProject(c *gin.Context) {
	var p models.Project
	id := c.Param("id")
	if err := db.DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "项目不存在"})
		return
	}
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if v, ok := raw["name"].(string); ok {
		p.Name = v
	}
	if v, ok := raw["description"].(string); ok {
		p.Description = v
	}
	if v, ok := raw["due_date"].(string); ok {
		p.DueDate = parseDate(v)
	}
	if err := db.DB.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// 删除项目
func DeleteProject(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

// 更新项目排序
func SortProjects(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for idx, id := range req.IDs {
		db.DB.Model(&models.Project{}).Where("id = ?", id).Update("sort_order", idx)
	}
	c.JSON(http.StatusOK, gin.H{"message": "排序已更新"})
}
