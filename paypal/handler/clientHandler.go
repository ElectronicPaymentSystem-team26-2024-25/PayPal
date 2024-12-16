package handler

import (
	"encoding/json"
	"net/http"

	"paypay.xws.com/paypal/model"
	"paypay.xws.com/paypal/service"
)

type ClientHandler struct {
	Service *service.ClientService
}

func (handler *ClientHandler) CreateClient(writer http.ResponseWriter, req *http.Request) {
	var client = model.Client{}
	err := json.NewDecoder(req.Body).Decode(&client)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	response, err := handler.Service.CreateClient(&client)
	if err != nil {
		println("Error while creating a new Tour")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(response)
}
