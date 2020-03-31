package database

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// User is the model of user
type User struct {
	gorm.Model
	UserName    string `json:"user_name" gorm:"not null;unique"`
	Password    string `json:"password" gorm:"not null"`
	Email       string `json:"email" gorm:"not null;type:varchar(100);unique"`
	Name        string `json:"name" gorm:"size:255"`
	Description string `json:"description" gorm:"size:255"`
	Notes       []Note `gorm:"foreignkey:ID"`
	Sns         []Sns
}

// Sns is
type Sns struct {
	gorm.Model
	User       User   `gorm:"foreignkey:UserName"`
	SnsAccount string `json:"sns_account" gorm:"type:varchar(100)"`
}

// Session is
type Session struct {
	gorm.Model
	UUID     string
	UserName string
}

// Note is
type Note struct {
	gorm.Model
	User     User `gorm:"foreignkey:UserName"`
	UserName string
	Title    string `json:"title" gorm:"not null;type:varchar(100)"`
	Content  string `json:"content" gorm:"not null;type:varchar(100)"`
}

// GetSession returns the signed-in user's session
func GetSession(sessionID string) (session Session, err error) {
	err = Db.Where("uuid = ?", sessionID).Find(&session).Error
	return
}

// GetNote is
func GetNote(ID uint) (note Note, err error) {
	err = Db.First(&note, ID).Error
	return
}

// GetMyNotes returns the signed-in user's note
func GetMyNotes(sessionUserName string) (notes []Note, err error) {
	query := Db.Table("notes").
		Select("notes.*, users.user_name").
		Joins("inner join users on notes.user_name = users.user_name").
		Where("notes.user_name = ? and notes.deleted_at is NULL", sessionUserName)
	rows, err := query.Rows()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var note Note
	var user User
	for rows.Next() {
		err = query.ScanRows(rows, &note)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = query.ScanRows(rows, &user)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		note.User = user
		notes = append(notes, note)
	}
	return
}

func (note *Note) UpdateNote(updateNote Note) error {
	return Db.Model(&note).Updates(updateNote).Error
}

// AuthValidate is
func (user *User) AuthValidate() error {
	if user.UserName == "" {
		return errors.New(fmt.Sprintf("User Name is required"))
	}
	if user.Email == "" {
		return errors.New(fmt.Sprintf("Email is required"))
	}
	if user.Password == "" {
		return errors.New(fmt.Sprintf("Password is required"))
	}
	return nil
}

// CreateSession is
func (user *User) CreateSession() (session Session) {
	id := make([]byte, 32)
	io.ReadFull(rand.Reader, id)
	session = Session{
		UUID:     base64.URLEncoding.EncodeToString(id),
		UserName: user.UserName, //gormを上手く使えばSessionのほうで上手くできそう
	}
	return
}

// ContentValidata is
func (note *Note) ContentValidate() error {
	if note.Title == "" {
		return errors.New(fmt.Sprintf("タイトルを入力してください"))
	}
	if note.Content == "" {
		return errors.New(fmt.Sprintf("本文を入力してください"))
	}
	return nil
}
