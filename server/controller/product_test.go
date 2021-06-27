package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"coding-challenge-go/pkg/product"
)

type (
	input struct {
		productUUID string
	}

	handlerErr struct {
		E string `json:"error"`
	}
)

func setupRouter(path string, handlerFunc func(c *gin.Context)) *gin.Engine {
	r := gin.Default()
	r.GET(path, handlerFunc)
	return r
}

func Test_GetProductV1(t *testing.T) {

	tests := []struct {
		name             string
		input            input
		expected         *ProductResponseV1
		statusCode       int
		err              *handlerErr
		DoGetProductFunc func(uuid string) (*product.ProductInfo, error)
	}{
		{
			name: "test get product success",
			input: input{
				productUUID: "e6461ea4-d698-11eb-890b-0242ac1a0002",
			},
			DoGetProductFunc: func(uuid string) (*product.ProductInfo, error) {
				return &product.ProductInfo{
					Product: &product.Product{
						ProductID:  1,
						UUID:       uuid,
						Name:       "product1",
						Brand:      "GFG",
						Stock:      1,
						SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
					},
					Seller: &product.SellerInfo{
						UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
						Links: &product.SellerLinks{
							Self: &product.SelfSellerLink{
								Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
							},
						},
					},
				}, nil
			},
			statusCode: 200,
			expected: &ProductResponseV1{
				ProductID:  1,
				UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0002",
				Name:       "product1",
				Brand:      "GFG",
				Stock:      1,
				SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
			},
		},
		{
			name: "test get product is not found",
			input: input{
				productUUID: "e6461ea4-d698-11eb-890b-0242ac1a0002",
			},
			DoGetProductFunc: func(uuid string) (*product.ProductInfo, error) {
				return nil, product.NewProductNotFoundError(uuid)
			},
			statusCode: 404,
			expected:   nil,
			err: &handlerErr{
				E: product.NewProductNotFoundError("e6461ea4-d698-11eb-890b-0242ac1a0002").Error(),
			},
		},
		{
			name: "test get product server error",
			input: input{
				productUUID: "e6461ea4-d698-11eb-890b-0242ac1a0002",
			},
			DoGetProductFunc: func(uuid string) (*product.ProductInfo, error) {
				return nil, fmt.Errorf("Something was wrong")
			},
			statusCode: 500,
			expected:   nil,
			err: &handlerErr{
				E: "Fail to query product by uuid",
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			service := &productServiceMock{
				DoGetProductFunc: test.DoGetProductFunc,
			}
			productController := NewProductController(service)
			router := setupRouter("/api/v1/product", productController.Get)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/product?id=%s", test.input.productUUID), nil)
			router.ServeHTTP(w, req)

			if w.Code == 200 {
				b, err := json.Marshal(test.expected)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, string(b), w.Body.String())
			} else {
				b, e := json.Marshal(test.err)
				if e != nil {
					t.Fatal(e)
				}
				assert.Equal(t, string(b), w.Body.String())

			}
			assert.Equal(t, test.statusCode, w.Code)
		})
	}

}

func Test_GetProductV2(t *testing.T) {

	tests := []struct {
		name             string
		input            input
		expected         *product.ProductInfo
		statusCode       int
		err              *handlerErr
		DoGetProductFunc func(uuid string) (*product.ProductInfo, error)
	}{
		{
			name: "test get product success",
			input: input{
				productUUID: "e6461ea4-d698-11eb-890b-0242ac1a0002",
			},
			DoGetProductFunc: func(uuid string) (*product.ProductInfo, error) {
				return &product.ProductInfo{
					Product: &product.Product{
						ProductID:  1,
						UUID:       uuid,
						Name:       "product1",
						Brand:      "GFG",
						Stock:      1,
						SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
					},
					Seller: &product.SellerInfo{
						UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
						Links: &product.SellerLinks{
							Self: &product.SelfSellerLink{
								Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
							},
						},
					},
				}, nil
			},
			statusCode: 200,
			expected: &product.ProductInfo{
				Product: &product.Product{
					ProductID:  1,
					UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0002",
					Name:       "product1",
					Brand:      "GFG",
					Stock:      1,
					SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
				},
				Seller: &product.SellerInfo{
					UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
					Links: &product.SellerLinks{
						Self: &product.SelfSellerLink{
							Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
						},
					},
				},
			},
		},
		{
			name: "test get product is not found",
			input: input{
				productUUID: "e6461ea4-d698-11eb-890b-0242ac1a0002",
			},
			DoGetProductFunc: func(uuid string) (*product.ProductInfo, error) {
				return nil, product.NewProductNotFoundError(uuid)
			},
			statusCode: 404,
			expected:   nil,
			err: &handlerErr{
				E: product.NewProductNotFoundError("e6461ea4-d698-11eb-890b-0242ac1a0002").Error(),
			},
		},
		{
			name: "test get product server error",
			input: input{
				productUUID: "e6461ea4-d698-11eb-890b-0242ac1a0002",
			},
			DoGetProductFunc: func(uuid string) (*product.ProductInfo, error) {
				return nil, fmt.Errorf("Something was wrong")
			},
			statusCode: 500,
			expected:   nil,
			err: &handlerErr{
				E: "Fail to query product by uuid",
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			service := &productServiceMock{
				DoGetProductFunc: test.DoGetProductFunc,
			}
			productController := NewProductController(service)
			router := setupRouter("/api/v2/product", productController.GetV2)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v2/product?id=%s", test.input.productUUID), nil)
			router.ServeHTTP(w, req)

			if w.Code == 200 {
				b, err := json.Marshal(test.expected)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, string(b), w.Body.String())
			} else {
				b, e := json.Marshal(test.err)
				if e != nil {
					t.Fatal(e)
				}
				assert.Equal(t, string(b), w.Body.String())

			}
			assert.Equal(t, test.statusCode, w.Code)
		})
	}

}

func Test_ProductListV1(t *testing.T) {

	tests := []struct {
		name               string
		page               int
		expected           []*ProductResponseV1
		statusCode         int
		err                *handlerErr
		DoListProductsFunc func() ([]*product.ProductInfo, error)
	}{
		{
			name: "test get list product success",
			page: 1,
			DoListProductsFunc: func() ([]*product.ProductInfo, error) {
				return []*product.ProductInfo{
					{

						Product: &product.Product{
							ProductID:  1,
							UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0002",
							Name:       "product1",
							Brand:      "GFG",
							Stock:      1,
							SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
						},
						Seller: &product.SellerInfo{
							UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
							Links: &product.SellerLinks{
								Self: &product.SelfSellerLink{
									Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
								},
							},
						},
					},
					{

						Product: &product.Product{
							ProductID:  2,
							UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0001",
							Name:       "product2",
							Brand:      "GFG",
							Stock:      1,
							SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
						},
						Seller: &product.SellerInfo{
							UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
							Links: &product.SellerLinks{
								Self: &product.SelfSellerLink{
									Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
								},
							},
						},
					},
				}, nil
			},
			statusCode: 200,
			expected: []*ProductResponseV1{
				&ProductResponseV1{
					ProductID:  1,
					UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0002",
					Name:       "product1",
					Brand:      "GFG",
					Stock:      1,
					SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
				},
				&ProductResponseV1{
					ProductID:  2,
					UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0001",
					Name:       "product2",
					Brand:      "GFG",
					Stock:      1,
					SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
				},
			},
		},
		{
			name: "test get product server error",

			DoListProductsFunc: func() ([]*product.ProductInfo, error) {
				return nil, fmt.Errorf("Something was wrong")
			},
			statusCode: 500,
			expected:   nil,
			err: &handlerErr{
				E: "Fail to query product list",
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			service := &productServiceMock{
				DoListProductsFunc: test.DoListProductsFunc,
			}
			productController := NewProductController(service)
			router := setupRouter("/api/v1/products", productController.List)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/products", nil)
			router.ServeHTTP(w, req)

			if w.Code == 200 {
				b, err := json.Marshal(test.expected)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, string(b), w.Body.String())
			} else {
				b, e := json.Marshal(test.err)
				if e != nil {
					t.Fatal(e)
				}
				assert.Equal(t, string(b), w.Body.String())

			}
			assert.Equal(t, test.statusCode, w.Code)
		})
	}

}

func Test_ProductListV2(t *testing.T) {

	tests := []struct {
		name               string
		page               int
		expected           []*product.ProductInfo
		statusCode         int
		err                *handlerErr
		DoListProductsFunc func() ([]*product.ProductInfo, error)
	}{
		{
			name: "test get list product success",
			page: 1,
			DoListProductsFunc: func() ([]*product.ProductInfo, error) {
				return []*product.ProductInfo{
					{

						Product: &product.Product{
							ProductID:  1,
							UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0002",
							Name:       "product1",
							Brand:      "GFG",
							Stock:      1,
							SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
						},
						Seller: &product.SellerInfo{
							UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
							Links: &product.SellerLinks{
								Self: &product.SelfSellerLink{
									Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
								},
							},
						},
					},
					{

						Product: &product.Product{
							ProductID:  2,
							UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0001",
							Name:       "product2",
							Brand:      "GFG",
							Stock:      1,
							SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
						},
						Seller: &product.SellerInfo{
							UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
							Links: &product.SellerLinks{
								Self: &product.SelfSellerLink{
									Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
								},
							},
						},
					},
				}, nil
			},
			statusCode: 200,
			expected: []*product.ProductInfo{
				{

					Product: &product.Product{
						ProductID:  1,
						UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0002",
						Name:       "product1",
						Brand:      "GFG",
						Stock:      1,
						SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
					},
					Seller: &product.SellerInfo{
						UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
						Links: &product.SellerLinks{
							Self: &product.SelfSellerLink{
								Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
							},
						},
					},
				},
				{

					Product: &product.Product{
						ProductID:  2,
						UUID:       "e6461ea4-d698-11eb-890b-0242ac1a0001",
						Name:       "product2",
						Brand:      "GFG",
						Stock:      1,
						SellerUUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
					},
					Seller: &product.SellerInfo{
						UUID: "e6461ea4-d698-11eb-890b-0242ac1a0003",
						Links: &product.SellerLinks{
							Self: &product.SelfSellerLink{
								Href: "http://localhost:8080/api/v1/sellers/e6461ea4-d698-11eb-890b-0242ac1a0003",
							},
						},
					},
				},
			},
		},
		{
			name: "test get product server error",

			DoListProductsFunc: func() ([]*product.ProductInfo, error) {
				return nil, fmt.Errorf("Something was wrong")
			},
			statusCode: 500,
			expected:   nil,
			err: &handlerErr{
				E: "Fail to query product list",
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			service := &productServiceMock{
				DoListProductsFunc: test.DoListProductsFunc,
			}
			productController := NewProductController(service)
			router := setupRouter("/api/v2/products", productController.ListV2)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v2/products", nil)
			router.ServeHTTP(w, req)

			if w.Code == 200 {
				b, err := json.Marshal(test.expected)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, string(b), w.Body.String())
			} else {
				b, e := json.Marshal(test.err)
				if e != nil {
					t.Fatal(e)
				}
				assert.Equal(t, string(b), w.Body.String())

			}
			assert.Equal(t, test.statusCode, w.Code)
		})
	}

}

type productServiceMock struct {
	DoGetProductFunc   func(uuid string) (*product.ProductInfo, error)
	DoListProductsFunc func() ([]*product.ProductInfo, error)
}

func (m *productServiceMock) FindByUUID(ctx context.Context, uuid string) (*product.ProductInfo, error) {
	return m.DoGetProductFunc(uuid)
}

func (m *productServiceMock) Create(ctx context.Context, p *product.Product) error {
	return nil
}

func (m *productServiceMock) Update(ctx context.Context, p *product.Product) error {
	return nil
}

func (m *productServiceMock) Delete(ctx context.Context, uuid string) error {
	return nil
}

func (m *productServiceMock) List(ctx context.Context, params *product.FilterParams) ([]*product.ProductInfo, error) {
	return m.DoListProductsFunc()
}
