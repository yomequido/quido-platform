package models

import (
	"time"

	"github.com/jackc/pgtype"
)

type DBUser struct {
	Email         pgtype.Varchar   `json:"email,omitempty"`
	Phone         pgtype.Varchar   `json:"phone,omitempty"`
	CountryCode   pgtype.Varchar   `json:"country_code,omitempty"`
	GivenName     pgtype.Varchar   `json:"given_names,omitempty"`
	FamilyName    pgtype.Varchar   `json:"family_names,omitempty"`
	Birthdate     pgtype.Date      `json:"birthdate,omitempty"`
	GovernmentId  pgtype.Varchar   `json:"government_id,omitempty"`
	TaxId         pgtype.Varchar   `json:"tax_id,omitempty"`
	BirthSex      pgtype.Varchar   `json:"birth_sex,omitempty"`
	Gender        pgtype.Varchar   `json:"gender,omitempty"`
	CreatedDate   pgtype.Timestamp `json:"created_date,omitempty"`
	FullLegalName pgtype.Varchar   `json:"full_legal_name,omitempty"`
}

type User struct {
	Email         string    `json:"email,omitempty"`
	Phone         string    `json:"phone,omitempty"`
	CountryCode   string    `json:"country_code,omitempty"`
	GivenName     string    `json:"given_names,omitempty"`
	FamilyName    string    `json:"family_names,omitempty"`
	Birthdate     string    `json:"birthdate,omitempty"`
	GovernmentId  string    `json:"government_id,omitempty"`
	TaxId         string    `json:"tax_id,omitempty"`
	BirthSex      string    `json:"birth_sex,omitempty"`
	Gender        string    `json:"gender,omitempty"`
	CreatedDate   time.Time `json:"created_date,omitempty"`
	FullLegalName string    `json:"full_legal_name,omitempty"`
}
