package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var host string
var maxConcurrent int

type ByteSize float64

type result struct {
	size     int
	duration time.Duration
}

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

func init() {
	flag.StringVar(&host, "host", "http://localhost:8000/julia", "HTTP host")
	flag.IntVar(&maxConcurrent, "maxConcurrent", 1,
		"Max number of simultaneous requests")

	flag.Parse()
}

func urlBuilder(urls chan string) {
	hi_x := 6
	low_x := -hi_x

	hi_y := 6
	low_y := -hi_y

	hi_z := 10
	low_z := 1

	i := 50

	for z := low_z; z < hi_z; z++ {
		for x := low_x * z; x < hi_x*z; x++ {
			for y := low_y * z; y < hi_y*z; y++ {
				p := url.Values{}
				p.Add("h", "128")
				p.Add("i", strconv.Itoa(i))
				p.Add("method", "1")
				p.Add("mu_i", "0.32")
				p.Add("mu_r", "0.36237")
				p.Add("w", "128")
				p.Add("x", strconv.Itoa(x))
				p.Add("y", strconv.Itoa(y))
				p.Add("z", strconv.Itoa(z))

				urls <- host + "?" + p.Encode()
			}
		}
	}
	close(urls)
}

func main() {
	urls := make(chan string, maxConcurrent)
	go urlBuilder(urls)
	work(urls)
}

func fetch(url string) int {
	//log.Print("Fetching", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch %q: %s", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return len(body)
}

func timedFetch(url string) result {
	start := time.Now()
	size := fetch(url)
	d := time.Since(start)

	return result{size, d}
}

func work(urls chan string) {
	var total_request, total_size int
	var total_duration time.Duration

	urlCh := make(chan string, maxConcurrent)
	res := make(chan result, maxConcurrent)
	for i := 0; i < maxConcurrent; i++ {
		go func() {
			for url := range urlCh {
				res <- timedFetch(url)
			}
		}()
	}

	printStats := func() {
		log.Printf("Fetched %d urls %d bytes in %s", total_request,
			total_size, total_duration)
		log.Printf("Avg %.2f QPS", float64(total_request)/
			total_duration.Seconds())
		log.Printf("Avg %s/s", ByteSize(float64(total_size)/
			total_duration.Seconds()))
	}

	for url := range urls {
		if total_request != 0 && total_request%1000 == 0 {
			printStats()
		}
		urlCh <- url
		r := <-res
		total_size += r.size
		total_duration += r.duration
		total_request++
	}
	close(urlCh)

	printStats()
}
