package database

import (
	"fmt"
	"io"
	"errors"
	"crypto/rand"
	"encoding/base64"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
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

type Note struct {
	gorm.Model
	Title string `json:"title" gorm:"type:varchar(100)"`
	Content string `json:"content" gorm:"type:varchar(100)"`
}

func (user *User) AuthValidate() error {
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

func (user *User) CreateSession() (session Session) {
	id := make([]byte, 32)
	io.ReadFull(rand.Reader, id)
	session = Session{
		UUID: base64.URLEncoding.EncodeToString(id),
		UserId: user.UserId, //gormを上手く使えばSessionのほうで上手くできそう
	}
	return
}

func (note *Note) ContentValidate() error {
	if note.Title == "" {
		return errors.New(fmt.Sprintf("タイトルを入力してください"))
	}
	if note.Content == "" {
		return errors.New(fmt.Sprintf("本文を入力してください"))
	}
	return nil
}

