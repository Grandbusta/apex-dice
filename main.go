package main

import (
	"fmt"
	"log"

	"github.com/Grandbusta/apex-dice/config"
	"github.com/Grandbusta/apex-dice/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		fmt.Println(".env loaded")
	}
}

func main() {
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	db := config.NewDB()
	db.Debug().AutoMigrate()

	r.Run(":8080")

}
