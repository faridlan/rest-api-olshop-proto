package w_product

type UpdateRequest struct {
	IdProduct string `json:"id_product,omitempty"`
	Name      string `json:"name,omitempty" validate:"required,max=150,min=1"`
	Price     int    `json:"price,omitempty" validate:"required,number"`
	Quantity  int    `json:"quantity,omitempty" validate:"required,number"`
	ImageUrl  string `json:"image_url,omitempty" validate:"required,max=200,min=1"`
}
