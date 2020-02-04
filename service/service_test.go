package service

import (
	"fmt"
	"github.com/ogoyukin/testutils/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func BenchmarkExecute(b *testing.B) {
	requestCount := 1
	threadCount := 1

	result := makeResult(requestCount, threadCount)
	fmt.Println(fmt.Sprintf("Completed :%v", result.Completed))
	fmt.Println(fmt.Sprintf("Duration  :%v", result.Duration))

}

func TestRequestService_Execute(t *testing.T) {
	var requestCount = 1
	var threadCount = 6

	result := makeResult(requestCount, threadCount)
	assert.True(t, result.Completed+result.Failed == int64(requestCount))
	assert.Equal(t, int(result.Duration.Milliseconds()), 1)

	requestCount = 1000
	threadCount = 10
	result = makeResult(requestCount, threadCount)
	assert.True(t, result.Completed+result.Failed == int64(requestCount))
	assert.True(t, result.Duration.Milliseconds()/int64(requestCount) <= int64(1))
	assert.True(t, result.RequestDuration(threadCount).Milliseconds() <= 2)
	assert.True(t, result.RequestDuration(threadCount).Milliseconds() >= 1)
}

func makeResult(requestCount, threadCount int) models.Result {
	params := models.RequestParams{RequestsCount: requestCount, ThreadsCount: threadCount}
	service := RequestExecutor{params: params, repository: &randomRepositoryMock{}}
	return service.Execute()
}

type randomRepositoryMock struct {
}

func (r *randomRepositoryMock) Request() bool {
	time.Sleep(time.Millisecond)
	return rand.Int()%2 == 0
}

func BenchmarkAsyncExecute(b *testing.B) {
	wg := sync.WaitGroup{}
	var requestCount = 1
	var threadCount = 1

	weight := semaphore.NewWeighted(int64(threadCount))
	params := models.RequestParams{RequestsCount: requestCount, ThreadsCount: threadCount}

	var service = RequestExecutor{params: params, repository: &successRepositoryMock{}}
	var failedChan = make(chan bool, requestCount)

	for i := 0; i < requestCount; i++ {
		service.asyncExecute(&wg, weight, failedChan)
	}
	wg.Wait()
}

func TestRequestService_asyncExecute(t *testing.T) {
	wg := sync.WaitGroup{}
	var requestCount = 100
	var threadCount = 5

	weight := semaphore.NewWeighted(int64(threadCount))
	params := models.RequestParams{RequestsCount: requestCount, ThreadsCount: threadCount}

	var service = RequestExecutor{params: params, repository: &failRepositoryMock{}}
	var failedChan = make(chan bool, requestCount)

	for i := 0; i < requestCount; i++ {
		service.asyncExecute(&wg, weight, failedChan)
	}
	assert.Equal(t, 0, len(failedChan))
	wg.Wait()
	assert.Equal(t, requestCount, len(failedChan))

	requestCount = 1000
	threadCount = 5
	//weight = semaphore.NewWeighted(int64(threadCount))
	params = models.RequestParams{RequestsCount: requestCount, ThreadsCount: threadCount}

	service = RequestExecutor{params: params, repository: &successRepositoryMock{}}
	failedChan = make(chan bool, requestCount)

	for i := 0; i < requestCount; i++ {
		service.asyncExecute(&wg, weight, failedChan)
	}
	wg.Wait()
	assert.Equal(t, 0, len(failedChan))
}

type failRepositoryMock struct {
}

func (r *failRepositoryMock) Request() bool {
	time.Sleep(time.Millisecond)
	return false
}

type successRepositoryMock struct {
}

func (r *successRepositoryMock) Request() bool {
	return true
}
