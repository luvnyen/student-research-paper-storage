package service

import (
	"log"

	"github.com/luvnyen/student-research-paper-storage/pkg/dto"
	"github.com/luvnyen/student-research-paper-storage/pkg/models"
	"github.com/luvnyen/student-research-paper-storage/repository"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email, password string) interface{}
	CreateStudent(student dto.RegisterDTO) models.Student
	FindByEmail(email string) models.Student
	IsDuplicateEmail(email string) bool
}

type authService struct {
	studentRepository repository.StudentRepository
}

func NewAuthService(studentRepository repository.StudentRepository) AuthService {
	return &authService{studentRepository: studentRepository}
}

func (service *authService) VerifyCredential(email, password string) interface{} {
	res := service.studentRepository.VerifyCredential(email, password)
	if v, ok := res.(models.Student); ok {
		comparedPassword := comparePassword(v.Password, []byte(password))
		if v.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func (service *authService) CreateStudent(student dto.RegisterDTO) models.Student {
	studentToCreate := models.Student{}
	err := smapping.FillStruct(&studentToCreate, smapping.MapFields(&student))
	if err != nil {
		log.Fatalf("Fail to map student create dto to student model: %v", err)
	}
	res := service.studentRepository.InsertStudent(studentToCreate)
	return res
}

func (service *authService) FindByEmail(email string) models.Student {
	return service.studentRepository.FindByEmail(email)
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.studentRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func comparePassword(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
