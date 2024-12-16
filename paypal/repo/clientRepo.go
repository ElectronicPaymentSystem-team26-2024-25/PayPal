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
