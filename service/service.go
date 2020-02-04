package service

import (
	"context"
	"github.com/ogoyukin/testutils/models"
	"github.com/ogoyukin/testutils/repositories"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

func Request(params models.RequestParams) models.Result {
	repository := repositories.NewRepository(params.Details)
	executor := RequestExecutor{params: params, repository: repository}
	return executor.Execute()
}

type RequestExecutor struct {
	params     models.RequestParams
	repository repositories.RequestRepository
}

func (m *RequestExecutor) Execute() models.Result {
	var failedChan = make(chan bool, m.params.RequestsCount)

	wg := sync.WaitGroup{}
	weight := semaphore.NewWeighted(int64(m.params.ThreadsCount))
	start := time.Now()
	for i := 0; i < m.params.RequestsCount; i++ {
		m.asyncExecute(&wg, weight, failedChan)
	}
	wg.Wait()
	duration := time.Since(start)
	failed := int64(len(failedChan))
	completed := int64(m.params.RequestsCount) - failed
	return models.Result{
		Completed: completed,
		Failed:    failed,
		Duration:  duration,
	}
}

func (m *RequestExecutor) asyncExecute(wg *sync.WaitGroup, weight *semaphore.Weighted, failedChan chan<- bool) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := weight.Acquire(context.Background(), 1)
		if err != nil {
			weight.TryAcquire(1)
		}
		defer func() {
			if err == nil {
				weight.Release(1)
			}
		}()
		if !m.repository.Request() {
			failedChan <- true
		}
	}()
}
