package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	log.Print(ctx.Request.URL.Host)
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
