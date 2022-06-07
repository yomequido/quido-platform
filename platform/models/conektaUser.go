package models

import "database/sql"

type ConektaUser struct {
	PatientId   int            `db:"patient_id"`
	GivenName   sql.NullString `db:"givennames"`
	FamilyName  sql.NullString `db:"familynames"`
	Email       string         `db:"email"`
	CountryCode string         `db:"countrycode"`
	Phone       string         `db:"phone"`
	ConektaId   sql.NullString `db:"conekta_id"`
}
