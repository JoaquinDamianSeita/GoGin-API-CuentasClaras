package dto

import "time"

type TransformedOperation struct {
	ID       int                 `json:"id"`
	Type     string              `json:"type"`
	Amount   float64             `json:"amount"`
	Date     time.Time           `json:"date"`
	Category TransformedCategory `json:"category"`
}

type TransformedShowOperation struct {
	ID          int                     `json:"id"`
	Type        string                  `json:"type"`
	Amount      float64                 `json:"amount"`
	Date        time.Time               `json:"date"`
	Category    TransformedShowCategory `json:"category"`
	Description string                  `json:"description"`
}

type TransformedShowCategory struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type CreateOperationRequest struct {
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id"`
}
