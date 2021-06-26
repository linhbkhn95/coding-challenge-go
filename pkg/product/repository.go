package product

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

func (r *repository) Delete(ctx context.Context, product *Product) error {
	rows, err := r.db.Query("DELETE FROM product WHERE uuid = ?", product.UUID)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (r *repository) Create(ctx context.Context, product *Product) error {
	rows, err := r.db.Query(
		"INSERT INTO product (id_product, name, brand, stock, fk_seller, uuid) VALUES(?,?,?,?,(SELECT id_seller FROM seller WHERE uuid = ?),?)",
		product.ProductID, product.Name, product.Brand, product.Stock, product.SellerUUID, product.UUID,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (r *repository) Update(ctx context.Context, product *Product) error {
	rows, err := r.db.Query(
		"UPDATE product SET name = ?, brand = ?, stock = ? WHERE uuid = ?",
		product.Name, product.Brand, product.Stock, product.UUID,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (r *repository) List(ctx context.Context, offset int, limit int) ([]*Product, error) {
	rows, err := r.db.Query(
		"SELECT p.id_product, p.name, p.brand, p.stock, s.uuid, p.uuid FROM product p "+
			"INNER JOIN seller s ON(s.id_seller = p.fk_seller) LIMIT ? OFFSET ?",
		limit, offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*Product

	for rows.Next() {
		product := &Product{}

		err = rows.Scan(&product.ProductID, &product.Name, &product.Brand, &product.Stock, &product.SellerUUID, &product.UUID)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *repository) FindByUUID(ctx context.Context, uuid string) (*Product, error) {
	rows, err := r.db.Query(
		"SELECT p.id_product, p.name, p.brand, p.stock, s.uuid, p.uuid FROM product p "+
			"INNER JOIN seller s ON(s.id_seller = p.fk_seller) WHERE p.uuid = ?",
		uuid,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	product := &Product{}
	err = rows.Scan(&product.ProductID, &product.Name, &product.Brand, &product.Stock, &product.SellerUUID, &product.UUID)

	if err != nil {
		return nil, err
	}

	return product, nil
}
