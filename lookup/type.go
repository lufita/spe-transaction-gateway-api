package lookup

const (
	service_name = "spe-transaction-gateway-api"
)

type NotificationRequest struct {
	RequestId          string `json:"request_id" validate:"required,max=32"`
	CustomerPan        string `json:"customer_pan" validate:"required"`
	Amount             string `json:"amount" validate:"required"`
	TransactionDate    string `json:"transaction_datetime" validate:"required"`
	RetrievalRefNum    string `json:"rrn" validate:"required"`
	BillingNumber      string `json:"bill_number" validate:"required"`
	CustomerName       string `json:"customer_name"`
	MerchantId         string `json:"merchant_id" validate:"required"`
	MerchantName       string `json:"merchant_name" validate:"required"`
	MerchantCity       string `json:"merchant_city" validate:"required"`
	CurrencyCode       string `json:"currency_code" validate:"required"`
	PaymentStatus      string `json:"payment_status" validate:"required"`
	PaymentDescription string `json:"payment_description" validate:"required"`
}

type NotificationResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type InquiryRequest struct {
	RequestId     string `json:"request_id" validate:"required"`
	BillingNumber string `json:"bill_number" validate:"required"`
}

type InquiryResponse struct {
	Code               string  `json:"code"`
	Message            string  `json:"message"`
	Id                 string  `json:"id"`
	RequestId          string  `json:"request_id"`
	CustomerPan        string  `json:"customer_pan"`
	Amount             float64 `json:"amount"`
	TransactionDate    string  `json:"transaction_datetime"`
	RetrievalRefNum    string  `json:"rrn"`
	BillNumber         string  `json:"bill_number"`
	CustomerName       string  `json:"customer_name"`
	MerchantId         string  `json:"merchant_id"`
	MerchantName       string  `json:"merchant_name"`
	MerchantCity       string  `json:"merchant_city"`
	CurrencyCode       string  `json:"currency_code"`
	PaymentStatus      string  `json:"payment_status"`
	PaymentDescription string  `json:"payment_description"`
}
