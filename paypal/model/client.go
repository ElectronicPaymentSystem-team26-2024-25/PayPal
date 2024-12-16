package model

import "gorm.io/gorm"

type Client struct {
	Id           int64 `json:"id"`
	Email        string
	ClientId     string
	ClientSecret string
}

func (client *Client) BeforeCreate(scope *gorm.DB) error {
	return nil
}
