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
	http.HandleFunc("/post", post.CreatePost)
	http.HandleFunc("/getpost", post.GetMyNotesHandler)
	// 確認用
	http.HandleFunc("/getcookie", account.GetCookie)
	http.HandleFunc("/update", post.UpdataNoteHandler)
	http.ListenAndServe(":8080", nil)
}
