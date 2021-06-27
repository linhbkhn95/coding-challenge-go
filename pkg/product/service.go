package product

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"coding-challenge-go/pkg/seller"
)

var (
	serverAddress = "http://localhost:8080"
)

const (
	defaultListPageSize = 10
)

type (
	Service interface {
		List(ctx context.Context, params *FilterParams) ([]*ProductInfo, error)
		FindByUUID(ctx context.Context, uuid string) (*ProductInfo, error)
		Update(ctx context.Context, product *Product) error
		Create(ctx context.Context, product *Product) error
		Delete(ctx context.Context, uuid string) error
	}

	FilterParams struct {
		Paging *Paging
	}

	Paging struct {
		PageNumber int
	}

	Repository interface {
		// List to get list of products by offset and limit.
		List(ctx context.Context, offset int, limit int) ([]*Product, error)
		// FindByUUID return a product when found.
		FindByUUID(ctx context.Context, uuid string) (*Product, error)
		// Update to update product information.
		Update(ctx context.Context, product *Product) error
		Create(ctx context.Context, product *Product) error
		Delete(ctx context.Context, product *Product) error
	}

	service struct {
		repo         Repository
		sellerRepo   seller.Repository
		notiProvider seller.NotiProvider
	}

	ProductInfo struct {
		*Product
		Seller *SellerInfo `json:"seller"`
	}
	SellerInfo struct {
		UUID  string       `json:"uuid"`
		Links *SellerLinks `json:"_links"`
	}

	SellerLinks struct {
		Self *SelfSellerLink `json:"self"`
	}
	SelfSellerLink struct {
		Href string `json:"href"`
	}
)

func NewService(productRepo Repository, sellerRepo seller.Repository, notiProvider seller.NotiProvider) Service {
	return &service{
		repo:         productRepo,
		sellerRepo:   sellerRepo,
		notiProvider: notiProvider,
	}
}

func (s *service) List(ctx context.Context, params *FilterParams) ([]*ProductInfo, error) {
	if params.Paging == nil {
		params.Paging = &Paging{
			PageNumber: 0,
		}
	}
	products, err := s.repo.List(ctx, (params.Paging.PageNumber-1)*defaultListPageSize, defaultListPageSize)
	if err != nil {
		return nil, err
	}
	result := make([]*ProductInfo, len(products))
	for i, p := range products {
		result[i] = &ProductInfo{
			Product: p,
			Seller:  generateSellerInfo(p.SellerUUID),
		}
	}
	return result, nil
}

func (s *service) FindByUUID(ctx context.Context, uuid string) (*ProductInfo, error) {
	product, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, &ProductNotFoundError{id: uuid}
	}

	return &ProductInfo{
		Product: product,
		Seller:  generateSellerInfo(product.SellerUUID),
	}, nil
}

func generateSellerInfo(sellerUUID string) *SellerInfo {
	return &SellerInfo{
		UUID: sellerUUID,
		Links: &SellerLinks{
			Self: &SelfSellerLink{
				Href: generateSellerLink(sellerUUID),
			},
		},
	}
}

func generateSellerLink(sellerUUID string) string {
	return fmt.Sprintf("%s/api/v1/sellers/%s", serverAddress, sellerUUID)
}

func (s *service) Update(ctx context.Context, product *Product) error {
	p, err := s.repo.FindByUUID(ctx, product.UUID)
	if err != nil {
		return err
	}
	if p == nil {
		return &ProductNotFoundError{id: product.UUID}
	}

	oldStock := p.Stock

	err = s.repo.Update(ctx, product)
	if err != nil {
		return err
	}

	if oldStock != product.Stock {
		sl, err := s.sellerRepo.FindByUUID(ctx, product.SellerUUID)

		if err != nil {
			return err
		}
		if sl == nil {
			return SellerNotFoundError{id: product.SellerUUID}
		}
		s.notiProvider.StockChanged(oldStock, product.Stock, sl.Email)
		log.Info().Msg(fmt.Sprintf("%s Warning sent to %s (Phone: %s): %s Product stock changed", s.notiProvider.Type().String(), sl.UUID, sl.Phone, p.Name))

	}
	return nil
}

func (s *service) Create(ctx context.Context, product *Product) error {
	seller, err := s.sellerRepo.FindByUUID(ctx, product.SellerUUID)
	if err != nil {
		return err
	}
	if seller == nil {
		return &SellerNotFoundError{id: product.UUID}
	}
	return s.repo.Create(ctx, product)
}

func (s *service) Delete(ctx context.Context, uuid string) error {
	product, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if product == nil {
		return &ProductNotFoundError{id: uuid}
	}
	return s.repo.Delete(ctx, product)
}
