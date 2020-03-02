package account

import (
"encoding/json"
"fmt"
"net/http"
"errors"
"log"

"golang.org/x/crypto/bcrypt"
	"github.com/YosukeHoshi/media_platform_sunday/internal/database"
_ "github.com/go-sql-driver/mysql"
_ "github.com/jinzhu/gorm/dialects/mysql"
)


func Signup(w http.ResponseWriter, r *http.Request) {
	err := HandleOnlyPost(w, r)
	if err != nil {
		PrintLog(w, err.Error())
		return
	}

	var user database.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		PrintLog(w, "Failed To Decode Json.")
		return
	}

	if err = user.AuthValidate(); err != nil {
		PrintLog(w, err.Error())
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		PrintLog(w, err.Error())
		return
	}

	user.Password = string(passwordHash)

	session := user.CreateSession()
	cookie := &http.Cookie{
		Name: "session_id",
		Value: session.UUID,
	}
	http.SetCookie(w, cookie)

	if !database.Db.NewRecord(user) {
		return
	}
	database.Db.Create(&user)
	PrintLog(w,"add user")

	if !database.Db.NewRecord(session) {
		return
	}
	database.Db.Create(&session)
	PrintLog(w,"set cookie")
}

func Signin(w http.ResponseWriter ,r *http.Request) {
	err := HandleOnlyPost(w, r)
	if err != nil {
		PrintLog(w, err.Error())
		return
	}

	var user database.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		PrintLog(w, err.Error())
	}
	if err := user.AuthValidate(); err != nil {
		fmt.Println(err)
	}

	var savedUser database.User
	if database.Db.First(&savedUser, &database.User{UserId:user.UserId}).RecordNotFound() {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Please sign up."))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Password is wrong."))
		return
	}

	session := user.CreateSession()
	cookie := &http.Cookie{
		Name: "session_id",
		Value: session.UUID,
	}
	http.SetCookie(w, cookie)

	w.Write([]byte("Signed In."))
}

func Signout(w http.ResponseWriter, r *http.Request) {
	err := HandleOnlyPost(w, r)
	if err != nil {
		PrintLog(w, err.Error())
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
		fmt.Println("Session Record Not Found.")
		return
	}
	database.Db.Delete(&session)

	w.Write([]byte("Signed out."))
}

func HandleOnlyPost(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return errors.New(fmt.Sprintf("Method Not Allowed. Only POST Is Available."))
	}
	return nil
}

func PrintLog(w http.ResponseWriter, st string)  {
	w.Write([]byte(st))
	log.Println(st)
}

