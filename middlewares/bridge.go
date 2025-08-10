package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"spe-trx-gateway/config"
	"spe-trx-gateway/lookup"
	"strings"
)

func SignatureMiddleware() gin.HandlerFunc {
	secret := lookup.JWTTokenSecret

	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		raw := strings.TrimSpace(h[len("Bearer "):])

		tkn, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !tkn.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := tkn.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		//jti, _ := claims["jti"].(string)
		cid, _ := claims["sub"].(string)
		fpr, _ := claims["akh"].(string)
		if cid == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token missing fields"})
			return
		}

		if _, err := config.ReadRedis(raw); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired or revoked"})
			return
		}

		c.Set("id_data_client", cid)
		c.Set("id_data_fingerprint", fpr)
		c.Next()
	}
}

func SigNotificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := lookup.NotificationSignSecret
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "signature not provisioned"})
			return
		}

		sig := c.GetHeader("X-Signature")
		if sig == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing signature"})
			return
		}

		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(raw))

		var body struct {
			RequestID  string `json:"request_id"`
			RRN        string `json:"rrn"`
			MerchantID string `json:"merchant_id"`
		}
		if err := json.Unmarshal(raw, &body); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body: bill_number required"})
			return
		}
		// rekonstruksi payload sesuai spesifikasi
		payload := body.RequestID + ":" + body.RRN + ":" + body.MerchantID

		mac := hmac.New(sha512.New, []byte(secret))
		mac.Write([]byte(payload))
		expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		if subtle.ConstantTimeCompare([]byte(expected), []byte(sig)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}

		c.Next()
	}
}

func SigInquiryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := lookup.InquirySignSecret
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "signature not provisioned"})
			return
		}

		sig := c.GetHeader("X-Signature")
		if sig == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing signature"})
			return
		}

		raw, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "read body failed"})
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(raw))

		var body struct {
			BillNumber string `json:"bill_number"`
		}
		if err := json.Unmarshal(raw, &body); err != nil || strings.TrimSpace(body.BillNumber) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid body: bill_number required"})
			return
		}

		payload := body.BillNumber

		mac := hmac.New(sha512.New, []byte(secret))
		mac.Write([]byte(payload))
		expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		if subtle.ConstantTimeCompare([]byte(expected), []byte(sig)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}
		c.Next()
	}
}
