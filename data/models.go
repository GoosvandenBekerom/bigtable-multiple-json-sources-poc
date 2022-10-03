package data

import "encoding/json"

// Product contains product information
type Product struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// Offer contains offer information. Offer has a "many-to-one" relation to Product.
type Offer struct {
	ProductID    string `json:"-"`
	ID           string `json:"id"`
	PriceInCents int    `json:"price_in_cents"`
}

// Review contains review information. Review has a "many-to-one" relation to Product.
type Review struct {
	ProductID string `json:"-"`
	ID        string `json:"id"`
	Rating    int    `json:"rating"`
	Message   string `json:"message,omitempty"`
}

type AggregatedProduct struct {
	Product json.RawMessage   `json:"product"`
	Offers  []json.RawMessage `json:"offers"`
	Reviews []json.RawMessage `json:"reviews"`
}
