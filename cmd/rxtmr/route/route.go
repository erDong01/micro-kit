package route

import (
	"rxt/internal/http/middleware"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init(g *gin.Engine) *gin.Engine {
	g.Use(middleware.Cors())

	userGroup(g.Group("/api", middleware.Auth))

	return g
}
