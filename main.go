package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type TestResults struct {
	TotalRequests    int
	SuccessRequests  int
	StatusCodeCounts map[int]int
	TotalDuration    time.Duration
}

func main() {
	url := flag.String("url", "", "URL do serviço a ser testado")
	requests := flag.Int("requests", 100, "Número total de requests")
	concurrency := flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	flag.Parse()

	if *url == "" {
		log.Fatal("A URL deve ser fornecida")
	}

	results := runLoadTest(*url, *requests, *concurrency)

	printReport(results)
}

func runLoadTest(url string, totalRequests, concurrency int) *TestResults {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	results := &TestResults{
		TotalRequests:    totalRequests,
		StatusCodeCounts: make(map[int]int),
	}

	startTime := time.Now()

	requestFunc := func() {
		defer wg.Done()
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		resp, err := http.Get(url)
		if err != nil {
			log.Println("Erro ao fazer requisição:", err)
			return
		}
		defer resp.Body.Close()

		results.StatusCodeCounts[resp.StatusCode]++
		if resp.StatusCode == http.StatusOK {
			results.SuccessRequests++
		}
	}

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go requestFunc()
	}

	wg.Wait()

	results.TotalDuration = time.Since(startTime)
	return results
}

func printReport(results *TestResults) {
	fmt.Printf("\nRelatório de Teste de Carga:\n")
	fmt.Printf("----------------------------\n")
	fmt.Printf("Total de requisições: %d\n", results.TotalRequests)
	fmt.Printf("Requisições com sucesso (HTTP 200): %d\n", results.SuccessRequests)
	fmt.Println("Distribuição de códigos HTTP:")
	for code, count := range results.StatusCodeCounts {
		fmt.Printf("  HTTP %d: %d\n", code, count)
	}
	fmt.Printf("\nTempo total de execução: %v\n", results.TotalDuration)
}
