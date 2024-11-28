package server

import (
	"dbServer/server/handlers"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	Router *gin.Engine
}

var (
	Server *GinServer = &GinServer{}
)

func New() {
	Server.Router = gin.Default()

	Server.Router.GET("/api/disaster/list", handlers.GetDisasterList)
}
