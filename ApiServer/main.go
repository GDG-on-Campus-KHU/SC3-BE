package main

import (
	cfg "apiServer/config"
	"apiServer/db"
	"apiServer/server"
	"context"
)

func main() {
	cfg.Init()

	ctx := context.Background()
	ctx, cancel := db.New(ctx)
	defer func() {
		db.DisconnectDB()
		cancel()
	}()

	server.New()
	server.Server.Router.Run(":" + cfg.AppConfig.Srv.Port)
}
