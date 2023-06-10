package controllers

import (
	"log"
	"net/http"

	"github.com/Grandbusta/apex-dice/config"
	"github.com/Grandbusta/apex-dice/models"
	"github.com/Grandbusta/apex-dice/utils"
	"github.com/gin-gonic/gin"
)

func FundWallet(ctx *gin.Context) {
	db := config.NewDB()
	user := models.User{}
	user.Email = "olaifabolu@gmail.com"
	err := user.AddToWallet(db, config.DEFAULT_FUND_AMOUNT)
	if err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithMessage(ctx, http.StatusOK, "Wallet funded")
}

func GetWalletBalance(ctx *gin.Context) {
	db := config.NewDB()
	user := models.User{}
	user.Email = "olaifabolu@gmail.com"
	publicWallet, err := user.GetWalletBalance(db)
	if err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithData(ctx, http.StatusOK, publicWallet)
}

func StartGame(ctx *gin.Context) {
	db := config.NewDB()
	user := models.User{}
	user.Email = "olaifabolu@gmail.com"
	if err := user.DeductWallet(db, config.COMMITMENT_COST); err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	game := models.Game{}
	game.Cost = config.COMMITMENT_COST
	game.Target = utils.GenerateTargetNumber()
	game.UserID = 1
	_, err := game.SaveGame(db)
	if err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	utils.SuccessWithMessage(ctx, http.StatusOK, "Game started")
}
