package http_client

import (
	"io/ioutil"
	"net/http"
	"net/url"
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
	fetchUrl *url.URL,
	data []byte,
	headers http.Header,
	timeout time.Duration,
) (int, []byte, error) {

	httpReq := &http.Request{
		Method: http.MethodPost,
		URL:    fetchUrl,
		Header: headers,
	}

	return h.sendReq(httpReq, timeout)
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
