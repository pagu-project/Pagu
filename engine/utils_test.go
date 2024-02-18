package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckArgs(t *testing.T) {
	err := CheckArgs(2, []string{"1", "2", "3"})
	assert.EqualError(t, err, "incorrect number of arguments, expected 2 but got 3")

	err = CheckArgs(2, []string{"1"})
	assert.EqualError(t, err, "incorrect number of arguments, expected 2 but got 1")
}

func TestBoosterPrice(t *testing.T) {
	for i := 0; i < 501; i++ {
		if i < 100 {
			price := boosterPrice(i)
			assert.Equal(t, 30, price)
		}

		if i > 100 && i < 200 {
			price := boosterPrice(i)
			assert.Equal(t, 40, price)
		}

		if i > 200 {
			price := boosterPrice(i)
			assert.Equal(t, 50, price)
		}
	}
}
