package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	usermodel "oravue_backend/internal/http/model"
	"oravue_backend/internal/utils/response"
)

func GetOtp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("sending otp")

		vars := mux.Vars(r)
		phoneNumber := vars["phone_number"]
		response.WriteJson(w, http.StatusAccepted, fmt.Sprintf("success %s", phoneNumber))
	}
}
func VerifyUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating User")
		var userModel usermodel.UserModel
		err := json.NewDecoder(r.Body).Decode(&userModel)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//validate request
		if err := validator.New().Struct(userModel); err != nil {
			var validateErrs validator.ValidationErrors
			errors.As(err, &validateErrs)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		w.Write([]byte("User created"))
	}
}
