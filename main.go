package main

import (
	"fmt"
	"log"

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
)

func connectGorm() *gorm.DB {
	connectTemplate := "%s:%s@%s/%s"
	connect := fmt.Sprintf(connectTemplate, DBUser, DBPass, DBProtocol, DBName)
	db, err := gorm.Open(Dialect, connect)

	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("成功")
	}

	return db
}

func main() {
	db := connectGorm()

	//db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&User{})

	defer db.Close()
}

/*
type User struct {
	gorm.Model
	Name string `gorm:"size:255"`
	Age  int
	Sex  string `gorm:"size:255"`
}

func (u User) String() string {
	return fmt.Sprintf("%s(%d)", u.Name, u.Age)
}
*/
