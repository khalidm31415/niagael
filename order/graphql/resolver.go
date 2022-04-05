package graphql

import (
	"context"
	"order/service"
)

type resolver struct {
	orderService service.OrderService
}

func NewResolver(orderService service.OrderService) *resolver {
	return &resolver{orderService: orderService}
}

func (r resolver) MyOrders(ctx context.Context, args struct {
	IDs *[]string
}) (*[]OrderResolver, error) {
	orders, err := r.orderService.MyOrders(ctx, args.IDs)
	if err != nil {
		return nil, err
	}
	orderResolvers := []OrderResolver{}
	for _, o := range *orders {
		orderResolver := OrderResolver{order: o}
		orderResolvers = append(orderResolvers, orderResolver)
	}
	return &orderResolvers, nil
}

func (r resolver) CheckoutCart(ctx context.Context) (*OrderResolver, error) {
	order, err := r.orderService.CheckoutCart(ctx)
	if err != nil {
		return nil, err
	}
	return &OrderResolver{order: *order}, nil
}

func (r resolver) PayOrder(ctx context.Context, args struct {
	OrderID string
}) (*bool, error) {
	if err := r.orderService.PayOrder(ctx, args.OrderID); err != nil {
		ok := false
		return &ok, err
	}
	ok := true
	return &ok, nil
}

func (r resolver) CancelOrder(ctx context.Context, args struct {
	OrderID string
}) (*bool, error) {
	if err := r.orderService.CancelOrder(ctx, args.OrderID); err != nil {
		ok := false
		return &ok, err
	}
	ok := true
	return &ok, nil
}
