package main

import (
	"github.com/YosukeHoshi/media_platform_sunday/internal/account"
	"github.com/YosukeHoshi/media_platform_sunday/internal/post"
	"net/http"
)

func main() {
	http.HandleFunc("/signup", account.Signup)
	http.HandleFunc("/signin", account.Signin)
	http.HandleFunc("/signout", account.Signout)
	http.HandleFunc("/post", account.SignInRequired(post.CreatePost))
	http.HandleFunc("/getpost", account.SignInRequired(post.GetMyNotesHandler))
	http.HandleFunc("/getall", account.SignInRequired(post.GetAllNotesHandler))
	http.HandleFunc("/delete", account.SignInRequired(post.DeleteNoteHandler))
	// 確認用
	http.HandleFunc("/getcookie", account.GetCookie)
	http.HandleFunc("/update", account.SignInRequired(post.UpdataNoteHandler))
	http.ListenAndServe(":8080", nil)
}
