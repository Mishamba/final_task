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
	log.Println("server is listening")
	router := mux.NewRouter()
	router.HandleFunc("/", handleHomePage).Methods(http.MethodGet)
	router.HandleFunc("/subscribe", handleSubscribeForm).Methods(http.MethodGet)
	router.HandleFunc("/subscribe{data}", handleSubscribe).Methods(http.MethodPost)
	router.HandleFunc("/register", handleRegistrationForm).Methods(http.MethodGet)
	router.HandleFunc("/register{data}", handleRegistration).Methods(http.MethodPost)
	router.HandleFunc("/login", handleLoginForm).Methods(http.MethodGet)
	router.HandleFunc("/login{data}", handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/tweets", handleTweetForm).Methods(http.MethodGet)
	router.HandleFunc("/tweets/post", handleTweetCreate).Methods(http.MethodPost)
	router.HandleFunc("/tweets", handleTweetView).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8083", router))
}

//	Just ginig form to choose user to subscribe
func handleSubscribeForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "cmd/subscribe_form.html")
}

//	In this function we've got data of user to follow and we will save this info to database
func handleSubscribe(w http.ResponseWriter, r *http.Request) {

}

//	In this finction we serve html form to register user
func handleRegistrationForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "register.html")
}

//	In this function user will create his account. After that we will send for him JWT token.
func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
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
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "error.html")
		return
	}

	jwtToSend, expireTime, err := jwt.GeneratedToken(newUser.Name)
	if errorCheck(w, err) {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "error.html")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   jwtToSend,
		Expires: expireTime,
	})
}

//	Just giving form to login
func handleLoginForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "login.html")
}

//	In this function user will login to his account. If login passed succesfuly we will send for him JWT token.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("jwt")
	if errorCheck(w, err) {
		if err == http.ErrNoCookie {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				http.ServeFile(w, r, "error.html")
				return
			}

			username := r.FormValue("username")
			password := r.FormValue("password")

			_, err := db.FindUser(model.User{Name: username, Password: password}, conn)
			if errorCheck(w, err) {
				log.Println(err)
				w.WriteHeader(http.StatusUnauthorized)
				http.ServeFile(w, r, "error.html")
				return
			}

			jwtToSend, expireTime, err := jwt.GeneratedToken(username)
			if errorCheck(w, err) {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				http.ServeFile(w, r, "error.html")
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
	if errorCheck(w, err) {
		if err == http.ErrNoCookie {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			http.ServeFile(w, r, "error.html")
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

	if err := r.ParseForm(); errorCheck(w, err) {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
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
	if errorCheck(w, err) {
		if err == http.ErrNoCookie {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			http.ServeFile(w, r, "error.html")
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
	if errorCheck(w, err) {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		http.ServeFile(w, r, "error.html")
	}

	if err := r.ParseForm(); errorCheck(w, err) {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		http.ServeFile(w, r, "error.html")
		return
	}

	tweetsString, err := db.GetTweets(user.ID, conn)
	if errorCheck(w, err) {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "error.html")
		return
	}

	tweets, err := json.Marshal(tweetsString)
	if errorCheck(w, err) {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "error.html")
		return
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
