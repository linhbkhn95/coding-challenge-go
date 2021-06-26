package controller

import (
	"coding-challenge-go/pkg/product"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func NewProductController(productSvc product.Service) *productController {
	return &productController{
		productSvc: productSvc,
	}
}

type productController struct {
	productSvc product.Service
}

func (pc *productController) List(c *gin.Context) {
	request := &struct {
		Page int `form:"page,default=1"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := pc.productSvc.List(c.Request.Context(), &product.FilterParams{Paging: &product.Paging{
		PageNumber: request.Page,
	}})

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Fail to query product list with err=%s", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product list"})
		return
	}

	productsJson, err := json.Marshal(products)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal products"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productsJson)
}

func (pc *productController) Get(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := pc.productSvc.FindByUUID(c.Request.Context(), request.UUID)

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Fail to get product with err=%s", err.Error()))

		if _, ok := err.(*product.ProductNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	productJson, err := json.Marshal(p)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productJson)
}

func (pc *productController) Post(c *gin.Context) {
	request := &struct {
		Name   string `form:"name"`
		Brand  string `form:"brand"`
		Stock  int    `form:"stock"`
		Seller string `form:"seller"`
	}{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &product.Product{
		UUID:       uuid.New().String(),
		Name:       request.Name,
		Brand:      request.Brand,
		Stock:      request.Stock,
		SellerUUID: request.Seller,
	}

	err := pc.productSvc.Create(c.Request.Context(), p)
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Fail to create product with err=%s", err.Error()))
		if _, ok := err.(*product.SellerNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to insert product"})
		return
	}

	jsonData, err := json.Marshal(p)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

func (pc *productController) Put(c *gin.Context) {
	queryRequest := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(queryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	request := &struct {
		Name  string `form:"name"`
		Brand string `form:"brand"`
		Stock int    `form:"stock"`
	}{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &product.Product{
		UUID:  queryRequest.UUID,
		Name:  request.Name,
		Brand: request.Brand,
		Stock: request.Stock,
	}
	err := pc.productSvc.Update(c.Request.Context(), p)

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Fail to update product with err=%s", err.Error()))

		if _, ok := err.(*product.ProductNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if _, ok := err.(*product.SellerNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	jsonData, err := json.Marshal(p)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

func (pc *productController) Delete(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := pc.productSvc.Delete(c.Request.Context(), request.UUID)

	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Fail to delete product with err=%s", err.Error()))

		if _, ok := err.(*product.ProductNotFoundError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
