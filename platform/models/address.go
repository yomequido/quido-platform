package models

type Address struct {
	ID             int    `json:"id"`
	FkPatient      int    `db:"fk_patient"`
	Street         string `json:"street"`
	InteriorNumber string `json:"interior_number"`
	ExteriorNumber string `json:"exterior_number"`
	Neighborhood   string `json:"neighborhood"`
	PostalCode     string `json:"postal_code"`
	Country        string `json:"country"`
	Reference      string `json:"reference"`
}
