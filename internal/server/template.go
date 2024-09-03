package server

import (
	"bytes"
	"log"
	"net/http"
	"text/template"
)

type tmplEditor struct {
	IsDemo bool
}

type tmplApp struct {
	Body       string
	IsLoggedIn bool
}

func renderTemplate(w http.ResponseWriter, name string, param any) {
	// TODO: フロントエンドが書き終わったらグローバル変数に移して、リクエストごとに回さなくてよくする
	var t = template.Must(template.ParseGlob("templates/*.html"))

	var b bytes.Buffer
	if err := t.ExecuteTemplate(&b, name+".html", param); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, "app.html", tmplApp{
		Body: b.String(),
	}); err != nil {
		panic(err)
	}
}
