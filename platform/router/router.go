// platform/router/router.go

package router

import (
	"encoding/gob"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/yomequido/quido-platform/platform/authenticator"
	"github.com/yomequido/quido-platform/platform/middleware"
	"github.com/yomequido/quido-platform/web/app/callback"
	"github.com/yomequido/quido-platform/web/app/chat"
	"github.com/yomequido/quido-platform/web/app/login"
	"github.com/yomequido/quido-platform/web/app/logout"
	"github.com/yomequido/quido-platform/web/app/paymentMethods"
	"github.com/yomequido/quido-platform/web/app/user"
	"github.com/yomequido/quido-platform/web/app/websocket"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	router.Static("/public", "web/static")
	router.LoadHTMLGlob("web/template/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/user", middleware.IsAuthenticated, user.Handler)
	router.GET("/chat", middleware.IsAuthenticated, chat.Handler)
	router.GET("/logout", middleware.IsAuthenticated, logout.Handler)
	router.GET("/ws", middleware.IsAuthenticated, websocket.Handler)
	router.GET("/paymentMethods", middleware.IsAuthenticated, paymentMethods.Handler)

	return router
}
