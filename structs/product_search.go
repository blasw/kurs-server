package structs

type ProductsSearch struct {
	Amount     int                `json:"amount" binding:"required"`
	Page       int                `json:"page" binding:"required"`
	Sort       string             `json:"sort"`
	Brand      string             `json:"brand"`
	Name       string             `json:"name"`
	Categories []CategoriesSearch `json:"categories"`
	MinPrice   float32            `json:"minprice"`
	MaxPrice   float32            `json:"maxprice"`
}

type CategoriesSearch struct {
	ID      uint            `json:"id"`
	Details []DetailsSearch `json:"details"`
}

type DetailsSearch struct {
	ID     uint     `json:"id"`
	Values []string `json:"values"`
}
