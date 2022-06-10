package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/yomequido/quido-platform/platform/models"
	"github.com/yomequido/quido-platform/platform/tools"
)

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("Warning: %s environment variable not set.\n", k)
	}
	return v
}

func initSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_postgres_databasesql_create_socket]
	var (
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'my-db-user'
		dbPwd                  = mustGetenv("DB_PASS")                  // e.g. 'my-db-password'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", dbUser, dbPwd, dbName, socketDir, instanceConnectionName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_postgres_databasesql_create_socket]
}

func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_postgres_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_postgres_databasesql_limit]

	// [START cloud_sql_postgres_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800 * time.Second)

	// [END cloud_sql_postgres_databasesql_lifetime]
}

func GetTest() {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

	var id int
	err = db.QueryRow("select patient_id from patients where patient_id=$1", 1).Scan(&id)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("ID: %d", id)
}

func GetUserMessages(authId string) []models.Message {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

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
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

	message_id := -1

	sqlStatement := `INSERT INTO messages (sent_by_user, channel, fk_patient, message, sentDate) VALUES ($1, $2, (SELECT patient_id FROM patients WHERE $3 = ANY(authId)), $4, $5) RETURNING message_id;`

	err = db.QueryRow(sqlStatement, message.SentByUser, message.Channel, authId, message.Message, message.SentDate).Scan(&message_id)
	if err != nil {
		log.Panic(err)
	}

	return message_id != -1
}

//to-do remove test function
func InserNewUser(profile models.Profile) {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

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
			err = db.QueryRow(sqlStatement, authId, profile.GivenName, profile.FamilyName, profile.Email).Scan(&givenName)
			if err != nil {
				log.Panic(err)
			}
		} else {
			log.Print("Patient with same email has two user: " + profile.Email)
		}
	}
}

func GetConektaPayments(profile models.Profile) string {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

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

	rows.Close()

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
