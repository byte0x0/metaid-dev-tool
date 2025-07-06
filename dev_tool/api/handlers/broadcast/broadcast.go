package broadcast

import (
	"bytes"
	"dev_tool/api/response"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Request struct {
	URL     string `json:"url" binding:"required,url"`
	Content string `json:"content" binding:"required"`
}

type BroadcastRequest struct {
	RawTx string `json:"rawtx"`
}

// ProxyBroadcast 处理广播代理请求
func ProxyBroadcast(c *gin.Context) {
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CustomError(c, http.StatusBadRequest, "Invalid request: "+err.Error(), nil)
		return
	}

	// 构造广播请求体
	broadcastReq := BroadcastRequest{
		RawTx: req.Content,
	}

	// 将请求体转换为 JSON
	jsonData, err := json.Marshal(broadcastReq)
	if err != nil {
		response.CustomError(c, http.StatusInternalServerError, "Failed to marshal request: "+err.Error(), nil)
		return
	}

	// 创建客户端，允许自动重定向
	client := &http.Client{}

	// 发起 POST 请求
	resp, err := client.Post(req.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		response.CustomError(c, http.StatusInternalServerError, "Failed to broadcast transaction: "+err.Error(), nil)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.CustomError(c, http.StatusInternalServerError, "Failed to read response: "+err.Error(), nil)
		return
	}

	// 如果响应不是 JSON 格式，直接返回原始响应
	if !json.Valid(body) {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
		return
	}

	// 返回 JSON 响应
	c.JSON(resp.StatusCode, json.RawMessage(body))
}
