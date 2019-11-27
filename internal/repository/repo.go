package repository

import (
	"context"

	"github.com/Krylphi/rockspoon-cart/internal/domain"
)

// CartWriteRepository used for wring data to DB
type CartWriteRepository interface {
	CreateCart(ctx context.Context) (*domain.Cart, error)
	AddItem(ctx context.Context, cartID string, product string, quantity int) (*domain.CartItem, error)
	RemoveItem(ctx context.Context, cartID string, itemID string) error
	DeleteCart(ctx context.Context, id string) error
}

//CartReadRepository used to read data from db
type CartReadRepository interface {
	Cart(ctx context.Context, id string) (*domain.Cart, error)
}

//CartRepository is repository capable of reading and writing
type CartRepository interface {
	CartWriteRepository
	CartReadRepository
}
