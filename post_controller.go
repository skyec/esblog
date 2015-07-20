package main

import (
	"html/template"
	"net/http"
)

var postHistory []interface{}

func init() {
	postHistory = []interface{}{
		postAdded{postEvent{newRandomPostId()}, postContent{"test@test.com", "First Post", "Your blog is all setup. Start adding posts!"}},
	}
}

type postController struct {
	posts     *posts
	templates *template.Template
}

func newPostController(t *template.Template) *postController {
	return &postController{
		posts:     postsFromHistory(postHistory),
		templates: helperMustLoadTemplates(),
	}
}

func (ctlr *postController) index(w http.ResponseWriter, req *http.Request) {
	err := ctlr.templates.ExecuteTemplate(w, "index", struct{ Posts []interface{} }{helperMapPosts(ctlr.posts.mostRecent(20))})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ctlr *postController) viewAdd(w http.ResponseWriter, req *http.Request) {
	err := ctlr.templates.ExecuteTemplate(w, "viewAdd", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ctrl *postController) add(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form:"+err.Error(), http.StatusBadRequest)
	}
	ctrl.posts.apply(
		postAdded{
			postEvent{newRandomPostId()},
			postContent{req.FormValue("author"), req.FormValue("title"), req.FormValue("content")},
		})
	http.Redirect(w, req, "/posts", http.StatusSeeOther)
}

func (ctlr *postController) reloadTemplates() {
	ctlr.templates = helperMustLoadTemplates()
}

func helperMustLoadTemplates() *template.Template {
	t, err := template.ParseGlob(config.resourceDir + "/views/posts/*.html")
	if err != nil {
		panic(err.Error())
	}
	return t
}

func helperMapPosts(posts []*post) []interface{} {
	type vm struct {
		Author, Title, Content string
	}
	newps := []interface{}{}
	for _, p := range posts {
		newps = append(newps, vm{p.content.author, p.content.title, p.content.content})
	}
	return newps
}
