package repositories

import (
	"github.com/ogoyukin/testutils/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

const GitHub = "https://github.com"

func TestGetRequestRepository_Request(t *testing.T) {
	assert.True(t, initGetRequestRepository(models.RequestDetails{Host: GitHub}).Request())
	assert.False(t, initGetRequestRepository(models.RequestDetails{Host: "ascasvsd:3232"}).Request())
	assert.False(t, initGetRequestRepository(models.RequestDetails{Host: "https://ascasvsd.com:3232"}).Request())
}

func BenchmarkGetRequestRepository_Request(b *testing.B) {
	repository := initGetRequestRepository(models.RequestDetails{Host: GitHub})
	for i := 0; i < b.N; i++ {
		repository.Request()
	}
}
