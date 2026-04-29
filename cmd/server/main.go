package main

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"sms-dashboard/internal/config"
	"sms-dashboard/internal/database"
	"sms-dashboard/internal/handler"
	"sms-dashboard/internal/middleware"
	"sms-dashboard/web"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	database.InitDB(cfg.DatabasePath)

	r := setupRouter(cfg)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func setupRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Token")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := handler.NewAuthHandler(cfg)
	smsHandler := handler.NewSMSHandler()

	// Initialize default user
	authHandler.InitDefaultUser()

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/sms", middleware.APITokenMiddleware(cfg), smsHandler.Receive)

		// Private (Protected by JWT)
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware(cfg))
		{
			authorized.POST("/change-password", authHandler.ChangePassword)
			authorized.GET("/sms/list", smsHandler.List)
			authorized.GET("/sms/grouped", smsHandler.GroupedList)
			authorized.GET("/sms/load-more", smsHandler.LoadMore)
			authorized.GET("/sms/load-all", smsHandler.LoadAll)
			authorized.DELETE("/sms/:id", smsHandler.Delete)
			authorized.POST("/sms/batch-delete", smsHandler.BatchDelete)
			authorized.GET("/sms/search", smsHandler.Search)
		}
	}

	// Serve Frontend
	distFS, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatal(err)
	}

	r.NoRoute(func(c *gin.Context) {
		path := strings.TrimPrefix(c.Request.URL.Path, "/")

		// Handle root path or frontend routes
		if path == "" || !strings.Contains(path, ".") {
			file, err := distFS.Open("index.html")
			if err != nil {
				c.String(http.StatusNotFound, "index.html not found")
				return
			}
			defer file.Close()
			content, err := io.ReadAll(file)
			if err != nil {
				c.String(http.StatusInternalServerError, "error reading index.html")
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", content)
			return
		}

		// Serve static files (assets etc)
		f, err := distFS.Open(path)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		f.Close()

		c.FileFromFS(path, http.FS(distFS))
	})

	return r
}
