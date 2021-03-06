<div align="center">
    <h1>~ ratelimit ~</h1>
    <strong>A simple token bucket based rate limiter.</strong><br><br>
    <a href="https://godoc.org/github.com/zekroTJA/ratelimit"><img src="https://godoc.org/github.com/zekroTJA/ratelimit?status.svg" /></a>&nbsp;
    <a href="https://travis-ci.org/zekroTJA/ratelimit" ><img src="https://travis-ci.org/zekroTJA/ratelimit.svg?branch=master" /></a>&nbsp;
    <a href="https://coveralls.io/github/zekroTJA/ratelimit"><img src="https://coveralls.io/repos/github/zekroTJA/ratelimit/badge.svg" /></a>&nbsp;
    <a href="https://goreportcard.com/report/github.com/zekroTJA/ratelimit"><img src="https://goreportcard.com/badge/github.com/zekroTJA/ratelimit"/></a>
<br>
</div>

---

<div align="center">
    <code>go get github.com/zekroTJA/ratelimit</code>
</div>

---

## Intro

This package provides a verry simple, [token bucket](https://en.wikipedia.org/wiki/Token_bucket) based rate limiter with the ability to return the status of the limiter on reserving a token.

[Here](https://godoc.org/github.com/zekroTJA/ratelimit) you can read the docs of this package, generated by godoc.org.

---

## Usage Example

In [examples/webserver](examples/webserver) you can find a simple HTTP REST API design limiting accesses with this rate limit package:

Taking a look on just these two functions in [webserver.go](examples/webserver/webserver.go), for example:

```go
// ...

const (
	limiterLimit = 10 * time.Second
	limiterBurst = 3
)

// ...

// checkLimit is a helper function checking
// the availability of a limiter for the connections
// address or creating it if not existing. Then,
// the availability of tokens will be checked
// to perform an action. This state will be returned
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

	// Getting the limiter for the current connections address
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
```

Our limiter has a total token volume (= `burst`) of 3 tokens and a limit of 1 token per 10 seconds, which means, that every 10 seconds a new token will be added to the token bucket *(until the bucket has is "full")*.

So, if you send 4 HTTP GET requests in a time span of under 10 Seconds to `/api/test`, you will get following result:

```
///////////////////////////
// REQUEST #1

*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/test HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.61.1
> Accept: */*
>
< HTTP/1.1 200 OK
< X-Ratelimit-Limit: 3
< X-Ratelimit-Remaining: 2
< X-Ratelimit-Reset: 0
< Date: Thu, 21 Mar 2019 09:03:11 GMT
< Content-Length: 0
<


///////////////////////////
// REQUEST #2

*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/test HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.61.1
> Accept: */*
>
< HTTP/1.1 200 OK
< X-Ratelimit-Limit: 3
< X-Ratelimit-Remaining: 1
< X-Ratelimit-Reset: 0
< Date: Thu, 21 Mar 2019 09:03:12 GMT
< Content-Length: 0
<


///////////////////////////
// REQUEST #3

*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/test HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.61.1
> Accept: */*
>
< HTTP/1.1 200 OK
< X-Ratelimit-Limit: 3
< X-Ratelimit-Remaining: 0
< X-Ratelimit-Reset: 1553159004253249800
< Date: Thu, 21 Mar 2019 09:03:14 GMT
< Content-Length: 0
<


///////////////////////////
// REQUEST #4

*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /api/test HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.61.1
> Accept: */*
>
< HTTP/1.1 429 Too Many Requests
< X-Ratelimit-Limit: 3
< X-Ratelimit-Remaining: 0
< X-Ratelimit-Reset: 1553159004253249800
< Date: Thu, 21 Mar 2019 09:03:15 GMT
< Content-Length: 0
<
```

As you can see, we are receiving status code `200 OK` exactly 3 times. Also, you can see that the ammount of remaining tokens is returned in a `X-Ratelimit-Remaining` header and also the time until a new token will be generated in the `X-Ratelimit-Reset` header *(as UnixNano timestamp)* when no tokens are remaining. 

This concept ofreturning these information as headers were inspired by the design of the [REST API of discordapp.com](https://discordapp.com/developers/docs/topics/rate-limits).

---

Copyright (c) 2019 zekro Development (Ringo Hoffmann).  
Covered by MIT licence.
