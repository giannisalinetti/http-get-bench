package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	respTimes []int64
	tr        = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
)

// doGet returns a channel withe the request body
func doGet(url string, ch chan<- string, printRes bool, counter int) {

	beginGet := time.Now()
	c := &http.Client{Transport: tr}
	resp, err := c.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respTime := time.Since(beginGet).Nanoseconds()
	respTimes = append(respTimes, respTime)

	if printRes {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}
		ch <- fmt.Sprintf("| Attempt: %d\t| Url: %s\t| Status: %s\t| Time: %d |\n%s\n", counter, url, resp.Status, respTime, string(body))
	} else {
		ch <- fmt.Sprintf("| Attempt: %d\t| Url: %s\t| Status: %s\t| Time: %d |\n", counter, url, resp.Status, respTime)
	}
}

// getAvg calculates the average response times
func getAvg(s []int64) int64 {
	var tot int64
	for _, i := range s {
		tot = tot + int64(i)
	}
	return tot / int64(len(s))
}

// verifyUrl performs a DNS lookup
func verifyUrl(rawUrl string) error {

	fmt.Printf("Performing host DNS lookup...\t")
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	res, err := net.LookupHost(parsedUrl.Hostname())
	if err != nil {
		return err
	}
	if len(res) == 0 {
		return err
	}

	fmt.Printf("Done\n")
	return nil
}

func main() {

	customUrl := flag.StringP("url", "u", "", "Custom url")
	numGet := flag.IntP("num", "n", 1, "Number of GET")
	printRes := flag.BoolP("print", "p", false, "Print results")
	flag.Parse()

	// Test if url field is empty
	if *customUrl == "" {
		fmt.Println("Syntax error: url cannot be empty.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// If protocol scheme is missing add "http://"
	r, _ := regexp.Compile("http[s]?://")
	if match := r.FindString(*customUrl); match == "" {
		*customUrl = fmt.Sprintf("http://%s", *customUrl)
	}

	// Verify hostname resolution
	err := verifyUrl(*customUrl)
	if err != nil {
		fmt.Println("Url cannot be resolved", err)
		os.Exit(1)
	}

	// Benchmark section begins here
	fmt.Printf("Beginning benchmark...\n\n")

	ch := make(chan string, *numGet)
	var count = 0

	// Start goroutines loop
	start := time.Now()
	for i := 0; i < *numGet; i++ {
		count++
		go doGet(*customUrl, ch, *printRes, count)
	}

	// Print results in order of completion
	for i := 0; i < count; i++ {
		fmt.Print(<-ch)
	}
	close(ch)

	// Store elapsed time
	elapsed := time.Since(start).Nanoseconds()

	// Store average time
	average := getAvg(respTimes)

	// Print results
	fmt.Printf("\nBenchmark completed.\n")
	fmt.Printf("Total number of requests:\t\t\t%d\n", count)
	fmt.Printf("Total elapsed time in nanoseconds:\t\t%d\n", elapsed)
	fmt.Printf("Average time in nanoseconds:\t\t\t%d\n", average)
}
