package main

import (
	"GoGin-API-CuentasClaras/api"
	"GoGin-API-CuentasClaras/config"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	config.InitLog()
}

func main() {
	port := os.Getenv("PORT")

	init := config.Init()
	app := api.Init(init)

	app.Run(":" + port)
}
