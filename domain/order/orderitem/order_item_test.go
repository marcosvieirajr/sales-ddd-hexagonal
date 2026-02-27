package orderitem_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/marcosvieirajr/sales/domain"
	"github.com/marcosvieirajr/sales/domain/order/orderitem"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidOrderItem(t *testing.T, unitPrice float64, quantity int) *orderitem.OrderItem {
	t.Helper()
	return domain.Must(orderitem.NewOrderItem("prod-123", "Test Product", unitPrice, quantity))
}

func TestNewOrderItem(t *testing.T) {
	t.Run("should successfully create a new order item with valid input", func(t *testing.T) {
		got, err := orderitem.NewOrderItem("prod-123", "Product Name", 10.0, 2)

		require.NoError(t, err)
		want := &orderitem.OrderItem{
			ProductID:       "prod-123",
			ProductName:     "Product Name",
			UnitPrice:       10.0,
			Quantity:        2,
			DiscountApplied: 0.0,
			TotalPrice:      20.0,
		}
		ignoreFields := cmpopts.IgnoreFields(orderitem.OrderItem{}, "ID", "CreatedAt") // ignore ID and CreatedAt since they are generated and not predictable
		assert.True(t, cmp.Equal(got, want, ignoreFields), "got and want should be equal ignoring ID and createdAt: %v", cmp.Diff(got, want, ignoreFields))
	})

	t.Run("should return an error when invalid input is provided", func(t *testing.T) {
		type args struct {
			productID   string
			productName string
			unitPrice   float64
			quantity    int
		}
		tests := []struct {
			name        string
			args        args
			wantErr error
		}{
			{
				name:        "should return an error if product ID is invalid",
				args:        args{productID: "", productName: "Product Name", unitPrice: 10.0, quantity: 2},
				wantErr: orderitem.ErrInvalidProductID,
			},
			{
				name:        "should return an error if product name is empty",
				args:        args{productID: "prod-123", productName: "", unitPrice: 10.0, quantity: 2},
				wantErr: orderitem.ErrInvalidProductName,
			},
			{
				name:        "should return an error if unit price is zero",
				args:        args{productID: "prod-123", productName: "Product Name", unitPrice: 0.0, quantity: 2},
				wantErr: orderitem.ErrInvalidUnitPrice,
			},
			{
				name:        "should return an error if unit price is negative",
				args:        args{productID: "prod-123", productName: "Product Name", unitPrice: -0.1, quantity: 2},
				wantErr: orderitem.ErrInvalidUnitPrice,
			},
			{
				name:        "should return an error if quantity is zero",
				args:        args{productID: "prod-123", productName: "Product Name", unitPrice: 10.0, quantity: 0},
				wantErr: orderitem.ErrInvalidQuantity,
			},
			{
				name:        "should return an error if quantity is negative",
				args:        args{productID: "prod-123", productName: "Product Name", unitPrice: 10.0, quantity: -1},
				wantErr: orderitem.ErrInvalidQuantity,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := orderitem.NewOrderItem(tt.args.productID, tt.args.productName, tt.args.unitPrice, tt.args.quantity)

				assert.Nil(t, got)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}

func TestOrderItem_ApplyDiscount(t *testing.T) {
	t.Run("should successfully apply discount", func(t *testing.T) {
		oi := createValidOrderItem(t, 10.0, 2)

		err := oi.ApplyDiscount(5.0)

		require.NoError(t, err)
		assert.Equal(t, 5.0, oi.DiscountApplied, "DiscountApplied should be set correctly")
		assert.Equal(t, 15.0, oi.TotalPrice, "TotalPrice should be (10 * 2) - 5 = 15")
		assert.NotNil(t, oi.UpdatedAt, "UpdatedAt should be set on success")
	})

	t.Run("should return an error when discount is invalid", func(t *testing.T) {
		type fields struct {
			unitPrice float64
			quantity  int
		}
		tests := []struct {
			name               string
			fields             fields
			discount           float64
			wantTotalPrice float64
			wantErr        error
		}{
			{
				name:               "should return an error when discount is negative",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				discount:           -1.0,
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrNegativeDiscount,
			},
			{
				name:               "should return an error when discount is greater than unit price",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				discount:           11.0,
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrDiscountExceedsUnitPrice,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				oi := createValidOrderItem(t, tt.fields.unitPrice, tt.fields.quantity)

				err := oi.ApplyDiscount(tt.discount)

				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, 0.0, oi.DiscountApplied, "DiscountApplied should remain zero on error")
				assert.Equal(t, tt.wantTotalPrice, oi.TotalPrice, "TotalPrice should not change on error")
				assert.Nil(t, oi.UpdatedAt, "UpdatedAt should remain nil on error")
			})
		}
	})
}

func TestOrderItem_AddUnits(t *testing.T) {
	t.Run("should successfully add units when valid units are provided", func(t *testing.T) {
		type fields struct {
			unitPrice float64
			quantity  int
		}
		tests := []struct {
			name               string
			fields             fields
			units              int
			wantQuantity   int
			wantTotalPrice float64
		}{
			{
				name:               "should add units when valid units are provided",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				units:              3,
				wantQuantity:   5,
				wantTotalPrice: 50.0, // 10 * 5 = 50
			},
			{
				name:               "should add single unit",
				fields:             fields{unitPrice: 15.0, quantity: 1},
				units:              1,
				wantQuantity:   2,
				wantTotalPrice: 30.0, // 15 * 2 = 30
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				oi := createValidOrderItem(t, tt.fields.unitPrice, tt.fields.quantity)

				err := oi.AddUnits(tt.units)

				require.NoError(t, err)
				assert.Equal(t, tt.wantQuantity, oi.Quantity, "Quantity should be updated correctly: actual %v, expected %v", oi.Quantity, tt.wantQuantity)
				assert.Equal(t, tt.wantTotalPrice, oi.TotalPrice, "TotalPrice should be recalculated correctly: actual %v, expected %v", oi.TotalPrice, tt.wantTotalPrice)
				assert.NotNil(t, oi.UpdatedAt, "UpdatedAt should be set on success")
			})
		}
	})

	t.Run("should return an error when units is invalid", func(t *testing.T) {
		type fields struct {
			unitPrice float64
			quantity  int
		}
		tests := []struct {
			name               string
			fields             fields
			units              int
			wantQuantity   int
			wantTotalPrice float64
			wantErr        error
		}{
			{
				name:               "should return an error when units is zero",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				units:              0,
				wantQuantity:   2,
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrInvalidUnits,
			},
			{
				name:               "should return an error when units is negative",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				units:              -1,
				wantQuantity:   2,
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrInvalidUnits,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				oi := createValidOrderItem(t, tt.fields.unitPrice, tt.fields.quantity)

				err := oi.AddUnits(tt.units)

				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.wantQuantity, oi.Quantity, "Quantity should not change on error: actual %v, expected %v", oi.Quantity, tt.wantQuantity)
				assert.Equal(t, tt.wantTotalPrice, oi.TotalPrice, "TotalPrice should not change on error: actual %v, expected %v", oi.TotalPrice, tt.wantTotalPrice)
				assert.Nil(t, oi.UpdatedAt, "UpdatedAt should remain nil on error")
			})
		}
	})
}

func TestOrderItem_RemoveUnits(t *testing.T) {
	t.Run("should successfully remove units when valid units are provided", func(t *testing.T) {
		type fields struct {
			unitPrice float64
			quantity  int
		}
		tests := []struct {
			name               string
			fields             fields
			units              int
			wantQuantity   int
			wantTotalPrice float64
		}{
			{
				name:               "should remove units when valid units are provided",
				fields:             fields{unitPrice: 10.0, quantity: 5},
				units:              2,
				wantQuantity:   3,
				wantTotalPrice: 30.0, // 10 * 3 = 30
			},
			{
				name:               "should remove single unit",
				fields:             fields{unitPrice: 15.0, quantity: 3},
				units:              1,
				wantQuantity:   2,
				wantTotalPrice: 30.0, // 15 * 2 = 30
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				oi := createValidOrderItem(t, tt.fields.unitPrice, tt.fields.quantity)

				err := oi.RemoveUnits(tt.units)

				require.NoError(t, err)
				assert.Equal(t, tt.wantQuantity, oi.Quantity, "Quantity should be updated correctly: actual %v, expected %v", oi.Quantity, tt.wantQuantity)
				assert.Equal(t, tt.wantTotalPrice, oi.TotalPrice, "TotalPrice should be recalculated correctly: actual %v, expected %v", oi.TotalPrice, tt.wantTotalPrice)
				assert.NotNil(t, oi.UpdatedAt, "UpdatedAt should be set on success")
			})
		}
	})

	t.Run("should return an error when units is invalid", func(t *testing.T) {
		type fields struct {
			unitPrice float64
			quantity  int
		}
		tests := []struct {
			name               string
			fields             fields
			units              int
			wantQuantity   int
			wantTotalPrice float64
			wantErr        error
		}{
			{
				name:               "should return an error when units is zero",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				units:              0,
				wantQuantity:   2,
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrInvalidUnits,
			},
			{
				name:               "should return an error when units is negative",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				units:              -1,
				wantQuantity:   2,
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrInvalidUnits,
			},
			{
				name:               "should return an error when units to remove equals current quantity",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				units:              2,
				wantQuantity:   2,
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrInsufficientQuantity,
			},
			{
				name:               "should return an error when units to remove is greater than current quantity",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				units:              5,
				wantQuantity:   2,
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrInsufficientQuantity,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				oi := createValidOrderItem(t, tt.fields.unitPrice, tt.fields.quantity)

				err := oi.RemoveUnits(tt.units)

				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.wantQuantity, oi.Quantity, "Quantity should not change on error: actual %v, expected %v", oi.Quantity, tt.wantQuantity)
				assert.Equal(t, tt.wantTotalPrice, oi.TotalPrice, "TotalPrice should not change on error: actual %v, expected %v", oi.TotalPrice, tt.wantTotalPrice)
				assert.Nil(t, oi.UpdatedAt, "UpdatedAt should remain nil on error")
			})
		}
	})
}

func TestOrderItem_UpdateUnitPrice(t *testing.T) {
	t.Run("should successfully update unit price when valid price is provided", func(t *testing.T) {
		type fields struct {
			unitPrice float64
			quantity  int
		}
		tests := []struct {
			name               string
			fields             fields
			value              float64
			wantUnitPrice  float64
			wantTotalPrice float64
		}{
			{
				name:               "should update unit price when valid price is provided",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				value:              15.0,
				wantUnitPrice:  15.0,
				wantTotalPrice: 30.0, // 15 * 2 = 30
			},
			{
				name:               "should update unit price with decimal value",
				fields:             fields{unitPrice: 10.0, quantity: 3},
				value:              12.50,
				wantUnitPrice:  12.50,
				wantTotalPrice: 37.50, // 12.50 * 3 = 37.50
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				oi := createValidOrderItem(t, tt.fields.unitPrice, tt.fields.quantity)

				err := oi.UpdateUnitPrice(tt.value)

				require.NoError(t, err)
				assert.Equal(t, tt.wantUnitPrice, oi.UnitPrice, "UnitPrice should be updated correctly: actual %v, expected %v", oi.UnitPrice, tt.wantUnitPrice)
				assert.Equal(t, tt.wantTotalPrice, oi.TotalPrice, "TotalPrice should be recalculated correctly: actual %v, expected %v", oi.TotalPrice, tt.wantTotalPrice)
				assert.NotNil(t, oi.UpdatedAt, "UpdatedAt should be set on success")
			})
		}
	})

	t.Run("should return an error when unit price is invalid", func(t *testing.T) {
		type fields struct {
			unitPrice float64
			quantity  int
		}
		tests := []struct {
			name               string
			fields             fields
			value              float64
			wantUnitPrice  float64
			wantTotalPrice float64
			wantErr        error
		}{
			{
				name:               "should return an error when unit price is zero",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				value:              0.0,
				wantUnitPrice:  10.0, // no change
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrInvalidUnitPrice,
			},
			{
				name:               "should return an error when unit price is negative",
				fields:             fields{unitPrice: 10.0, quantity: 2},
				value:              -5.0,
				wantUnitPrice:  10.0, // no change
				wantTotalPrice: 20.0, // no change
				wantErr:        orderitem.ErrInvalidUnitPrice,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				oi := createValidOrderItem(t, tt.fields.unitPrice, tt.fields.quantity)

				err := oi.UpdateUnitPrice(tt.value)

				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.wantUnitPrice, oi.UnitPrice, "UnitPrice should not change on error: actual %v, expected %v", oi.UnitPrice, tt.wantUnitPrice)
				assert.Equal(t, tt.wantTotalPrice, oi.TotalPrice, "TotalPrice should not change on error: actual %v, expected %v", oi.TotalPrice, tt.wantTotalPrice)
				assert.Nil(t, oi.UpdatedAt, "UpdatedAt should remain nil on error")
			})
		}
	})
}

func TestOrderItem_Equals(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T) (oi, other *orderitem.OrderItem)
		want     bool
	}{
		// ==================== Success cases ==================== //
		{
			name: "should return true when comparing same order item instance",
			setup: func(t *testing.T) (*orderitem.OrderItem, *orderitem.OrderItem) {
				oi := createValidOrderItem(t, 10.0, 2)
				return oi, oi
			},
			want:     true,
		},
		{
			name: "should return true when order items have same ID",
			setup: func(t *testing.T) (*orderitem.OrderItem, *orderitem.OrderItem) {
				return &orderitem.OrderItem{ID: "same-id", ProductID: "prod-1", ProductName: "Product A", UnitPrice: 10.0, Quantity: 2},
					&orderitem.OrderItem{ID: "same-id", ProductID: "prod-2", ProductName: "Product B", UnitPrice: 20.0, Quantity: 5}
			},
			want:     true,
		},
		// ==================== Failure cases ==================== //
		{
			name: "should return false when comparing with nil",
			setup: func(t *testing.T) (*orderitem.OrderItem, *orderitem.OrderItem) {
				return createValidOrderItem(t, 10.0, 2), nil
			},
			want:     false,
		},
		{
			name: "should return false when order items have different IDs",
			setup: func(t *testing.T) (*orderitem.OrderItem, *orderitem.OrderItem) {
				return &orderitem.OrderItem{ID: "id-1", ProductID: "prod-1", ProductName: "Product A", UnitPrice: 10.0, Quantity: 2},
					&orderitem.OrderItem{ID: "id-2", ProductID: "prod-1", ProductName: "Product A", UnitPrice: 10.0, Quantity: 2}
			},
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oi, other := tt.setup(t)

			result := oi.Equals(other)

			assert.Equal(t, tt.want, result, "Equals should return %v but got %v", tt.want, result)
		})
	}
}
