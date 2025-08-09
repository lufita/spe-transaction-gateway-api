package middlewares

import (
	"github.com/gin-gonic/gin"
)

func SignatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//token := c.GetHeader("Authorization")
		//if token == "" || !strings.HasPrefix(token, "Bearer ") {
		//	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		//	return
		//}
		//sig := c.GetHeader("X-Signature")
		//if sig == "" {
		//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing signature"})
		//	return
		//}
		// Payload & key sesuai tipe request
		// Contoh untuk notification: request_id:rrn:merchant_id
		// Untuk check-status: bill_number
		c.Next()
	}
}
