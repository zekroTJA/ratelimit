package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zekroTJA/ratelimit"
)

const (
	limiterLimit = 10 * time.Second
	limiterBurst = 3
)

// webServer contains the actual HTTP server
// instance, the ServeMux handling roots and
// a map of limiters.
type webServer struct {
	s   *http.Server
	mux *http.ServeMux
	// limiters is a map binding a limiter
	// to a remote address, so a limiter only
	// counts for one connection
	limiters map[string]*ratelimit.Limiter
}

// newWebServer creates a new instance of
// a webServer and adds the root handlers
func newWebServer(addr string) *webServer {
	// Creating a new webServer witn a new
	// Server, a new ServeMux and the
	// initialized limiters map.
	ws := &webServer{
		s: &http.Server{
			Addr: addr,
		},
		mux:      http.NewServeMux(),
		limiters: make(map[string]*ratelimit.Limiter),
	}

	// Setting handlerTest as handler for
	// root /api/test
	ws.mux.HandleFunc("/api/test", ws.handlerTest)

	// Adding ServeMux instance as Handler
	// for the HTTP Server
	ws.s.Handler = ws.mux

	return ws
}

// start the web server blocking the current
// thread and retuning an error if it fails.
func (ws *webServer) start() error {
	return ws.s.ListenAndServe()
}

// checkLimit is a helper function checking
// the availability of a limiter for the connections
// address or creating it if not existing. Then,
// the availability of tokens will be checked
// to perform an action. This state will be retunred
// as boolean.
func (ws *webServer) checkLimit(w http.ResponseWriter, r *http.Request) bool {
	// Getting the address of the incomming connection.
	// Because you will likely test this with a local connection,
	// the local port number will be attached and differ on every
	// request. So, we need to cut away everything behind the last
	// ":", if existent.
	addr := r.RemoteAddr
	if strings.Contains(addr, ":") {
		split := strings.Split(addr, ":")
		addr = strings.Join(split[0:len(split)-1], ":")
	}

	// Getting the limiter for the current connections addres
	// or create one if not existent.
	limiter, ok := ws.limiters[addr]
	if !ok {
		limiter = ratelimit.NewLimiter(limiterLimit, limiterBurst)
		ws.limiters[addr] = limiter
	}

	// Reserve a token from the limiter.
	a, res := limiter.Reserve()

	// Attach the reservation result to the three headers
	// "X-RateLimit-Limit"
	//    - containing the absolute burst rate
	//      of the limiter,
	// "X-RateLimit-Remaining"
	//    - the number of remaining tickets after
	//      the request
	// "X-RateLimit-Reset"
	//    - the UnixNano timestamp until a new token
	//      will be generated (only if remaining == 0)
	w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", res.Burst))
	w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
	w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", res.Reset.UnixNano()))

	// Return the succeed status of
	// the token request
	return a
}

// handlerTest is the handler for /api/test root
func (ws *webServer) handlerTest(w http.ResponseWriter, r *http.Request) {
	// Check and consume a token from the limiter,
	// if available. If succeed, return status 200,
	// else status 429.
	if ws.checkLimit(w, r) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
}
