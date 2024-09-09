package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

		renderTemplate(s, w, "editor", "記事を作成", tmplEditor{
			SubmitTarget: "/post/create",
		})
	})

	http.HandleFunc("POST /post/create", func(w http.ResponseWriter, r *http.Request) {
		s := resumeSession(r, kvs)
		if !s.isLoggedIn() {
			if !s.isLoggedIn() {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var p post
		if err := json.Unmarshal(b, &p); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u, err := randomString(32)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p.URLKey = u

		now := time.Now().Format(time.DateTime)
		p.CreatedDatetime = now
		p.UpdatedDatetime = now

		con, err := getConnection()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := con.createPost(r.Context(), p); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		j, _ := json.Marshal(redirectResponse{Location: p.getURL()})
		w.Write(j)
	})

	http.HandleFunc("GET /posts/private/{url_key}", func(w http.ResponseWriter, r *http.Request) {
		s := resumeSession(r, kvs)
		if !s.isLoggedIn() {
			if !s.isLoggedIn() {
				w.WriteHeader(http.StatusUnauthorized)
				renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Unauthorized", Message: "ログインが必要."})
				return
			}
		}

		key := r.PathValue("url_key")

		p, err := getPost(r.Context(), key)
		if err != nil && errors.Is(err, errNotFound) {
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(s, w, "not-found", "Not Found", nil)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Internal Server Error", Message: "記事の取得に失敗"})
			return
		}

		if p.Visibility != postVisibilityPrivate {
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(s, w, "not-found", "Not Found", nil)
			return
		}

		renderTemplate(s, w, "post", p.Title, tmpPost{Post: *p})
	})

	http.HandleFunc("GET /manage/posts", func(w http.ResponseWriter, r *http.Request) {
		s := resumeSession(r, kvs)
		if !s.isLoggedIn() {
			if !s.isLoggedIn() {
				w.WriteHeader(http.StatusUnauthorized)
				renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Unauthorized", Message: "ログインが必要."})
				return
			}
		}

		con, err := getConnection()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(s, w, "Internal Server Error", "エラー", nil)
			return
		}

		p, err := con.getPosts(r.Context())
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(s, w, "Internal Server Error", "エラー", nil)
			return
		}

		renderTemplate(s, w, "manage-posts", "記事一覧", tmplManagePosts{Posts: p})
	})

	// === 誰でもアクセス可能 ===

	http.HandleFunc("GET /editor/demo", func(w http.ResponseWriter, r *http.Request) {
		s := resumeSession(r, kvs)

		f, err := os.Open("static/demo.md")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(s, w, "Internal Server Error", "エラー", nil)
			return
		}
		defer f.Close()

		c, err := io.ReadAll(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(s, w, "Internal Server Error", "エラー", nil)
			return
		}

		renderTemplate(s, w, "editor", "エディタ", tmplEditor{
			IsDemo: true,
			Post: post{
				Title:      "Demo",
				Text:       string(c),
				Visibility: postVisibilityPublic,
			},
		})
	})

	http.HandleFunc("GET /posts/limited/{url_key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("url_key")
		s := resumeSession(r, kvs)

		p, err := getPost(r.Context(), key)
		if err != nil && errors.Is(err, errNotFound) {
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(s, w, "not-found", "Not Found", nil)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Internal Server Error", Message: "記事の取得に失敗"})
			return
		}

		if p.Visibility != postVisibilityLimited {
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(s, w, "not-found", "Not Found", nil)
			return
		}

		renderTemplate(s, w, "post", p.Title, tmpPost{Post: *p})
	})

	http.HandleFunc("GET /posts/public/{url_key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("url_key")
		s := resumeSession(r, kvs)

		p, err := getPost(r.Context(), key)
		if err != nil && errors.Is(err, errNotFound) {
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(s, w, "not-found", "Not Found", nil)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			renderTemplate(nil, w, "error", "エラー", tmplError{Title: "Internal Server Error", Message: "記事の取得に失敗"})
			return
		}

		if p.Visibility != postVisibilityPublic {
			w.WriteHeader(http.StatusNotFound)
			renderTemplate(s, w, "not-found", "Not Found", nil)
			return
		}

		renderTemplate(s, w, "post", p.Title, tmpPost{Post: *p})
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
	ID              uint64         `json:"id"`
	URLKey          string         `json:"url_key"`
	CreatedDatetime string         `json:"-"`
	UpdatedDatetime string         `json:"-"`
	Title           string         `json:"title"`
	Text            string         `json:"text"`
	Visibility      postVisibility `json:"visibility"`
	HTML            string         `json:"-"`
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

func (p *post) getURL() string {
	switch p.Visibility {
	case postVisibilityPublic:
		return fmt.Sprintf("/posts/public/%s", p.URLKey)
	case postVisibilityLimited:
		return fmt.Sprintf("/posts/limited/%s", p.URLKey)
	case postVisibilityPrivate:
		return fmt.Sprintf("/posts/private/%s", p.URLKey)
	}

	panic("unknown visibility")
}

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

type redirectResponse struct {
	Location string `json:"location"`
}
