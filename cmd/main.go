package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "github.com/Mishamba/final_task/data_process"
	jwt "github.com/Mishamba/final_task/jwt_process"
	"github.com/Mishamba/final_task/model"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var conn *sql.DB

func main() {
	err := db.SQLStart(conn)
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

//	In this finction we serve html form to register user
func handleRegistrationForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "register.html")
}

//	In this function user will create his account. After that we will send for him JWT token.
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

	jwtToSend, expireTime, err := jwt.GeneratedToken(newUser.Name)
	if errorCheck(w, err) {
		http.ServeFile(w, r, "error.html")

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   jwtToSend,
		Expires: expireTime,
	})
}

//	In this function user will login to his account. If login passed succesfuly we will send for him JWT token.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			username := r.FormValue("username")
			password := r.FormValue("password")

			_, err := db.FindUser(model.User{Name: username, Password: password}, conn)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
			}

			jwtToSend, expireTime, err := jwt.GeneratedToken(username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "jwt",
				Value:   jwtToSend,
				Expires: expireTime,
			})

			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := c.Value

	ok, status, username, err := jwt.DecodeToken(tokenString)
	if !ok || err != nil {
		w.WriteHeader(status)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(username + " authorized succesfuly. now u have jwt token"))
}

func handleTweetForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "post_form.html")
}

//	In this function user will create and post tweet
func handleTweetCreate(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := c.Value

	ok, status, username, err := jwt.DecodeToken(tokenString)
	if !ok || err != nil {
		w.WriteHeader(status)
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")
		return
	}

	user, err := db.FindUser(model.User{Name: username}, conn)

	db.PostTweet(r.FormValue("post"), user.ID, conn)
	http.ServeFile(w, r, "post_answer.html")
}

//	In this function user will view tweets
func handleTweetView(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := c.Value

	ok, status, username, err := jwt.DecodeToken(tokenString)
	if !ok || err != nil {
		w.WriteHeader(status)
		return
	}

	user, err := db.FindUser(model.User{Name: username}, conn)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprint(w, err)
		http.ServeFile(w, r, "error.html")

		return
	}

	tweetsString, err := db.GetTweets(user.ID, conn)
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
