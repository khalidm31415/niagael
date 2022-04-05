package graphql

import (
	"context"
	"product/service"
)

type resolver struct {
	productService service.ProductService
}

func NewResolver(productService service.ProductService) *resolver {
	return &resolver{productService: productService}
}

func (r resolver) CreateProduct(ctx context.Context, args struct {
	Title string
	Price int32
}) (*bool, error) {
	if err := r.productService.CreateProduct(ctx, args.Title, args.Price); err != nil {
		ok := false
		return &ok, err
	}
	ok := true
	return &ok, nil
}

func (r resolver) UpdateProduct(ctx context.Context, args struct {
	ID    string
	Title *string
	Price *int32
}) (*bool, error) {
	if err := r.productService.UpdateProduct(ctx, args.ID, args.Title, args.Price); err != nil {
		ok := false
		return &ok, err
	}
	ok := true
	return &ok, nil

}

func (r resolver) DeleteProduct(ctx context.Context, args struct {
	ID string
}) (*bool, error) {
	if err := r.productService.DeleteProduct(ctx, args.ID); err != nil {
		ok := false
		return &ok, err
	}
	ok := true
	return &ok, nil
}

func (r resolver) SearchProducts(ctx context.Context, args struct {
	Query string
}) (*[]ProductResolver, error) {
	products, err := r.productService.SearchProducts(ctx, args.Query)
	if err != nil {
		return nil, err
	}
	productResolvers := []ProductResolver{}
	for _, p := range *products {
		productResolvers = append(productResolvers, ProductResolver{product: p})
	}
	return &productResolvers, nil
}
