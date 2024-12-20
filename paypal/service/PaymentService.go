package service

import (
	"context"
	"time"

	"github.com/plutov/paypal/v4"
	"paypay.xws.com/paypal/model"
	"paypay.xws.com/paypal/repo"
)

type PaymentService struct {
	ClientRepo *repo.ClientRepo
	OrderRepo  *repo.OrderRepo
}

func (service *PaymentService) ProcessPayment(paymentReq *model.PaymentRequest) (*model.PaymentResponse, error) {
	paymentRes := model.PaymentResponse{}

	c, err := service.generatePayPalClient(paymentReq.MerchantId)
	if err != nil {
		return nil, err
	}
	ppOrder, err := createOrder(c, paymentReq)
	if err != nil {
		return nil, err
	}

	service.saveOrder(paymentReq, ppOrder)
	approveLink := getApproveLink(ppOrder)
	paymentRes.Message = approveLink
	return &paymentRes, nil
}

func (service *PaymentService) CapturePayment(ppOrderId string) (*model.PaymentResponse, error) {
	var paymentRes = model.PaymentResponse{}
	order, err := service.OrderRepo.GetOrder(ppOrderId)
	if err != nil {
		return nil, err
	}

	c, err := service.generatePayPalClient(order.MerchantId)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	orderReq := paypal.CaptureOrderRequest{
		PaymentSource: &paypal.PaymentSource{},
	}
	orderRes, err := c.CaptureOrder(ctx, order.PaypalOrderId, orderReq)
	if err != nil {
		return nil, err
	}
	if orderRes.Status == "COMPLETED" {
		err := service.updateOrder(order, model.Success)
		if err != nil {
			paymentRes.Message = "Could not update the order"
		} else {
			paymentRes.Message = "Order completed successfully"
		}
	} else {
		paymentRes.Message = "Order couldn't be completed"
	}
	return &paymentRes, nil
}

func (service *PaymentService) generatePayPalClient(merchantId string) (*paypal.Client, error) {
	client, err := service.ClientRepo.GetClient(merchantId)
	if err != nil {
		return nil, err
	}

	c, err := paypal.NewClient(client.ClientId, client.ClientSecret, paypal.APIBaseSandBox)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getApproveLink(order *paypal.Order) string {
	var approvalURL string
	for _, link := range order.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}
	return approvalURL
}

func (service *PaymentService) saveOrder(paymentReq *model.PaymentRequest, ppOrder *paypal.Order) {
	order := &model.Order{
		OrderId:       paymentReq.OrderId,
		MerchantId:    paymentReq.MerchantId,
		Amount:        paymentReq.Amount,
		PaypalOrderId: ppOrder.ID,
		TimeStamp:     time.Now(),
	}
	service.OrderRepo.CreateOrder(order)
}

func (service *PaymentService) updateOrder(order *model.Order, newStatus model.Status) error {
	order.OrderStatus = newStatus
	if _, err := service.OrderRepo.UpdateOrder(order); err != nil {
		return err
	}
	return nil
}

func createOrder(c *paypal.Client, paymentReq *model.PaymentRequest) (*paypal.Order, error) {
	ctx := context.Background()
	units := []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Currency: paymentReq.Currency, //TODO: Change to USD
				Value:    paymentReq.Amount},
		},
	}
	//TODO: Add ReturnURL, CancelURL
	//TODO: Change returnUrl to be less coupled with PSP
	//TODO: CancelUrl should be to merchant's webshop
	returnUrl := "http://localhost:4200/success/" + paymentReq.OrderId

	appCtx := &paypal.ApplicationContext{
		UserAction: paypal.UserActionPayNow,
		BrandName:  paymentReq.BrandName,
		ReturnURL:  returnUrl,
	}
	order, err := c.CreateOrder(ctx, paypal.OrderIntentCapture, units, nil, appCtx)

	if err != nil {
		return nil, err
	}
	return order, nil
}
