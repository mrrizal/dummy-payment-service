package chaos

import (
	"errors"
	"math/rand"
	"time"
)

var ErrChaosInjected = errors.New("chaos: injected error")

func MaybeDelay(prob float64, maxDelay time.Duration) {
	// #nosec G404
	if rand.Float64() < prob {
		// #nosec G404
		delay := time.Duration(rand.Int63n(int64(maxDelay)))
		time.Sleep(delay)
	}
}

func MaybeError(prob float64) error {
	// #nosec G404
	if rand.Float64() < prob {
		return ErrChaosInjected
	}
	return nil
}
