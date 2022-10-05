package models

type SSPResponse struct {
	Id      string  `json:"id"`
	Tiles   []Tiles `json:"tiles"`
	Context Context `json:"context"`
}

type Tiles struct {
	Id    uint    `json:"id"`
	Width uint    `json:"width"`
	Ratio float32 `json:"ratio"`
}
