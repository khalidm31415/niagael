package entity

type CartItem struct {
	UserID    string `gorm:"type:VARCHAR(36); NOT NULL; PRIMARY_KEY"`
	ProductID string `gorm:"type:VARCHAR(36); NOT NULL; PRIMARY_KEY"`
	Quantity  int32  `gorm:"type:INT; NOT NULL"`
}
