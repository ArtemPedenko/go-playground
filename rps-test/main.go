package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	url      = "https://xn----8sbfyhjh4c.xn--p1ai/directus/items/houses?limit=3&fields=*,house_news.*,house_reports.*&page=1&meta=*&" // URL
	duration = 200 * time.Second                                                                                                      // Длительность теста
	rps      = 100                                                                                                                    // Количество запросов в секунду
	workers  = 500                                                                                                                    // Количество горутин (регулируется)

)

var totalRequests int64

func worker(client *http.Client, requests <-chan struct{}, results chan<- time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	for range requests {

		reqNum := atomic.AddInt64(&totalRequests, 1)

		start := time.Now()
		resp, err := client.Get(url)

		elapsed := time.Since(start)

		if err != nil || resp == nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(reqNum, resp.StatusCode)

		if err == nil && resp.StatusCode == http.StatusOK {
			results <- elapsed
			resp.Body.Close()
		}
	}
}

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: workers,
		},
	}

	requests := make(chan struct{}, rps) // Канал для запросов
	results := make(chan time.Duration, rps*int(duration.Seconds()))

	var wg sync.WaitGroup

	// Запуск воркеров
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(client, requests, results, &wg)
	}

	fmt.Printf("Запуск теста: %d RPS на %s в течение %v\n", rps, url, duration)
	startTime := time.Now()

	// Генерация запросов
	ticker := time.NewTicker(time.Second / time.Duration(rps))
	defer ticker.Stop()

	go func() {
		for i := 0; i < rps*int(duration.Seconds()); i++ {
			<-ticker.C
			requests <- struct{}{}
		}
		close(requests)
	}()

	wg.Wait()
	close(results)

	totalTime := time.Since(startTime)
	var sum time.Duration
	var count int
	for res := range results {
		sum += res
		count++
	}

	fmt.Printf("Тест завершен за %v\n", totalTime)
	fmt.Printf("Успешные запросы: %d\n", count)
	if count > 0 {
		fmt.Printf("Среднее время ответа: %v\n", sum/time.Duration(count))
	}
}
