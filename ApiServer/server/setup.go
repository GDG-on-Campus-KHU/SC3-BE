package server

import (
	"apiServer/server/handlers"

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

	Server.Router.GET("/db/missing/exist", handlers.NewDataExist)
	Server.Router.GET("/db/missing/person/:sn", handlers.SearchBySN)
}
