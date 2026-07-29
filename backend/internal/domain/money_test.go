package domain

import (
	"encoding/json"
	"testing"
)

func TestMoneyFromFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected Money
	}{
		{"zero", 0.0, 0},
		{"um centavo", 0.01, 1},
		{"dez centavos", 0.10, 10},
		{"um real", 1.00, 100},
		{"dez reais", 10.00, 1000},
		{"dez reais e cinquenta centavos", 10.50, 1050},
		{"valor grande", 999999.99, 99999999},
		{"arredondamento para cima", 10.505, 1051},
		{"arredondamento para baixo", 10.504, 1050},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromFloat64(tt.input)
			if result != tt.expected {
				t.Errorf("FromFloat64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMoneyToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    Money
		expected float64
	}{
		{"zero", 0, 0.0},
		{"um centavo", 1, 0.01},
		{"dez centavos", 10, 0.10},
		{"um real", 100, 1.00},
		{"dez reais", 1000, 10.00},
		{"dez reais e cinquenta centavos", 1050, 10.50},
		{"valor grande", 99999999, 999999.99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToFloat64()
			if result != tt.expected {
				t.Errorf("Money(%v).ToFloat64() = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMoneyAdd(t *testing.T) {
	tests := []struct {
		name     string
		a        Money
		b        Money
		expected Money
	}{
		{"zero + zero", 0, 0, 0},
		{"zero + valor", 0, 100, 100},
		{"valor + zero", 100, 0, 100},
		{"valor + valor", 100, 50, 150},
		{"overflow teste", 9223372036854775807, 1, -9223372036854775808}, // overflow esperado
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Add(tt.b)
			if result != tt.expected {
				t.Errorf("Money(%v).Add(%v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMoneySub(t *testing.T) {
	tests := []struct {
		name     string
		a        Money
		b        Money
		expected Money
	}{
		{"zero - zero", 0, 0, 0},
		{"valor - zero", 100, 0, 100},
		{"zero - valor", 0, 100, -100},
		{"valor - valor", 150, 50, 100},
		{"negativo", 50, 100, -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Sub(tt.b)
			if result != tt.expected {
				t.Errorf("Money(%v).Sub(%v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMoneyMul(t *testing.T) {
	tests := []struct {
		name     string
		a        Money
		b        int64
		expected Money
	}{
		{"zero * zero", 0, 0, 0},
		{"valor * zero", 100, 0, 0},
		{"zero * valor", 0, 5, 0},
		{"valor * valor", 100, 3, 300},
		{"valor * negativo", 100, -2, -200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Mul(tt.b)
			if result != tt.expected {
				t.Errorf("Money(%v).Mul(%v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMoneyMulFloat(t *testing.T) {
	tests := []struct {
		name     string
		a        Money
		b        float64
		expected Money
	}{
		{"zero * zero", 0, 0.0, 0},
		{"valor * zero", 100, 0.0, 0},
		{"zero * valor", 0, 0.5, 0},
		{"valor * 0.5", 100, 0.5, 50},
		{"valor * 2.0", 100, 2.0, 200},
		{"valor * 1.5", 100, 1.5, 150},
		{"desconto 10%", 1000, 0.9, 900},
		{"desconto 25%", 1000, 0.75, 750},
		{"imposto 20%", 1000, 1.2, 1200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.MulFloat(tt.b)
			if result != tt.expected {
				t.Errorf("Money(%v).MulFloat(%v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMoneyDiv(t *testing.T) {
	tests := []struct {
		name     string
		a        Money
		b        int64
		expected Money
	}{
		{"zero / valor", 0, 5, 0},
		{"valor / 1", 100, 1, 100},
		{"valor / 2", 100, 2, 50},
		{"valor / 4", 100, 4, 25},
		{"valor / 10", 100, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Div(tt.b)
			if result != tt.expected {
				t.Errorf("Money(%v).Div(%v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMoneyCmp(t *testing.T) {
	tests := []struct {
		name     string
		a        Money
		b        Money
		expected int
	}{
		{"igual", 100, 100, 0},
		{"a < b", 50, 100, -1},
		{"a > b", 100, 50, 1},
		{"zero < valor", 0, 100, -1},
		{"valor > zero", 100, 0, 1},
		{"negativo < positivo", -100, 100, -1},
		{"negativo < zero", -100, 0, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Cmp(tt.b)
			if result != tt.expected {
				t.Errorf("Money(%v).Cmp(%v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMoneyIsZero(t *testing.T) {
	tests := []struct {
		name     string
		input    Money
		expected bool
	}{
		{"zero", 0, true},
		{"positivo", 100, false},
		{"negativo", -100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsZero()
			if result != tt.expected {
				t.Errorf("Money(%v).IsZero() = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMoneyIsNegative(t *testing.T) {
	tests := []struct {
		name     string
		input    Money
		expected bool
	}{
		{"zero", 0, false},
		{"positivo", 100, false},
		{"negativo", -100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsNegative()
			if result != tt.expected {
				t.Errorf("Money(%v).IsNegative() = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMoneyIsPositive(t *testing.T) {
	tests := []struct {
		name     string
		input    Money
		expected bool
	}{
		{"zero", 0, false},
		{"positivo", 100, true},
		{"negativo", -100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.IsPositive()
			if result != tt.expected {
				t.Errorf("Money(%v).IsPositive() = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMoneyMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    Money
		expected string
	}{
		{"zero", 0, "0"},
		{"um centavo", 1, "0.01"},
		{"um real", 100, "1"},
		{"dez reais e cinquenta centavos", 1050, "10.5"},
		{"valor grande", 99999999, "999999.99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Money(%v).MarshalJSON() error = %v", tt.input, err)
				return
			}
			if string(result) != tt.expected {
				t.Errorf("Money(%v).MarshalJSON() = %v, want %v", tt.input, string(result), tt.expected)
			}
		})
	}
}

func TestMoneyUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Money
	}{
		{"zero", "0", 0},
		{"um centavo", "0.01", 1},
		{"um real", "1", 100},
		{"dez reais e cinquenta centavos", "10.5", 1050},
		{"valor grande", "999999.99", 99999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result Money
			err := json.Unmarshal([]byte(tt.input), &result)
			if err != nil {
				t.Errorf("Money.UnmarshalJSON(%v) error = %v", tt.input, err)
				return
			}
			if result != tt.expected {
				t.Errorf("Money.UnmarshalJSON(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMoneyRoundTrip(t *testing.T) {
	tests := []Money{0, 1, 10, 100, 1050, 99999999}

	for _, original := range tests {
		t.Run("", func(t *testing.T) {
			// Marshal
			jsonBytes, err := json.Marshal(original)
			if err != nil {
				t.Errorf("Marshal error: %v", err)
				return
			}

			// Unmarshal
			var result Money
			err = json.Unmarshal(jsonBytes, &result)
			if err != nil {
				t.Errorf("Unmarshal error: %v", err)
				return
			}

			// Verificar round-trip
			if result != original {
				t.Errorf("Round-trip failed: %v -> %v -> %v", original, string(jsonBytes), result)
			}
		})
	}
}
