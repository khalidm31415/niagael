package graphql

import "order/entity"

type OrderItemResolver struct {
	orderItem entity.OrderItem
}

func (o OrderItemResolver) ID() *string {
	return &o.orderItem.ID
}

func (o OrderItemResolver) ProductID() *string {
	return &o.orderItem.ProductID
}

func (o OrderItemResolver) ProductTitle() *string {
	return &o.orderItem.ProductTitle
}

func (o OrderItemResolver) UnitPrice() *int32 {
	return &o.orderItem.UnitPrice
}

func (o OrderItemResolver) Quantity() *int32 {
	return &o.orderItem.Quantity
}
