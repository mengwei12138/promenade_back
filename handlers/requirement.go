package handlers

import (
	"fmt"
	"net/http"
	"promanage/backend/db"
	"promanage/backend/models"

	"github.com/gin-gonic/gin"
)

// 获取项目下所有需求
func GetRequirements(c *gin.Context) {
	projectID := c.Param("id")
	var requirements []models.Requirement
	if err := db.DB.Where("project_id = ?", projectID).Order("created_at desc").Find(&requirements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": requirements})
}

// 新建需求
func CreateRequirement(c *gin.Context) {
	projectID := c.Param("id")
	var req models.Requirement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	if req.Content == "" || req.Proposer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "需求内容和提出人不能为空"})
		return
	}
	req.ProjectID = parseUintLocal(projectID)
	req.Status = "Pending"
	if err := db.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": req})
}

// 处理需求（改状态）
func UpdateRequirementStatus(c *gin.Context) {
	projectID := c.Param("id")
	reqID := c.Param("rid")
	var req models.Requirement
	if err := db.DB.Where("project_id = ? AND id = ?", projectID, reqID).First(&req).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到该需求"})
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || (body.Status != "Pending" && body.Status != "Resolved") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误或状态非法"})
		return
	}
	req.Status = body.Status
	if err := db.DB.Save(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": req})
}

// 工具函数：字符串转uint
func parseUintLocal(s string) uint {
	var n uint
	_, _ = fmt.Sscan(s, &n)
	return n
}
