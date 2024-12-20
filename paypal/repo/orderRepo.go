package repo

import (
	"gorm.io/gorm"
	"paypay.xws.com/paypal/model"
)

type OrderRepo struct {
	DbConnection *gorm.DB
}

func (repo *OrderRepo) CreateOrder(order *model.Order) error {
	dbResult := repo.DbConnection.Create(order)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("CreateOrder: Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *OrderRepo) GetOrder(orderId string) (*model.Order, error) {
	var order = model.Order{}
	dbResult := repo.DbConnection.First(&order, "paypal_order_id = ?", orderId)
	if dbResult.Error != nil {
		return &order, dbResult.Error
	}
	return &order, nil
}

func (repo *OrderRepo) UpdateOrder(order *model.Order) (*model.Order, error) {
	dbResult := repo.DbConnection.Save(order)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return order, nil
}
