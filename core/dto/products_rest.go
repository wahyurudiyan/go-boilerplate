package dto

type AddProductRequest struct {
	Id          uint64 `json:"id,omitempty"`
	Qty         uint64 `json:"stock,omitempty"`
	Name        string `json:"name,omitempty"`
	Price       int64  `json:"price,omitempty"` // store price with the smallest unit like cents
	Description string `json:"description,omitempty"`
}

type AddProductResponse struct {
	Id   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type RetrieveProductRequest struct {
	Ids []uint64 `json:"ids,omitempty"`
}
