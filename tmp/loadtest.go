package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Language string `json:"language"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"token"`
	} `json:"data"`
}

func main() {
	baseURL := "https://localhost:8888/api"
	loginURL := baseURL + "/auth/login"
	
	// Configuration
	concurrency := 50
	requestsPerUser := 20
	totalRequests := concurrency * requestsPerUser
	
	username := "admin"
	passwords := []string{"123456", "admin", "password", "root"}
	
	var token string
	var err error

	fmt.Printf("Starting Load Test...\n")
	fmt.Printf("Target: %s\n", baseURL)

	// 1. Authenticate
	fmt.Println("\n[1] Authenticating...")
	
	for _, p := range passwords {
		fmt.Printf("Trying password: %s ... ", p)
		token, err = login(loginURL, username, p)
		if err == nil {
			fmt.Println("Success!")
			break
		}
		fmt.Printf("Failed (%v)\n", err)
	}

	targetURL := "" // Decides which endpoint to test
	if token != "" {
		fmt.Println("Login Successful! Testing Protected Endpoint: /auth/me")
		targetURL = baseURL + "/auth/me"
	} else {
		fmt.Println("Login Failed. Falling back to Public Endpoint: / (Static)")
		targetURL = "https://localhost:8888/"
	}

	// 2. Load Test
	fmt.Printf("\n[2] Running Load Test on %s...\n", targetURL)
	fmt.Printf("Concurrency: %d users\n", concurrency)
	fmt.Printf("Total Requests: %d\n", totalRequests)
	
	start := time.Now()
	var wg sync.WaitGroup
	results := make(chan int64, totalRequests)
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			client := &http.Client{
				Timeout:   10 * time.Second,
				Transport: tr,
			}
			
			for j := 0; j < requestsPerUser; j++ {
				req, _ := http.NewRequest("GET", targetURL, nil)
				if token != "" {
					req.Header.Set("Authorization", "Bearer "+token)
				}
				req.Header.Set("Accept-Encoding", "gzip")

				reqStart := time.Now()
				resp, err := client.Do(req)
				latency := time.Since(reqStart).Milliseconds()
				
				if err == nil && resp.StatusCode == 200 {
					mu.Lock()
					successCount++
					mu.Unlock()
					results <- latency
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
				} else {
					mu.Lock()
					failCount++
					mu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	close(results)
	totalDuration := time.Since(start)

	// 3. Analyze Results
	var totalLatency int64 = 0
	var minLatency int64 = 100000
	var maxLatency int64 = 0

	for lat := range results {
		totalLatency += lat
		if lat < minLatency { minLatency = lat }
		if lat > maxLatency { maxLatency = lat }
	}

	avgLatency := 0.0
	if successCount > 0 {
		avgLatency = float64(totalLatency) / float64(successCount)
	}
	rps := float64(successCount) / totalDuration.Seconds()

	fmt.Println("\n[3] Results Summary")
	fmt.Printf("Total Duration: %v\n", totalDuration)
	fmt.Printf("Successful Requests: %d\n", successCount)
	fmt.Printf("Failed Requests: %d\n", failCount)
	fmt.Printf("Throughput: %.2f req/sec\n", rps)
	fmt.Printf("Latency (ms): Min=%d, Max=%d, Avg=%.2f\n", minLatency, maxLatency, avgLatency)
}

func login(url, user, pass string) (string, error) {
	reqBody, _ := json.Marshal(LoginRequest{Username: user, Password: pass, Language: "en"})
	
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	var res LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	
	if res.Data.Token == "" {
		return "", fmt.Errorf("no token in response")
	}

	return res.Data.Token, nil
}
