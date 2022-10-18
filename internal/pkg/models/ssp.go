package models

type SspRequest struct {
	Id      string  `valid:"required,fullwidth" json:"id"`
	Tiles   []Tiles `valid:"required" json:"tiles"`
	Context Context `valid:"required" json:"context"`
}

type Tiles struct {
	Id    uint    `valid:"required,numeric" json:"id"`
	Width uint    `valid:"required,numeric" json:"width"`
	Ratio float32 `valid:"required,float" json:"ratio"`
}

type SspResponse struct {
	Id  string   `json:"id"`
	Imp []SspImp `json:"imp"`
}

type SspImp struct {
	Id     uint   `json:"id"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}
