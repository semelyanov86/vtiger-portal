package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/logger"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/webhook"
	"io"
	"log"
	"net/http"
	"os"
)

type PaymentIntentResponse struct {
	ClientSecret string `json:"clientSecret"`
}

func (h *Handler) initPaymentRoutes(api *gin.RouterGroup) {
	products := api.Group("/payments")
	{
		products.GET("/config", h.getPaymentConfig)
		products.POST("/create-payment-intent", h.createPaymentIntent)
		products.POST("/webhook", h.handleWebhook)
	}
}

func (h *Handler) getPaymentConfig(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	if userModel == nil {
		return
	}
	paymentConfig := domain.PaymentConfig{PublishableKey: h.config.Payment.StripePublic}
	res := AloneDataResponse[domain.PaymentConfig]{
		Data: paymentConfig,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) createPaymentIntent(c *gin.Context) {
	userModel := h.getValidatedUser(c)

	if userModel == nil {
		return
	}

	req := service.PaymentIntent{}
	if err := c.BindJSON(&req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "/", "message": "Please pass correct data"})
		return
	}
	if req.PaymentMethodType == "" || req.Currency == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "paymentMethodType", "message": "Payment method type and currency should be not empty"})
		return
	}
	if req.SoId == "" && req.InvoiceId == "" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "invoiceId", "message": "You should pass sales order ID or invoice id"})
		return
	}
	if req.Amount < 0.1 {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "amount", "message": "Please pass correct amount"})
		return
	}
	req.UserId = userModel.Crmid

	pi, err := h.services.Payments.CreatePaymentIntent(c.Request.Context(), req)

	if err != nil {
		// Try to safely cast a generic error to a stripe.Error so that we can get at
		// some additional Stripe-specific information about what went wrong.
		if stripeErr, ok := err.(*stripe.Error); ok {
			logger.Error(logger.GenerateErrorMessageFromString("Other Stripe error occurred: " + stripeErr.Error()))
			newResponse(c, http.StatusBadRequest, stripeErr.Error())
		} else {
			logger.Error(logger.GenerateErrorMessageFromString("Other Stripe error occurred: " + err.Error()))
			newResponse(c, http.StatusInternalServerError, "Unknown server error")
		}

		return
	}

	res := AloneDataResponse[PaymentIntentResponse]{
		Data: PaymentIntentResponse{ClientSecret: pi.ClientSecret},
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) handleWebhook(c *gin.Context) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		log.Printf("ioutil.ReadAll: %v", err)
		return
	}
	event, err := webhook.ConstructEvent(b, c.Request.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		log.Printf("webhook.ConstructEvent: %v", err)
		return
	}
	if event.Type == "checkout.session.completed" {
		fmt.Println("Checkout Session completed!")
	}
}
