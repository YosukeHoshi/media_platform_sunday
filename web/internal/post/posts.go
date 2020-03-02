package post

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/YosukeHoshi/media_platform_sunday/internal/database"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	err := HandleOnlyPost(w, r)
	if err != nil {
		// err = xerrors.Errorf("Error in request method: %v", err)
		log.Println(err.Error())
		return
	}

	var note database.Note
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		log.Println(err)
		return
	}

	if err = note.ContentValidate(); err != nil {
		log.Println(err)
		// w.WriteHeader()
		return
	}

	database.Db.Create(&note)
}

func HandleOnlyPost(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return errors.New(fmt.Sprintf("Method Not Allowed. Only POST Is Available."))
	}
	return nil
}

func PrintLog(w http.ResponseWriter, st string) {
	w.Write([]byte(st))
	log.Println(st)
}
