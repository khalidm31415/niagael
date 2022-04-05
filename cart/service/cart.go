package service

import (
	"cart/entity"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CartService interface {
	UpsertCartItem(ctx context.Context, productID string, quantity int32) error
	MyCart(ctx context.Context) (*[]entity.CartItem, error)
	ClearCart(ctx context.Context) error
}

type cartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) CartService {
	return cartService{db: db}
}

func (c cartService) UpsertCartItem(ctx context.Context, productID string, quantity int32) error {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return errors.New("unauthorized")
	}

	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	if quantity == 0 {
		if err := c.db.Where("user_id = ?", userID).Where("product_id = ?", productID).Delete(&entity.CartItem{}).Error; err != nil {
			zap.S().Error(err)
			return err
		}
		return nil
	}

	res := c.db.Model(&entity.CartItem{}).Where("user_id = ?", userID).Where("product_id = ?", productID).Update("quantity", quantity)
	if res.Error != nil {
		zap.S().Error(res.Error)
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}

	cartItem := entity.CartItem{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}
	db := c.db.Session(&gorm.Session{NewDB: true})
	if err := db.Create(&cartItem).Error; err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}

func (c cartService) MyCart(ctx context.Context) (*[]entity.CartItem, error) {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	var cartItems []entity.CartItem
	if err := c.db.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		zap.S().Error(err)
		return nil, err
	}

	return &cartItems, nil
}

func (c cartService) ClearCart(ctx context.Context) error {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return errors.New("unauthorized")
	}

	if err := c.db.Where("user_id = ?", userID).Delete(&entity.CartItem{}).Error; err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}
