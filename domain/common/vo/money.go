package vo

import (
	"errors"
	"math"
)

// Money is a value object that represents an amount of currency
type Money struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// NewMoney creates a new Money value object with validation
func NewMoney(amount float64, currency string) (Money, error) {
	if currency == "" {
		return Money{}, errors.New("currency cannot be empty")
	}
	
	// Round to 2 decimal places
	roundedAmount := math.Round(amount*100) / 100
	
	return Money{
		Amount:   roundedAmount,
		Currency: currency,
	}, nil
}

// Add adds the amount of another Money object and returns a new Money object
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot add different currencies")
	}
	
	newAmount := m.Amount + other.Amount
	return NewMoney(newAmount, m.Currency)
}

// Subtract subtracts the amount of another Money object and returns a new Money object
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("cannot subtract different currencies")
	}
	
	newAmount := m.Amount - other.Amount
	return NewMoney(newAmount, m.Currency)
}

// Multiply multiplies the amount by a factor and returns a new Money object
func (m Money) Multiply(factor float64) (Money, error) {
	newAmount := m.Amount * factor
	return NewMoney(newAmount, m.Currency)
}

// IsZero checks if the money amount is zero
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsNegative checks if the money amount is negative
func (m Money) IsNegative() bool {
	return m.Amount < 0
}

// IsPositive checks if the money amount is positive
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// Equals checks if two money objects are equal
func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}
