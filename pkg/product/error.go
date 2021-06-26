package product

import "fmt"

type ProductNotFoundError struct {
	id string
}

func (e ProductNotFoundError) Error() string {
	return fmt.Sprintf("Product is not found with id=%s", e.id)
}

type SellerNotFoundError struct {
	id string
}

func (e SellerNotFoundError) Error() string {
	return fmt.Sprintf("Seller is not found with id=%s", e.id)
}
