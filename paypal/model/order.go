package model

import (
	"time"

	"gorm.io/gorm"
)

type Status int

const (
	InProgress Status = iota
	Failed
	Success
)

type Order struct {
	Id            int64 `json:"id"`
	OrderId       string
	PaypalOrderId string
	MerchantId    string
	Amount        string
	OrderStatus   Status
	TimeStamp     time.Time
}

func (order *Order) BeforeCreate(scope *gorm.DB) error {
	return nil
}
