package main

import (
 "bufio"
 "context"
 "fmt"
 "net/http"
 "net/url"
 "os"
 "sync"
 "time"
)

const (
 numThreads = 100
 timeout    = 5 * time.Second
)
 
func checkProxy(ctx context.Context, proxy string) error {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Host: proxy,
			}),
		},
		Timeout: timeout,
	}
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://www.google.com", nil)
	if err != nil {
		return fmt.Errorf("err: %w", err)
		}
		
		resp, err := httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("err: %d", resp.StatusCode)
		}
		
		fmt.Printf("Proxy %s is working\n", proxy)
		return nil
	}

func main() {
	file, err := os.Open("proxy.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	
	var wg sync.WaitGroup
	wg.Add(numThreads)
	
	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()
			for scanner.Scan() {
				proxy := scanner.Text()
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				if err := checkProxy(ctx, proxy); err != nil {
					fmt.Printf("Proxy %s is not working: %v\n", proxy, err)
				}
			}
		}()
	}
		
		wg.Wait()
}
