package handler

import (
	"net/http"
	"payment-service/internal/core/usecase"

	"github.com/gin-gonic/gin"
)

type createPaymentRequest struct {
	OrderID  string `json:"order_id" binding:"required"`
	PayerID  int    `json:"payer_id" binding:"required"`
	Amount   int    `json:"amount" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Provider string `json:"provider" binding:"required"`
	Method   string `json:"method" binding:"required"`
}

type createPaymentResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}

type PaymentHandler struct {
	createPaymentUC *usecase.CreatePaymentUsecase
}

func NewPaymentHandler(
	createPaymentUC *usecase.CreatePaymentUsecase,
) *PaymentHandler {
	return &PaymentHandler{
		createPaymentUC: createPaymentUC,
	}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 🔑 Idempotency-Key wajib dari header
	idempotencyKey := c.GetHeader("Idempotency-Key")
	if idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Idempotency-Key header is required",
		})
		return
	}

	output, err := h.createPaymentUC.Execute(
		c.Request.Context(),
		usecase.CreatePaymentInput{
			OrderID:        req.OrderID,
			PayerID:        req.PayerID,
			Amount:         req.Amount,
			Currency:       req.Currency,
			Provider:       req.Provider,
			Method:         req.Method,
			IdempotencyKey: idempotencyKey,
		},
	)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, createPaymentResponse{
		PaymentID: output.PaymentID,
		Status:    string(output.Status),
	})
}
