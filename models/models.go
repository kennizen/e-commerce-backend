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
	Created_at string `json:"createdAt"`
	Updated_at string `json:"updatedAt"`
}
