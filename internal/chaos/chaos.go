package chaos

import (
	"errors"
	"math/rand"
	"time"
)

var ErrChaosInjected = errors.New("chaos: injected error")

func MaybeDelay(prob float64, maxDelay time.Duration) {
	// nolint:gosec // G404: math/rand is sufficient for chaos injection
	if rand.Float64() < prob {
		// nolint:gosec // G404: math/rand is sufficient for chaos injection
		delay := time.Duration(rand.Int63n(int64(maxDelay)))
		time.Sleep(delay)
	}
}

func MaybeError(prob float64) error {
	// nolint:gosec // G404: math/rand is sufficient for chaos injection
	if rand.Float64() < prob {
		return ErrChaosInjected
	}
	return nil
}
