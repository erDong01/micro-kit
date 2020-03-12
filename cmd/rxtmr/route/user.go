package route

import (
	"rxt/cmd/rxtmr/app/handler/user"

	"github.com/gin-gonic/gin"
)

func userGroup(group *gin.RouterGroup) {
	g := group.Group("/auth")
	{
		g.GET("/login", user.Login)
	}
}
