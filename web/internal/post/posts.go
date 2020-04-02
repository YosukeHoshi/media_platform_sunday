package post

import (
	"encoding/json"

	// "errors"

	"log"
	"net/http"

	"github.com/YosukeHoshi/media_platform_sunday/internal/database"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// CreatePost is
func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed. Only POST Is Available."))
		log.Println("Method Not Allowed. Only POST Is Available.")
		return
	}

	var note database.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		log.Println(err)
		return
	}

	if err = note.ContentValidate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		log.Println(err)
		return
	}

	// cookieからsessionIDの取得
	cookie, err := r.Cookie(("session_id"))
	if err != nil {
		// check later
		log.Println("Cookie: ", err)
	}

	// sessionテーブルからcookieにあるsessionIDとuuidが一致するものを取得
	session, err := database.GetSession(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		return
	}
	note.UserName = session.UserName

	database.Db.Create(&note)
	w.WriteHeader(http.StatusOK)
	log.Println("posted")
	return
}

// GetMyNotesHandker is
func GetMyNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed. Only POST Is Available."))
		log.Println("Method Not Allowed. Only POST Is Available.")
		return
	}

	cookie, err := r.Cookie(("session_id"))
	if err != nil {
		// check later
		log.Println("Cookie: ", err)
	}

	session, err := database.GetSession(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		return
	}

	notes, err := database.GetMyNotes(session.UserName)
	if err != nil {
		log.Println(err.Error())
	}

	bytes, err := json.Marshal(&notes)
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	log.Println("get your notes")
	return
}

// GetAllNotesHandler returns all notes.(フォロー機能を付けたらフォローしているユーザーに絞る)
func GetAllNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed. Only POST Is Available."))
		log.Println("Method Not Allowed. Only POST Is Available.")
		return
	}

	_, err := r.Cookie(("session_id"))
	if err != nil {
		// check later
		log.Println("Cookie: ", err)
	}

	notes, err := database.GetAllNotes()
	if err != nil {
		log.Println(err.Error())
	}

	bytes, err := json.Marshal(&notes)
	if err != nil {
		log.Println(err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	log.Println("get  notes")
	return
}

// UpdateNote is
func UpdataNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed. Only POST Is Available."))
		log.Println("Method Not Allowed. Only POST Is Available.")
		return
	}

	var updateNote database.Note
	err := json.NewDecoder(r.Body).Decode(&updateNote)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if updateNote.ID == 0 {
		log.Println("id is requested")
		return
	}

	note, err := database.GetNote(updateNote.ID)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = note.UpdateNote(updateNote)
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("update successful"))
	log.Println("update successful")
	return
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed. Only POST Is Available."))
		log.Println("Method Not Allowed. Only POST Is Available.")
		return
	}

	_, err := r.Cookie(("session_id"))
	if err != nil {
		// check later
		log.Println("Cookie: ", err)
	}

	var note database.Note
	err = json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if note.ID == 0 {
		log.Println("id is requested")
		return
	}

	deleteNote, err := database.GetNote(note.ID)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = deleteNote.Delete()
	if err != nil {
		log.Println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete successful"))
	log.Println("delete successful")
	return
}
