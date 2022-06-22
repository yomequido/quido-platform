package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgtype"
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

func InserNewUser(profile models.Profile) bool {
	var userExists = true
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

	//Check if user exists through auth0 ID
	rows, err := db.Query(`SELECT patient_id FROM patients WHERE $1 =ANY(auth_id)`, profile.Sub)
	if err != nil {
		log.Panic(err)
	}
	//User doesn't exist by auth0 ID we should check if using other login method through email matching
	if !rows.Next() {
		//check if email exists
		rows, err = db.Query(`SELECT patient_id FROM patients WHERE $1 = email`, profile.Email)
		if err != nil {
			log.Panic(err)
		}
		//if email exists insert user auth0 into existing user
		if rows.Next() {
			var patientId = ""
			err := rows.Scan(&patientId)
			if err != nil {
				log.Panic(err)
			}
			sqlStatement := `UPDATE patients SET auth_id = array_append(auth_id, $1) WHERE $2 = patient_id ;`
			//authId is a string array
			authId := []string{profile.Sub}
			_, err = db.Exec(sqlStatement, authId, patientId)
			if err != nil {
				log.Panic(err)
			}
		} else {
			userExists = false
			//if email doesn't exist then insert user
			sqlStatement := `INSERT INTO patients (auth_id, email) VALUES ($1, $2) RETURNING email;`
			//authId is a string array
			authId := []string{profile.Sub}
			email := ""
			err = db.QueryRow(sqlStatement, authId, profile.Email).Scan(&email)
			if err != nil {
				log.Panic(err)
			}
		}
	}
	return userExists
}

func GetUser(authId string) models.DBUser {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

	rows, err := db.Query(`SELECT created_date, given_names, family_names, email, country_code, phone, CAST(birth_sex as VARCHAR(1)) as birth_sex, gender, birthdate, tax_id, government_id, full_legal_name FROM patients WHERE $1 =ANY(auth_id)`, authId)
	if err != nil {
		log.Panic(err)
	}

	var user models.DBUser

	if rows.Next() {
		err = rows.Scan(
			&user.CreatedDate,
			&user.GivenName,
			&user.FamilyName,
			&user.Email,
			&user.CountryCode,
			&user.Phone,
			&user.BirthSex,
			&user.Gender,
			&user.Birthdate,
			&user.TaxId,
			&user.GovernmentId,
			&user.FullLegalName,
		)
	}

	if err != nil {
		log.Panic(err)
	}

	return user

}

func SetUser(authId string, user models.User) {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Getting existing user data from db: " + authId)
	rows, err := db.Query(`SELECT patient_id, given_names, family_names, country_code, phone, CAST(birth_sex as VARCHAR(1)) as birth_sex, gender, birthdate, tax_id, government_id, full_legal_name FROM patients WHERE $1 =ANY(auth_id)`, authId)
	if err != nil {
		log.Panic(err)
	}

	var id int
	var currentUser models.DBUser
	log.Println("Extracting existing user data from resulting row: " + authId)
	if rows.Next() {
		err = rows.Scan(
			&id,
			&currentUser.GivenName,
			&currentUser.FamilyName,
			&currentUser.CountryCode,
			&currentUser.Phone,
			&currentUser.BirthSex,
			&currentUser.Gender,
			&currentUser.Birthdate,
			&currentUser.TaxId,
			&currentUser.GovernmentId,
			&currentUser.FullLegalName,
		)
	}

	//Indicator that there are changes betweeen the old values and the new values, so as to not waste an update in postgres for data that won't change
	updateUser := false

	//Add columns to update only if they have values
	log.Println("Given name: " + user.GivenName)
	if user.GivenName != "" {
		currentUser.GivenName.Set(user.GivenName)
		currentUser.GivenName.Status = pgtype.Present
		updateUser = true
	}
	log.Println("Family name: " + user.FamilyName)
	if user.FamilyName != "" {
		currentUser.FamilyName.Set(user.FamilyName)
		updateUser = true
	}
	log.Println("Country: " + user.CountryCode)
	if user.CountryCode != "" {
		currentUser.CountryCode.Set(user.CountryCode)
		updateUser = true
	}
	if user.Phone != "" {
		currentUser.Phone.Set(user.Phone)
		updateUser = true
	}
	if user.BirthSex != "" {
		currentUser.BirthSex.Set(user.BirthSex)
		updateUser = true
	}
	if user.Gender != "" {
		currentUser.Gender.Set(user.Gender)
		updateUser = true
	}
	if user.Birthdate != "" {
		currentUser.Birthdate.Set(user.Birthdate)
		updateUser = true
	}
	if user.TaxId != "" {
		currentUser.TaxId.Set(user.TaxId)
		updateUser = true
	}
	if user.GovernmentId != "" {
		currentUser.GovernmentId.Set(user.GovernmentId)
		updateUser = true
	}
	if user.FullLegalName != "" {
		currentUser.FullLegalName.Set(user.FullLegalName)
		updateUser = true
	}

	if updateUser {
		log.Println("Starting user update: " + authId)
		result, err := db.Exec(`UPDATE patients SET given_names = $1, family_names = $2, country_code = $3, phone =$4, birth_sex = $5, gender = $6, birthdate = $7, tax_id = $8, government_id = $9, full_legal_name = $10 WHERE patient_id = $11;`,
			currentUser.GivenName,
			currentUser.FamilyName,
			currentUser.CountryCode,
			currentUser.Phone,
			currentUser.BirthSex,
			currentUser.Gender,
			currentUser.Birthdate,
			currentUser.TaxId,
			currentUser.GovernmentId,
			currentUser.FullLegalName,
			id)
		if err != nil {
			log.Panic(err)
		}

		res, _ := result.RowsAffected()
		log.Printf("rows affected from %d", res)
	} else {
		log.Printf("User %s didn't have any new values to update", authId)
	}
}

func GetUserMessages(authId string) []models.Message {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

	rows, err := db.Query(`SELECT * FROM messages WHERE fk_patient = (SELECT patient_id FROM patients WHERE $1 =ANY(auth_id))`, authId)
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

	sqlStatement := `INSERT INTO messages (sent_by_user, channel, fk_patient, message, sent_date) VALUES ($1, $2, (SELECT patient_id FROM patients WHERE $3 = ANY(auth_id)), $4, $5) RETURNING message_id;`

	err = db.QueryRow(sqlStatement, message.SentByUser, message.Channel, authId, message.Message, message.SentDate).Scan(&message_id)
	if err != nil {
		log.Panic(err)
	}

	return message_id != -1
}

func GetConektaPayments(profile models.Profile) models.PaymentMethods {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

	rows, err := db.Query(`SELECT patient_id, given_names, family_names, email, country_code, phone, conekta_id FROM patients LEFT JOIN conekta_id ON patient_id = fk_patient WHERE $1 = ANY(auth_id)`, profile.Sub)
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

	if conektaUser.ConektaId.Status == pgtype.Null {
		conektaid := tools.CreateCustomer(conektaUser)

		sqlStatement := `INSERT INTO conekta_id (fk_patient, conekta_id) VALUES ($1, $2) RETURNING conekta_id;`

		err = db.QueryRow(sqlStatement, conektaUser.PatientId, conektaid.ID).Scan(&conektaUser.ConektaId)
		if err != nil {
			log.Panic(err)
		}

	}

	return tools.GetConektaPaymentMethods(conektaUser.ConektaId.String)
}

func GetConektaUser(profile models.Profile) string {
	db, err := initSocketConnectionPool()
	if err != nil {
		log.Panic(err)
	}

	rows, err := db.Query(`SELECT conekta_id FROM patients LEFT JOIN conekta_id ON patient_id = fk_patient WHERE $1 = ANY(auth_id)`, profile.Sub)
	if err != nil {
		log.Panic(err)
	}

	var conektaUser models.ConektaUser
	if rows.Next() {
		err = rows.Scan(
			&conektaUser.ConektaId)
		if err != nil {
			log.Panic(err)
		}
	}

	rows.Close()

	return conektaUser.ConektaId.String

}
