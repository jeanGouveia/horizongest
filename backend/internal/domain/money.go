package domain

import (
	"encoding/json"
	"math"
)

// Money representa um valor monetário em centavos.
// Exemplo: R$ 10,50 = 1050 centavos
type Money int64

// ToFloat64 converte centavos para reais (apenas para display/compatibilidade).
// Exemplo: 1050 → 10.50
func (m Money) ToFloat64() float64 {
	return float64(m) / 100.0
}

// FromFloat64 converte reais para centavos (apenas para input de legado).
// Exemplo: 10.50 → 1050
func FromFloat64(f float64) Money {
	return Money(math.Round(f * 100))
}

// MarshalJSON serializa Money como float64 para compatibilidade com API.
// Exemplo: 1050 → 10.50 (no JSON)
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.ToFloat64())
}

// UnmarshalJSON deserializa float/string para Money.
// Exemplo: 10.50 → 1050
func (m *Money) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*m = FromFloat64(f)
	return nil
}

// Add soma dois valores Money.
func (m Money) Add(other Money) Money {
	return m + other
}

// Sub subtrai dois valores Money.
func (m Money) Sub(other Money) Money {
	return m - other
}

// Mul multiplica Money por int64.
func (m Money) Mul(other int64) Money {
	return Money(int64(m) * other)
}

// Div divide Money por int64.
func (m Money) Div(other int64) Money {
	return Money(int64(m) / other)
}

// MulFloat multiplica Money por float64 (para cálculos como 30%).
func (m Money) MulFloat(f float64) Money {
	return Money(math.Round(float64(m) * f))
}

// Cmp compara dois valores Money.
// Retorna -1 se m < other, 0 se m == other, 1 se m > other.
func (m Money) Cmp(other Money) int {
	if m < other {
		return -1
	}
	if m > other {
		return 1
	}
	return 0
}

// IsZero retorna true se o valor é zero.
func (m Money) IsZero() bool {
	return m == 0
}

// IsNegative retorna true se o valor é negativo.
func (m Money) IsNegative() bool {
	return m < 0
}

// IsPositive retorna true se o valor é positivo.
func (m Money) IsPositive() bool {
	return m > 0
}
