package service

import (
	"context"
	"errors"
	"product/entity"

	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(ctx context.Context, title string, price int32) error
	UpdateProduct(ctx context.Context, id string, title *string, price *int32) error
	DeleteProduct(ctx context.Context, id string) error
	SearchProducts(ctx context.Context, query string) (*[]entity.Product, error)
	GetProducts(ctx context.Context, ids []string) (*[]entity.Product, error)
}

type productService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) ProductService {
	return productService{db: db}
}

func (p productService) CreateProduct(ctx context.Context, title string, price int32) error {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return errors.New("unauthorized")
	}
	product, err := entity.NewProduct(userID, title, price)
	if err != nil {
		return err
	}
	if err := p.db.Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func (p productService) UpdateProduct(ctx context.Context, id string, title *string, price *int32) error {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return errors.New("unauthorized")
	}
	var productUserID string
	if err := p.db.Table("products").Select("user_id").Where("id = ?", id).Scan(&productUserID).Error; err != nil {
		return err
	}
	if userID != productUserID {
		return errors.New("cannot update other user's product")
	}

	db := p.db.Session(&gorm.Session{NewDB: true})
	updates := entity.Product{}
	if title != nil {
		updates.Title = *title
	}
	if price != nil {
		updates.Price = *price
	}
	if err := db.Model(entity.Product{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (p productService) DeleteProduct(ctx context.Context, id string) error {
	userID := ctx.Value("userID").(string)
	if userID == "" {
		return errors.New("unauthorized")
	}
	var productUserID string
	if err := p.db.Table("products").Select("user_id").Where("id = ?", id).Scan(&productUserID).Error; err != nil {
		return err
	}
	if userID != productUserID {
		return errors.New("cannot delete other user's product")
	}

	db := p.db.Session(&gorm.Session{NewDB: true})
	if err := db.Where("id = ?", id).Delete(&entity.Product{}).Error; err != nil {
		return err
	}
	return nil
}

func (p productService) SearchProducts(ctx context.Context, query string) (*[]entity.Product, error) {
	var products []entity.Product
	if err := p.db.Where("MATCH(title) AGAINST(? IN NATURAL LANGUAGE MODE)", query).Find(&products).Error; err != nil {
		return nil, err
	}
	return &products, nil
}

func (p productService) GetProducts(ctx context.Context, ids []string) (*[]entity.Product, error) {
	var products []entity.Product
	if err := p.db.Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}
	return &products, nil
}
