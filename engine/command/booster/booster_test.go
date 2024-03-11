package booster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

		if i > 200 && i < 300 {
			price := boosterPrice(i)
			assert.Equal(t, 50, price)
		}

		if i > 400 {
			price := boosterPrice(i)
			assert.Equal(t, 100, price)
		}
	}
}
