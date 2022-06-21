// web/app/user/user.go

package user

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yomequido/quido-platform/platform/database"
	"github.com/yomequido/quido-platform/platform/models"
	"github.com/yomequido/quido-platform/platform/tools"
)

// Handler for our logged-in user page.
func Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	ctx.HTML(http.StatusOK, "user.html", profile)
}

// Handler for getting the logged in user
func Get(ctx *gin.Context) {
	sub := tools.GetProfile(ctx).Sub

	user := database.GetUser(sub)

	ctx.JSON(http.StatusOK, user)

}

// Handler for updating the logged in user
func Post(ctx *gin.Context) {
	sub := tools.GetProfile(ctx).Sub
	var user models.User

	err := ctx.BindJSON(&user)
	if err != nil {
		log.Panic(err)
	}

	database.SetUser(sub, user)

	ctx.Status(http.StatusOK)
}
