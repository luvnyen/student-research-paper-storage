package repository

import (
	"log"

	"github.com/luvnyen/student-research-paper-storage/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StudentRepository interface {
	InsertStudent(student models.Student) models.Student
	UpdateUser(student models.Student) models.Student
	VerifyCredential(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) models.Student
	ProfileStudent(id int64) models.Student
}

type studentConnection struct {
	connection *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentConnection{connection: db}
}

func (db *studentConnection) InsertStudent(student models.Student) models.Student {
	student.Password = hashAndSalt([]byte(student.Password))
	db.connection.Save(&student)
	return student
}

func (db *studentConnection) UpdateUser(student models.Student) models.Student {
	if student.Password != "" {
		student.Password = hashAndSalt([]byte(student.Password))
	} else {
		var tempStudent models.Student
		db.connection.Find(&tempStudent, student.ID)
		student.Password = tempStudent.Password
	}
	db.connection.Save(&student)
	return student
}

func (db *studentConnection) VerifyCredential(email string, password string) interface{} {
	var student models.Student
	res := db.connection.Where("email = ?", email).Take(&student)
	if res.Error == nil {
		return student
	}
	return res.Error
}

func (db *studentConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var student models.Student
	return db.connection.Where("email = ?", email).Take(&student)
}

func (db *studentConnection) FindByEmail(email string) models.Student {
	var student models.Student
	db.connection.Where("email = ?", email).Take(&student)
	return student
}

func (db *studentConnection) ProfileStudent(id int64) models.Student {
	var student models.Student
	db.connection.Find(&student, id)
	return student
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash the password")
	}
	return string(hash)
}
