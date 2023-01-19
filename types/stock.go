package types

import (
	"math"
)

// Stock represents a stock value
type Stock struct {
	value           Cash // The cash value in the stock currency
	conversionPrice Cash // conversion Cash from stock currency to CAD (1$ | 0$ == stock in CAD)
}

// Create a new stock given its value and conversion price
// Ex: ("10USD", "1.23CAD/USD")
// Note: conversion prices 0 and 1 mean the value is in CAD
// i.e. ("10CAD", "0") == ("10CAD", "1")
func ParseStock(valueStr string, conversionStr string) (Stock, error) {
	value, err := ParseCash(valueStr, false)
	if err != nil {
		return Stock{}, err
	}
	conversion, err := ParseCash(conversionStr, false)
	if err != nil {
		return Stock{}, err
	}
	if conversion == 0 {
		conversion = CashDollar
	}
	return Stock{
		value:           value,
		conversionPrice: conversion,
	}, nil
}

// Multiply the stock for a given percentage
func (s Stock) Percentage(p Percentage) (Stock, error) {
	newValue, err := s.value.Percentage(p)
	if err != nil {
		return Stock{}, err
	}
	s.value = newValue
	return s, nil
}

// Returns Cash in CAD, Cash in the original value,
// a bool indicating if original currency is CAD
func (s Stock) Value() (cad Cash, orig Cash, isCAD bool, err error) {
	// If CAD, return value
	isCAD = s.conversionPrice == CashDollar
	if isCAD {
		return s.value, s.value, true, nil
	}
	// Detect potential overflows
	if math.MaxInt64/int64(s.conversionPrice) < int64(s.value) {
		return 0, 0, false, ErrorOverflow
	}
	return s.value * s.conversionPrice / CashDollar, s.value, false, nil
}
