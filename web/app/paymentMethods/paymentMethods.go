package paymentMethods

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yomequido/quido-platform/platform/database"
	"github.com/yomequido/quido-platform/platform/tools"
)

// Handler for user payment methods.
func Get(ctx *gin.Context) {
	profile := tools.GetProfile(ctx)
	conektaPaymentMethods := database.GetConektaPayments(profile)
	ctx.JSON(http.StatusOK, conektaPaymentMethods)

}

func Post(ctx *gin.Context) {
	profile := tools.GetProfile(ctx)
	conektaUser := database.GetConektaUser(profile)
	var cardToken CardToken

	err := ctx.BindJSON(&cardToken)
	if err != nil {
		log.Panic(err)
	}
	success := tools.CreateCard(conektaUser, cardToken.Token)
	if success {
		ctx.Status(http.StatusCreated)
	} else {
		ctx.Status(http.StatusBadRequest)
	}

}

type CardToken struct {
	Token string `json:"token"`
}
