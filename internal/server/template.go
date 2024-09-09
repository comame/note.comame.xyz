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

type tmpPost struct {
	Post post
}

type tmplEditor struct {
	IsDemo       bool
	SubmitTarget string
	Post         post
}

type tmplManagePosts struct {
	Posts []post
}

type tmplApp struct {
	Title      string
	Body       string
	IsLoggedIn bool
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
	})
	template.Must(t.ParseGlob("templates/*.html"))
	return t
}

func renderTemplate(s *session, w http.ResponseWriter, name, title string, param any) {
	t := setupTemplate()

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
