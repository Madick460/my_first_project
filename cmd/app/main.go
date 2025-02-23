package main

import (
	"context"
	"github.com/Rasikrr/my_project/internal/app"
	"log"
)

const (
	appName = "my_project"
)

func main() {
	ctx := context.Background()
	a := app.InitApp(ctx, appName)
	if err := a.Start(ctx); err != nil {
		log.Fatal(err)
	}
}
