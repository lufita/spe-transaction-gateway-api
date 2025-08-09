package controllers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	tp "spe-trx-gateway/lookup"
	model "spe-trx-gateway/models"
)

func (s *Server) InquiryController(c *gin.Context) {
	req := tp.InquiryRequest{}
	res := tp.InquiryResponse{}

	if err := c.ShouldBind(&req); err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	//if err := helper.ValidateInquiry(&req, c.GetHeader("X-Signature")); err != nil {
	//	res.Code = tp.UNAUTHORIZED_CODE
	//	res.Message = err.Error()
	//	c.JSON(http.StatusUnauthorized, res)
	//	return
	//}

	httpResp, respBody, err := s.CheckTransaction(c.Request.Context(), req)
	if err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(httpResp, respBody)
	return
}

func (s *Server) CheckTransaction(ctx context.Context, req tp.InquiryRequest) (httpCode int, res tp.InquiryResponse, err error) {
	data := model.InquiryModel{}
	err = s.DB.QueryRow(ctx, `
		SELECT 
		    t.id, 
		    t.number_billing,
		    t.request_id,
		    t.customer_pan,
		    t.amount,
		    t.transaction_datetime::text as trx_date,
		    t.retrieval_reference_number,
		    t.customer_name,
		    t.merchant_id,
		    t.merchant_name,
		    t.merchant_city,
		    t.currency_code,
		    t.payment_status,
		    t.payment_description
		FROM transactions t
		WHERE t.request_id = $1 AND t.number_billing = $2 LIMIT 1`, req.RequestId, req.BillingNumber).
		Scan(
			&data.Id,
			&data.NumberBilling,
			&data.RequestId,
			&data.CustomerPan,
			&data.Amount,
			&data.TransactionDate,
			&data.RetrievalRefNum,
			&data.CustomerName,
			&data.MerchantId,
			&data.MerchantName,
			&data.MerchantCity,
			&data.CurrencyCode,
			&data.PaymentStatus,
			&data.PaymentDescription,
		)

	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			res.Code = tp.INTERNAL_SERVER_ERROR
			res.Message = "data-not-found"
			return http.StatusInternalServerError, res, errors.New("data-not-found")
		}
	}

	res.Code = tp.SUCCESS_CODE
	res.Message = "success"
	res.RequestId = data.RequestId.String
	res.CustomerPan = data.CustomerPan.String
	res.Amount = data.Amount.Float64
	res.TransactionDate = data.TransactionDate.String
	res.RetrievalRefNum = data.RetrievalRefNum.String
	res.BillNumber = data.NumberBilling.String
	res.CustomerName = data.CustomerName.String
	res.MerchantId = data.MerchantId.String
	res.MerchantName = data.MerchantName.String
	res.MerchantCity = data.MerchantCity.String
	res.CurrencyCode = data.CurrencyCode.String
	res.PaymentStatus = data.PaymentStatus.String
	res.PaymentDescription = data.PaymentDescription.String

	return http.StatusOK, res, nil
}
