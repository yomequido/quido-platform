// web/app/user/user.go

package user

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yomequido/quido-platform/platform/database"
)

// Handler for our logged-in user page.
func Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	database.GetTest()

	ctx.HTML(http.StatusOK, "user.html", profile)
}
