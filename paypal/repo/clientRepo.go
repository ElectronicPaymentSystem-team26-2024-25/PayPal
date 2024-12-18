package repo

import (
	"gorm.io/gorm"
	"paypay.xws.com/paypal/model"
)

type ClientRepo struct {
	DbConnection *gorm.DB
}

func (repo *ClientRepo) CreateClient(client *model.Client) error {
	dbResult := repo.DbConnection.Create(client)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("CreateClient: Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *ClientRepo) GetClient(merchantId string) (*model.Client, error) {
	var client = model.Client{}
	dbResult := repo.DbConnection.First(&client, "merchant_id = ?", merchantId)
	if dbResult.Error != nil {
		return &client, dbResult.Error
	}
	return &client, nil
}
