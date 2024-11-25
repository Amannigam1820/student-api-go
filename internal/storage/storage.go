package storage

import "github.com/Amannigam1820/student-api-go/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudent() ([]types.Student, error)
	DeleteStudent(id int64) (string, error)
}
