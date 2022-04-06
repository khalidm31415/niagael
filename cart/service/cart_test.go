package service

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestCartService(t *testing.T) {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	cartService := NewCartService(db)

	ctx := context.WithValue(context.Background(), "userID", "")
	err = cartService.UpsertCartItem(ctx, "test-product-id", 1)
	if err == nil {
		t.Errorf("expected unauthorized error if userID is empty string, got %s", err)
	}

	ctx = context.WithValue(context.Background(), "userID", "test-user-id")
	err = cartService.UpsertCartItem(ctx, "test-product-id", -1)
	if err == nil {
		t.Errorf("expected an error if quantity is negative, got %s", err)
	}
}
