package main

import (
	"github.com/damirqa/shortener/internal/app"
)

func main() {
	a := &app.App{}
	a.Init()
	go a.Start()
	a.Listen()
	a.Stop()
}
