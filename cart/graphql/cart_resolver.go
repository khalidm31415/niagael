package graphql

import "cart/entity"

type CartResolver struct {
	cartItems   []entity.CartItem
	productsMap map[string]entity.Product
}

func (c CartResolver) CartItems() *[]CartItemResolver {
	cartItemResolvers := []CartItemResolver{}
	for _, cartItem := range c.cartItems {
		product := c.productsMap[cartItem.ProductID]
		cartItemResolver := CartItemResolver{
			cartItem: cartItem,
			product:  product,
		}
		cartItemResolvers = append(cartItemResolvers, cartItemResolver)
	}
	return &cartItemResolvers
}

func (c CartResolver) TotalPrice() *int32 {
	totalPrice := int32(0)
	for _, cartItem := range c.cartItems {
		product := c.productsMap[cartItem.ProductID]
		totalPrice += (cartItem.Quantity * product.Price)
	}
	return &totalPrice
}
