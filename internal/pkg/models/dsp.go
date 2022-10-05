package models

type DspBidRequest struct {
	Id      string       `json:"id"`
	Imp     []RequestImp `json:"imp"`
	Context Context      `json:"context"`
}

type RequestImp struct {
	Id        uint `json:"id"`
	Minwidth  uint `json:"minwidth"`
	Minheight uint `json:"minheight"`
}

type DspBidResponse struct {
	Id  string        `json:"id"`
	Imp []ResponseImp `json:"imp"`
}

type ResponseImp struct {
	Id     uint    `json:"id"`
	Width  uint    `json:"width"`
	Height uint    `json:"height"`
	Title  string  `json:"title"`
	Url    string  `json:"url"`
	Price  float32 `json:"price"`
}
