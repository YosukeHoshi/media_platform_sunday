package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const DRIVER = "mysql"
const DSN = "root:root-mps@tcp(mysql-container:3306)/golang-mps-database" //envファイル

func main() {
	db, err := sql.Open(DRIVER, DSN)
	if err != nil {
		fmt.Println("Openエラー")
	} else {
		fmt.Println("OpenOK！")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("接続失敗！")
	} else {
		fmt.Println("接続OK！")
	}

	db.Close()
}