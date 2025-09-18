package sqlite

import (
	"database/sql" // it provides all the interfaces to
	"fmt"

	"github.com/Shubham1-Kumar/students-apis/internal/config"
	"github.com/Shubham1-Kumar/students-apis/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

// work with databases

type Sqlite struct {
	Db *sql.DB
}

// it's convention to make the function by name New
func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
       id integer primary key autoincrement,
	   name text,
	   email text,
	   age integer
   )`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (S *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	// here to prevent sql injection first we prepare the query and then bind the values with the data
	statement, err := S.Db.Prepare("INSERT INTO students (name, email, age)VALUES (?,?,?) ")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	// now we need to execute the query
	result, err := statement.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}


func (s * Sqlite) GetStudentById(id int64)(types.Student, error){
	statement, err := s.Db.Prepare("select id, name, email, age from students where id = ? limit 1")
	if err != nil{
		return types.Student{}, err
	}
	defer statement.Close()

	var student types.Student
    
	// order should be same in which you've added into db
	err = statement.QueryRow(id).Scan(&student.Id,&student.Name, &student.Email, &student.Age)
	if err != nil{
	   if err == sql.ErrNoRows{
		return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
	   }
       return types.Student{}, fmt.Errorf("query error: %w" , err)
	}
    
	return student,nil
}