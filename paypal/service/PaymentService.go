package service

import (
	"context"
	"fmt"
	"time"

	"github.com/plutov/paypal/v4"
	"paypay.xws.com/paypal/model"
	"paypay.xws.com/paypal/repo"
	"paypay.xws.com/paypal/security"
)

type PaymentService struct {
	ClientRepo   *repo.ClientRepo
	OrderRepo    *repo.OrderRepo
	SecurityRepo *repo.SecurityRepo
}

func (service *PaymentService) ProcessPayment(paymentReq *model.PaymentRequest) (*model.PaymentApproveLink, error) {
	paymentRes := model.PaymentApproveLink{}

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
	fmt.Print("Trying to capture payment")
	var paymentRes = model.PaymentResponse{}
	order, err := service.OrderRepo.GetOrder(ppOrderId)
	if err != nil {
		return nil, err
	}

	secCtx, err1 := service.SecurityRepo.GetOrder(order.Id)
	if err1 != nil {
		return nil, err1
	}
	decMerchantOrderId, _ := security.Decrypt(order.OrderId, secCtx.IVector)
	decMerchantId, _ := security.Decrypt(order.MerchantId, secCtx.IVector)

	paymentRes.MerchantOrderId = decMerchantOrderId
	paymentRes.AcquirerOrderId = order.PaypalOrderId
	paymentRes.AcquirerTimestamp = order.TimeStamp
	paymentRes.PaymentId = 2 //TODO: Change this

	c, err := service.generatePayPalClient(decMerchantId)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	orderReq := paypal.CaptureOrderRequest{
		PaymentSource: &paypal.PaymentSource{},
	}
	orderRes, err := c.CaptureOrder(ctx, order.PaypalOrderId, orderReq)
	if err != nil {
		//TODO: Order is not successfull, update paymentRes
		return nil, err
	}
	if orderRes.Status == "COMPLETED" {
		fmt.Print("Order completed")
		err := service.updateOrder(order, model.Success)
		if err != nil {
			paymentRes.FailReason = "Could not update the order"
		}
	} else {
		paymentRes.FailReason = "Order couldn't be completed"
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
	bytes := security.CreateIV()

	encOrderId, _ := security.Encrypt(paymentReq.MerchantOrderId, bytes)
	encMerchantId, _ := security.Encrypt(paymentReq.MerchantId, bytes)

	order := &model.Order{
		OrderId:       encOrderId,
		MerchantId:    encMerchantId,
		Amount:        paymentReq.Amount,
		PaypalOrderId: ppOrder.ID,
		TimeStamp:     time.Now(),
	}
	service.OrderRepo.CreateOrder(order)

	orderSecCtx := &model.OrderSecurityContext{
		OrderId: order.Id,
		IVector: bytes,
	}
	service.SecurityRepo.CreateOrder(orderSecCtx)
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
				Currency: "USD",
				Value:    paymentReq.Amount},
		},
	}
	returnUrl := paymentReq.SucessUrl + paymentReq.MerchantOrderId
	fmt.Print(returnUrl)
	//TODO: CancelUrl should be to merchant's webshop
	cancelUrl := "https://webshop-client:5173/"

	appCtx := &paypal.ApplicationContext{
		UserAction: paypal.UserActionPayNow,
		BrandName:  paymentReq.BrandName,
		ReturnURL:  returnUrl,
		CancelURL:  cancelUrl,
	}
	order, err := c.CreateOrder(ctx, paypal.OrderIntentCapture, units, nil, appCtx)

	if err != nil {
		return nil, err
	}
	return order, nil
}
