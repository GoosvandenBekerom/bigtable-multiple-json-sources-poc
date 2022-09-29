package data

import "encoding/json"

type Product struct {
	Id          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type Offer struct {
	ProductID    string `json:"-"`
	PriceInCents int    `json:"price_in_cents"`
}

type Review struct {
	ProductID string `json:"-"`
	Rating    int    `json:"rating"`
	Message   string `json:"message,omitempty"`
}

type AggregatedProduct struct {
	Product json.RawMessage `json:"product"`
	Offer   json.RawMessage `json:"offer"`
	Review  json.RawMessage `json:"review"`
}
