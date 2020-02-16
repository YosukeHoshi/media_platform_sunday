package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"errors"
	"io"
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	// データベース
	Dialect = "mysql"

	// ユーザー名
	DBUser = "root"

	// パスワード
	DBPass = "root-mps"

	// プロトコル
	DBProtocol = "tcp(mysql-container:3306)"

	// DB名
	DBName = "data_base"

	// DB文字コード
	DBChar = "charset=utf8mb4&parseTime=True&loc=Local"
)

func connectGorm() *gorm.DB {
	connectTemplate := "%s:%s@%s/%s?%s"
	connect := fmt.Sprintf(connectTemplate, DBUser, DBPass, DBProtocol, DBName, DBChar)
	fmt.Println(connect)
	db, err := gorm.Open(Dialect, connect)

	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("データベースへの接続成功")
	}

	return db
}

var db *gorm.DB

func init()  {
	db = connectGorm()
	db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&User{}, &Note{}, &Sns{}, &Session{})
	// defer db.Close() <= what's this?
}

func main() {
	http.HandleFunc("/signup", signup)
	// http.HandleFunc("/signin", signin)
	http.ListenAndServe(":8080", nil)
}

func handleOnlyPost(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return errors.New(fmt.Sprintf("Method Not Allowed. Only POST Is Available."))
	}
	return nil
}

func printLog(w http.ResponseWriter, st string)  {
	w.Write([]byte(st))
	log.Println(st)
}

func signup(w http.ResponseWriter, r *http.Request) {
	err := handleOnlyPost(w, r)
	if err != nil {
		printLog(w, err.Error())
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		printLog(w, "Failed To Decode Json.")
		return
	}

	if err = user.authValidate(); err != nil {
		printLog(w, err.Error())
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		printLog(w, err.Error())
		return
	}

	user.Password = string(passwordHash)


	// fmt.Printf("(%%#v) %#v\n", user)
	if !db.NewRecord(user) {
		return
	}
	db.Create(&user)
	printLog(w,"add user")

	session := user.CreateSession()
	cookie := &http.Cookie{
		Name: "session_id",
		Value: session.UUID,
	}
	http.SetCookie(w, cookie)
	printLog(w,"set cookie")
}

func (user User) CreateSession() (session Session) {
	id := make([]byte, 32)
	io.ReadFull(rand.Reader, id)
	session = Session{
		UUID: base64.URLEncoding.EncodeToString(id),
		UserId: user.UserId, //gormを上手く使えばSessionのほうで上手くできそう
	}
	return
}

/*
func signin(w http.ResponseWriter ,r *http.Response) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	return
}
*/
func (user User) authValidate() error {
	if user.UserId == "" {
		return errors.New(fmt.Sprintf("User Id is required"))
	}
	if user.Email == "" {
		return errors.New(fmt.Sprintf("Email is required"))
	}
	if user.Password == "" {
		return errors.New(fmt.Sprintf("Email is required"))
	}
	return nil
}

type User struct {
	gorm.Model
	UserId string `json:"user_id" gorm:"primary_key,not null;unique"`
	Password string `json:"password" gorm:"not null"`
	Email string `json:"email" gorm:"not null;type:varchar(100);unique"`
	UserName string `json:"user_name" gorm:"size:255"`
	Description string `json:"description" gorm:"size:255"`
	Notes []Note
	Sns []Sns
}

type Note struct {
	gorm.Model
	User User
	Text string `json:"email" gorm:"type:varchar(100)"`
}

type Sns struct {
	gorm.Model
	User User
	SnsAccount string `json:"sns_account" gorm:"type:varchar(100)"`
}


type Session struct {
	gorm.Model
	UUID string
	UserId string
}

