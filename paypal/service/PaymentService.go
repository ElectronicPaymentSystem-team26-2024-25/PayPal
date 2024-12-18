package service

import (
	"context"

	"github.com/plutov/paypal/v4"
	"paypay.xws.com/paypal/model"
	"paypay.xws.com/paypal/repo"
)

type PaymentService struct {
	ClientRepo *repo.ClientRepo
}

func (service *PaymentService) ProcessPayment(paymentReq *model.PaymentRequest) (*model.PaymentResponse, error) {
	paymentRes := model.PaymentResponse{}
	client, err := service.ClientRepo.GetClient(paymentReq.MerchantId)
	if err != nil {
		return nil, err
	}

	c, err := paypal.NewClient(client.ClientId, client.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		return nil, err
	}
	order, err := CreateOrder(c, paymentReq)
	if err != nil {
		return nil, err
	}

	approveLink := GetApproveLink(order)
	paymentRes.Message = approveLink
	return &paymentRes, nil
}

func CreateOrder(c *paypal.Client, paymentReq *model.PaymentRequest) (*paypal.Order, error) {
	ctx := context.Background()
	units := []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Currency: paymentReq.Currency, //TODO: Change to USD
				Value:    paymentReq.Amount},
		},
	}
	//TODO: Add ReturnURL, CancelURL
	appCtx := &paypal.ApplicationContext{
		UserAction: paypal.UserActionPayNow,
		BrandName:  paymentReq.BrandName,
	}
	order, err := c.CreateOrder(ctx, paypal.OrderIntentCapture, units, nil, appCtx)

	if err != nil {
		return nil, err
	}
	return order, nil
}

func GetApproveLink(order *paypal.Order) string {
	var approvalURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}
	return approvalURL
}
