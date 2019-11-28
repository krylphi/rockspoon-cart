package domain

import (
	"errors"
)

var (
	// ErrCartNotFound is returned, when cart is not found.
	ErrCartNotFound error = errors.New("no such cart")
	// ErrNoSuchCartItem is returned when cart has no requested item
	ErrNoSuchCartItem error = errors.New("this cart doesn't contain such item")
	// ErrZeroQuantity is returned when quantity is 0 for new CartItem
	ErrZeroQuantity error = errors.New("product quantity can not be 0")
	// ErrEmptyProduct is returned when Product field of CartItem is not specified
	ErrEmptyProduct error = errors.New("product can not be empty")
)
