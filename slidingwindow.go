package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Applying singleton to this variable as it creates the in-memory storage
var initiated *SlidingWindow
var once sync.Once

// A SlidingWindow represents the Sliding window strategy settings to configure the
// rate limiter function
type SlidingWindow struct {
	// Maximum amount of requests the requester can do
	MaxRequests uint64

	// WindowDuration is to determine how long the user can make X
	// requests per duration defined
	WindowDuration time.Duration

	// Prevents race conditions for concurrent requests
	mtx sync.Mutex

	// CacheInterface responsible for creating in-memory storage
	// for caching purposes
	CacheInterface CacheInterface
}

/*
	Creates a new entry in cache if given a given user doesn't exist
	Parameters:
      userId: a unique key string identifying a user. (IP, uuid, id, etc)
      limit: limit of request the requester can make during a specific time window
	Returns error
      errors when there was a problem inserting/retrieving an entry from memory.
*/
func (sw *SlidingWindow) CreateNewEntry(userId string, limit time.Duration) error {
	userCache, err := sw.CacheInterface.Get(userId)

	if err != nil {
		return fmt.Errorf("error: %v", err.Error())
	}

	if userCache != nil && err == nil {
		return fmt.Errorf(
			"user %s has an entry in cache", userId,
		)
	}

	cache := &CacheBody{
		UserID:    userId,
		Requests:  1,
		Limit:     limit,
		ExpiresAt: time.Now().Add(sw.WindowDuration),
	}

	// CacheInterface.Set inserts the CacheBody to in-memory storage
	err = sw.CacheInterface.Set(userId, &cache, sw.WindowDuration)

	if err != nil {
		return fmt.Errorf("error: %v", err.Error())
	}

	return nil
}

/*
	Returns the cached user parameters from in-memory storage
	Parameters:
      userId: a unique key string identifying a user. (IP, uuid, id, etc)
	Returns (*CacheBody, error)
      *CacheBody contains user parameters data.
	  error is returned when there was an issue retrieving it from in-memory storage
*/
func (sw *SlidingWindow) GetUserParameters(userId string) (*CacheBody, error) {
	var cacheBody CacheBody
	userCache, err := sw.CacheInterface.Get(userId)

	if err != nil {
		return nil, fmt.Errorf("error: %v", err.Error())
	}

	if userCache == nil {
		return nil, nil
	} else {
		err = json.Unmarshal([]byte(userCache), &cacheBody)

		if err != nil {
			return nil, fmt.Errorf("error formatting userCache: %v", err.Error())
		}

		return &cacheBody, nil
	}

}

/*
	Updates the request count from a user parameter
	Parameters:
      cache: contains user parameters data that will be used the update itself.
	  n: the number of which the request should be incremented
	Returns error
	  error: is returned when there was an issue updating the in-memory storage
*/
func (sw *SlidingWindow) UpdateRequestCount(cache *CacheBody, n uint64) error {
	cache.Requests += n
	err := sw.CacheInterface.Set(cache.UserID, &cache, sw.WindowDuration)
	if err != nil {
		return fmt.Errorf("error: %v", err.Error())
	}
	return nil
}

/*
	Resets the request count for a specific user
	Parameters:
      cache: contains user parameters data that will be used the update itself.
	Returns error
	  error: is returned when there was an issue updating the in-memory storage
*/
func (sw *SlidingWindow) ResetRequestCount(cache *CacheBody) error {
	cache.Requests = 1
	err := sw.CacheInterface.Set(cache.UserID, &cache, sw.WindowDuration)
	if err != nil {
		return fmt.Errorf("error: %v", err.Error())
	}
	return nil
}

/*
	Creates a new SlidingWindow object with in-memory storage initiated
	Parameters:
      maxReq: maximum amount of requests it will allow a requester to do.
      duration: the window duration of the requests. E.g. 100 requests per 60 minutes
	Returns *SlidingWindow
	  *SlidingWindow: object is returned
*/
func NewRateLimiter(sw *SlidingWindow, maxReq uint64, duration time.Duration) *SlidingWindow {
	sw.MaxRequests = maxReq
	sw.WindowDuration = duration
	return sw
}

/*
	Singleton function
	Parameters:
	Returns *SlidingWindow
	  *SlidingWindow: object is returned
*/
func GetSWInstance() *SlidingWindow {
	once.Do(func() {
		initiated = &SlidingWindow{CacheInterface: InitMemCache()}
	})
	return initiated
}
