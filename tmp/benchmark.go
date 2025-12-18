package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	url := "https://localhost:8888/"
	totalRequests := 1000
	concurrency := 50

	fmt.Printf("Starting benchmark for %s\n", url)

	// 1. Check Connectivity & Compression
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept-Encoding", "gzip")
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("ERROR: Failed to connect: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Content-Encoding: %s\n", resp.Header.Get("Content-Encoding"))

	if resp.Header.Get("Content-Encoding") != "gzip" && resp.Header.Get("Content-Encoding") != "br" {
		fmt.Println("WARNING: Compression not detected!")
	} else {
		fmt.Println("SUCCESS: Compression detected.")
	}

	// 2. Latency Benchmark
	fmt.Printf("\nBenchmarking %d requests with %d concurrency...\n", totalRequests, concurrency)
	
	start := time.Now()
	var wg sync.WaitGroup
	sem := make(chan bool, concurrency)
	
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		sem <- true
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{Transport: tr} // Re-create or reuse, here valid but slightly inefficient, ideally reuse
			
			r, e := client.Get(url)
			if e == nil && r.StatusCode == 200 {
				mu.Lock()
				successCount++
				mu.Unlock()
				r.Body.Close()
			}
		}()
	}
	
	wg.Wait()
	duration := time.Since(start)
	
	rps := float64(totalRequests) / duration.Seconds()
	avgLatency := duration.Seconds() * 1000 / float64(totalRequests)

	fmt.Printf("\nResults:\n")
	fmt.Printf("Time Taken: %v\n", duration)
	fmt.Printf("Successful Requests: %d/%d\n", successCount, totalRequests)
	fmt.Printf("Requests Per Second: %.2f\n", rps)
	fmt.Printf("Avg Latency (Theoretical): %.2f ms\n", avgLatency)
}
