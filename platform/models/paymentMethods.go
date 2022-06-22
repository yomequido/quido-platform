package models

import "fmt"

type PaymentMethod struct {
	Type       string `json:"type"`
	CardEnding string `json:"card_ending,omitempty"`
	CardToken  string `json:"card_token,omitempty"`
	Default    bool   `json:"default,omitempty"`
	Reference  string `json:"reference,omitempty"`
	BarcodeUrl string `json:"barcode_url,omitempty"`
	Id         string `json:"id"`
}

type PaymentMethods struct {
	CardPaymentMethods []PaymentMethod `json:"card_payment_methods"`
	OxxoPaymentMethod  PaymentMethod   `json:"oxxo_payment_method"`
	SpeiPaymentMethod  PaymentMethod   `json:"spei_payment_method"`
}

type ConektaCustomer struct {
	PaymentSources PaymentSources `json:"payment_sources"`
}

type PaymentSources struct {
	Data []PaymentMethod `json:"data"`
}

func (cus ConektaCustomer) String() string {
	return fmt.Sprintf("%s", cus.PaymentSources.Data[0].Type)
}
