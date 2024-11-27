package storage

import "github.com/Amannigam1820/student-api-go/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudent() ([]types.Student, error)
	DeleteStudent(id int64) (string, error)
	UpdateStudent(id int64, name string, age int, email string) (string, types.Student, error)

	// interface for searching and sorting for student function
	//GetStudentByFilter(name string, sortOrder string) ([]types.Student, error)

	// USer Operation

	RegisterUser(username string, password string) (int64, error)
	GetUserByUsername(username string) (types.User, error)
	GetLoggedInUserDetail(username string) (types.User, error)
}
