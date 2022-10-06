package models

type SSPResponse struct {
	Id      string  `valid:"required,alpha" json:"id"`
	Tiles   []Tiles `valid:"required" json:"tiles"`
	Context Context `valid:"required" json:"context"`
}

type Tiles struct {
	Id    uint    `valid:"required,numeric" json:"id"`
	Width uint    `valid:"required,numeric" json:"width"`
	Ratio float32 `valid:"required,numeric" json:"ratio"`
}
