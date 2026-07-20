package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements rate limiting per IP and per user using a simple token bucket algorithm
type RateLimiter struct {
	ips       map[string]*ipTracker
	users     map[uint]*userTracker
	mu        sync.RWMutex
	ipLimit   int // requests per minute per IP
	userLimit int // requests per hour per user
}

type ipTracker struct {
	requests []time.Time
	mu       sync.Mutex
}

type userTracker struct {
	requests []time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
// ipLimit: requests per minute per IP (e.g., 5 = 5 requests per minute)
// userLimit: requests per hour per user (e.g., 30 = 30 requests per hour)
func NewRateLimiter(ipLimit, userLimit int) *RateLimiter {
	return &RateLimiter{
		ips:       make(map[string]*ipTracker),
		users:     make(map[uint]*userTracker),
		ipLimit:   ipLimit,
		userLimit: userLimit,
	}
}

// checkIPRateLimit checks if the IP is within rate limits
func (rl *RateLimiter) checkIPRateLimit(ip string) bool {
	rl.mu.Lock()
	tracker, exists := rl.ips[ip]
	if !exists {
		tracker = &ipTracker{requests: []time.Time{}}
		rl.ips[ip] = tracker
	}
	rl.mu.Unlock()

	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	now := time.Now()
	// Remove requests older than 1 minute
	var validRequests []time.Time
	for _, req := range tracker.requests {
		if now.Sub(req) < time.Minute {
			validRequests = append(validRequests, req)
		}
	}
	tracker.requests = validRequests

	// Check if under limit
	if len(tracker.requests) >= rl.ipLimit {
		return false
	}

	// Add current request
	tracker.requests = append(tracker.requests, now)
	return true
}

// checkUserRateLimit checks if the user is within rate limits
func (rl *RateLimiter) checkUserRateLimit(userID uint) bool {
	rl.mu.Lock()
	tracker, exists := rl.users[userID]
	if !exists {
		tracker = &userTracker{requests: []time.Time{}}
		rl.users[userID] = tracker
	}
	rl.mu.Unlock()

	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	now := time.Now()
	// Remove requests older than 1 hour
	var validRequests []time.Time
	for _, req := range tracker.requests {
		if now.Sub(req) < time.Hour {
			validRequests = append(validRequests, req)
		}
	}
	tracker.requests = validRequests

	// Check if under limit
	if len(tracker.requests) >= rl.userLimit {
		return false
	}

	// Add current request
	tracker.requests = append(tracker.requests, now)
	return true
}

// RateLimitByIP limits requests by IP address
func (rl *RateLimiter) RateLimitByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if !rl.checkIPRateLimit(ip) {
			http.Error(w, `{"error":"rate_limit_exceeded","message":"Too many requests from this IP. Please try again later."}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimitByUser limits requests by user ID (requires authentication)
func (rl *RateLimiter) RateLimitByUser(getUserID func(r *http.Request) (uint, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := getUserID(r)
			if err != nil {
				// If we can't get user ID, skip user rate limiting
				next.ServeHTTP(w, r)
				return
			}

			if !rl.checkUserRateLimit(userID) {
				http.Error(w, `{"error":"rate_limit_exceeded","message":"Too many requests from this user. Please try again later."}`, http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitByIPAndUser limits requests by both IP and user ID
func (rl *RateLimiter) RateLimitByIPAndUser(getUserID func(r *http.Request) (uint, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check IP rate limit
			ip := r.RemoteAddr
			if !rl.checkIPRateLimit(ip) {
				http.Error(w, `{"error":"rate_limit_exceeded","message":"Too many requests from this IP. Please try again later."}`, http.StatusTooManyRequests)
				return
			}

			// Check user rate limit if authenticated
			userID, err := getUserID(r)
			if err == nil {
				if !rl.checkUserRateLimit(userID) {
					http.Error(w, `{"error":"rate_limit_exceeded","message":"Too many requests from this user. Please try again later."}`, http.StatusTooManyRequests)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Cleanup removes stale limiters (should be called periodically)
func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// Clean up IP trackers
	for ip, tracker := range rl.ips {
		tracker.mu.Lock()
		var validRequests []time.Time
		for _, req := range tracker.requests {
			if now.Sub(req) < time.Minute {
				validRequests = append(validRequests, req)
			}
		}
		tracker.requests = validRequests
		if len(tracker.requests) == 0 {
			delete(rl.ips, ip)
		}
		tracker.mu.Unlock()
	}

	// Clean up user trackers
	for userID, tracker := range rl.users {
		tracker.mu.Lock()
		var validRequests []time.Time
		for _, req := range tracker.requests {
			if now.Sub(req) < time.Hour {
				validRequests = append(validRequests, req)
			}
		}
		tracker.requests = validRequests
		if len(tracker.requests) == 0 {
			delete(rl.users, userID)
		}
		tracker.mu.Unlock()
	}
}

// Stats returns current rate limiter statistics
func (rl *RateLimiter) Stats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"tracked_ips":   len(rl.ips),
		"tracked_users": len(rl.users),
		"ip_limit":      fmt.Sprintf("%d req/min", rl.ipLimit),
		"user_limit":    fmt.Sprintf("%d req/hour", rl.userLimit),
	}
}
