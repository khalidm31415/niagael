package graphql

import "cart/entity"

type CartItemResolver struct {
	cartItem entity.CartItem
	product  entity.Product
}

func (c CartItemResolver) ProductID() *string {
	return &c.cartItem.ProductID
}

func (c CartItemResolver) ProductTitle() *string {
	return &c.product.Title
}

func (c CartItemResolver) UnitPrice() *int32 {
	return &c.product.Price
}

func (c CartItemResolver) Quantity() *int32 {
	return &c.cartItem.Quantity
}
