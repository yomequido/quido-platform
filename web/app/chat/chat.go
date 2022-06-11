package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yomequido/quido-platform/platform/database"
	"github.com/yomequido/quido-platform/platform/tools"
)

// Handler for the chat page
func Handler(ctx *gin.Context) {

	ctx.HTML(http.StatusOK, "chat.html", gin.H{
		"historicMessages": fetchMessages(ctx),
	})
	/*
		ctx.JSON(http.StatusOK, gin.H{
			"historicMessages": fetchMessages(ctx),
		}) */
}

func fetchMessages(ctx *gin.Context) []string {
	profile := tools.GetProfile(ctx)

	messages := database.GetUserMessages(profile.Sub)

	var result []string

	for _, message := range messages {
		result = append(result, message.Message.String)
	}

	return result
}
