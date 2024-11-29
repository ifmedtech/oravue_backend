package user

import (
	"github.com/gorilla/mux"
)

func Routes(api *mux.Router) {
	router := api.PathPrefix("/user").Subrouter()
	router.Handle("/verify", VerifyUser()).Methods("POST")
	router.Handle("/otp/{phone_number}", GetOtp()).Methods("GET")
}
