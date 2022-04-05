package graphql

import "product/entity"

type ProductResolver struct {
	product entity.Product
}

func (p ProductResolver) ID() *string {
	return &p.product.ID
}

func (p ProductResolver) UserID() *string {
	return &p.product.UserID
}

func (p ProductResolver) Title() *string {
	return &p.product.Title
}

func (p ProductResolver) Price() *int32 {
	price := int32(p.product.Price)
	return &price
}
