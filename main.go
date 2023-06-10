package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/Grandbusta/apex-dice/config"
	"github.com/Grandbusta/apex-dice/controllers"
	"github.com/Grandbusta/apex-dice/middlewares"
	"github.com/Grandbusta/apex-dice/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
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
	if err := db.Debug().AutoMigrate(&models.User{}, &models.Game{}); err == nil && db.Migrator().HasTable(&models.User{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			//Insert seed data
			user := models.User{Email: "olaifabolu@gmail.com"}
			user.SaveUser(db)
		}
	}

	r.GET("/wallet/fund", controllers.FundWallet)
	r.GET("/wallet/get-balance", controllers.GetWalletBalance)

	r.GET("/game/start", controllers.StartGame)

	r.Run(":8080")

}
