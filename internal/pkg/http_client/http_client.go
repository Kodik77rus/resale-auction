package http_client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func InitHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{},
	}
}

func (h HttpClient) POST(
	fetchUrl *string,
	data []byte,
	timeout time.Duration,
) (int, []byte, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		*fetchUrl,
		bytes.NewReader(data),
	)
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return h.sendReq(req, timeout)
}

func (h *HttpClient) sendReq(
	req *http.Request,
	timeout time.Duration,
) (int, []byte, error) {
	h.client.Timeout = timeout

	resp, err := h.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}
