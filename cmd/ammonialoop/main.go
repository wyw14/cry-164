package main

import (
	"context"
	appRuntime "github.com/wyw14/cry-164/internal/runtime"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := appRuntime.LoadConfig()
	app := appRuntime.New(config)
	if err := appRuntime.Probe(app.Server); err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := app.Start(ctx); err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := app.Server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			log.Print(err)
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = app.Shutdown(shutdown)
}
