package model

type PaymentRequest struct {
	BrandName  string
	MerchantId string
	Currency   string
	Amount     string
}
