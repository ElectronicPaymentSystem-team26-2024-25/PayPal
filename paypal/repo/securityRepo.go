package repo

import (
	"gorm.io/gorm"
	"paypay.xws.com/paypal/model"
)

type SecurityRepo struct {
	DbConnection *gorm.DB
}

func (repo *SecurityRepo) CreateOrder(order *model.OrderSecurityContext) error {
	dbResult := repo.DbConnection.Create(order)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Order security context created: Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *SecurityRepo) GetOrder(orderId int64) (*model.OrderSecurityContext, error) {
	var order = model.OrderSecurityContext{}
	dbResult := repo.DbConnection.First(&order, "order_id = ?", orderId)
	if dbResult.Error != nil {
		return &order, dbResult.Error
	}
	return &order, nil
}
