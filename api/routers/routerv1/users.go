package routerv1

import (
	"github.com/decentorganization/topaz/api/authentication"
	"github.com/decentorganization/topaz/api/controllers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// SetUsersRoutes ...
func SetUsersRoutes(r *mux.Router) *mux.Router {
	s := r.PathPrefix("/users").Subrouter()

	s.Handle("", negroni.New(
		negroni.HandlerFunc(controllers.CreateUser),
	)).Methods("POST")

	s.Handle("/me", negroni.New(
		negroni.HandlerFunc(authentication.UserAuth),
		negroni.HandlerFunc(controllers.GetUser),
	)).Methods("GET")

	s.Handle("/me", negroni.New(
		negroni.HandlerFunc(authentication.UserAuth),
		negroni.HandlerFunc(authentication.DoubleAuth),
		negroni.HandlerFunc(controllers.EditUser),
	)).Methods("PUT", "PATCH")

	s.Handle("/me/update-password", negroni.New(
		negroni.HandlerFunc(authentication.UserAuth),
		negroni.HandlerFunc(authentication.DoubleAuth),
		negroni.HandlerFunc(controllers.UpdatePassword),
	)).Methods("PUT", "PATCH")

	s.Handle("/reset-password/generate-token", negroni.New(
		negroni.HandlerFunc(controllers.ResetPasswordGenerateToken),
	)).Methods("POST")

	s.Handle("/reset-password/verify-token", negroni.New(
		negroni.HandlerFunc(controllers.ResetPasswordVerifyToken),
	)).Methods("POST")

	return r
}
