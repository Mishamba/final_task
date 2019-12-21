package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	db "github.com/Mishamba/FinalTaks/data_process"
	"github.com/Mishamba/FinalTaks/model"
	"github.com/Mishamba/FinalTask/jwt_process/jwt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var conn *sql.DB

func main() {
	err := db.ConnectToSQL(conn)
	if err != nil {
		return
	}
	router := mux.NewRouter()
	router.HandleFunc("/", handleHomePage).Methods(http.MethodGet)
	router.HandleFunc("/register", handleRegistrationForm).Methods(http.MethodPost)
	router.HandleFunc("/register/data", handleRegistration).Methods(http.MethodPost)
	router.HandleFunc("/login/data", handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/tweets", handleTweetForm).Methods(http.MethodPost)
	router.HandleFunc("/tweets/post", handleTweetCreate).Methods(http.MethodPost)
	router.HandleFunc("/tweets", handleTweetView).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8083", router))
}

func handleRegistrationForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "register.html")
}

//	In this function user will create his account
func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	newUser := model.User{
		Name:     username,
		Email:    email,
		Password: password,
	}
	err := db.AddUser(newUser, conn)
	if errorCheck(w, err) {
		http.ServeFile(w, r, "error.html")

		return
	}

	jwtToSend, err := jwt.GeneratedToken(newUser)
	if errorCheck(w, err) {
		http.ServeFile(w, r, "error.html")

		return
	}

	w.Write(jwtToSend)

	//	TODO
	// need save token in some form in database
}

//	In this function user will connect to his account
func handleLogin(w http.ResponseWriter, r *http.Request) {

}

func handleTweetForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "post_form.html")
}

//	In this function user will create and post tweet
func handleTweetCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")

		return
	}

	iD, err := strconv.Atoi(r.FormValue("userID"))
	if err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")
		return
	}
	db.PostTweet(r.FormValue("post"), iD, conn) //add possibility, to recognize user by JWT
	http.ServeFile(w, r, "post_answer.html")
}

//	In this function user will view tweets
func handleTweetView(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")

		return
	}

	iD, err := strconv.Atoi(r.FormValue("userID")) //JWT
	if err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")

	}
	tweetsString, err := db.GetTweets(iD, conn)
	if err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")

	}

	tweets, err := json.Marshal(tweetsString)
	if err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(tweets)
}

func handleHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}

func errorCheck(w http.ResponseWriter, err error) bool {
	if err != nil {
		fmt.Fprint(w, err)
		return true
	}

	return false
}
