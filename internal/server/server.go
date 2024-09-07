package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/comame/note.comame.xyz/internal/md"
	oidc "github.com/comame/note.comame.xyz/internal/odic"
	_ "github.com/go-sql-driver/mysql"
)

func Start() {
	// === ログイン ===

	oidcIssuer := os.Getenv("OIDC_ISSUER")
	oidcClientID := os.Getenv("OIDC_CLIENT_ID")
	oidcClientSecret := os.Getenv("OIDC_CLIENT_SECRET")
	oidcRedirectURI := fmt.Sprintf("%s/login/oidc-callback", os.Getenv("ORIGIN"))
	oidcAud := "note.comame.xyz"

	oidc.InitializeDiscovery(oidcIssuer)
	kvs := initKVS()

	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		u, s, err := oidc.GenerateAuthenticationRequestUrl(oidcClientID, oidcRedirectURI, kvs)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Internal Server Error", Message: "ログインに失敗."})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "state",
			Value:    s,
			MaxAge:   600,
			Secure:   true,
			HttpOnly: true,
		})

		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc("GET /logout", func(w http.ResponseWriter, r *http.Request) {
		destroySession(w, r, kvs)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("GET /login/oidc-callback", func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("state")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Bad Request", Message: "ログインに失敗."})
			return
		}
		p, err := oidc.CallbackCode(c.Value, r.URL.Query(), oidcClientID, oidcClientSecret, oidcRedirectURI, kvs, oidcAud)
		if err != nil {
			renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Bad Request", Message: "ログインに失敗."})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		startNewSession(w, p.Sub, kvs)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// === ログイン専用 ===

	http.HandleFunc("GET /post/new", func(w http.ResponseWriter, r *http.Request) {
		s := resumeSession(r, kvs)
		if !s.isLoggedIn() {
			w.WriteHeader(http.StatusUnauthorized)
			renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Unauthorized", Message: "ログインが必要."})
			return
		}

		renderTemplate(s, w, "editor", "記事を作成", tmplEditor{})
	})

	// === 誰でもアクセス可能 ===

	http.HandleFunc("GET /editor/demo", func(w http.ResponseWriter, r *http.Request) {
		s := resumeSession(r, kvs)
		renderTemplate(s, w, "editor", "エディタ", tmplEditor{
			IsDemo: true,
		})
	})

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

	http.Handle("GET /static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static")),
	))
	http.Handle("GET /out/dist/", http.StripPrefix("/out/dist/",
		http.FileServer(http.Dir("out/dist")),
	))

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			s := resumeSession(r, kvs)
			renderTemplate(s, w, "not-found", "Not Found", nil)
			return
		}

		w.Write([]byte("Not found"))
	})

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
