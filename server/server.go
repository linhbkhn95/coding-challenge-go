package server

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"coding-challenge-go/pkg/product"
	"coding-challenge-go/pkg/seller"
	"coding-challenge-go/server/config"
	"coding-challenge-go/server/controller"
)

func Server(cfg *config.AppConfig) {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	db, err := sql.Open("mysql", cfg.MySQLConfig.DSN())

	if err != nil {
		log.Error().Err(err).Msg("Fail to create server")
		return
	}

	defer db.Close()

	r := gin.New()
	v1 := r.Group("api/v1")
	productRepository := product.NewRepository(db)
	sellerRepository := seller.NewRepository(db)
	emailProvider := seller.NewEmailProvider()

	productSvc := product.NewService(productRepository, sellerRepository, emailProvider)
	sellerSvc := seller.NewService(sellerRepository)
	productController := controller.NewProductController(productSvc)
	sellerController := controller.NewSellerController(sellerSvc)

	// path for product
	v1.GET("products", productController.List)
	v1.GET("product", productController.Get)
	v1.POST("product", productController.Post)
	v1.PUT("product", productController.Put)
	v1.DELETE("product", productController.Delete)

	// Path for seller
	v1.GET("sellers", sellerController.List)

	log.Info().Msg("Start server")
	log.Fatal().Err(r.Run(fmt.Sprintf(":%d", cfg.HTTPPort))).Msg("Fail to listen and serve")

}