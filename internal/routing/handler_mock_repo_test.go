package routing

import (
	"context"
	"rockspoon-cart/internal/domain"
	"rockspoon-cart/internal/repository"
	"strconv"
)

type (
	cartMockRepository struct {
		carts map[string]*domain.Cart
	}
)

func InitMockRepo() repository.CartRepository {
	res := &cartMockRepository{
		carts: make(map[string]*domain.Cart, 0),
	}

	ctx := context.Background()

	cart, _ := res.CreateCart(ctx)
	_, _ = res.AddItem(ctx, cart.ID, "Product1", 9)
	_, _ = res.AddItem(ctx, cart.ID, "Product2", 8)
	_, _ = res.AddItem(ctx, cart.ID, "Product3", 7)
	_, _ = res.AddItem(ctx, cart.ID, "Product4", 6)
	_, _ = res.AddItem(ctx, cart.ID, "Product5", 5)
	_, _ = res.CreateCart(ctx)
	_, _ = res.CreateCart(ctx)
	return res
}

func (c *cartMockRepository) CreateCart(ctx context.Context) (*domain.Cart, error) {
	id := strconv.Itoa(len(c.carts) + 1)
	cart := &domain.Cart{
		ID:    id,
		Items: make([]*domain.CartItem, 0),
	}

	c.carts[id] = cart
	return cart, nil
}

func (c *cartMockRepository) AddItem(ctx context.Context, cartID string, product string, quantity int) (*domain.CartItem, error) {
	cart, ok := c.carts[cartID]
	if !ok {
		return nil, domain.ErrCartNotFound
	}

	item, err := cart.AddItem(strconv.Itoa(len(cart.Items)+1), product, quantity)

	if err != nil {
		return nil, err
	}

	c.carts[cartID] = cart

	return item, nil
}

func (c *cartMockRepository) RemoveItem(ctx context.Context, cartID string, itemID string) error {
	cart, ok := c.carts[cartID]
	if !ok {
		return domain.ErrCartNotFound
	}

	err := cart.RemoveItem(itemID)

	if err != nil {
		return err
	}

	c.carts[cartID] = cart

	return nil
}

func (c *cartMockRepository) DeleteCart(ctx context.Context, id string) error {
	_, ok := c.carts[id]
	if !ok {
		return domain.ErrCartNotFound
	}

	delete(c.carts, id)

	return nil
}

func (c *cartMockRepository) Cart(ctx context.Context, id string) (*domain.Cart, error) {
	cart, ok := c.carts[id]
	if !ok {
		return nil, domain.ErrCartNotFound
	}

	return cart, nil
}
