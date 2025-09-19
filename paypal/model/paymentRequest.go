package model

import "time"

type PaymentRequest struct {
	BrandName         string
	MerchantId        string
	MerchantTimeStamp time.Time
	MerchantOrderId   string
	Amount            string
	ErrorUrl          string
	SucessUrl         string
	FailedUrl         string
}
