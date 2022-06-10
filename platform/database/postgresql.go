package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
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

func getDbConnection() *pgx.Conn {
	//conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	conn, err := pgx.Connect(context.Background(), "postgres://termis:@localhost:5432/postgres")
	if err != nil {
		log.Panic(err)
	}
	return conn
}

func GetTest() {
	conn := getDbConnection()
	defer conn.Close(context.Background())

	var id int

	err := conn.QueryRow(context.Background(), "select patient_id from patients where patient_id=$1", 1).Scan(&id)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("ID: %d", id)
}

func GetUserMessages(authId string) []models.Message {
	//Cretae DB connection
	db := getDbConnection()
	background := context.Background()
	defer db.Close(background)

	rows, err := db.Query(background, `SELECT * FROM messages WHERE fk_patient = (SELECT patient_id FROM patients WHERE $1 =ANY(authId))`, authId)
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
	db := getDbConnection()
	background := context.Background()
	defer db.Close(background)

	message_id := -1

	sqlStatement := `INSERT INTO messages (sent_by_user, channel, fk_patient, message, sentDate) VALUES ($1, $2, (SELECT patient_id FROM patients WHERE $3 = ANY(authId)), $4, $5) RETURNING message_id;`

	err := db.QueryRow(background, sqlStatement, message.SentByUser, message.Channel, authId, message.Message, message.SentDate).Scan(&message_id)
	if err != nil {
		log.Panic(err)
	}

	return message_id != -1
}

//to-do remove test function
func InserNewUser(profile models.Profile) {
	//Cretae DB connection
	db := getDbConnection()
	background := context.Background()
	defer db.Close(background)

	rows, err := db.Query(background, `SELECT patient_id FROM patients WHERE $1 =ANY(authId)`, profile.Sub)
	if err != nil {
		log.Panic(err)
	}

	if !rows.Next() {
		rows, err = db.Query(background, `SELECT patient_id FROM patients WHERE $1 = email`, profile.Email)
		if err != nil {
			log.Panic(err)
		}
		if !rows.Next() {
			givenName := ""

			sqlStatement := `INSERT INTO patients (authId, givenNames, familyNames, email) VALUES ($1, $2, $3, $4) RETURNING givenNames;`
			//authId is a string array
			authId := []string{profile.Sub}
			err = db.QueryRow(background, sqlStatement, pq.Array(authId), profile.GivenName, profile.FamilyName, profile.Email).Scan(&givenName)
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
	db := getDbConnection()
	background := context.Background()
	defer db.Close(background)

	rows, err := db.Query(background, `SELECT patient_id, givenNames, familyNames, email, countryCode, phone, conekta_id FROM patients LEFT JOIN conekta_id USING (patient_id) WHERE $1 = ANY(authId)`, profile.Sub)
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

	rows.Close()

	if !conektaUser.ConektaId.Valid {
		conektaid := tools.CreateCustomer(conektaUser)

		sqlStatement := `INSERT INTO conekta_id (patient_id, conekta_id) VALUES ($1, $2) RETURNING conekta_id;`

		err = db.QueryRow(background, sqlStatement, conektaUser.PatientId, conektaid).Scan(&conektaUser.ConektaId)
		if err != nil {
			log.Panic(err)
		}

	}

	return conektaUser.ConektaId.String
}
