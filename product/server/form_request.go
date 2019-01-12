package server

import (
	"errors"
)

type CreateProductFormRequest struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	ImageURL string `json:"imageurl"`
}

func (r CreateProductFormRequest) Validate() map[string][]error {
	errs := make(map[string][]error)

	if r.Name == "" {
		errs["name"] = append(errs["name"], errors.New("REQUIRED"))
	}
	if r.Price == 0 {
		errs["price"] = append(errs["price"], errors.New("REQUIRED"))
	}
	if r.ImageURL == "" {
		errs["imageurl"] = append(errs["imageurl"], errors.New("REQUIRED"))
	}

	return errs
}

type UpdateProductFormRequest struct {
	Name     *string `json:"name,omitempty"`
	Price    *int    `json:"price,omitempty"`
	ImageURL *string `json:"imageurl,omitempty"`
}
