package entity

type Product struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Price int32  `json:"price"`
}
