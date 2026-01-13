package model

type Catalog struct {
	Id        int    `json:"id"`
	Parid     int    `json:"parId"`
	Kodcat    int    `json:"kodcat"`
	Kod       int    `json:"kod"`
	Name      string `json:"name"`
	Comm      string `json:"comm"`
	DopInt1   int    `json:"dopInt1"`
	DopInt2   int    `json:"dopInt2"`
	DopFloat1 int    `json:"dopFloat1"`
	DopFloat2 int    `json:"dopFloat2"`
	DopBit1   bool   `json:"dopBit1"`
	DopBit2   bool   `json:"dopBit2"`
	Archive   bool   `json:"archive"`
}
