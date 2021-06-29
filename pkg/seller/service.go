package seller

import (
	"context"
)

type (
	Service interface {
		List(ctx context.Context) ([]*Seller, error)
		Top10ByProduct(ctx context.Context) ([]*Seller, error)
	}

	service struct {
		repo Repository
	}

	Repository interface {
		List(ctx context.Context) ([]*Seller, error)
		FindByUUID(ctx context.Context, uuid string) (*Seller, error)
		TopByProduct(ctx context.Context, limit int) ([]*Seller, error)
	}
)

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) List(ctx context.Context) ([]*Seller, error) {
	return s.repo.List(ctx)
}

func (s *service) Top10ByProduct(ctx context.Context) ([]*Seller, error) {
	return s.repo.TopByProduct(ctx, 10)
}
