package graphql

import (
	"cart/entity"
	"cart/service"
	"context"
)

type resolver struct {
	cartService    service.CartService
	productService service.ProductService
}

func NewResolver(cartService service.CartService, productService service.ProductService) *resolver {
	return &resolver{cartService: cartService, productService: productService}
}

func GetProductsMap(products []entity.Product) map[string]entity.Product {
	productsMap := map[string]entity.Product{}
	for _, p := range products {
		productsMap[p.ID] = p
	}
	return productsMap
}

func (r resolver) MyCart(ctx context.Context) (*CartResolver, error) {
	cartItems, err := r.cartService.MyCart(ctx)
	if err != nil {
		return nil, err
	}
	productIDs := []string{}
	for _, c := range *cartItems {
		productIDs = append(productIDs, c.ProductID)
	}
	products, err := r.productService.GetProducts(ctx, productIDs)
	if err != nil {
		return nil, err
	}
	productsMap := GetProductsMap(*products)
	return &CartResolver{cartItems: *cartItems, productsMap: productsMap}, nil
}

func (r resolver) UpsertCartItem(ctx context.Context, args struct {
	ProductID string
	Quantity  int32
}) (*bool, error) {
	if err := r.cartService.UpsertCartItem(ctx, args.ProductID, args.Quantity); err != nil {
		ok := false
		return &ok, err
	}
	ok := true
	return &ok, nil

}
