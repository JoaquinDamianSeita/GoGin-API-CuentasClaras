package dto

type TransformedCategory struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TransformedShowCategory struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
}

type TransformedIndexCategory struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
}

type CategoryCreateRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
