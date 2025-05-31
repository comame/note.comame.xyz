package server

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"text/template"
)

type breadcrumb struct {
	Label    string
	Location string
}

type page string

const (
	pageTop        page = "TopPage"
	pageAllPosts   page = "AllPostsPage"
	pageNewPost    page = "NewPostPage"
	pageEditPost   page = "EditPostPage"
	pagePost       page = "PostPage"
	pageDemoEditor page = "DemoEditorPage"
	pageNotFound   page = "NotFoundPage"
)

type pageProps struct {
	IsLoggedIn    bool
	Breadcrumbs   []breadcrumb
	Title         string
	Page          page
	PageData      any
	PagePropsJSON string `json:"-,omitempty"`
}

type allPostsPageData struct {
	Posts []post `json:"posts"`
}

type postPageData struct {
	Post post `json:"post"`
}

type editPostPageData struct {
	Post post `json:"post"`
}

func setupTemplate() *template.Template {
	// TODO: フロントエンドが書き終わったらグローバル変数に移して、リクエストごとに回さなくてよくする
	t := template.New("_")
	template.Must(t.ParseGlob("templates/*.html"))
	return t
}

func renderTemplate(s *session, w http.ResponseWriter, name page, title string, pageData any) {
	t := setupTemplate()

	props := pageProps{
		IsLoggedIn: s.isLoggedIn(),
		Breadcrumbs: []breadcrumb{
			{Label: "Top", Location: "/"},
		},
		Title:    html.EscapeString(title),
		Page:     name,
		PageData: pageData,
	}

	pagePropsJSON, err := json.Marshal(props)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	props.PagePropsJSON = html.EscapeString(string(pagePropsJSON))

	if err := t.ExecuteTemplate(w, "app_index.html", props); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
