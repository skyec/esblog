package main

import (
	"log"
	"net/http"
)

func main() {

	reloadable := NewTemplateReloader()
	mux := http.NewServeMux()
	mid := NewMiddleware(mux).
		AddPost(http.HandlerFunc(logRequests)).
		AddPre(reloadable)

	posts := newPostController(nil)
	reloadable.Add(posts)

	mr := newMethodRouter(mux, "/posts")
	mr.Add("GET", http.HandlerFunc(posts.index))
	mr.Add("POST", http.HandlerFunc(posts.add))

	mux.HandleFunc("/posts/add", posts.viewAdd)

	log.Print("Listening on http://localhost:7777 ...")
	http.ListenAndServe(":7777", mid)
}
