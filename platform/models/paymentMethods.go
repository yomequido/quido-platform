package models

type CardPaymentMethod struct {
	Type       string `json:"type"`
	CardEnding int    `json:"card_ending"`
	CardToken  string `json:"card_token"`
	Default    bool   `json:"default"`
}

type OxxoPaymentMethod struct {
	Type       string `json:"type"`
	Reference  string `json:"reference"`
	BarcodeUrl string `json:"barcode_url"`
}

type SpeiPaymentMethod struct {
	Type      string `json:"type"`
	Reference string `json:"reference"`
}
