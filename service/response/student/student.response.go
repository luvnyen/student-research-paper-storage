package _student

import "github.com/luvnyen/student-research-paper-storage/pkg/models"

type StudentResponse struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	NRP   string `json:"nrp"`
	Email string `json:"email"`
}

func NewStudentResponse(student models.Student) StudentResponse {
	return StudentResponse{
		ID:    student.ID,
		Name:  student.Name,
		NRP:   student.NRP,
		Email: student.Email,
	}
}
