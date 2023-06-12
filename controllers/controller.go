package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Grandbusta/apex-dice/config"
	"github.com/Grandbusta/apex-dice/models"
	"github.com/Grandbusta/apex-dice/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FundWallet(ctx *gin.Context) {
	db := config.NewDB()
	user := models.User{}
	user.Email = "olaifabolu@gmail.com"
	currentUser, err := user.GetWalletBalance(db)
	walletBalance, _ := strconv.Atoi(currentUser.WalletBalance)
	if walletBalance > 35 {
		utils.SuccessWithMessage(ctx, http.StatusBadRequest, "Balance still greater than 35")
	}
	err = user.AddToWallet(db, config.DEFAULT_FUND_AMOUNT)
	tlog := models.TransactionLog{}
	tlog.Action = "credit"
	tlog.UserID = 1
	tlog.Amount = config.DEFAULT_FUND_AMOUNT
	tlog.Add(db)
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
	tlog := models.TransactionLog{}
	tlog.Action = "debit"
	tlog.UserID = 1
	tlog.Amount = config.COMMITMENT_COST
	tlog.Add(db)
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

func RollDice(ctx *gin.Context) {
	db := config.NewDB()
	game := models.Game{}
	tlog := models.TransactionLog{}
	game.UserID = 1
	user := models.User{}
	user.Email = "olaifabolu@gmail.com"
	activeGame, err := game.GetActiveGame(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusNotFound, "No active game")
		return
	}
	if err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}

	if activeGame.FirstRoll != 0 && activeGame.SecondRoll != 0 {
		utils.SuccessWithMessage(ctx, http.StatusOK, "Session complete, you can try again")
		return
	}
	// if first is 0, insert first and deduct 5 from wallet
	firstRoll := false
	if activeGame.FirstRoll == 0 {
		firstRoll = true
		activeGame.FirstRoll = utils.GenerateDiceNumber()
		user.DeductWallet(db, 5)
		tlog.Action = "debit"
		tlog.UserID = 1
		tlog.Amount = 5
		tlog.Add(db)
	}
	// if not, and second is zero, insert second
	if !firstRoll {
		activeGame.SecondRoll = utils.GenerateDiceNumber()
	}
	updatedGame, err := activeGame.UpdateGame(db)
	if err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}

	if firstRoll {
		utils.SuccessWithMessage(
			ctx,
			http.StatusOK,
			fmt.Sprintf("You just rolled %v. Roll the second time", activeGame.FirstRoll),
		)
		return
	}

	if !firstRoll && updatedGame.FirstRoll+updatedGame.SecondRoll == updatedGame.Target {
		user.AddToWallet(db, 10)
		tlog.Action = "credit"
		tlog.UserID = 1
		tlog.Amount = 10
		tlog.Add(db)
		utils.SuccessWithMessage(ctx, http.StatusOK, "You won. Awarded 10 sats")
		return
	}

	utils.SuccessWithMessage(ctx, http.StatusOK, "You did not win, please try again")

}

func EndGame(ctx *gin.Context) {
	db := config.NewDB()
	game := models.Game{}
	game.UserID = 1
	activeGame, err := game.GetActiveGame(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusNotFound, "No active game")
		return
	}
	if err != nil {
		log.Println(err)
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}
	err = activeGame.EndSession(db)
	if err != nil {
		utils.ServerResponse(ctx, http.StatusInternalServerError, "An error occured")
		return
	}

	utils.SuccessWithMessage(ctx, http.StatusOK, "Game ended")
}
