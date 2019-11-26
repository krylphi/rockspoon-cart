package domain

import (
	"testing"
)

func TestCartItem_Validate(t *testing.T) {
	type fields struct {
		ID       string
		CartID   string
		Product  string
		Quantity int
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		expectedErr error
	}{
		// TODO: Add test cases.
		{
			name:    "Test empty Product",
			fields:  fields{
				ID:       "ID",
				CartID:   "CartID",
				Product:  "",
				Quantity: 1,
			},
			wantErr:     true,
			expectedErr: ErrEmptyProduct,
		},
		{
			name:    "Test Zero Quantity",
			fields:  fields{
				ID:       "ID",
				CartID:   "CartID",
				Product:  "Product",
				Quantity: 0,
			},
			wantErr:     true,
			expectedErr: ErrZeroQuantity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ci := CartItem{
				ID:       tt.fields.ID,
				CartID:   tt.fields.CartID,
				Product:  tt.fields.Product,
				Quantity: tt.fields.Quantity,
			}
			if err := ci.Validate(); (err != nil) != tt.wantErr || (tt.wantErr && tt.expectedErr != nil && err != tt.expectedErr) {
				t.Errorf("Validate() error = %v, wantErr: %v, expectedErr: %v", err, tt.wantErr, tt.expectedErr)
			}
		})
	}
}
