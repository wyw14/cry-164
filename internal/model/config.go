package model

import "time"

type Config struct {
	Port            int
	StageTimeout    time.Duration
	LubeSoakSeconds int
	MaxPressure     float64
	MinLevel        float64
}

func DefaultConfig() Config {
	return Config{Port: 21264, StageTimeout: 250 * time.Millisecond, LubeSoakSeconds: 45, MaxPressure: 180, MinLevel: 0.25}
}

func (c Config) LubeSoak() time.Duration {
	return time.Duration(c.LubeSoakSeconds) * time.Second
}

func (c Config) Valid() bool {
	return c.Port > 0 && c.StageTimeout > 0 && c.LubeSoakSeconds > 0 && c.MaxPressure > 0 && c.MinLevel >= 0
}
