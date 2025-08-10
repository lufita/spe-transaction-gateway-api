package controllers

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"spe-trx-gateway/config"
	"spe-trx-gateway/lookup"
	"time"
)

type GetAuthReq struct {
	ApiKey   string `json:"api_key" binding:"required"`
	ClientId string `json:"client_id" binding:"required"`
}

type GetAuthResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (s *Server) RegisterController(c *gin.Context) {
	var req GetAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Minute)
	defer cancel()

	const q = `
		SELECT id, api_key_hash, is_active, COALESCE(rate_limit,60)
		FROM api_clients
		WHERE client_id=$1
		LIMIT 1;
	`
	var (
		clientID   string
		apiKeyHash string
		isActive   bool
		rateLimit  int
	)
	if err := s.DB.QueryRow(ctx, q, req.ClientId).
		Scan(&clientID, &apiKeyHash, &isActive, &rateLimit); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "client id not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if !isActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "client disabled"})
		return
	}

	mac := hmac.New(sha512.New, []byte(lookup.RegisterSecret))
	mac.Write([]byte(req.ApiKey))
	hashedkey := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if apiKeyHash != hashedkey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
		return
	}

	secret := lookup.JWTTokenSecret
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server not configured"})
		return
	}
	exp := time.Now().Add(1 * time.Hour)
	jti := uuid.NewString()

	claims := jwt.MapClaims{
		"jti":           jti,
		"sub":           clientID,
		"client_id":     req.ClientId,
		"api_client_id": clientID,
		"akh":           apiKeyHash,
		"rl":            rateLimit,
		"iat":           time.Now().Unix(),
		"exp":           exp.Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tkn.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sign token failed"})
		return
	}

	ttl := time.Until(exp)
	if _, err := config.WriteRedis(signed, signed, ttl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis set failed"})
		return
	}

	c.JSON(http.StatusOK, GetAuthResp{
		AccessToken: signed,
		TokenType:   "Bearer",
		ExpiresIn:   int64(ttl.Seconds()),
	})
}
