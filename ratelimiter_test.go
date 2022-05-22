package main

import (
	"sync"
	"testing"
	"time"
)

/* TestValidIsRequestAllowed simple single request to verify one request can succeed
Variables:
	duration: the period of time the requester can send X amount of requests before resetting
*/
func TestValidIsRequestAllowed(t *testing.T) {
	var duration, _ = time.ParseDuration("1m")

	var sw = new(SlidingWindow).NewRateLimiter(100, duration)

	var limiterInterface RateLimiterInterface = sw

	result, err := limiterInterface.IsRequestAllowed("user123")

	if err != nil {
		t.Fatalf("%v", err)
	}

	if result != true {
		t.Fatalf("IsRequestAllowed() failed for a single request, should've returned true")
	}
}

/* TestRefusedIsRequestAllowed to get rejected calls from rate limiter
Variables:
	runs: change this to increase/decrease amount of requests
	maxReq: determines the max amount of requests before rejecting calls
	duration: the period of time the requester can send X amount of requests before resetting
*/
func TestRefusedIsRequestAllowed(t *testing.T) {
	var runs = 101

	var count = 0

	var maxReq uint64 = 100

	var duration, _ = time.ParseDuration("1m")

	var sw = new(SlidingWindow).NewRateLimiter(maxReq, duration)

	var limiterInterface RateLimiterInterface = sw

	var wg sync.WaitGroup

	var result = true

	for i := 0; i < runs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count += 1
			res, err := limiterInterface.IsRequestAllowed("123")

			if err != nil {
				t.Errorf("%v", err)
				return
			}

			if res == false {
				result = false
			}
		}()
	}
	wg.Wait()
	if result != false {
		t.Fatalf("IsRequestAllowed() failed to return false when limit is 2 and request is 3. Number of requests: %v", count)
	}
}

/* TestPerformanceIsRequestAllowed to evaluate the performance of the rate limiter
Variables:
	runs: change this to increase/decrease amount of requests
	durationLimit: determines how fast the application should return
	maxReq: determines the max amount of requests before rejecting calls
	duration: the period of time the requester can send X amount of requests before resetting
*/
func TestPerformanceIsRequestAllowed(t *testing.T) {
	var runs = 100

	var durationLimit = 10 * time.Millisecond

	var count = 0

	var maxReq uint64 = 100

	var duration, _ = time.ParseDuration("1m")

	var started = time.Now()

	var sw = new(SlidingWindow).NewRateLimiter(maxReq, duration)

	var limiterInterface RateLimiterInterface = sw

	var wg sync.WaitGroup

	for i := 0; i < runs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			count += 1
			_, err := limiterInterface.IsRequestAllowed("123")

			if err != nil {
				t.Errorf("%v", err)
				return
			}
		}()
	}
	wg.Wait()

	perfResults := time.Now().Sub(started)

	if perfResults > durationLimit {
		t.Errorf("Results: %v", perfResults)
	}
}
