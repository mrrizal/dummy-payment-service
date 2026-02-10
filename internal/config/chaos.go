package config

import "time"

type ChaosConfig struct {
	Enabled          bool
	ErrorProbability float64       // 0.0 - 1.0
	DelayProbability float64       // 0.0 - 1.0
	MaxDelay         time.Duration // contoh: 500ms
}
