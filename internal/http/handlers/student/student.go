package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Shubham1-Kumar/students-apis/internal/storage"
	"github.com/Shubham1-Kumar/students-apis/internal/types"
	"github.com/Shubham1-Kumar/students-apis/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("creating a student")
		// here we created a student struct to serialize the incoming information
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student) // decoding the incoming req into the struct student
		// here we need to check that the error we have is of type io.EOF or not ---> when input would be empty
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		// if it's not an eof error treat it as a general error and process
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// creating student
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
		}

		response.WriteJson(w, http.StatusCreated, map[string](int64){
			"id": lastId,
		})

	}
}

func GetById(storage storage.Storage)http.HandlerFunc{
	return func(w http.ResponseWriter,r*http.Request){
		id := r.PathValue("id") // name should be same as you gave in rout
		slog.Info("getting a student", slog.String("id", id))
        
		intId, err := strconv.ParseInt(id,10,64)
		if err != nil{
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil{
			slog.Error("error in getting user ", slog.String("id",id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}
       
		response.WriteJson(w,http.StatusOK,student)

	}
}