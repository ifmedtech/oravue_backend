package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log/slog"
	"math/rand"
	"net/http"
	"oravue_backend/internal/config"
	"oravue_backend/internal/utils/response"
	"oravue_backend/pkg/jwt"
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
		expiry := time.Now().Add(10 * time.Minute)

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

func VerifyOtp(repository UserRepository, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var request struct {
			PhoneNumber string `json:"phone_number" validate:"required,len=10,numeric"`
			OTP         string `json:"otp" validate:"required,len=6,numeric"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			slog.Error("Failed to decode request", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate request fields
		if err := validator.New().Struct(request); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		// Verify OTP
		userID, err := repository.VerifyOtpRepository(request.PhoneNumber, request.OTP)
		if err != nil {
			slog.Error("Failed to verify OTP", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		if userID == "" {
			response.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"status":  "error",
				"message": "Invalid or expired OTP",
			})
			return
		}

		// Generate JWT token
		token, err := jwt.GenerateJWT(userID, config)
		if err != nil {
			slog.Error("Failed to generate JWT", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		// Send success response with token
		response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"status":  "success",
			"message": "OTP verified successfully",
			"token":   token,
		})
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
