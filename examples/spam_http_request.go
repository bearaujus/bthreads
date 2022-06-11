package examples

import (
	"log"
	"net/http"
	"time"

	"github.com/Bearaujus/bthreads"
)

func SpamHttpRequest() {
	// Initialize bthread instance
	bt, err := bthreads.NewInstance(&bthreads.Config{
		// Add delay between goroutines before it revoked by endless loop
		GoroutinesDelay: 60 * time.Millisecond,
		// Put your defined funcs into X goroutine
		FuncGoroutinesCount: 2,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Add funcs
	bt.AddFuncs(
		func() bool {
			req, err := http.NewRequest("GET", "https://google.com/search", nil)
			if err != nil {
				return false
			}

			queryParam := req.URL.Query()
			queryParam.Add("q", "test")
			req.URL.RawQuery = queryParam.Encode()

			resp, err := http.Get(req.URL.String())
			if err != nil {
				return false
			}

			return resp.StatusCode == 200
		},
		func() bool {
			req, err := http.NewRequest("PUT", "https://google.com", nil)
			if err != nil {
				return false
			}

			queryParam := req.URL.Query()
			queryParam.Add("test", "test")
			req.URL.RawQuery = queryParam.Encode()

			resp, err := http.Get(req.URL.String())
			if err != nil {
				return false
			}

			return resp.StatusCode == 200
		},
		func() bool {
			req, err := http.NewRequest("POST", "https://google.com", nil)
			if err != nil {
				return false
			}

			queryParam := req.URL.Query()
			queryParam.Add("test", "test")
			req.URL.RawQuery = queryParam.Encode()

			resp, err := http.Get(req.URL.String())
			if err != nil {
				return false
			}

			return resp.StatusCode == 200
		},
	)

	// Start bthread instance
	bt.Start()
}
