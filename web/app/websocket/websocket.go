package websocket

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yomequido/quido-platform/platform/database"
	"github.com/yomequido/quido-platform/platform/models"
	"github.com/yomequido/quido-platform/platform/tools"
)

var upgrader = websocket.Upgrader{
	//check origin will check the cross region source (note : please not using in production)
	CheckOrigin: func(r *http.Request) bool {
		//Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
		return true
	},
}

// Handler for websockets
func Handler(ctx *gin.Context) {
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Panic(err)
	}

	defer ws.Close()

	for {
		//Read Message from client
		messageType, message, err := ws.ReadMessage()
		log.Print(string(message))
		if err != nil {
			log.Panic(err)
		}

		profile := tools.GetProfile(ctx)

		var newMessage models.Message

		newMessage.SentByUser = true
		newMessage.Channel = "livechat"
		newMessage.Message = sql.NullString{String: string(message), Valid: true}
		newMessage.SentDate = sql.NullTime{Time: time.Now(), Valid: true}

		database.InsertUserMessage(profile.Sub, message)

		//If client message is ping will return pong
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//Response message to client
		err = ws.WriteMessage(messageType, message)
		if err != nil {
			log.Panic(err)
		}
	}
}
