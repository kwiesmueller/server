package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/GruffDebate/server/api"
	"github.com/GruffDebate/server/config"
)

func main() {

	config.Init()
	api.RW_DB_POOL = config.InitDB()
	api.RW_DB_POOL.LogMode(true)

	root := api.SetUpRouter(false, api.RW_DB_POOL)
	addr := ":" + os.Getenv("PORT")

	go func() {
		if err := root.Start(addr); err != nil {
			root.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := root.Shutdown(ctx); err != nil {
		root.Logger.Fatal(err)
	}
}
