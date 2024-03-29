// web/app/callback/callback.go

package callback

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/yomequido/quido-platform/platform/authenticator"
	"github.com/yomequido/quido-platform/platform/database"
	"github.com/yomequido/quido-platform/platform/tools"
)

// Handler for our callback.
func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		if ctx.Query("state") != session.Get("state") {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
		if err != nil {
			ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		log.Print("Inserting user into DB")

		//Insert new user into DB if doesn't exist
		userProfile := tools.GetProfile(ctx)
		//Insert new user returns true if user previously existed in the DB
		userExists := strconv.FormatBool(database.InserNewUser(userProfile))

		ctx.SetCookie("auth-session", token.AccessToken, -1, "/", ".quido.mx", true, false)

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "https://quido.mx/callbackLogin?returning_user="+userExists)
	}
}
