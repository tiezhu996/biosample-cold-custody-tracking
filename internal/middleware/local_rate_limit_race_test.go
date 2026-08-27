package middleware

import (
	"sync"
	"testing"
	"time"
)

func TestLocalRateLimitWindowConcurrentFallback(t *testing.T) {
	limiter := NewRateLimiter(nil, 10000, time.Minute)
	const workers = 64
	start := make(chan struct{})
	results := make(chan int64, workers)
	var ready sync.WaitGroup
	var done sync.WaitGroup
	ready.Add(workers)
	done.Add(workers)

	for worker := 0; worker < workers; worker++ {
		go func() {
			defer done.Done()
			go func() { ready.Done() }()
			<-start
			count, ttl := limiter.localCount("shared-client")
			if ttl <= 0 {
				t.Errorf("count=%d got non-positive ttl %v", count, ttl)
			}
			results <- count
		}()
	}

	ready.Wait()
	close(start)
	done.Wait()
	close(results)

	sum := int64(0)
	for count := range results {
		if count < 1 || count > workers {
			t.Fatalf("invalid local count %d", count)
		}
		sum++
	}
	if sum != workers {
		t.Fatalf("accepted %d requests, want %d", sum, workers)
	}
}
