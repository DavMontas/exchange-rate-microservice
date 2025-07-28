package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/domain"
	"github.com/davmontas/exchange-rate-offers/internal/exchangerate/service"
)

func RegisterRoutes(r *gin.Engine, svc *service.ExchangeService, log *zap.SugaredLogger) {
	r.POST("/best-quote", func(c *gin.Context) {
		var pair domain.CurrencyPair
		if err := c.ShouldBindJSON(&pair); err != nil {
			log.Warnw("invalid JSON payload", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
			return
		}

		if pair.From == "" || pair.To == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "both 'from' and 'to' currencies must be provided"})
			return
		}
		if pair.Amount <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "'amount' must be greater than zero"})
			return
		}

		quote := svc.BestQuote(c.Request.Context(), pair)
		if quote.Err != nil {
			log.Errorw("all providers failed", "error", quote.Err)
			c.JSON(http.StatusBadGateway, gin.H{"error": quote.Err.Error()})
			return
		}
		c.JSON(http.StatusOK, quote)
	})
}
