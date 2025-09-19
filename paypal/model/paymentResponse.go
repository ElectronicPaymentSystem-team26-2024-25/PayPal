package model

import "time"

type PaymentResponse struct {
	MerchantOrderId   string    //This is from PaymentRequest field OrderId
	AcquirerOrderId   string    //This is PayPal's orderId
	AcquirerTimestamp time.Time //This is PayPal's timestamp from order
	PaymentId         int64     //This is from PSP's PaymentMethod
	RedirectUrl       string    //This is from PSP, where to go after payment
	FailReason        string
}

// TODO: This should be expanded in case something wrong happens
type PaymentApproveLink struct {
	Message string
}
