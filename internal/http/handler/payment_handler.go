package handler

import (
	"net/http"
	"payment-service/internal/core/usecase"
	"payment-service/internal/observability"

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

type getPaymentResponse struct {
	PaymentID string `json:"payment_id"`
	OrderID   string `json:"order_id"`
	PayerID   int    `json:"payer_id"`
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	Provider  string `json:"provider"`
	Method    string `json:"method"`
	CreatedAt string `json:"created_at"`
	PaidAt    string `json:"paid_at,omitempty"`
}

type PaymentHandler struct {
	createPaymentUC *usecase.CreatePaymentUsecase
	getPaymentUC    *usecase.GetPaymentUsecase
}

func NewPaymentHandler(
	createPaymentUC *usecase.CreatePaymentUsecase,
	getPaymentUC *usecase.GetPaymentUsecase,
) *PaymentHandler {
	return &PaymentHandler{
		createPaymentUC: createPaymentUC,
		getPaymentUC:    getPaymentUC,
	}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := observability.Tracer().Start(ctx, "PaymentHandler.Create")
	defer span.End()

	var req createPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Set("payment_method", req.Method)

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

func (h *PaymentHandler) Get(c *gin.Context) {
	ctx := c.Request.Context()
	ctx, span := observability.Tracer().Start(ctx, "PaymentHandler.GET")
	defer span.End()

	publicID := c.Param("public_id")
	payment, err := h.getPaymentUC.Execute(
		ctx,
		publicID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var paidAt string
	if payment.PaidAt != nil {
		paidAt = payment.PaidAt.Format("2006-01-02T15:04:05Z07:00")
	}

	c.JSON(http.StatusOK, getPaymentResponse{
		PaymentID: payment.PublicID,
		OrderID:   payment.OrderID,
		PayerID:   payment.PayerID,
		Amount:    payment.Amount,
		Currency:  payment.Currency,
		Status:    string(payment.Status),
		Provider:  payment.Provider,
		Method:    payment.Method,
		CreatedAt: payment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		PaidAt:    paidAt,
	})
}
