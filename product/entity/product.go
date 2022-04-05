package entity

import "github.com/google/uuid"

type Product struct {
	Base
	UserID string `gorm:"type:VARCHAR(36); NOT NULL"`
	Title  string `gorm:"type:TEXT; NOT NULL; index:,class:FULLTEXT"`
	Price  int32  `gorm:"type:INT; NOT NULL"`
}

func NewProduct(userID, title string, price int32) (*Product, error) {
	product := Product{
		UserID: userID,
		Title:  title,
		Price:  price,
	}
	productID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	product.ID = productID.String()
	return &product, nil
}
