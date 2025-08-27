package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Address           string `json:"address"`
	Port              int    `json:"port"`
	Error             string `json:"error,omitempty"`
	TLSVersion        string `json:"tlsVersion"`
	CipherSuite       string `json:"cipherSuite"`
	ServerCertKeyAlgo string `json:"serverCertKeyAlgo"`
	CurveName         string `json:"ecCurve"`
	IsPQCCurve        bool   `json:"isPQCCurve"`
}

func scanTLSPort(host string, port int, timeout time.Duration) (result ScanResult) {
	defer func() {
		if r := recover(); r != nil {
			result = ScanResult{
				Address: host,
				Port:    port,
				Error:   fmt.Sprintf("panic: %v", r),
			}
		}
	}()

	r, err := scanTLSPortWithErr(host, port, timeout)
	if err != nil {
		r.Error = err.Error()
	}
	return r
}

// scanPorts scans a range of ports using multiple goroutines
func scanPorts(host string, start, end int, timeout time.Duration, concurrency int) []ScanResult {
	var results []ScanResult
	var wg sync.WaitGroup

	// Create a buffered channel to collect results
	resultChan := make(chan ScanResult, end-start+1)

	// Create a semaphore to limit concurrent goroutines
	semaphore := make(chan struct{}, concurrency)

	// Launch goroutines for each port
	for port := start; port <= end; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := scanTLSPort(host, p, timeout)
			resultChan <- result
		}(port)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from channel
	for result := range resultChan {
		if result.Error == "" {
			results = append(results, result)
		}
	}

	// Sort results by port number for better readability
	sort.Slice(results, func(i, j int) bool {
		if results[i].Address < results[j].Address {
			return true
		}
		return results[i].Port < results[j].Port
	})

	return results
}

func outputText(w *os.File, results []ScanResult, elapsed time.Duration, verbose bool) {
	fmt.Fprintf(w, "\nScan completed in %s\n", elapsed)
	fmt.Fprintf(w, "Found %d open ports:\n\n", len(results))

	if len(results) == 0 {
		fmt.Fprintf(w, "No open ports found.\n")
		return
	}

	fmt.Fprintf(w, "HOST            \t  PORT\tTLS VER\tCIPHER   \tERROR \n")
	fmt.Fprintf(w, "----------------\t------\t-------\t---------\t------\n")

	for _, result := range results {
		fmt.Fprintf(w, "%-16s\t%6d\t%s\t%s\t%s\n",
			result.Address,
			result.Port,
			result.TLSVersion,
			result.CipherSuite,
			result.Error)
	}
}

func outputJSON(w *os.File, results []ScanResult, elapsed time.Duration) {
	output := struct {
		ScanTime    string       `json:"scan_time"`
		ElapsedTime string       `json:"elapsed_time"`
		OpenPorts   int          `json:"open_tls_ports"`
		Results     []ScanResult `json:"results"`
	}{
		ScanTime:    time.Now().Format(time.RFC3339),
		ElapsedTime: elapsed.String(),
		OpenPorts:   len(results),
		Results:     results,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(output)
}

func outputCSV(w *os.File, results []ScanResult, elapsed time.Duration, verbose bool) {
	fmt.Fprintf(w, "Host,Port,CipherSuite,TLS Version\n")

	for _, result := range results {
		fmt.Fprintf(w, "%s,%d,%s,%s\n",
			result.Address,
			result.Port,
			escapeCSV(result.CipherSuite),
			result.TLSVersion)
	}

	fmt.Fprintf(w, "\n# Scan completed in %s, found %d open ports\n",
		elapsed, len(results))
}

func escapeCSV(s string) string {
	if strings.Contains(s, ",") || strings.Contains(s, "\"") || strings.Contains(s, "\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

func main2() {
	hostPtr := flag.String("host", "", "Target host to scan (required)")
	startPortPtr := flag.Int("start", 1, "Starting port number")
	endPortPtr := flag.Int("end", 1024, "Ending port number")
	timeoutPtr := flag.Int("timeout", 100, "Timeout in milliseconds")
	concurrencyPtr := flag.Int("concurrency", 100, "Number of concurrent scans")
	formatPtr := flag.String("format", "text", "Output format: text, json, or csv")
	verbosePtr := flag.Bool("verbose", false, "Show verbose output including banners")
	outputFilePtr := flag.String("output", "", "Output file (default is stdout)")

	flag.Parse()

	if *hostPtr == "" {
		fmt.Println("Error: host is required")
		fmt.Println("\nUsage examples:")
		fmt.Println("  goscan -host example.com")
		fmt.Println("  goscan -host 192.168.1.1 -start 80 -end 443")
		fmt.Println("  goscan -host example.com -format json -output results.json")
		fmt.Println("\nFor more options:")
		flag.Usage()
		os.Exit(1)
	}

	if *startPortPtr < 1 || *startPortPtr > 65535 {
		fmt.Println("Error: starting port must be between 1 and 65535")
		os.Exit(1)
	}
	if *endPortPtr < 1 || *endPortPtr > 65535 {
		fmt.Println("Error: ending port must be between 1 and 65535")
		os.Exit(1)
	}
	if *startPortPtr > *endPortPtr {
		fmt.Println("Error: starting port must be less than or equal to ending port")
		os.Exit(1)
	}

	if *formatPtr != "text" && *formatPtr != "json" && *formatPtr != "csv" {
		fmt.Println("Error: format must be one of: text, json, or csv")
		os.Exit(1)
	}

	timeout := time.Duration(*timeoutPtr) * time.Millisecond

	var outputFile *os.File
	var err error

	if *outputFilePtr != "" {
		outputFile, err = os.Create(*outputFilePtr)
		if err != nil {
			fmt.Printf("Error creating output file: %v\n", err)
			os.Exit(1)
		}
		defer outputFile.Close()
	} else {
		outputFile = os.Stdout
	}

	fmt.Fprintf(outputFile, "PQC TLS Scan - Scans whether a TLS server is post quantum safe\n")
	fmt.Fprintf(outputFile, "======================================\n")
	fmt.Fprintf(outputFile, "Target: %s\n", *hostPtr)
	fmt.Fprintf(outputFile, "Port range: %d-%d\n", *startPortPtr, *endPortPtr)
	fmt.Fprintf(outputFile, "Timeout: %d ms\n", *timeoutPtr)
	fmt.Fprintf(outputFile, "Concurrency: %d\n", *concurrencyPtr)
	fmt.Fprintf(outputFile, "======================================\n")

	fmt.Fprintf(outputFile, "Scanning %s from port %d to %d...\n", *hostPtr, *startPortPtr, *endPortPtr)
	startTime := time.Now()

	results := scanPorts(*hostPtr, *startPortPtr, *endPortPtr, timeout, *concurrencyPtr)

	elapsed := time.Since(startTime)

	switch *formatPtr {
	case "json":
		outputJSON(outputFile, results, elapsed)
	case "csv":
		outputCSV(outputFile, results, elapsed, *verbosePtr)
	default:
		outputText(outputFile, results, elapsed, *verbosePtr)
	}

	if *outputFilePtr != "" {
		fmt.Printf("Scan complete! Results saved to %s\n", *outputFilePtr)
	}
}
