package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

type templateName string

const (
	templateNameEditor      templateName = "editor"
	templateNameError       templateName = "error"
	templateNameManagePosts templateName = "manage-posts"
	templateNameNotFound    templateName = "not-found"
	templateNamePost        templateName = "post"
	templateNameTop         templateName = "top"
)

type templateError struct {
	Title   string
	Message string
}

type templatePost struct {
	Post       post
	EditLink   string
	IsLoggedIn bool
}

type templateEditor struct {
	IsDemo       bool
	SubmitTarget string
	Post         post
}

type templateManagePosts struct {
	Posts []post
}

type templateTop struct{}

type templateApp struct {
	Title         string
	Body          string
	IsLoggedIn    bool
	OgDescription string
}

func setupTemplate() *template.Template {
	// TODO: フロントエンドが書き終わったらグローバル変数に移して、リクエストごとに回さなくてよくする
	t := template.New("_")
	t.Funcs(map[string]any{
		"toYMDString": func(datetime string) string {
			l := len("2024-09-01")
			if len(datetime) < l {
				return datetime
			}
			return datetime[:l]
		},
		"postURL": func(p post) string {
			return p.getURL()
		},
		"editURL": func(p post) string {
			return p.editURL()
		},
		"visibilityLabel": func(p post) string {
			return p.visibilityLabel()
		},
	})
	template.Must(t.ParseGlob("templates/*.html"))
	return t
}

func renderTemplate(s *session, w http.ResponseWriter, name templateName, title string, param any) {
	t := setupTemplate()

	switch p := param.(type) {
	case templatePost:
		p.IsLoggedIn = s.isLoggedIn()
		param = p
	}

	var b bytes.Buffer
	if err := t.ExecuteTemplate(&b, string(name)+".html", param); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ogDescription := "note.comame.xyz"
	if name == templateNamePost {
		p := param.(templatePost)
		ogDescription = fmt.Sprintf("%d字", len(p.Post.Text))
	}

	if err := t.ExecuteTemplate(w, "app.html", templateApp{
		Title:         title,
		Body:          b.String(),
		IsLoggedIn:    s.isLoggedIn(),
		OgDescription: ogDescription,
	}); err != nil {
		panic(err)
	}
}
