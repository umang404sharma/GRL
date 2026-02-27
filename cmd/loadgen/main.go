package main

import (
	"net/http"
	"sync"
	"sync/atomic"
)

func main() {
	var success int64
	var dropped int64

	total := 10000
	wg := sync.WaitGroup{}

	client := &http.Client{}

	for range total {
		wg.Go(func() {
			resp, err := client.Get("http://localhost:8001/api")
			if err != nil {
				return
			}

			if resp.StatusCode == 200 {
				atomic.AddInt64(&success, 1)
			} else {
				atomic.AddInt64(&dropped, 1)
			}
		})
	}

	wg.Wait()

	println("Success:", success)
	println("Dropped:", dropped)
}
