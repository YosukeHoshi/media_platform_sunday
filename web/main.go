package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

func main() {
	http.HandleFunc("/signup", signup)
	http.ListenAndServe(":8080", nil)
}

func signup(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&User{})

	defer db.Close()

	var user User
	json.NewDecoder(r.Body).Decode(&user)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	user.Password = passwordHash

	fmt.Printf("(%%#v) %#v\n", user)
	db.NewRecord(user)
	db.Create(&user)
	fmt.Fprintf(w, "add user")
}

type User struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"size:255"`
	Password []byte `json:"password" gorm:"size:255"`
	Email string `json:"email" gorm:"type:varchar(100);unique_index"`
	Description string `json:"description" gorm:"size:255"`
	Notes []Note
	Sns []Sns
}

type Note struct {
	User User
	Text string `json:"email" gorm:"type:varchar(100)"`
}

type Sns struct {
	User User
	SnsAccount string `json:"sns_account" gorm:"type:varchar(100)"`
}

/*
func (u User) String() string {
	return fmt.Sprintf("%s(%d)", u.Name, u.Age)
}
*/
