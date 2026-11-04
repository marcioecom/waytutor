package ratelimit

import (
	"sync"
	"time"
)

const sweepInterval = 10 * time.Minute

// Limiter is per-instance: with more than one machine a client gets one bucket
// per machine. The quota, not this, is the money guard.
type Limiter struct {
	mu      sync.Mutex
	rate    float64
	burst   float64
	buckets map[string]*bucket
	lastGC  time.Time
}

type bucket struct {
	tokens float64
	seen   time.Time
}

func New(perMinute int) *Limiter {
	return &Limiter{
		rate:    float64(perMinute) / 60,
		burst:   float64(perMinute),
		buckets: make(map[string]*bucket),
		lastGC:  time.Now(),
	}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if now.Sub(l.lastGC) > sweepInterval {
		l.sweep(now)
	}

	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.burst, seen: now}
		l.buckets[key] = b
	} else {
		b.tokens = min(l.burst, b.tokens+now.Sub(b.seen).Seconds()*l.rate)
		b.seen = now
	}

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func (l *Limiter) sweep(now time.Time) {
	l.lastGC = now
	for key, b := range l.buckets {
		if now.Sub(b.seen) > sweepInterval {
			delete(l.buckets, key)
		}
	}
}
