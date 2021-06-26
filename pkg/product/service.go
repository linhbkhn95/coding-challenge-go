package product

import (
	"context"

	"coding-challenge-go/pkg/seller"
)

const (
	defaultListPageSize = 10
)

type (
	Service interface {
		List(ctx context.Context, params *FilterParams) ([]*Product, error)
		FindByUUID(ctx context.Context, uuid string) (*Product, error)
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
		List(ctx context.Context, offset int, limit int) ([]*Product, error)
		FindByUUID(ctx context.Context, uuid string) (*Product, error)
		Update(ctx context.Context, product *Product) error
		Create(ctx context.Context, product *Product) error
		Delete(ctx context.Context, product *Product) error
	}

	service struct {
		repo          Repository
		sellerRepo    seller.Repository
		emailProvider seller.EmailProvider
	}
)

func NewService(productRepo Repository, sellerRepo seller.Repository, emailProvider seller.EmailProvider) Service {
	return &service{
		repo:          productRepo,
		sellerRepo:    sellerRepo,
		emailProvider: emailProvider,
	}
}

func (s *service) List(ctx context.Context, params *FilterParams) ([]*Product, error) {
	if params.Paging == nil {
		params.Paging = &Paging{
			PageNumber: 0,
		}
	}
	return s.repo.List(ctx, (params.Paging.PageNumber-1)*defaultListPageSize, defaultListPageSize)
}

func (s *service) FindByUUID(ctx context.Context, uuid string) (*Product, error) {
	product, err := s.repo.FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, &ProductNotFoundError{id: uuid}
	}
	return product, nil
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

	product.Name = p.Name
	product.Brand = p.Brand
	product.Stock = p.Stock

	err = s.repo.Update(ctx, product)
	if err != nil {
		return err
	}

	if oldStock != product.Stock {
		seller, err := s.sellerRepo.FindByUUID(ctx, product.SellerUUID)

		if err != nil {
			return err
		}
		s.emailProvider.StockChanged(oldStock, product.Stock, seller.Email)
	}
	return s.repo.Update(ctx, product)
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
