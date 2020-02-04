package repositories

import (
	"github.com/ogoyukin/testutils/models"
	"net/http"
	"sync"
	"time"
)

func NewRepository(details models.RequestDetails) RequestRepository {
	return initGetRequestRepository(details)
}

type RequestRepository interface {
	Request() bool
}

func initGetRequestRepository(details models.RequestDetails) RequestRepository {
	return &getRequestRepository{
		requestDetails: details,
		httpClientPool: sync.Pool{New: func() interface{} {
			transport := http.Transport{TLSHandshakeTimeout: time.Millisecond * 500}
			client := http.Client{Timeout: 2 * time.Second, Transport: &transport}
			return &client
		}},
	}
}

type getRequestRepository struct {
	requestDetails models.RequestDetails
	httpClientPool sync.Pool
}

func (r *getRequestRepository) Request() bool {
	client := r.httpClientPool.Get().(*http.Client)
	defer r.httpClientPool.Put(client)
	resp, err := client.Get(r.requestDetails.Host)
	if err != nil {
		return false
	}
	err = resp.Body.Close()
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}
