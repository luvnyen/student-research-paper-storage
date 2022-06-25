package repository

import (
	"github.com/luvnyen/student-research-paper-storage/pkg/models"
	"gorm.io/gorm"
)

type PaperRepository interface {
	All() ([]models.Paper, error)
	FindByID(paperID string) (models.Paper, error)
	FindByTitleAuthorAbstract(title string, author string, abstract string) ([]models.Paper, error)
	InsertPaper(paper models.Paper) (models.Paper, error)
}

type paperConnection struct {
	connection *gorm.DB
}

func NewPaperRepository(db *gorm.DB) PaperRepository {
	return &paperConnection{connection: db}
}

func (db *paperConnection) All() ([]models.Paper, error) {
	papers := []models.Paper{}
	err := db.connection.Preload("Student").Find(&papers)
	if err.Error != nil {
		return papers, err.Error
	}
	return papers, nil
}

func (db *paperConnection) FindByID(paperID string) (models.Paper, error) {
	var paper models.Paper
	err := db.connection.Preload("Student").Where("id = ?", paperID).Take(&paper)
	if err.Error != nil {
		return paper, err.Error
	}
	return paper, nil
}

func (db *paperConnection) FindByTitleAuthorAbstract(title string, author string, abstract string) ([]models.Paper, error) {
	papers := []models.Paper{}

	title = "%" + title + "%"
	author = "%" + author + "%"
	abstract = "%" + abstract + "%"

	err := db.connection.Preload("Student").Where("title LIKE ? AND author LIKE ? AND abstract LIKE ?", title, author, abstract).Find(&papers)
	if err.Error != nil {
		return papers, err.Error
	}

	return papers, nil
}

func (db *paperConnection) InsertPaper(paper models.Paper) (models.Paper, error) {
	errSave := db.connection.Save(&paper)
	if errSave.Error != nil {
		return paper, errSave.Error
	}

	errFetch := db.connection.Preload("Student").Find(&paper)
	if errFetch.Error != nil {
		return paper, errFetch.Error
	}

	return paper, nil
}
