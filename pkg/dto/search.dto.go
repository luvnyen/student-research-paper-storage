package dto

type SearchDTO struct {
	Title    string `json:"title" form:"title"`
	Author   string `json:"author" form:"author"`
	Abstract string `json:"abstract" form:"abstract"`
}
