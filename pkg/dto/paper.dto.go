package dto

import "mime/multipart"

type PaperDTO struct {
	Title    string               `json:"title" form:"title" binding:"required"`
	Author   string               `json:"author" form:"author" binding:"required"`
	Year     int64                `json:"year" form:"year" binding:"required"`
	Abstract string               `json:"abstract" form:"abstract" binding:"required"`
	File     multipart.FileHeader `json:"file" form:"file" binding:"required"`
}
