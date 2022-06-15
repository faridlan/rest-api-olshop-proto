package w_product

type Response struct {
	IdProduct string `json:"id_product,omitempty"`
	Name      string `json:"name,omitempty"`
	Price     int    `json:"price,omitempty"`
	Quantity  int    `json:"quantity,omitempty"`
	ImageUrl  string `json:"image_url,omitempty"`
	CreatedAt int64  `json:"created_at,omitempty"`
}
