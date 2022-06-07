package tools

import (
	"encoding/json"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/yomequido/quido-platform/platform/models"
)

func GetProfile(ctx *gin.Context) models.Profile {
	profileInterface := sessions.Default(ctx).Get("profile")

	jsonString, err := json.Marshal(profileInterface)
	if err != nil {
		log.Panic(err)
	}

	var profile models.Profile

	json.Unmarshal(jsonString, &profile)

	return profile
}
