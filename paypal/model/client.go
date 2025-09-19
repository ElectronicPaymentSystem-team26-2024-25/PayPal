package model

import "gorm.io/gorm"

//TODO: Encrypt everything, except ID and ClientId
type Client struct {
	Id           int64 `json:"id"`
	MerchantId   string
	Email        string
	ClientId     string
	ClientSecret string
}

func (client *Client) BeforeCreate(scope *gorm.DB) error {
	return nil
}
