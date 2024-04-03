package circuitbreaker

import (
	"sync"
	"time"
)

type State int

const (
	Open State = iota
	HalfOpen
	Closed
)

type CircuitBreaker struct {
	mu            sync.Mutex
	state         State
	failureCount  int
	resetTimeout  time.Duration
	maxFailures   int
	openStartTime time.Time
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:         Closed,
		maxFailures:   maxFailures,
		resetTimeout:  resetTimeout,
		openStartTime: time.Now(),
	}
}

func (cb *CircuitBreaker) AllowRequest() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case Open:
		if time.Since(cb.openStartTime) >= cb.resetTimeout {
			cb.state = HalfOpen
			return true
		}
		return false
	case HalfOpen:
		return true
	case Closed:
		return true
	default:
		return false
	}
}
