package models

type User struct {
	Id         int    `json:"id"`
	Firstname  string `json:"firstname"`
	Middlename string `json:"middlename"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Age        int    `json:"age"`
	Avatar     string `json:"avatar"`
	Password   string `json:"-"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

type Product struct {
	Id           int     `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	Price        float32 `json:"price"`
	Stock        int     `json:"stock"`
	Image        string  `json:"image"`
	Thumbnail    string  `json:"thumbnail"`
	Rating       float32 `json:"rating"`
	Weight       int     `json:"weight"`
	Width        float32 `json:"width"`
	Height       float32 `json:"height"`
	Depth        float32 `json:"depth"`
	Warranty     string  `json:"warranty"`
	Shipping     string  `json:"shipping"`
	Availability string  `json:"availability"`
	ReturnPolicy string  `json:"returnPolicy"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

type ProductReview struct {
	Id        int     `json:"id"`
	ReviewBy  int     `json:"reviewBy"`
	ProductId int     `json:"productId"`
	Review    string  `json:"review"`
	Rating    float32 `json:"rating"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}
