package controllers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrClientInactive      = errors.New("client disabled")
	ErrFingerprintMismatch = errors.New("token no longer valid")
)

func (s *Server) ValidateAccess(ctx context.Context, clientID string, fpr string) (int, error) {
	const query = `
		SELECT api_key_hash, is_active, COALESCE(rate_limit, 60)
		FROM api_clients
		WHERE id = $1
		LIMIT 1`
	var (
		hash    string
		active  bool
		rateLim int
	)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	if err := s.DB.QueryRow(ctx, query, clientID).Scan(&hash, &active, &rateLim); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrClientNotFound
		}
		return 0, err
	}
	if !active {
		return 0, ErrClientInactive
	}
	if hash != fpr {
		return 0, ErrFingerprintMismatch
	}
	return rateLim, nil
}

func getClaims(c *gin.Context) (jti, clientID, akh string) {
	if v, ok := c.Get("jti"); ok {
		if s, _ := v.(string); s != "" {
			jti = s
		}
	}
	if v, ok := c.Get("id_data_client"); ok {
		if s, _ := v.(string); s != "" {
			clientID = s
		}
	}
	if v, ok := c.Get("id_data_fingerprint"); ok {
		if s, _ := v.(string); s != "" {
			akh = s
		}
	}
	return
}
