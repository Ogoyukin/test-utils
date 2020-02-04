package models

type RequestParams struct {
	RequestsCount int
	ThreadsCount  int
	Details       RequestDetails
}

type RequestDetails struct {
	Host string
}
