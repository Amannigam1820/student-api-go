package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Amannigam1820/student-api-go/internal/config"
	"github.com/Amannigam1820/student-api-go/internal/types"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
id INTEGER PRIMARY KEY AUTOINCREMENT,
name TEXT,
email TEXT,
age INTEGER

	
	)`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES(?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, nil
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil

}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("select * from students where id = ?")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Age, &student.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}
	return student, nil
}

func (s *Sqlite) GetAllStudent() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("select id, name, age, email from students")
	if err != nil {
		return []types.Student{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return []types.Student{}, err
	}

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Age, &student.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				return []types.Student{}, fmt.Errorf("no student found")
			}
			return []types.Student{}, fmt.Errorf("query error: %w", err)
		}

		students = append(students, student)
	}
	return students, nil

}

func (s *Sqlite) DeleteStudent(id int64) (string, error) {
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return "", err
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	if rowAffected == 0 {
		return "No student found with the given ID", nil

	}
	return "Student deleted successfully", nil
}

func (s *Sqlite) UpdateStudent(id int64, name string, age int, email string) (string, types.Student, error) {
	if id <= 0 {
		return "", types.Student{}, fmt.Errorf("invalid ID: %d", id)
	}

	var existingStudent types.Student
	query := "select id,name,age,email from students where id = ?"
	err := s.Db.QueryRow(query, id).Scan(&existingStudent.Id, &existingStudent.Name, &existingStudent.Age, &existingStudent.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "Student not found", types.Student{}, nil
		}
		return "", types.Student{}, err
	}

	updateQuery := "UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?"

	res, err := s.Db.Exec(updateQuery, name, age, email)
	if err != nil {
		return "", types.Student{}, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", types.Student{}, err
	}
	if rowsAffected == 0 {
		return "No updates were made", existingStudent, nil
	}

	updatedStudent := types.Student{
		Id:    int(id),
		Name:  name,
		Email: email,
		Age:   age,
	}

	return "Student updated successfully", updatedStudent, nil

}
