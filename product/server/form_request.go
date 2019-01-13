package server

type CreateProductFormRequest struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	ImageURL string `json:"imageurl"`
}

func (r CreateProductFormRequest) Validate() map[string]interface{} {
	errs := make(map[string]interface{})

	nameErrs := []string{}
	if r.Name == "" {
		nameErrs = append(nameErrs, "Cannot be blank")
	}

	priceErrs := []string{}
	if r.Price == 0 {
		priceErrs = append(priceErrs, "Cannot be blank")
	}

	imageURLErrs := []string{}
	if r.ImageURL == "" {
		imageURLErrs = append(imageURLErrs, "Cannot be blank")
	}

	if len(nameErrs) > 0 {
		errs["name"] = nameErrs
	}
	if len(priceErrs) > 0 {
		errs["price"] = priceErrs
	}
	if len(imageURLErrs) > 0 {
		errs["imageurl"] = imageURLErrs
	}

	return errs
}

type UpdateProductFormRequest struct {
	Name     *string `json:"name,omitempty"`
	Price    *int    `json:"price,omitempty"`
	ImageURL *string `json:"imageurl,omitempty"`
}
