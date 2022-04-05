package service

import (
	"cart/entity"
	"context"

	"github.com/machinebox/graphql"
	"go.uber.org/zap"
)

type ProductService interface {
	GetProducts(ctx context.Context, ids []string) (*[]entity.Product, error)
}

type productService struct {
	client *graphql.Client
}

func NewProductService(url string) ProductService {
	client := graphql.NewClient(url)
	// client.Log = func(s string) { log.Println(s) }
	return productService{client: client}
}

type GetProductsResponse struct {
	Products []entity.Product `json:"getProducts"`
}

func (p productService) GetProducts(ctx context.Context, ids []string) (*[]entity.Product, error) {
	req := graphql.NewRequest(`
		query getProducts($ids: [String!]!) {
			getProducts(ids: $ids) {
				id
				title
				price
			}
		}
	`)

	req.Var("ids", ids)

	var respData GetProductsResponse
	if err := p.client.Run(ctx, req, &respData); err != nil {
		zap.S().Error(err)
		return nil, err
	}
	products := respData.Products
	return &products, nil
}
