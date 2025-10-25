package payment

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

var coreApiClient coreapi.Client

func RandomKeyTrans() string {
	time.Sleep(500 * time.Millisecond)
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func NewMidtrans() {
	serverKey := os.Getenv("PAYMENT_SERVER_KEY")
	paymentSandbox := os.Getenv("PAYMENT_SANDBOX")
	if paymentSandbox == "true" {
		coreApiClient.New(serverKey, midtrans.Sandbox)
	} else {
		coreApiClient.New(serverKey, midtrans.Production)
	}
}

func CheckTransaction(orderID string) error {
	res, err := coreApiClient.CheckTransaction(orderID)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", res)
	return nil
}

func CheckStatusB2B(orderID string) error {
	res, err := coreApiClient.GetStatusB2B(orderID)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", res)
	return nil
}

func ApproveTransaction(orderID string) error {
	res, err := coreApiClient.ApproveTransaction(orderID)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", res)
	return nil
}

func DenyTransaction(orderID string) error {
	res, err := coreApiClient.DenyTransaction(orderID)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", res)
	return nil
}

func CancelTransaction(orderID string) error {
	res, err := coreApiClient.CancelTransaction(orderID)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", res)
	return nil
}

func ExpireTransaction(orderID string) error {
	res, err := coreApiClient.ExpireTransaction(orderID)
	if err != nil {
		// do something on error handle
	}
	fmt.Println("Response: ", res)
	return nil
}

func CaptureTransaction(orderID string, grossAmount float64) error {
	refundRequest := &coreapi.CaptureReq{
		TransactionID: orderID,
		GrossAmt:      grossAmount,
	}
	res, err := coreApiClient.CaptureTransaction(refundRequest)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", res)
	return nil
}

func RefundTransaction(orderID string, amount int64, reason string) error {
	refundRequest := &coreapi.RefundReq{
		Amount: amount,
		Reason: reason,
	}

	res, err := coreApiClient.RefundTransaction(orderID, refundRequest)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", res)
	return nil
}

func DirectRefundTransaction(orderID string, amount int64, reason string, refundKey string) error {
	refundRequest := &coreapi.RefundReq{
		RefundKey: refundKey,
		Amount:    amount,
		Reason:    reason,
	}

	// Optional: set payment idempotency key to prevent duplicate request
	coreApiClient.Options.SetPaymentIdempotencyKey(RandomKeyTrans())

	res, err := coreApiClient.DirectRefundTransaction(orderID, refundRequest)
	if err != nil {
		return err
	}
	fmt.Println("Response: ", res)
	return nil
}
