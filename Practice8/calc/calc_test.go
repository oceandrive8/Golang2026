package calc

import (
	"testing"
)

func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 5, 3, 2},
		{"positive - zero", 5, 0, 5},
		{"negative - positive", -1, 4, -5},
		{"both negative", -2, -3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {

	t.Run("valid division", func(t *testing.T) {
		got, err := Divide(10, 2)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got != 5 {
			t.Errorf("Divide(10, 2) = %d; want 5", got)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		got, err := Divide(10, 0)

		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if got != 0 {
			t.Errorf("Divide(10, 0) = %d; want 0", got)
		}
	})
}
