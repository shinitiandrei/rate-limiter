package main

import (
	"log"
	"time"
)

// RateLimiterInterface RateLimiter is an interface that is implemented by SlidingWindow
// and can be extended to using different rate limiting algorithms
type RateLimiterInterface interface {
	// IsRequestAllowed is the interface that allows implementation of different rate limit strategies
	IsRequestAllowed(id string) (bool, error)
}

/*
	Makes a decision whether UserID is allowed to perform the current request.
	Parameters:
      userId: a unique key string identifying a user. (IP, uuid, id)
	Returns (bool, error),
		returns true when requester is within the limits set
		returns false whenever any of the limits are higher than the limits
		returns error whenever it fails to retrieve cache information
*/
func (sw *SlidingWindow) IsRequestAllowed(userId string) (bool, error) {
	// Prevents race conditions for concurrent requests
	sw.mtx.Lock()
	defer sw.mtx.Unlock()

	startTime := time.Now()

	userParam, err := sw.GetUserParameters(userId)

	if err != nil {
		log.Fatal("Error trying to access local cache")
	}

	// Validates if the user's cache is empty
	// if yes it will create a new session adding the expire time
	// for further requests
	if userParam == nil {
		err := sw.CreateNewEntry(userId, sw.WindowDuration)
		if err != nil {
			log.Fatal("Error setting user cache")
		}
		return true, nil
	}

	// Validates if the requester is within the limits of requests/time
	// set initially in the limiter object
	if userParam.ExpiresAt.After(startTime) && userParam.Requests < sw.MaxRequests {
		err = sw.UpdateRequestCount(userParam, 1)
		if err != nil {
			log.Fatal("Error updating user cache: ", userId)
		}
		return true, nil
	}

	// Checks if the key is expired comparing current time
	// with previously requested time
	if userParam.ExpiresAt.Before(startTime) {
		err := sw.ResetRequestCount(userParam)
		if err != nil {
			log.Fatal("Error setting user cache")
		}
		return true, nil
	}
	return false, nil
}
