package service

import (
	"strconv"

	"github.com/luvnyen/student-research-paper-storage/pkg/dto"
	"github.com/luvnyen/student-research-paper-storage/pkg/models"
	"github.com/luvnyen/student-research-paper-storage/repository"
	_paper "github.com/luvnyen/student-research-paper-storage/service/response/paper"
)

type PaperService interface {
	All(studentID string) (*[]_paper.PaperResponse, error)
	FindByID(paperID string) (*_paper.PaperResponse, error)
	InsertPaper(paperRequest dto.PaperDTO, studentID string, newFileName string) (*_paper.PaperResponse, error)
	FindByTitleAuthorAbstract(searchRequest dto.SearchDTO, studentID string) (*[]_paper.PaperResponse, error)
}

type paperService struct {
	paperRepository repository.PaperRepository
}

func NewPaperService(paperRepository repository.PaperRepository) PaperService {
	return &paperService{paperRepository: paperRepository}
}

func (service *paperService) All(studentID string) (*[]_paper.PaperResponse, error) {
	papers, err := service.paperRepository.All()
	if err != nil {
		return nil, err
	}
	papers_res := _paper.NewPaperArrayResponse(papers)
	return &papers_res, nil
}

func (service *paperService) FindByID(paperID string) (*_paper.PaperResponse, error) {
	paper, err := service.paperRepository.FindByID(paperID)
	if err != nil {
		return nil, err
	}
	paper_res := _paper.NewPaperResponse(paper)
	return &paper_res, nil
}

func (service *paperService) FindByTitleAuthorAbstract(searchRequest dto.SearchDTO, studentID string) (*[]_paper.PaperResponse, error) {
	papers, err := service.paperRepository.FindByTitleAuthorAbstract(searchRequest.Title, searchRequest.Author, searchRequest.Abstract)
	if err != nil {
		return nil, err
	}
	papers_res := _paper.NewPaperArrayResponse(papers)
	return &papers_res, nil
}

func (service *paperService) InsertPaper(paperRequest dto.PaperDTO, studentID string, newFileName string) (*_paper.PaperResponse, error) {
	paper := models.Paper{}
	paper.Title = paperRequest.Title
	paper.Author = paperRequest.Author
	paper.Year = paperRequest.Year
	paper.Abstract = paperRequest.Abstract
	paper.File = newFileName

	id, _ := strconv.ParseInt(studentID, 0, 64)
	paper.StudentID = id

	paper, err := service.paperRepository.InsertPaper(paper)
	if err != nil {
		return nil, err
	}

	paper_res := _paper.NewPaperResponse(paper)
	return &paper_res, nil
}
