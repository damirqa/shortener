package main

import (
	"context"
	"github.com/damirqa/shortener/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	a := &app.App{}
	a.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go a.Start()

	<-ctx.Done()

	a.Shutdown()
}
