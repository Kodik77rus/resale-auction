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
	DspInfo        string
	DspBidRequest  *DspBidRequest
	DspBidResponse *DspBidResponse
}

type DspBidRequest struct {
	Id      string
	Imp     []*RequestDspImp
	Context Context
}

type RequestDspImp struct {
	Id        uint `valid:"required,numeric" json:"id"`
	Minwidth  uint `valid:"required,numeric" json:"minwidth"`
	Minheight uint `valid:"required,numeric" json:"minheight"`
}

type DspBidResponse struct {
	Id  string           `valid:"required,halfwidth" json:"id"`
	Imp []ResponseDspImp `valid:"required" json:"imp"`
}

type ResponseDspImp struct {
	Id     uint    `valid:"required,numeric" json:"id"`
	Width  uint    `valid:"required,numeric" json:"width"`
	Height uint    `valid:"required,numeric" json:"height"`
	Title  string  `valid:"required,halfwidth" json:"title"`
	Url    string  `valid:"required,halfwidth" json:"url"`
	Price  float32 `valid:"required,float" json:"price"`
}
