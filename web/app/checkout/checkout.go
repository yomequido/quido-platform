package checkout

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yomequido/quido-platform/platform/tools"
)

func Handler(ctx *gin.Context) {
	checkout_id, public_key := tools.CreateCheckout()
	ctx.JSON(http.StatusCreated, gin.H{
		"checkout_id": checkout_id,
		"public_key":  public_key,
	})
}
