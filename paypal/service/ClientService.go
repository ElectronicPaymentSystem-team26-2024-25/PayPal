package service

import (
	"paypay.xws.com/paypal/model"
	"paypay.xws.com/paypal/repo"
)

type ClientService struct {
	Repo *repo.ClientRepo
}

func (service *ClientService) CreateClient(client *model.Client) (*model.Client, error) {
	err := service.Repo.CreateClient(client)
	if err != nil {
		return nil, err
	}
	return client, nil
}
