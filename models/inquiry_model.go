package models

import "database/sql"

type InquiryModel struct {
	Id                 sql.NullString  `json:"id"`
	NumberBilling      sql.NullString  `json:"number_billing"`
	RequestId          sql.NullString  `json:"request_id"`
	CustomerPan        sql.NullString  `json:"customer_pan"`
	Amount             sql.NullFloat64 `json:"amount"`
	TransactionDate    sql.NullString  `json:"transaction_datetime"`
	RetrievalRefNum    sql.NullString  `json:"retrieval_reference_number"`
	CustomerName       sql.NullString  `json:"customer_name"`
	MerchantId         sql.NullString  `json:"merchant_id"`
	MerchantName       sql.NullString  `json:"merchant_name"`
	MerchantCity       sql.NullString  `json:"merchant_city"`
	CurrencyCode       sql.NullString  `json:"currency_code"`
	PaymentStatus      sql.NullString  `json:"payment_status"`
	PaymentDescription sql.NullString  `json:"payment_description"`
}
