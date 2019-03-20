package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zekroTJA/ratelimit"
)

const (
	limiterLimit = 10 * time.Second
	limiterBurst = 3
)

type WebServer struct {
	s        *http.Server
	mux      *http.ServeMux
	limiters map[string]*ratelimit.Limiter
}

func NewWebServer(addr string) *WebServer {
	ws := &WebServer{
		s: &http.Server{
			Addr: addr,
		},
		mux:      http.NewServeMux(),
		limiters: make(map[string]*ratelimit.Limiter),
	}

	ws.mux.HandleFunc("/api/test", ws.handlerTest)

	ws.s.Handler = ws.mux

	return ws
}

func (ws *WebServer) Start() error {
	return ws.s.ListenAndServe()
}

func (ws *WebServer) checkLimit(w http.ResponseWriter, r *http.Request) bool {
	addr := r.RemoteAddr
	if strings.Contains(addr, ":") {
		split := strings.Split(addr, ":")
		addr = strings.Join(split[0:len(split)-1], ":")
	}

	limiter, ok := ws.limiters[addr]
	if !ok {
		limiter = ratelimit.NewLimiter(limiterLimit, limiterBurst)
		ws.limiters[addr] = limiter
	}

	a, res := limiter.Reserve()

	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(res.Burst))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
	w.Header().Set("X-RateLimit-Reset", res.Reset.Format(time.RFC3339))

	return a
}

func (ws *WebServer) handlerTest(w http.ResponseWriter, r *http.Request) {
	if ws.checkLimit(w, r) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
	}
}
