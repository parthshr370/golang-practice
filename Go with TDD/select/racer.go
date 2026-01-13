package racer

import (
	"fmt"
	"net/http"
	"time"
)

func Racer(a, b string, timout time.Duration) (winner string, errror error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(a):
		return a, nil
	case <-time.After(timout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)

	}
}
func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}
func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}
