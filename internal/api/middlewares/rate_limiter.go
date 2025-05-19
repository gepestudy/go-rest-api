package middlewares

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mutex     sync.Mutex
	visitor   map[string]int
	limit     int
	resetTime time.Duration
}

func NewRatelimiter(limit int, resetTime time.Duration) *rateLimiter {
	rateLimit := &rateLimiter{
		visitor:   make(map[string]int),
		mutex:     sync.Mutex{},
		limit:     limit,
		resetTime: resetTime,
	}

	// start the reset routine
	go rateLimit.resetVisitorCount()

	return rateLimit
}

// misal rate limit 2x dan reset time 5detik
// method ini akan reset semua visitor, jadi ketika visitor a hit 2x maka akan bisa kembali setelah di reset
// tapi ga adil nya, ketika visitor b hit di detik ke 4. dan dia hit lagi 2 detik kemudian maka itunganya dia sudah 1x hit di cycle terbaru, jadi sisa 1 hit lg di sisa 3 detik itu
// dan  jika ada 1jt visitor maka map nya juga akan berisi 1jt dalam kurun waktu 5 detik
func (rl *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mutex.Lock()
		rl.visitor = make(map[string]int)
		fmt.Println("visitor count cleared")
		rl.mutex.Unlock()
	}
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	fmt.Println("Rate Limiter Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Rate Limiter Middleware being returned...")
		rl.mutex.Lock()
		defer rl.mutex.Unlock()

		// remove port dari ip visitor
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		visitorIP := host

		rl.visitor[visitorIP]++
		fmt.Printf("visitor count from %v is %v\n", visitorIP, rl.visitor[visitorIP])

		if rl.visitor[visitorIP] > rl.limit {
			http.Error(w, "Too many request", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("Rate Limiter ends...")
	})
}
