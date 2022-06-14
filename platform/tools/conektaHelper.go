package tools

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	conekta "github.com/conekta/conekta-go"
	"github.com/conekta/conekta-go/customer"
	"github.com/conekta/conekta-go/paymentsource"
	"github.com/yomequido/quido-platform/platform/models"
)

func CreateCustomer(conektaUser models.ConektaUser) *conekta.Customer {
	conekta.APIKey = os.Getenv("CONEKTA_API")

	cus := &conekta.CustomerParams{}
	cus.Name = conektaUser.GivenName.String + " " + conektaUser.FamilyName.String
	cus.Email = conektaUser.Email
	cus.Phone = conektaUser.CountryCode + conektaUser.Phone

	payment := &conekta.PaymentSourceCreateParams{
		PaymentType: "oxxo_recurrent",
	}

	cus.PaymentSources = append(cus.PaymentSources, payment)

	res, err := customer.Create(cus)
	if err != nil {
		log.Panic(err)
	}

	payment = &conekta.PaymentSourceCreateParams{
		PaymentType: "spei_recurrent",
	}

	paymentsource.Create(res.ID, payment)

	log.Print(res)
	return res
}

func CreateCheckout() (string, string) {
	var urlStr = "https://api.conekta.io/tokens"

	key := os.Getenv("CONEKTA_API")

	bearerToken := base64.StdEncoding.EncodeToString([]byte(key + ":"))

	log.Print(bearerToken)

	var bodyString = `{
		"checkout": {
		"returns_control_on": "Token"
		}
	}`

	body := strings.NewReader(bodyString)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, body)
	r.Header.Add("Authorization", "Basic "+bearerToken)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/vnd.conekta-v2.0.0+json")

	resp, err := client.Do(r)
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()

	var checkoutToken CheckoutToken

	log.Print(resp.Status)

	json.NewDecoder(resp.Body).Decode(&checkoutToken)

	public_key := os.Getenv("CONEKTA_PUBLIC_KEY")

	return checkoutToken.Checkout.ID, public_key
}

type CheckoutToken struct {
	Checkout Checkout `json:"checkout"`
}

type Checkout struct {
	ID string `json:"id"`
}

/*
func GetConektaPayments(conektaId string) map[string]interface{} {

	conekta.APIKey = "key_zyaxzY5JAjNAGTv8f8TroA"

	paymenthMethods, err := paymentsource.All(conektaId)
	if err != nil {
		log.Panic(err)
	}

	log.Print(paymenthMethods)
	var paymentMethodMap models.PaymentMethods

	for _, paymentMethod := range paymenthMethods.Data {
		if paymentMethod.PaymentType == "card" {
			newCard := models.CardPaymentMethod{Type="card", CardEnding=paymentMethod.Last4, paymentMethod.ID, Default=paymentMethod.Default}
			paymentMethodMap.CardPaymentMethods = append(paymentMethodMap.CardPaymentMethods, newCard)
		}else if paymentMethod.PaymentType == "oxxo_recurrent" {
			newOxxo := models.OxxoPaymentMethod{Type="oxxo_recurrent", Reference=paymentMethod.Reference, BarcodeUrl=paymenthMethods.Barcode
			}
		}
	}

	return make(map[string]interface{})
}

*/

/*{
	"card_payment_methods": []models.CardPaymentMethod{
		{Type: "card", CardEnding: 1432, CardToken: "test_3ed98d239dn9238", Default: true},
		{Type: "card", CardEnding: 4352, CardToken: "test_3edr32r32432432", Default: false},
		{Type: "card", CardEnding: 8032, CardToken: "test_3ed98d239dn9258", Default: false},
	},
	"oxxo_payment_method": models.OxxoPaymentMethod{Type: "oxxo", Reference: "0000-0000-0000-0000", BarcodeUrl: "test.net"},
	"spei_payment_method": models.SpeiPaymentMethod{Type: "spei", Reference: "16537213202193820183"},
},*/
