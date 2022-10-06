package models

import (
	"net/http"
	"net/url"
)

type DspConfig struct {
	Endpoint       *url.URL
	RequestHeaders http.Header
}

type DspBidRequestInfo struct {
	DspEndpoint    *url.URL
	DspBidRequest  DspBidRequest
	DspBidResponse DspBidResponse
	HttStatusCode  int
}

type DspBidRequest struct {
	Id      string          `valid:"required,alpha" json:"id"`
	Imp     []RequestDspImp `valid:"required" json:"imp"`
	Context Context         `valid:"required" json:"context"`
}

type RequestDspImp struct {
	Id        uint `valid:"required,numeric" json:"id"`
	Minwidth  uint `valid:"required,numeric" json:"minwidth"`
	Minheight uint `valid:"required,numeric" json:"minheight"`
}

type DspBidResponse struct {
	Id  string           `valid:"required,alpha" json:"id"`
	Imp []ResponseDspImp `valid:"required" json:"imp"`
}

type ResponseDspImp struct {
	Id     uint    `valid:"required,numeric" json:"id"`
	Width  uint    `valid:"required,numeric" json:"width"`
	Height uint    `valid:"required,numeric" json:"height"`
	Title  string  `valid:"required,alpha" json:"title"`
	Url    string  `valid:"required,alpha" json:"url"`
	Price  float32 `valid:"required,numeric" json:"price"`
}
