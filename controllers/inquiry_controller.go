package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"spe-trx-gateway/config"
	tp "spe-trx-gateway/lookup"
	model "spe-trx-gateway/models"
	"time"
)

func cacheKeyInquiry(req tp.InquiryRequest) string {
	return "trx:status:" + req.BillingNumber + ":" + req.RequestId
}

func (s *Server) InquiryController(c *gin.Context) {
	req := tp.InquiryRequest{}
	res := tp.InquiryResponse{}

	_, IDclient, fpr := getClaims(c) // get IdClient and Fingerprint from Claims Set
	if IDclient == "" || fpr == "" {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = "missing jwt claims"
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if _, err := s.ValidateAccess(c.Request.Context(), IDclient, fpr); err != nil {
		switch {
		case errors.Is(err, ErrClientNotFound):
			res.Code = tp.UNAUTHORIZED_CODE
			res.Message = "client not found"
			c.JSON(http.StatusUnauthorized, res)
			return
		case errors.Is(err, ErrClientInactive):
			res.Code = tp.UNAUTHORIZED_CODE
			res.Message = "client has been inactive"
			c.JSON(http.StatusUnauthorized, res)
			return
		case errors.Is(err, ErrFingerprintMismatch):
			res.Code = tp.UNAUTHORIZED_CODE
			res.Message = "token no longer valid"
			c.JSON(http.StatusUnauthorized, res)
			return
		default:
			res.Code = tp.INTERNAL_SERVER_ERROR
			res.Message = "db error"
			c.JSON(http.StatusInternalServerError, res)
			return
		}
	}

	if err := c.ShouldBind(&req); err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	key := cacheKeyInquiry(req)
	if cached, err := config.ReadRedis(key); err == nil && cached != "" {
		var resp tp.InquiryResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	httpResp, respBody, err := s.CheckTransaction(c.Request.Context(), req)
	if err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if b, err := json.Marshal(respBody); err == nil {
		_, _ = config.WriteRedis(key, string(b), 3*time.Hour)
	}

	go func(trxID string, msg tp.InquiryResponse) {

		msgBytes, _ := json.Marshal(msg)
		msgStr := string(msgBytes)

		pctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = s.PublishTrxEvent(pctx, TrxEvent{
			TransactionID: trxID,
			Message:       msgStr,
			Source:        "transaction-inquiry",
			Timestamp:     time.Now().UTC().Format(time.RFC3339),
		})
	}(respBody.Id, respBody)

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
		WHERE t.request_id = $1 AND t.number_billing = $2
		AND t.transaction_datetime::date = now()::date
		LIMIT 1`, req.RequestId, req.BillingNumber).
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
	res.Id = data.Id.String
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
