package entity

type CartItem struct {
	ProductID    string `json:"productId"`
	ProductTitle string `json:"productTitle"`
	UnitPrice    int32  `json:"unitPrice"`
	Quantity     int32  `json:"quantity"`
}

type Cart struct {
	CartItems  []CartItem `json:"cartItems"`
	TotalPrice int32      `json:"totalPrice"`
}
