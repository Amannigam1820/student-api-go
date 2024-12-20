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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
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
func (s *Sqlite) RegisterUser(username, password string) (int64, error) {
	stmt, err := s.Db.Prepare("insert into users(username,password) values(?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password)
	if err != nil {
		return 0, nil
	}
	lastd, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastd, nil
}
func (s *Sqlite) GetUserByUsername(username string) (types.User, error) {
	query := ("select Id,username,password from  users where username = ?")
	row := s.Db.QueryRow(query, username)
	var user types.User
	err := row.Scan(&user.Id, &user.Username, &user.Password)
	// fmt.Println(user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil

}

func (s *Sqlite) GetLoggedInUserDetail(username string) (types.User, error) {
	var user types.User
	query := "SELECT id, username, email FROM users WHERE username = ?"
	err := s.Db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return types.User{}, err
	}
	return user, nil
}

// function to apply searching and sorting feature in backend api

// func (s *Sqlite) GetStudentByFilter(name string, sortOrder string) ([]types.Student, error) {
// 	query := "SELECT id, name, age FROM students WHERE 1=1"
// 	var args []interface{}

// 	// Apply name filter if provided
// 	if name != "" {
// 		query += " AND name LIKE ?"
// 		args = append(args, "%"+name+"%")
// 	}

// 	// Validate and apply sorting order
// 	if sortOrder != "asc" && sortOrder != "desc" {
// 		sortOrder = "asc" // Default to ascending if invalid or not provided
// 	}
// 	query += fmt.Sprintf(" ORDER BY age %s", sortOrder)

// 	// Execute query
// 	rows, err := s.Db.Query(query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("query execution failed: %w", err)
// 	}
// 	defer rows.Close()

// 	// Parse query results
// 	var students []types.Student
// 	for rows.Next() {
// 		var student types.Student
// 		if err := rows.Scan(&student.Id, &student.Name, &student.Age); err != nil {
// 			return nil, fmt.Errorf("failed to scan row: %w", err)
// 		}
// 		students = append(students, student)
// 	}

// 	// Check for errors encountered during iteration
// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("row iteration error: %w", err)
// 	}

// 	return students, nil
// }
