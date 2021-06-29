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

	productRepository := product.NewRepository(db)
	sellerRepository := seller.NewRepository(db)
	notiProvider := getNotiProvider(cfg.NotiProdiverType)
	if notiProvider == nil {
		log.Fatal().Msg("NotiProvider is nil")
	}
	productSvc := product.NewService(productRepository, sellerRepository, notiProvider)
	sellerSvc := seller.NewService(sellerRepository)
	productController := controller.NewProductController(productSvc)
	sellerController := controller.NewSellerController(sellerSvc)

	v1 := r.Group("api/v1")
	{
		// path for product
		v1.GET("products", productController.List)
		v1.GET("product", productController.Get)
		v1.POST("product", productController.Post)
		v1.PUT("product", productController.Put)
		v1.DELETE("product", productController.Delete)

		// Path for seller
		v1.GET("sellers", sellerController.List)
	}

	v2 := r.Group("api/v2")
	{
		v2.GET("products", productController.ListV2)
		v2.GET("product", productController.GetV2)

		v2.GET("sellers/top10", sellerController.Top10ByProduct)
	}

	log.Info().Msg("Start server")
	log.Fatal().Err(r.Run(fmt.Sprintf(":%d", cfg.HTTPPort))).Msg("Fail to listen and serve")

}

func getNotiProvider(providerType string) seller.NotiProvider {
	NotiProvider, err := seller.ProviderTypeString(providerType)
	if err != nil {
		log.Fatal().Err(err).Msg("Unsupport type")
	}
	switch NotiProvider {
	case seller.Email:
		return seller.NewEmailProvider()

	case seller.SMS:
		return seller.NewSMSProvider()

	default:
		log.Fatal().Err(err).Msg("Unsupport type")
	}
	return nil
}
