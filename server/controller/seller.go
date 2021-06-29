package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"coding-challenge-go/pkg/seller"
)

func NewSellerController(sellerSvc seller.Service) *sellerController {
	return &sellerController{
		sellerSvc: sellerSvc,
	}
}

type sellerController struct {
	sellerSvc seller.Service
}

func (sc *sellerController) List(c *gin.Context) {
	sellers, err := sc.sellerSvc.List(c.Request.Context())

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller list"})
		return
	}

	sellersJson, err := json.Marshal(sellers)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal sellers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal sellers"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", sellersJson)
}

func (sc *sellerController) Top10ByProduct(c *gin.Context) {
	sellers, err := sc.sellerSvc.Top10ByProduct(c.Request.Context())

	if err != nil {
		log.Error().Err(err).Msg("Fail to query top 10 seller by product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query top 10 seller"})
		return
	}

	sellersJson, err := json.Marshal(sellers)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal sellers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal sellers"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", sellersJson)
}
