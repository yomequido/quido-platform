package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/yomequido/quido-platform/platform/models"
	"github.com/yomequido/quido-platform/platform/tools"
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
			&message.FkPatient,
			&message.FkEmployee,
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

	err = db.QueryRow(sqlStatement, message.SentByUser, message.Channel, authId, message.Message, message.SentDate).Scan(&message_id)
	if err != nil {
		log.Panic(err)
	}

	return message_id != ""
}

//to-do remove test function
func InserNewUser(profile models.Profile) {
	//Cretae DB connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT patient_id FROM patients WHERE $1 =ANY(authId)`, profile.Sub)
	if err != nil {
		log.Panic(err)
	}

	if !rows.Next() {
		rows, err = db.Query(`SELECT patient_id FROM patients WHERE $1 = email`, profile.Email)
		if err != nil {
			log.Panic(err)
		}
		if !rows.Next() {
			givenName := ""

			sqlStatement := `INSERT INTO patients (authId, givenNames, familyNames, email) VALUES ($1, $2, $3, $4) RETURNING givenNames;`
			//authId is a string array
			authId := []string{profile.Sub}
			err = db.QueryRow(sqlStatement, pq.Array(authId), profile.GivenName, profile.FamilyName, profile.Email).Scan(&givenName)
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Print("Patient with same email has two user: " + profile.Email)
		}
	}
}

func GetConektaPayments(profile models.Profile) string {
	//Cretae DB connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT patient_id, givenNames, familyNames, email, countryCode, phone, conekta_id FROM patients LEFT JOIN conekta_id USING (patient_id) WHERE $1 = ANY(authId)`, profile.Sub)
	if err != nil {
		log.Panic(err)
	}

	var conektaUser models.ConektaUser
	if rows.Next() {
		err = rows.Scan(
			&conektaUser.PatientId,
			&conektaUser.GivenName,
			&conektaUser.FamilyName,
			&conektaUser.Email,
			&conektaUser.CountryCode,
			&conektaUser.Phone,
			&conektaUser.ConektaId)
		if err != nil {
			log.Panic(err)
		}
	}

	if !conektaUser.ConektaId.Valid {
		conektaid := tools.CreateCustomer(conektaUser)

		sqlStatement := `INSERT INTO conekta_id (patient_id, conekta_id) VALUES ($1, $2) RETURNING conekta_id;`

		err = db.QueryRow(sqlStatement, conektaUser.PatientId, conektaid).Scan(&conektaUser.ConektaId)
		if err != nil {
			log.Panic(err)
		}

	}

	return conektaUser.ConektaId.String
}
