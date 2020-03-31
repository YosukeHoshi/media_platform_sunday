package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
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

func ConnectGorm() *gorm.DB {
	connectTemplate := "%s:%s@%s/%s?%s"
	connect := fmt.Sprintf(connectTemplate, DBUser, DBPass, DBProtocol, DBName, DBChar)
	fmt.Println(connect)
	Db, err := gorm.Open(Dialect, connect)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("データベースへの接続成功")
	}

	return Db
}

var Db *gorm.DB

func init() {
	Db = ConnectGorm()
	Db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&User{}, &Note{}, &Sns{}, &Session{})
	Db.LogMode(true)

	// defer Db.Close()
}
