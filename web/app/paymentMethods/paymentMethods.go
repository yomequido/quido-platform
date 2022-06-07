package paymentMethods

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yomequido/quido-platform/platform/database"
	"github.com/yomequido/quido-platform/platform/tools"
)

// Handler for user payment methods.
func Handler(ctx *gin.Context) {
	profile := tools.GetProfile(ctx)
	conektaId := database.GetConektaPayments(profile)
	ctx.JSON(http.StatusOK, gin.H{"ConektaId": conektaId})

}
