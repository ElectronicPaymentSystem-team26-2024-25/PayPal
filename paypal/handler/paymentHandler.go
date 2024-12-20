package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(response)
}

func (handler *PaymentHandler) CapturePayment(writer http.ResponseWriter, req *http.Request) {
	orderId := mux.Vars(req)["orderId"]
	if orderId == "" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := handler.Service.CapturePayment(orderId)
	if err != nil {
		println("Error capturing payment")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(response)
}
