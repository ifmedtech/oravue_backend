package user

import (
	"github.com/gorilla/mux"
	"oravue_backend/internal/config"
)

func Routes(api *mux.Router, userRepository UserRepository, config *config.Config) {
	router := api.PathPrefix("/user").Subrouter()
	router.Handle("/otp/{phone_number}", GetOtp(userRepository, config)).Methods("GET")
	router.Handle("/verify", VerifyOtp(userRepository, config)).Methods("POST")

}
