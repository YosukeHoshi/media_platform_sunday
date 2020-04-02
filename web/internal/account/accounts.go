package account

import (
	"encoding/json"
	// "errors"
	// "fmt"
	"log"
	"net/http"

	"github.com/YosukeHoshi/media_platform_sunday/internal/database"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

// SignInRequired
func SignInRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isSignIn(r) {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Please Sign In."))
		}
	}
}

// isSignIn
func isSignIn(r *http.Request) bool {
	cookie, err := r.Cookie("session_id")
	// if err == http.ErrNoCookie {
	// 	log.Println("Cookie: ", err)
	// 	return false
	// }
	if err != nil {
		log.Println("Cookie: ", err)
		return false
	}

	_, err = database.GetSession(cookie.Value)
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	return true
}

// GetCookie is get session id for check
func GetCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(("session_id"))
	if err != nil {
		log.Fatal("Cookie: ", err)
	}
	v := cookie.Value
	log.Println(v)
}

// Signup is sign up
func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed. Only POST Is Available."))
		return
	}

	var user database.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = user.AuthValidate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = string(passwordHash)

	session := user.CreateSession()
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: session.UUID,
	}
	http.SetCookie(w, cookie)

	if database.Db.NewRecord(user) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("this user infomation is already used."))
		log.Println("this user infomation is already used.")
		return
	}
	database.Db.Create(&user)
	log.Println("add user")

	if !database.Db.NewRecord(session) {
		return
	}
	database.Db.Create(&session)
	log.Println("set cookie")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("sign up"))
	log.Println("sign up")
	return
}

// Signin is sign in
func Signin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user database.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if err := user.AuthValidate(); err != nil {
		log.Println(err)
	}

	var savedUser database.User
	if database.Db.First(&savedUser, &database.User{UserName: user.UserName}).RecordNotFound() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Please sign up."))
		log.Println("Please sign up.")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Password is wrong."))
		log.Println("Password is wrong.")
		return
	}

	session := user.CreateSession()
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: session.UUID,
	}
	http.SetCookie(w, cookie)
	log.Println("set cookie")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("sign in."))
	log.Println("sign in")
	return
}

// Signout is sign out
func Signout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil || cookie.Value == "" {
		return
	}
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	var session database.Session
	if database.Db.First(&session, &database.Session{UUID: cookie.Value}).RecordNotFound() {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Session Record Not Found.")
		return
	}
	database.Db.Delete(&session)
	log.Println("set cookie")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("sign out"))
	log.Println("sign out")
	return
}
