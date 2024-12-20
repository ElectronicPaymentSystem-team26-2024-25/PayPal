package model

type PaymentRequest struct {
	BrandName  string
	MerchantId string
	OrderId    string
	Currency   string
	Amount     string
}
