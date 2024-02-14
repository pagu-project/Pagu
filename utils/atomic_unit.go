package utils

import "strconv"

const changeFactor = float64(1_000_000_000)

// CoinToChange converts a coin amount to its corresponding change value.
// Example: CoinToChange(2.75) returns 2750000000.
func CoinToChange(coin float64) int64 {
	return int64(coin * changeFactor)
}

// ChangeToCoin converts a change value to its corresponding coin amount.
// Example: ChangeToCoin(2750000000) returns 2.75.
func ChangeToCoin(change int64) float64 {
	return float64(change) / changeFactor
}

// StringToChange converts a string representation of a coin amount to its corresponding change value.
// It returns the change value as an int64 and an error if the string cannot be parsed.
// Example: StringToChange("2.75") returns 2750000000, nil.
func StringToChange(amount string) (int64, error) {
	coin, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0, err
	}

	return CoinToChange(coin), nil
}

// ChangeToStringWithTrailingZeros converts a change value to its string representation with trailing zeros.
// The returned string will have up to 9 decimal places.
// Example: ChangeToStringWithTrailingZeros(2750000000) returns "2.750000000".
func ChangeToStringWithTrailingZeros(change int64) string {
	coin := ChangeToCoin(change)

	return strconv.FormatFloat(coin, 'f', 9, 64)
}

// ChangeToString converts a change value to its string representation.
// Example: ChangeToString(2750000000) returns "2.75".
func ChangeToString(change int64) string {
	coin := ChangeToCoin(change)

	return strconv.FormatFloat(coin, 'f', -1, 64)
}

// ChangeToString converts a change value to its string representation.
// Example: ChangeToString(2750000000) returns "2".
func ChangeToStringNormal(change int64) string {
	coin := ChangeToCoin(change)

	return strconv.FormatFloat(coin, 'f', 0, 64)
}
