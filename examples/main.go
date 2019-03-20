package main

import (
	"log"
	"time"

	"github.com/zekroTJA/ratelimit"
)

func main() {
	l := ratelimit.NewLimiter(3*time.Second, 5)

	for {
		log.Println(l.ReserveN(1))
		time.Sleep(1000 * time.Millisecond)
	}
}
