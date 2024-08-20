package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/comame/note.comame.xyz/internal/md"
	_ "github.com/go-sql-driver/mysql"
)

func Start() {
	http.HandleFunc("GET /posts/limited/{url_key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("url_key")

		p, err := getPost(r.Context(), key)
		if err != nil && errors.Is(err, errNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		if p.Visibility != postVisibilityLimited {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		out := fmt.Sprintf(`
			<h1>%s</h1>
			<div class='post'>%s</div>
			<style></style>
		`, p.Title, p.HTML)

		w.Write([]byte(out))
	})

	http.HandleFunc("GET /", http.FileServerFS(os.DirFS("./out/dist")).ServeHTTP)

	http.ListenAndServe(":8080", nil)
}

type post struct {
	ID              uint64         `json:"-"`
	URLKey          string         `json:"url_key"`
	CreatedDatetime string         `json:"created_datetime"`
	UpdatedDatetime string         `json:"updated_datetime"`
	Title           string         `json:"title"`
	Text            string         `json:"text"`
	Visibility      postVisibility `json:"visibility"`
	HTML            string         `json:"html"`
}

type postVisibility int

const (
	// 非公開
	postVisibilityPrivate postVisibility = 0
	// 限定公開
	postVisibilityLimited = 1
	// 全体公開
	postVisibilityPublic = 2
)

func getPost(ctx context.Context, urlKey string) (*post, error) {
	c, err := getConnection()
	if err != nil {
		return nil, err
	}

	p, err := c.findPost(ctx, urlKey)
	if errors.Is(err, errNotFound) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}

	pv, err := c.findVisibility(ctx, *p)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to query postVisibility"))
	}

	pv.HTML = md.ToHTML(pv.Text)

	return pv, nil
}
