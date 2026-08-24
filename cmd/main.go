package main

import (
	"apg105/di"
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	if err := fx.ValidateApp(di.CommonModule); err != nil {
		panic(err)
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env: %v", err)
	}
	app := fx.New(di.CommonModule)
	app.Run()
}
