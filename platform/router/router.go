// platform/router/router.go

package router

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/yomequido/quido-platform/platform/authenticator"
	"github.com/yomequido/quido-platform/platform/middleware"
	"github.com/yomequido/quido-platform/platform/models"
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
		AllowMethods:     []string{"GET", "DELETE", "OPTIONS", "PUT", "PATCH"},
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
	router.GET("/user", middleware.IsAuthenticated, user.Handler)
	router.GET("/chat", middleware.IsAuthenticated, chat.Handler)
	router.GET("/logout", middleware.IsAuthenticated, logout.Handler)
	router.GET("/ws", middleware.IsAuthenticated, websocket.Handler)
	router.GET("/paymentMethods", middleware.IsAuthenticated, paymentMethods.Handler)

	v1 := router.Group("/v1")

	//Authentication login and callback endpoints
	v1.GET("login", login.Handler(auth))
	v1.GET("/callback", callback.Handler(auth))

	//Get a users data and update users data
	v1.GET("/user", middleware.IsAuthenticated, user.Get)
	v1.POST("/user", middleware.IsAuthenticated, user.Post)

	//Create a checkout id and public key for creating a card tokenizer
	v1.GET("/checkout", middleware.IsAuthenticated, checkout.Handler)

	v1.GET("/paymentMethods", middleware.IsAuthenticated, paymentMethods.Handler)

	//to-do
	/*
		v1.GET("/paymentMethods", func(ctx *gin.Context) {

			ctx.JSON(http.StatusOK, gin.H{
				"card_payment_methods": []models.CardPaymentMethod{
					{Type: "card", CardEnding: 1432, CardToken: "test_3ed98d239dn9238", Default: true},
					{Type: "card", CardEnding: 4352, CardToken: "test_3edr32r32432432", Default: false},
					{Type: "card", CardEnding: 8032, CardToken: "test_3ed98d239dn9258", Default: false},
				},
				"oxxo_payment_method": models.OxxoPaymentMethod{Type: "oxxo", Reference: "0000-0000-0000-0000", BarcodeUrl: "test.net"},
				"spei_payment_method": models.SpeiPaymentMethod{Type: "spei", Reference: "16537213202193820183"},
			},
			)
		})
	*/

	v1.POST("/paymentMethods", func(ctx *gin.Context) {
		var card models.CardPaymentMethod
		err := ctx.BindJSON(&card)
		if err != nil {
			log.Panic(err)
		}

		ctx.Status(http.StatusCreated)
	})

	return router
}
