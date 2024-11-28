package main

import (
	"context"
	cfg "dbServer/config"
	"dbServer/db"
	"dbServer/server"
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
