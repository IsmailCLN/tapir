package httpclient

import (
	"context"
	"net"
	"net/http"
	"time"
)

var transport = &http.Transport{
	MaxIdleConns:        100,
	MaxIdleConnsPerHost: 10,
	IdleConnTimeout:     90 * time.Second,
	DialContext: (&net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	TLSHandshakeTimeout:   5 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var client = &http.Client{
	Transport: transport,
	Timeout:   15 * time.Second,
}

func Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return client.Do(req.WithContext(ctx))
}
