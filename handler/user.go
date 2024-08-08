package handler

import (
	"awesomeProject/dbhelper"
	"awesomeProject/models"
	"awesomeProject/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

func UserRoutes(router *mux.Router) *mux.Router {
	router.HandleFunc("/user", UserRoute).Methods("GET")
	router.HandleFunc("/sign-up", SignUp).Methods("POST")
	//signup
	router.HandleFunc("/login", Login).Methods("POST")
	//login
	router.HandleFunc("/forget-password", ForgotPassword).Methods("POST")
	//forget-password
	router.HandleFunc("/reset-password", ResetPassword).Methods("POST")
	//reset-password

	//logout

	return router
}
func SignUp(w http.ResponseWriter, r *http.Request) {
	//todo!!! SIGNUP
	//signup:
	//create user
	//email and user_name should be unique.
	//store password using using hash

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "failed to parsed user body")
		return
	}
	v := validator.New()
	err := v.Struct(user)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "failed to validate user")
		return
	}
	var available bool
	available, err = util.EmailAvailable(user.Email)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "failed to validate email")
		return
	}
	if available {
		util.RespondError(w, http.StatusBadRequest, err, "email is already registered")
		return
	}
	available, err = util.UsernameAvailable(user.Username)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "failed to validate username")
		return
	}

	if available {
		util.RespondError(w, http.StatusBadRequest, err, "username is already registered")
		return
	}

	user.Password, err = util.HashPassword(user.Password)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "failed to hash password")
		return
	}
	err = dbhelper.CreateUser(user)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err, "failed to create user")
		return
	}
	util.ResponseJSON(w, http.StatusOK, "user created successfully")

}

//todo!!! LOGIN
//payload: input{username,password}
//response:{token,message,name,email,is_admin,[]permission}
//check username is available if yes then proceed
//get user info with hash password that is stored in db

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	v := validator.New()
	err := v.Struct(input)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "failed to validate body")
		return
	}

	var available bool
	available, err = util.UsernameAvailable(input.Username)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "failed to validate username")
		return
	}
	if !available {
		util.RespondError(w, http.StatusInternalServerError, err, "username is not exist ")
		return
	}

	user, err := dbhelper.GetUserInfo(input.Username)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err, "failed to get user")
		return
	}
	if user.Username == "" {
		util.RespondError(w, http.StatusInternalServerError, err, "user is not exist")
		return
	}
	matchPassword := util.CheckPasswordHash(input.Password, user.Password)
	if !matchPassword {
		util.RespondError(w, http.StatusUnauthorized, errors.New("incorrect password"), "password is not match")
		return
	}

	token := util.GenerateSession()
	err = util.CreateSession(token, user.ID)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err, "failed to create session")
		return
	}
	err = util.UpdateUserLastLogin(user.ID)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err, "failed to update last login")
		return
	}
	util.ResponseJSON(w, http.StatusOK, map[string]string{
		"token":    token,
		"username": user.Username,
	})

}

//todo!!compare hash-password(input,db) if match then proceed
//todo!!generate session token which is session id
//insert into sessions(id,user_id) values(token,userID)
//get permissions of that users
//update users last login

//todo!!! forget password

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	//
}
func ResetPassword(w http.ResponseWriter, r *http.Request) {

}
func UserRoute(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "hello test")
	if err != nil {
		return
	}
}
