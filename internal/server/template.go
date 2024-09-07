package server

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
)

type tmplError struct {
	Title   string
	Message string
}

type tmplEditor struct {
	IsDemo bool
}

type tmplApp struct {
	Title      string
	Body       string
	IsLoggedIn bool
}

func renderTemplate(s *session, w http.ResponseWriter, name, title string, param any) {
	// TODO: フロントエンドが書き終わったらグローバル変数に移して、リクエストごとに回さなくてよくする
	var t = template.Must(template.ParseGlob("templates/*.html"))

	var b bytes.Buffer
	if err := t.ExecuteTemplate(&b, name+".html", param); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, "app.html", tmplApp{
		Title:      title,
		Body:       b.String(),
		IsLoggedIn: s.isLoggedIn(),
	}); err != nil {
		panic(err)
	}
}
