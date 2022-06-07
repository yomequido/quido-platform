package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/yomequido/quido-platform/platform/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "termis"
	password = ""
	dbname   = "postgres"
)

func GetUserMessages(authId string) []models.Message {
	//Cretae DB connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM messages WHERE fk_patient = (SELECT patient_id FROM patients WHERE $1 =ANY(authId))`, authId)
	if err != nil {
		log.Panic(err)
	}

	var messages []models.Message

	for rows.Next() {
		var message models.Message
		err = rows.Scan(
			&message.MessageId,
			&message.ExternalMessagesId,
			&message.SentByUser,
			&message.Channel,
			&message.ChannelUserId,
			&message.FkEmployee,
			&message.FkPatient,
			&message.Title,
			&message.Message,
			&message.UrlReference,
			&message.SentDate,
			&message.ReadDate)
		if err != nil {
			log.Panic(err)
		}

		messages = append(messages, message)
	}

	return messages
}
