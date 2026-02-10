package chaos

import (
	"errors"
	"math/rand"
	"time"
)

var ErrChaosInjected = errors.New("chaos: injected error")

func MaybeDelay(prob float64, maxDelay time.Duration) {
	if rand.Float64() < prob {
		delay := time.Duration(rand.Int63n(int64(maxDelay)))
		time.Sleep(delay)
	}
}

func MaybeError(prob float64) error {
	if rand.Float64() < prob {
		return ErrChaosInjected
	}
	return nil
}
