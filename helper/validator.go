package helper

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"github.com/go-playground/validator/v10"
	"os"
	tp "spe-trx-gateway/lookup"
	"strings"
)

var v = validator.New()

func ValidateNotification(req *tp.NotificationRequest, signature string) error {
	if err := v.Struct(req); err != nil {
		return err
	}
	if signature == "" {
		return errors.New("missing X-Signature")
	}
	key := os.Getenv("SIGNING_KEY")
	msg := strings.Join([]string{req.RequestId, req.RetrievalRefNum, req.MerchantId}, ":")
	if !verifyHMAC512B64(msg, key, signature) {
		return errors.New("invalid signature")
	}
	return nil
}

func ValidateInquiry(req *tp.InquiryRequest, signature string) error {
	if err := v.Struct(req); err != nil {
		return err
	}
	if signature == "" {
		return errors.New("missing X-Signature")
	}
	key := os.Getenv("SIGNING_KEY")
	if !verifyHMAC512B64(req.BillingNumber, key, signature) {
		return errors.New("invalid signature")
	}
	return nil
}

func verifyHMAC512B64(message, key, sig string) bool {
	mac := hmac.New(sha512.New, []byte(key))
	mac.Write([]byte(message))
	exp := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	// gunakan constant-time compare
	return hmac.Equal([]byte(exp), []byte(sig))
}
