package main

import (
	"log"
	"os"

	"github.com/v4run/ws/internal/app"
)

func main() {
	a := app.InitApp(app.AppState{
		ShowRelativePath:     true,
		SwitchToWorktreeRoot: false,
	})
	switch err := a.Run(); err {
	case app.ErrCancelled:
		os.Exit(1)
	case nil:
	default:
		log.Fatal(err)
	}
}
