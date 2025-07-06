package proxy

import (
	"dev_tool/api/response"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProxyRequest struct {
	URL string `json:"url" binding:"required,url"`
}

// ProxyUTXO 处理 UTXO 代理请求
func ProxyUTXO(c *gin.Context) {
	var req ProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.CustomError(c, http.StatusBadRequest, "Invalid request: "+err.Error(), nil)
		return
	}

	// 发起请求
	resp, err := http.Get(req.URL)
	if err != nil {
		response.CustomError(c, http.StatusInternalServerError, "Failed to fetch UTXO data: "+err.Error(), nil)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.CustomError(c, http.StatusInternalServerError, "Failed to read response: "+err.Error(), nil)
		return
	}

	// 设置响应头
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Status(resp.StatusCode)

	// 返回响应
	c.Writer.Write(body)
}
