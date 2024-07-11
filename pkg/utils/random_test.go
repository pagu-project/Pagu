package utils_test

import (
	"testing"

	"github.com/pagu-project/Pagu/pkg/utils"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		n      int
		input  string
		length int
	}{
		{10, utils.Alphabet, 10},
		{5, utils.Numbers, 5},
		{15, utils.CapitalAlphanumerical, 15},
		{0, utils.Alphabet, 0},
		{20, "abc123", 20},
	}

	for _, tt := range tests {
		result := utils.RandomString(tt.n, tt.input)
		if len(result) != tt.length {
			t.Errorf("expected length %d, got %d", tt.length, len(result))
		}
		for _, char := range result {
			if !contains(tt.input, char) {
				t.Errorf("result contains invalid character: %c", char)
			}
		}
	}
}

func TestRandomStringNonCollision(t *testing.T) {
	const numStrings = 10000
	const strLength = 10
	input := utils.Alphabet

	generatedStrings := make(map[string]struct{}, numStrings)

	for i := 0; i < numStrings; i++ {
		str := utils.RandomString(strLength, input)
		if _, exists := generatedStrings[str]; exists {
			t.Errorf("collision detected: %s", str)
		}
		generatedStrings[str] = struct{}{}
	}

	if len(generatedStrings) != numStrings {
		t.Errorf("expected %d unique strings, got %d", numStrings, len(generatedStrings))
	}
}

func contains(s string, c rune) bool {
	for _, char := range s {
		if char == c {
			return true
		}
	}
	return false
}

func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.RandomString(10, utils.Alphabet)
	}
}
