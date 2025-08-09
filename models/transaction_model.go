package models

import (
	"time"
)

type TransactionModel struct {
	Id                 string    `json:"id"`
	NumberBilling      string    `json:"number_billing"`
	RequestId          string    `json:"request_id"`
	CustomerPan        string    `json:"customer_pan"`
	Amount             float64   `json:"amount"`
	TransactionDate    time.Time `json:"transaction_date"`
	RetrievalRefNum    string    `json:"retrieval_ref_num"`
	CustomerName       string    `json:"customer_name"`
	MerchantId         string    `json:"merchant_id"`
	MerchantName       string    `json:"merchant_name"`
	MerchantCity       string    `json:"merchant_city"`
	CurrencyCode       string    `json:"currency_code"`
	PaymentStatus      string    `json:"payment_status"`
	PaymentDescription string    `json:"payment_description"`
	CreatedBy          string    `json:"created_by"`
}
