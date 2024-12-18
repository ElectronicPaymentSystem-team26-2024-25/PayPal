package handler

import (
	"encoding/json"
	"net/http"

	"paypay.xws.com/paypal/model"
	"paypay.xws.com/paypal/service"
)

type PaymentHandler struct {
	Service *service.PaymentService
}

func (handler *PaymentHandler) ProcessPayment(writer http.ResponseWriter, req *http.Request) {
	var paymentReq = model.PaymentRequest{}
	err := json.NewDecoder(req.Body).Decode(&paymentReq)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := handler.Service.ProcessPayment(&paymentReq)
	if err != nil {
		println("Error processing payment")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}
