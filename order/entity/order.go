package entity

import "github.com/google/uuid"

type OrderItem struct {
	ID           string `gorm:"type:VARCHAR(36)"`
	OrderID      string `gorm:"type:VARCHAR(36); NOT NULL"`
	Order        Order
	ProductID    string `gorm:"type:VARCHAR(36); NOT NULL"`
	ProductTitle string `gorm:"type:TEXT; NOT NULL"`
	UnitPrice    int32  `gorm:"type:INT; NOT NULL"`
	Quantity     int32  `gorm:"type:INT; NOT NULL"`
}

func NewOrderItem(orderID string, cartItem CartItem) (*OrderItem, error) {
	orderItem := OrderItem{
		OrderID:      orderID,
		ProductID:    cartItem.ProductID,
		ProductTitle: cartItem.ProductTitle,
		UnitPrice:    cartItem.UnitPrice,
		Quantity:     cartItem.Quantity,
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	orderItem.ID = id.String()
	return &orderItem, nil
}

type Order struct {
	Base
	UserID     string `gorm:"type:VARCHAR(36); NOT NULL"`
	OrderItems []OrderItem
	Status     string `gorm:"type:VARCHAR(20); NOT NULL"`
}

func NewOrder(userID string) (*Order, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	order := Order{
		UserID: userID,
		Status: "WAITING_FOR_PAYMENT",
	}
	order.ID = id.String()
	return &order, nil
}
