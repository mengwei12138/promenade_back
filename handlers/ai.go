package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AIChatRequest struct {
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func AIChatHandler(c *gin.Context) {
	var req AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API Key 未配置"})
		return
	}

	// 拼接context为自然语言
	contextText := ""
	if req.Context != nil {
		if projects, ok := req.Context["projects"].([]interface{}); ok {
			contextText += "你当前管理的项目有："
			for _, p := range projects {
				if proj, ok := p.(map[string]interface{}); ok {
					contextText += proj["name"].(string)
					if desc, ok := proj["description"].(string); ok && desc != "" {
						contextText += "（" + desc + "）"
					}
					if tasks, ok := proj["tasks"].([]interface{}); ok && len(tasks) > 0 {
						contextText += "，任务："
						for _, t := range tasks {
							if task, ok := t.(map[string]interface{}); ok {
								contextText += task["name"].(string) + ","
							}
						}
					}
					if reqs, ok := proj["requirements"].([]interface{}); ok && len(reqs) > 0 {
						contextText += "，需求："
						for _, r := range reqs {
							if req, ok := r.(map[string]interface{}); ok {
								contextText += req["content"].(string)
								if proposer, ok := req["proposer"].(string); ok && proposer != "" {
									contextText += "（提出人:" + proposer + ")"
								}
								if status, ok := req["status"].(string); ok && status != "" {
									contextText += "（状态:" + status + ")"
								}
								contextText += ","
							}
						}
					}
					contextText += "; "
				}
			}
		}
	}
	// 优化AI风格：要求用简洁中文、非Markdown格式直接回答，不要反问用户，并结合项目的任务和需求信息。
	systemPrompt := "你是一个项目管理智能助手。请用简洁中文、非Markdown格式直接回答，不要反问用户。请结合项目的任务和需求信息进行回答。" + contextText

	requestBody := RequestBody{
		Model: "qwen-plus",
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: req.Message},
		},
	}
	jsonData, _ := json.Marshal(requestBody)
	httpReq, _ := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var respData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respData)

	// 解析通义千问返回的内容
	reply := ""
	if choices, ok := respData["choices"].([]interface{}); ok && len(choices) > 0 {
		if msg, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{}); ok {
			reply, _ = msg["content"].(string)
		}
	}
	if reply == "" {
		reply = "AI助手暂时无法回答，请稍后再试。"
	}
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}
