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

//to-do remove test function
func InsertUserMessage(authId string, message models.Message) bool {
	//Cretae DB connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	message_id := ""

	sqlStatement := `INSERT INTO messages (sent_by_user, channel, fk_patient, message, sentDate) VALUES ($1, $2, (SELECT patient_id FROM patients WHERE $3 = ANY(authId)), $4, $5) RETURNING message_id;`

	err = db.QueryRow(sqlStatement, message.SentByUser, message.Channel, authId, message.Message, message.SentDate, message.ExternalMessagesId).Scan(&message_id)
	if err != nil {
		log.Panic(err)
	}

	return message_id != ""
}
