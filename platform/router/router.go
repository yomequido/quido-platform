// platform/router/router.go

package router

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/yomequido/quido-platform/platform/authenticator"
	"github.com/yomequido/quido-platform/platform/middleware"
	"github.com/yomequido/quido-platform/web/app/callback"
	"github.com/yomequido/quido-platform/web/app/chat"
	"github.com/yomequido/quido-platform/web/app/checkout"
	"github.com/yomequido/quido-platform/web/app/login"
	"github.com/yomequido/quido-platform/web/app/logout"
	"github.com/yomequido/quido-platform/web/app/paymentMethods"
	"github.com/yomequido/quido-platform/web/app/user"
	"github.com/yomequido/quido-platform/web/app/websocket"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://api.quido.mx", "https://www.api.quido.mx", "https://quido.mx", "https://www.quido.mx"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
	router.GET("/user", user.Get)
	router.GET("/chat", middleware.IsAuthenticated, chat.Handler)
	router.GET("/logout", middleware.IsAuthenticated, logout.Handler)
	router.GET("/ws", middleware.IsAuthenticated, websocket.Handler)

	v1 := router.Group("/v1")

	//Authentication login and callback endpoints
	v1.GET("/login", login.Handler(auth))
	v1.GET("/callback", callback.Handler(auth))

	//Get a users data and update users data
	v1.GET("/user", middleware.IsAuthenticated, user.Get)
	v1.POST("/user", middleware.IsAuthenticated, user.Post)

	//Create a checkout id and public key for creating a card tokenizer
	v1.GET("/checkout", middleware.IsAuthenticated, checkout.Handler)

	//Get and post payment methods
	v1.GET("/paymentMethods", middleware.IsAuthenticated, paymentMethods.Get)
	v1.POST("/paymentMethods", middleware.IsAuthenticated, paymentMethods.Post)

	//Get and post products and prices
	v1.GET("/products", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	v1.POST("/products", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	// Get and post address

	v1.GET("/address", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	v1.POST("/address", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	//Get and post payment intent
	v1.GET("/paymentIntent", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	v1.POST("/paymentIntent", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})

	return router
}
