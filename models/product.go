package models

type Product struct {
	ID           int64  `json:"id" sql:"id" form:"id"`
	ProductName  string `json:"product_name" sql:"product_name" form:"product_name"`
	ProductNum   int64  `json:"product_num" sql:"product_num" form:"product_num"`
	ProductImage string `json:"product_image" sql:"product_image" form:"product_image"`
	ProductUrl   string `json:"product_url" sql:"product_url" form:"product_url"`
}
