package domain

import (
	"reflect"
	"testing"
)

func TestCart_AddItem(t *testing.T) {
	type (
		fields struct {
			ID    string
			Items []*CartItem
		}
		args struct {
			id       string
			product  string
			quantity int
		}
	)

	tests := []struct {
		name          string
		fields        fields
		args          args
		wantCi        *CartItem
		expectedCiCnt int
		wantErr       bool
		expectedErr   error
	}{
		{
			name: "Test 0 quantity add",
			fields: fields{
				ID:    "CartID",
				Items: make([]*CartItem, 0),
			},
			args: args{
				id:       "ProductID",
				product:  "Product",
				quantity: 0,
			},
			wantCi:        nil,
			expectedCiCnt: 0,
			wantErr:       true,
			expectedErr:   ErrZeroQuantity,
		},
		{
			name: "Test empty product",
			fields: fields{
				ID:    "CartID",
				Items: make([]*CartItem, 0),
			},
			args: args{
				id:       "ProductID",
				product:  "",
				quantity: 1,
			},
			wantCi:        nil,
			expectedCiCnt: 0,
			wantErr:       true,
			expectedErr:   ErrEmptyProduct,
		},
		{
			name: "Test regular case",
			fields: fields{
				ID:    "CartID",
				Items: make([]*CartItem, 0),
			},
			args: args{
				id:       "ProductID",
				product:  "Product",
				quantity: 1,
			},
			wantCi: &CartItem{
				ID:       "ProductID",
				CartID:   "CartID",
				Product:  "Product",
				Quantity: 1,
			},
			expectedCiCnt: 1,
			wantErr:       false,
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cart{
				ID:    tt.fields.ID,
				Items: tt.fields.Items,
			}

			gotCi, err := c.AddItem(tt.args.id, tt.args.product, tt.args.quantity)

			if (err != nil) != tt.wantErr || (tt.wantErr && tt.expectedErr != nil && err != tt.expectedErr) {
				t.Fatalf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotCi, tt.wantCi) {
				t.Fatalf("AddItem() gotCi = %v, want %v", gotCi, tt.wantCi)
			}

			if tt.expectedCiCnt != len(c.Items) {
				t.Fatalf("AddItem() gotCiCnt = %v, want %v", len(c.Items), tt.expectedCiCnt)
			}
		})
	}
}

func TestCart_RemoveItem(t *testing.T) {
	type (
		fields struct {
			ID    string
			Items []*CartItem
		}
		args struct {
			id string
		}
	)

	testCart := fields{
		ID: "CartID",
		Items: []*CartItem{
			{
				ID:       "Item1",
				CartID:   "CartID",
				Product:  "Product1",
				Quantity: 1,
			},
			{
				ID:       "Item2",
				CartID:   "CartID",
				Product:  "Product2",
				Quantity: 2,
			},
			{
				ID:       "Item3",
				CartID:   "CartID",
				Product:  "Product3",
				Quantity: 3,
			},
			{
				ID:       "Item4",
				CartID:   "CartID",
				Product:  "Product4",
				Quantity: 4,
			},
		},
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "No such item",
			fields: testCart,
			args: args{
				id: "NoExistingId",
			},
			wantErr:     true,
			expectedErr: ErrNoSuchCartItem,
		}, {
			name:   "Regular delete",
			fields: testCart,
			args: args{
				id: "Item2",
			},
			wantErr:     false,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cart{
				ID:    tt.fields.ID,
				Items: tt.fields.Items,
			}
			err := c.RemoveItem(tt.args.id)
			if (err != nil) != tt.wantErr || (tt.wantErr && tt.expectedErr != nil && err != tt.expectedErr) {
				t.Fatalf("RemoveItem() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && !tt.wantErr {
				for _, item := range c.Items {
					if item.ID == tt.args.id {
						t.Fatalf("RemoveItem() item with id %v was not removed", tt.args.id)
					}
				}
			}
		})
	}
}
