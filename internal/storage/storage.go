package storage

import "github.com/Shubham1-Kumar/students-apis/internal/types"

// we use interfaces which gives us leverage of switching
// databases with minimal changes

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
}
