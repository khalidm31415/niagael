package service

import (
	"context"
	"errors"
	"order/entity"

	"github.com/machinebox/graphql"
	"go.uber.org/zap"
)

type CartService interface {
	MyCart(ctx context.Context) (*entity.Cart, error)
	ClearCart(ctx context.Context) error
}

type cartService struct {
	client *graphql.Client
}

func NewCartService(url string) CartService {
	client := graphql.NewClient(url)
	// client.Log = func(s string) { log.Println(s) }
	return cartService{client: client}
}

type MyCartResponse struct {
	Cart entity.Cart `json:"myCart"`
}

func (c cartService) MyCart(ctx context.Context) (*entity.Cart, error) {
	req := graphql.NewRequest(`
		query myCart {
			myCart {
				cartItems {
					productId
					productTitle
					unitPrice
					quantity
				}
				totalPrice
			}
		}
	`)

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return nil, errors.New("unauthorized")
	}
	req.Header.Set("X-user-id", userID)

	var respData MyCartResponse
	if err := c.client.Run(ctx, req, &respData); err != nil {
		zap.S().Error(err)
		return nil, err
	}
	cart := respData.Cart
	return &cart, nil
}

type ClearCartResponse struct {
	OK bool `json:"clearCart"`
}

func (c cartService) ClearCart(ctx context.Context) error {
	req := graphql.NewRequest(`
		mutation clearCart {
			clearCart
		}
	`)

	userID := ctx.Value("userID").(string)
	if userID == "" {
		return errors.New("unauthorized")
	}
	req.Header.Set("X-user-id", userID)

	var respData ClearCartResponse
	if err := c.client.Run(ctx, req, &respData); err != nil {
		zap.S().Error(err)
		return err
	}
	return nil
}
