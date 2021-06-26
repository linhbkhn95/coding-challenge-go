package seller

import (
	"context"
	"database/sql"
)

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

type repository struct {
	db *sql.DB
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller WHERE uuid = ?", uuid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	seller := &Seller{}

	err = rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)

	if err != nil {
		return nil, err
	}

	return seller, nil
}

func (r *repository) List(ctx context.Context) ([]*Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sellers []*Seller

	for rows.Next() {
		seller := &Seller{}

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}
