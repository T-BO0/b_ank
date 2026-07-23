package util

import (
	"math/rand"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomDecimal(min, max int64, floatingNumber int32) decimal.Decimal {
	wholeValue := (random.Int63n(max) + min - 1)

	return decimal.New(wholeValue, floatingNumber)
}

func RandomString(length int) string {
	var sb strings.Builder
	sb.Grow(length)

	for range length {
		sb.WriteByte(charset[random.Intn(len(charset))])
	}

	return sb.String()
}

// Random balance
func RandomMonetaryAmount() decimal.Decimal {
	money := RandomDecimal(1000000000, 1000000000000, -6)
	return money
}

// Random Account owner name
func RandomOwnerName() string {
	owner := RandomString(10)
	return owner
}

// Random Currency Code
func RandomCurrencyCode() string {
	currecny := []string{"EUR", "USD", "CAD"}
	length := len(currecny)

	return currecny[random.Intn(length)]
}
