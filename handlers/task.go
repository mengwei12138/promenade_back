package handlers

import (
	"fmt"
	"net/http"
	"promanage/backend/db"
	"promanage/backend/models"
	"time"

	"github.com/gin-gonic/gin"
)

// 获取某项目所有任务
func GetTasks(c *gin.Context) {
	var tasks []models.Task
	pid := c.Param("id")
	if err := db.DB.Where("project_id = ?", pid).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func parseDate(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		return &t
	}
	t2, err2 := time.Parse(time.RFC3339, s)
	if err2 == nil {
		return &t2
	}
	return nil
}

// 新建任务
func CreateTask(c *gin.Context) {
	pid := c.Param("id")
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var t models.Task
	if v, ok := raw["name"].(string); ok {
		t.Name = v
	}
	if v, ok := raw["description"].(string); ok {
		t.Description = v
	}
	if v, ok := raw["status"].(string); ok {
		t.Status = v
	}
	t.ProjectID = parseUint(pid)
	if v, ok := raw["start_date"].(string); ok {
		t.StartDate = parseDate(v)
	}
	if v, ok := raw["end_date"].(string); ok {
		t.EndDate = parseDate(v)
	}
	if t.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "任务名称不能为空"})
		return
	}
	if t.Status == "" {
		t.Status = "Pending"
	}
	if err := db.DB.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

// 更新任务
func UpdateTask(c *gin.Context) {
	var t models.Task
	tid := c.Param("tid")
	if err := db.DB.First(&t, tid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}
	var raw map[string]interface{}
	if err := c.ShouldBindJSON(&raw); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if v, ok := raw["name"].(string); ok {
		t.Name = v
	}
	if v, ok := raw["description"].(string); ok {
		t.Description = v
	}
	if v, ok := raw["status"].(string); ok {
		t.Status = v
	}
	if v, ok := raw["start_date"].(string); ok {
		t.StartDate = parseDate(v)
	}
	if v, ok := raw["end_date"].(string); ok {
		t.EndDate = parseDate(v)
	}
	if err := db.DB.Save(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

// 删除任务
func DeleteTask(c *gin.Context) {
	tid := c.Param("tid")
	if err := db.DB.Delete(&models.Task{}, tid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已删除"})
}

func parseUint(s string) uint {
	var i uint64
	fmt.Sscanf(s, "%d", &i)
	return uint(i)
}
