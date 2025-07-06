package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"

	"dev_tool/api"
	"dev_tool/api/middleware"
	"dev_tool/config"
	"dev_tool/models"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var distFS embed.FS

// SpaHandler serves a Single Page Application.
// It serves static files from the embedded filesystem.
// If a file is not found, it serves the index.html file.
type SpaHandler struct {
	StaticFS  fs.FS
	IndexFile string
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the requested path relative to the StaticFS root
	reqPath := path.Clean(r.URL.Path)
	if reqPath == "/" {
		reqPath = "/" + h.IndexFile
	}
	filePath := reqPath[1:] // Remove leading '/' for fs.FS access

	// Check if the file exists in the embedded FS
	_, err := h.StaticFS.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File does not exist, serve index.html
			indexData, err := fs.ReadFile(h.StaticFS, h.IndexFile)
			if err != nil {
				http.Error(w, "Internal Server Error: Index file not found", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(indexData)
			return
		}
		// Other error opening file
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// File exists, serve it using http.FileServer
	http.FileServer(http.FS(h.StaticFS)).ServeHTTP(w, r)
}

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	models.InitDB()

	// 设置 gin 模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// --- API Server (Port 8080) ---
	apiRouter := gin.Default()
	api.SetupRouter(apiRouter)
	apiAddr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	go func() {
		log.Printf("API 服务器启动在端口 %d\n", config.GlobalConfig.Server.Port)
		if err := apiRouter.Run(apiAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API 服务器启动失败: %v", err)
		}
	}()

	// --- Frontend Server (Port 8081) ---
	subFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatalf("无法访问嵌入的前端文件: %v", err)
	}

	// Create the SPA handler
	spaHandler := SpaHandler{
		StaticFS:  subFS,
		IndexFile: "index.html", // Relative to subFS
	}

	// Create a standard net/http server mux
	mux := http.NewServeMux()
	mux.Handle("/", spaHandler) // Handle all paths

	// Adapt Gin CORS middleware for net/http
	corsMiddleware := middleware.Cors()
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a dummy Gin context to wrap the net/http request/response
		c, _ := gin.CreateTestContext(w)
		c.Request = r
		corsMiddleware(c)

		// If CORS middleware (like preflight) wrote a response and aborted, don't continue
		if c.Writer.Status() != 0 && c.Writer.Status() != http.StatusOK {
			// We assume the middleware handled writing the response
			return
		}

		// If not aborted, proceed to the actual handler
		mux.ServeHTTP(w, r)
	})

	// Use the configured frontend port
	frontendPort := config.GlobalConfig.Server.FrontendPort
	if frontendPort == 0 {
		// Provide a default if not set or set to 0 in config
		log.Println("Warning: frontend_port not set or is 0 in config, defaulting to API port + 1")
		frontendPort = config.GlobalConfig.Server.Port + 1
	}
	frontendAddr := fmt.Sprintf(":%d", frontendPort)
	log.Printf("前端服务器启动在端口 %d\n", frontendPort)

	// Start the frontend server using net/http
	frontendServer := &http.Server{
		Addr:    frontendAddr,
		Handler: corsHandler, // Use the CORS-wrapped handler
	}

	if err := frontendServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("前端服务器启动失败: %v", err)
	}
}
