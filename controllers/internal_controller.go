package controllers

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"spe-trx-gateway/lookup"
)

func InternalHashRegister(c *gin.Context) {
	secret := lookup.RegisterSecret
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CHECK_STATUS_SIGNING_SECRET not set"})
		return
	}

	var body struct {
		ApiKey string `json:"api_key"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ApiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ApiKey"})
		return
	}
	payload := body.ApiKey

	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(payload))
	hashedkey := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	c.JSON(http.StatusOK, gin.H{
		"payload":      payload,
		"api_key_hash": hashedkey,
	})
}

func InternalHashNotification(c *gin.Context) {
	secret := lookup.NotificationSignSecret
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "NOTIF_SIGNING_SECRET not set"})
		return
	}

	var body struct {
		RequestID  string `json:"request_id"`
		RRN        string `json:"rrn"`
		MerchantID string `json:"merchant_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	payload := body.RequestID + ":" + body.RRN + ":" + body.MerchantID

	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(payload))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	c.JSON(http.StatusOK, gin.H{
		"payload":   payload,
		"signature": signature,
	})
}

func InternalHashInquiry(c *gin.Context) {
	secret := lookup.InquirySignSecret
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CHECK_STATUS_SIGNING_SECRET not set"})
		return
	}

	var body struct {
		BillNumber string `json:"bill_number"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.BillNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	payload := body.BillNumber

	mac := hmac.New(sha512.New, []byte(secret))
	mac.Write([]byte(payload))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	c.JSON(http.StatusOK, gin.H{
		"payload":   payload,
		"signature": signature,
	})
}
