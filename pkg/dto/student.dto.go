package dto

type StudentDTO struct {
	ID    uint64 `json:"id" form:"id" binding:"required"`
	Name  string `json:"name" form:"name" binding:"required"`
	NRP   string `json:"nrp" form:"nrp" binding:"required"`
	Email string `json:"email" form:"email" binding:"required,email"`
}
