package main

import (
    "bufio"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "sync"
)

func checkProxy(proxy string, ch chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    transport := &http.Transport{
        Proxy: http.ProxyURL(&proxy),
    }
    client := &http.Client{Transport: transport}
    req, err := http.NewRequest("GET", "https://www.google.com/", nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
    resp, err := client.Do(req)
    if err != nil {
        ch <- fmt.Sprintf("Failed: %v", err)
        return
    }
    if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
        ch <- fmt.Sprintf("Success: %v", proxy)
    } else {
        ch <- fmt.Sprintf("Failed: %v", proxy)
    }
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: checkproxy <threads> <proxyfile>")
        os.Exit(1)
    }
    numThreads, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    proxyFile, err := os.Open(os.Args[2])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer proxyFile.Close()
    scanner := bufio.NewScanner(proxyFile)
    var wg sync.WaitGroup
    ch := make(chan string)
    for i := 0; i < numThreads; i++ {
        wg.Add(1)
        go func() {
            for proxy := range ch {
                checkProxy(proxy, ch, &wg)
            }
        }()
    }
    for scanner.Scan() {
        proxy := scanner.Text()
        if len(proxy) == 0 {
            continue
        }
        if _, err := http.ProxyFromEnvironment(nil); err == nil {
            if len(proxy) > 4 && proxy[:4] == "http" {
                ch <- proxy
            } else {
                ch <- "Failed: " + proxy
            }
        }
    }
    close(ch)
    wg.Wait()
}
