package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"math/rand"
	"net/http"
	"regexp"
	//"spe-trx-gateway/helper"
	tp "spe-trx-gateway/lookup"
	model "spe-trx-gateway/models"
	"strconv"
	"strings"
	"time"
)

var amountRegex = regexp.MustCompile(`^\d+\.\d{2}$`)

func (s *Server) NotificationController(c *gin.Context) {
	req := tp.NotificationRequest{}
	res := tp.NotificationResponse{}

	_, clientID, akh := getClaims(c)
	if clientID == "" || akh == "" {
		res.Code = tp.UNAUTHORIZED_CODE
		res.Message = "missing jwt claims"
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	_, err := s.ValidateAccess(c.Request.Context(), clientID, akh)
	if err != nil {
		switch {
		case errors.Is(err, ErrClientNotFound):
			res.Code = tp.UNAUTHORIZED_CODE
			res.Message = "client not found"
			c.JSON(http.StatusUnauthorized, res)
			return
		case errors.Is(err, ErrClientInactive):
			res.Code = tp.UNAUTHORIZED_CODE
			res.Message = "client disabled"
			c.JSON(http.StatusUnauthorized, res)
			return
		case errors.Is(err, ErrFingerprintMismatch):
			res.Code = tp.UNAUTHORIZED_CODE
			res.Message = "token no longer valid"
			c.JSON(http.StatusUnauthorized, res)
			return
		default:
			res.Code = tp.UNAUTHORIZED_CODE
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

	err = validateAmount(req.Amount)
	if err != nil {
		res.Code = tp.BAD_REQUEST
		res.Message = "invalid amount: " + err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	httpResp, respBody, err := s.ProcessPaymentNotification(c.Request.Context(), req, clientID)
	if err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	c.JSON(httpResp, respBody)
	return
}

func (s *Server) ProcessPaymentNotification(ctx context.Context, req tp.NotificationRequest, createdBy string) (httpCode int, res tp.NotificationResponse, err error) {
	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = "begin tx: " + err.Error()
		return http.StatusInternalServerError, res, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	amountFloat, err := convertToFloat64(req.Amount)
	if err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = err.Error()
		return http.StatusInternalServerError, res, err
	}

	trxDateTime, err := ValidateTimestampTZ(req.TransactionDate)
	if err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = err.Error()
		return http.StatusInternalServerError, res, err
	}

	data := model.TransactionModel{
		NumberBilling:      generateBillingNumber(),
		RequestId:          req.RequestId,
		CustomerPan:        req.CustomerPan,
		Amount:             amountFloat,
		TransactionDate:    trxDateTime,
		RetrievalRefNum:    req.RetrievalRefNum,
		CustomerName:       req.CustomerName,
		MerchantId:         req.MerchantId,
		MerchantName:       req.MerchantName,
		MerchantCity:       req.MerchantCity,
		CurrencyCode:       req.CurrencyCode,
		PaymentStatus:      req.PaymentStatus,
		PaymentDescription: req.PaymentDescription,
		CreatedBy:          createdBy,
	}
	err = s.CreateNewTransaction(ctx, tx, data)
	if err != nil {
		res.Code = tp.INTERNAL_SERVER_ERROR
		res.Message = "create new transaction: " + err.Error()
		return http.StatusInternalServerError, res, err
	}

	res.Code = tp.SUCCESS_CODE
	res.Message = "success"
	return http.StatusOK, res, nil
}

func (s *Server) CreateNewTransaction(ctx context.Context, tx pgx.Tx, data model.TransactionModel) error {
	err := tx.QueryRow(ctx, `
		INSERT INTO transactions (
		    number_billing, 
			request_id, 
			customer_pan, 
		    amount, 
			transaction_datetime, 
			retrieval_reference_number, 
			customer_name,
			merchant_id,
			merchant_name,
			merchant_city,
			currency_code,
			payment_status,
			payment_description,
			created_at,
			created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, CURRENT_TIMESTAMP , $14) RETURNING id
	`, data.NumberBilling, data.RequestId, data.CustomerPan, data.Amount, data.TransactionDate, data.RetrievalRefNum,
		data.CustomerName, data.MerchantId, data.MerchantName, data.MerchantCity, data.CurrencyCode, data.PaymentStatus,
		data.PaymentDescription, data.CreatedBy).Scan(&data.Id)
	if err != nil {
		return err
	}

	err = s.CreateTransactionHist(ctx, tx, data)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) CreateTransactionHist(ctx context.Context, tx pgx.Tx, data model.TransactionModel) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	jsonString := string(jsonData)

	_, err = tx.Exec(ctx, `
		INSERT INTO transaction_hist (
		    transaction_id,
			event_type,
			event_data,
			created_at,
			created_by
		) VALUES ($1, $2, $3, CURRENT_TIMESTAMP , $4)
	`, data.Id, "INSERT", jsonString, data.CreatedBy)
	if err != nil {
		return err
	}

	return nil
}

func validateAmount(s string) error {
	s = strings.TrimSpace(s)
	if !amountRegex.MatchString(s) {
		return errors.New("amount must be in format ###.## with exactly 2 decimals")
	}

	return nil
}

func convertToFloat64(s string) (float64, error) {
	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return amount, nil
}

func ValidateTimestampTZ(value string) (time.Time, error) {
	layout := "2006-01-02T15:04:05"
	t, err := time.Parse(layout, value)
	if err != nil {
		return time.Time{}, errors.New("invalid timestamptz format, must be YYYY-MM-DDTHH:MM:SS")
	}
	return t, nil
}

func generateBillingNumber() string {
	now := time.Now().UTC()
	timePart := now.Format("20060102150405") // 14 digit (tahun → detik)

	randPart := strconv.Itoa(rand.Intn(900000) + 100000) // 6 digit

	billingNumber := fmt.Sprintf("%s%s", timePart, randPart)

	return billingNumber
}
