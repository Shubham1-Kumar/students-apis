package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Shubham1-Kumar/students-apis/internal/types"
	"github.com/Shubham1-Kumar/students-apis/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("creating a student")
		// here we created a student struct to serialize the incoming information
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
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
		if err := validator.New().Struct(student); err != nil{
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{
			"success": "OK",
		})

		w.Write([]byte("Welcome to students api"))
	}
}
