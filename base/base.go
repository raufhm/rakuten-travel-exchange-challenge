package base

import (
	"encoding/xml"
	"time"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Cube    *Cube1   `xml:"Cube"`
}

type Cube1 struct {
	Cube []Cube2 `xml:"Cube"`
}

type Cube2 struct {
	Time string  `xml:"time,attr"`
	Cube []Cube3 `xml:"Cube"`
}

type Cube3 struct {
	Currency string `xml:"currency,attr" json:"currency"`
	Rate     string `xml:"rate,attr" json:"rate"`
}

type Output struct {
	Time     time.Time
	Currency string
	Rate     float64
}

type Analyze struct {
	Currency string  `json:"currency"`
	Min      float64 `json:"min"`
	Max      float64 `json:"max"`
	Avg      float64 `json:"avg"`
}
