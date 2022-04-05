package service

import (
	"context"
	"errors"
	"order/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderService interface {
	MyOrders(ctx context.Context, ids *[]string) (*[]entity.Order, error)
	CheckoutCart(ctx context.Context) (*entity.Order, error)
	PayOrder(ctx context.Context, orderID string) error
	CancelOrder(ctx context.Context, orderID string) error
}

type orderService struct {
	db          *gorm.DB
	cartService CartService
}

func NewOrderService(db *gorm.DB, cartService CartService) OrderService {
	return orderService{db: db, cartService: cartService}
}

func (o orderService) MyOrders(ctx context.Context, ids *[]string) (*[]entity.Order, error) {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	db := o.db.Preload("OrderItems")
	db = db.Where("user_id = ?", userID)
	if ids != nil {
		db = db.Where("id IN ?", *ids)
	}
	var orders []entity.Order
	if err := db.Find(&orders).Error; err != nil {
		zap.S().Error(err)
		return nil, err
	}
	return &orders, nil
}

func (o orderService) CheckoutCart(ctx context.Context) (*entity.Order, error) {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	cart, err := o.cartService.MyCart(ctx)
	if err != nil {
		zap.S().Error(err)
		return nil, err
	}
	order, err := entity.NewOrder(userID)
	if err != nil {
		return nil, err
	}
	orderItems := []entity.OrderItem{}
	for _, c := range cart.CartItems {
		orderItem, err := entity.NewOrderItem(order.ID, c)
		if err != nil {
			return nil, err
		}
		orderItems = append(orderItems, *orderItem)
	}

	err = o.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			zap.S().Error(err)
			return err
		}
		if err := tx.Create(&orderItems).Error; err != nil {
			zap.S().Error(err)
			return err
		}
		if err := o.cartService.ClearCart(ctx); err != nil {
			zap.S().Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	order.OrderItems = orderItems
	return order, nil
}

func (o orderService) PayOrder(ctx context.Context, orderID string) error {
	var orderStatus string
	if err := o.db.Table("orders").Select("status").Where("id = ?", orderID).Scan(&orderStatus).Error; err != nil {
		return err
	}
	if orderStatus != "WAITING_FOR_PAYMENT" {
		return errors.New("cannot pay settled order")
	}
	if err := o.db.Model(&entity.Order{}).Where("id = ?", orderID).Update("status", "PAYMENT_SUCCESS").Error; err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

func (o orderService) CancelOrder(ctx context.Context, orderID string) error {
	var orderStatus string
	if err := o.db.Table("orders").Select("status").Where("id = ?", orderID).Scan(&orderStatus).Error; err != nil {
		return err
	}
	if orderStatus != "WAITING_FOR_PAYMENT" {
		return errors.New("cannot cancel settled order")
	}
	if err := o.db.Model(&entity.Order{}).Where("id = ?", orderID).Update("status", "CANCELLED").Error; err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}
