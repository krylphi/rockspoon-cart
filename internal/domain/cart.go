package domain

type (
	//Cart is domain model for cart entity
	Cart struct {
		ID    string      `json:"id"`
		Items []*CartItem `json:"items"`
	}
)

func (c *Cart) AddItem(id string, product string, quantity int) (ci *CartItem, err error) {

	item := &CartItem{
		ID:       id,
		CartID:   c.ID,
		Product:  product,
		Quantity: quantity,
	}

	if err = item.Validate(); err != nil {
		return nil, err
	}

	//c.Items[id] = item
	c.Items = append(c.Items, item)
	return item, nil
}

func (c *Cart) RemoveItem(id string) (err error) {

	var ok bool
	var idx int

	for pos, item := range c.Items {
		if item.ID == id {
			ok = true
			idx = pos
			break
		}
	}

	if !ok {
		return ErrNoSuchCartItem
	}

	c.Items = append(c.Items[:idx], c.Items[idx+1:]...)

	return nil
}
