package main

import (
	"github.com/damirqa/shortener/internal/app"
)

func main() {
	a := &app.App{}
	a.Init()
	a.Start()
}
