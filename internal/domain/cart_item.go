package domain

type (
	//CartItem is item in Cart
	CartItem struct {
		ID       string `json:"id"`
		CartID   string `json:"cart_id"`
		Product  string `json:"product"`
		Quantity int    `json:"quantity"`
	}
)

func (ci CartItem) Validate() error {
	if len(ci.Product) == 0 {
		return ErrEmptyProduct
	}

	if ci.Quantity == 0 {
		return ErrZeroQuantity
	}

	return nil
}
