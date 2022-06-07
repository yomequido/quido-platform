package models

type Profile struct {
	Sub        string `json:"sub"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}
