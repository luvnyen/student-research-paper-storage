package _paper

import (
	"github.com/luvnyen/student-research-paper-storage/pkg/models"
	_student "github.com/luvnyen/student-research-paper-storage/service/response/student"
)

type PaperResponse struct {
	ID       uint64                   `json:"id"`
	Title    string                   `json:"title"`
	Author   string                   `json:"author"`
	Year     int64                    `json:"year"`
	Abstract string                   `json:"abstract"`
	File     string                   `json:"file"`
	Student  _student.StudentResponse `json:"student,omitempty"`
}

func NewPaperResponse(paper models.Paper) PaperResponse {
	return PaperResponse{
		ID:       paper.ID,
		Title:    paper.Title,
		Author:   paper.Author,
		Year:     paper.Year,
		Abstract: paper.Abstract,
		File:     paper.File,
		Student:  _student.NewStudentResponse(paper.Student),
	}
}

func NewPaperArrayResponse(papers []models.Paper) []PaperResponse {
	var papersResponse []PaperResponse
	for _, paper := range papers {
		papersResponse = append(papersResponse, NewPaperResponse(paper))
	}
	return papersResponse
}
