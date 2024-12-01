package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log/slog"
	"math/rand"
	"net/http"
	usermodel "oravue_backend/internal/http/model"
	"oravue_backend/internal/utils/response"
	"time"
)

func GetOtp(repository UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("sending OTP")

		// Extract phone number from URL variables
		vars := mux.Vars(r)
		phoneNumber := vars["phone_number"]

		// Generate OTP and expiry time
		otp := generateOTP(6)
		expiry := time.Now().Add(5 * time.Minute)

		// Store OTP in the repository
		_, err := repository.GetOtpRepository(phoneNumber, otp, expiry)
		if err != nil {
			slog.Error("Failed to store OTP in repository", slog.String("phone_number", phoneNumber), slog.String("error", err.Error()))
			err := response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			if err != nil {
				slog.Error("Failed to write JSON response", slog.String("error", err.Error()))
				return
			}
			return
		}

		// Optionally, send OTP using an external service
		/*
			err = sendOtpToExternalService(phoneNumber, otp)
			if err != nil {
				slog.Error("Failed to send OTP via external service", slog.String("phone_number", phoneNumber), slog.String("error", err.Error()))
				response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to send OTP")))
				return
			}
		*/

		// Write success response
		err = response.WriteJson(w, http.StatusAccepted, map[string]interface{}{
			"status":  "success",
			"message": fmt.Sprintf("OTP sent successfully to %s", phoneNumber),
		})
		if err != nil {
			slog.Error("Failed to write JSON response", slog.String("error", err.Error()))
		}
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

func generateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	otp := ""
	for i := 0; i < length; i++ {
		otp += string(digits[rand.Intn(len(digits))])
	}
	return otp
}
