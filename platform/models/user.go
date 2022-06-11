package models

import "time"

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

type InsertUser struct {
	Phone         string `json:"phone,omitempty"`
	CountryCode   string `json:"country_code,omitempty"`
	GivenName     string `json:"given_names,omitempty"`
	FamilyName    string `json:"family_names,omitempty"`
	Birthdate     string `json:"birthdate,omitempty"`
	GovernmentId  string `json:"government_id,omitempty"`
	TaxId         string `json:"tax_id,omitempty"`
	BirthSex      string `json:"birth_sex,omitempty"`
	Gender        string `json:"gender,omitempty"`
	FullLegalName string `json:"full_legal_name,omitempty"`
}
