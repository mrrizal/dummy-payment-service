package router

import (
	"github.com/gin-gonic/gin"

	"payment-service/internal/http/handler"
)

func Register(r *gin.Engine, paymentHandler *handler.PaymentHandler) {
	v1 := r.Group("/v1")
	{
		payments := v1.Group("/payments")
		{
			payments.POST("", paymentHandler.Create)
		}
	}
}
