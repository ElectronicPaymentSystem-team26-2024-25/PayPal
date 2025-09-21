package model

import "time"

type PaymentResponse struct {
	MerchantOrderId   string    `json:"merchantOrderId"`   //This is from PaymentRequest field OrderId
	AcquirerOrderId   string    `json:"acquirerOrderId"`   //This is PayPal's orderId
	AcquirerTimestamp time.Time `json:"acquirerTimestamp"` //This is PayPal's timestamp from order
	PaymentId         int64     `json:"paymentId"`         //This is from PSP's PaymentMethod
	RedirectUrl       string    `json:"redirectUrl"`       //This is from PSP, where to go after payment
	FailReason        string    `json:"failReason"`
}

// TODO: This should be expanded in case something wrong happens
type PaymentApproveLink struct {
	Message string `json:"message"`
}
