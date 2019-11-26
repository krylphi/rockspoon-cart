package domain

import (
	"errors"
)

var (
	ErrCartNotFound   error = errors.New("no such cart")
	ErrNoSuchCartItem error = errors.New("this cart doesn't contain such item")
	ErrZeroQuantity   error = errors.New("product quantity can not be 0")
	ErrEmptyProduct   error = errors.New("product can not be empty")
)
