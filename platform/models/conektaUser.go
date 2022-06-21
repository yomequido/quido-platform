package models

import "github.com/jackc/pgtype"

type ConektaUser struct {
	PatientId   int            `db:"patient_id"`
	GivenName   string         `db:"givennames"`
	FamilyName  string         `db:"familynames"`
	Email       string         `db:"email"`
	CountryCode string         `db:"countrycode"`
	Phone       string         `db:"phone"`
	ConektaId   pgtype.Varchar `db:"conekta_id"`
}
