package main

import (
	"flag"
	"fmt"
	"github.com/ogoyukin/testutils/models"
	"github.com/ogoyukin/testutils/service"
)

func main() {
	// # http://localhost:8080/load/assets/2020/1/1/1580500880123.png
	host := flag.String("h", "http://localhost:8089/api/user/1", "Request host")
	requestsCount := flag.Int("r", 10000, "Request count")
	threadsCount := flag.Int("t", 20, "Threads count")
	flag.Parse()
	params := models.RequestParams{
		RequestsCount: *requestsCount,
		ThreadsCount:  *threadsCount,
		Details:       models.RequestDetails{Host: *host},
	}
	result := service.Request(params)

	fmt.Println("		FINISH		")
	fmt.Println(fmt.Sprintf("Completed			: %v", result.Completed))
	fmt.Println(fmt.Sprintf("Failed				: %v", result.Failed))
	fmt.Println(fmt.Sprintf("Time per request	: %v", result.RequestDuration(*threadsCount)))
	fmt.Println(fmt.Sprintf("Time taken for tests: %v", result.Duration))
}
