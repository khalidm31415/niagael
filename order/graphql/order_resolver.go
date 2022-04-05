package graphql

import "order/entity"

type OrderResolver struct {
	order entity.Order
}

func (o OrderResolver) ID() *string {
	return &o.order.ID
}

func (o OrderResolver) OrderItems() *[]OrderItemResolver {
	orderItemResolvers := []OrderItemResolver{}
	for _, o := range o.order.OrderItems {
		orderItemResolver := OrderItemResolver{orderItem: o}
		orderItemResolvers = append(orderItemResolvers, orderItemResolver)
	}
	return &orderItemResolvers
}

func (o OrderResolver) TotalPrice() *int32 {
	totalPrice := int32(0)
	for _, o := range o.order.OrderItems {
		totalPrice += (o.Quantity * o.UnitPrice)
	}
	return &totalPrice
}

func (o OrderResolver) Status() *string {
	return &o.order.Status
}
