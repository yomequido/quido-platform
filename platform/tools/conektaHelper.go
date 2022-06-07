package tools

import (
	"log"

	conekta "github.com/conekta/conekta-go"
	"github.com/conekta/conekta-go/customer"
	"github.com/yomequido/quido-platform/platform/models"
)

func CreateCustomer(conektaUser models.ConektaUser) string {
	conekta.APIKey = "key_zyaxzY5JAjNAGTv8f8TroA"

	payment := &conekta.PaymentSourceCreateParams{
		PaymentType: "oxxo_recurrent",
	}

	cus := &conekta.CustomerParams{}
	cus.Name = conektaUser.GivenName.String + " " + conektaUser.FamilyName.String
	cus.Email = conektaUser.Email
	cus.Phone = conektaUser.CountryCode + conektaUser.Phone
	cus.PaymentSources = append(cus.PaymentSources, payment)

	res, err := customer.Create(cus)
	if err != nil {
		log.Panic(err)
	}

	return res.ID
}
