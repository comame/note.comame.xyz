package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

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

	// ログインを開始する
	http.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			renderBadRequest(nil, w)
			return
		}

		u, s, err := oidc.GenerateAuthenticationRequestUrl(oidcClientID, oidcRedirectURI, kvs)
		if err != nil {
			renderInternalServerError(nil, w)
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
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			renderBadRequest(nil, w)
			return
		}

		destroySession(w, r, kvs)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("GET /login/oidc-callback", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			renderBadRequest(nil, w)
			return
		}

		c, err := r.Cookie("state")
		if err != nil {
			renderInternalServerError(nil, w)
			return
		}
		p, err := oidc.CallbackCode(c.Value, r.URL.Query(), oidcClientID, oidcClientSecret, oidcRedirectURI, kvs, oidcAud)
		if err != nil {
			renderInternalServerError(nil, w)
			return
		}

		startNewSession(w, p.Sub, kvs)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	// === ログイン専用 ===

	http.HandleFunc("GET /post/new", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(true, r, kvs)
		if !ok {
			renderBadRequest(s, w)
			return
		}

		renderTemplate(s, w, templateNameEditor, "記事を作成", templateEditor{
			SubmitTarget: "/post/create",
		})
	})

	http.HandleFunc("POST /post/create", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(true, r, kvs)
		if !ok {
			renderBadRequest(s, w)
			return
		}

		var p1 post
		if err := readJSONFromBody(r, &p1); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p2, err := createPost(r.Context(), p1)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		j, _ := json.Marshal(redirectResponse{Location: p2.getURL()})
		w.Write(j)
	})

	http.HandleFunc("GET /posts/private/{url_key}", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(true, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}

		key := r.PathValue("url_key")

		p, err := getPost(r.Context(), key, postVisibilityPrivate)
		if err != nil && errors.Is(err, errNotFound) {
			renderNotFound(s, w)
			return
		}
		if err != nil {
			renderInternalServerError(s, w)
			return
		}

		renderTemplate(s, w, "post", p.Title+" | note.comame.xyz", templatePost{Post: *p, EditLink: fmt.Sprintf("/edit/post/%d", p.ID)})
	})

	http.HandleFunc("GET /manage/posts", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(true, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}

		con, err := GetConnection()
		if err != nil {
			log.Println(err)
			renderInternalServerError(s, w)
			return
		}

		p, err := con.getPosts(r.Context())
		if err != nil {
			log.Println(err)
			renderInternalServerError(s, w)
			return
		}

		renderTemplate(s, w, templateNameManagePosts, "記事一覧", templateManagePosts{Posts: p})
	})

	http.HandleFunc("GET /edit/post/{post_id}", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(true, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}

		idStr := r.PathValue("post_id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			renderBadRequest(s, w)
			return
		}

		con, err := GetConnection()
		if err != nil {
			log.Println(err)
			renderInternalServerError(s, w)
			return
		}

		p, err := con.findPostByID(r.Context(), id)
		if err != nil && errors.Is(err, errNotFound) {
			renderNotFound(s, w)
			return
		}

		renderTemplate(s, w, templateNameEditor, "記事を作成", templateEditor{
			SubmitTarget: "/edit/post/" + idStr,
			Post:         *p,
		})
	})

	http.HandleFunc("POST /edit/post/{post_id}", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(true, r, kvs); !ok {
			renderBadRequest(nil, w)
			return
		}

		var p post
		if err := readJSONFromBody(r, &p); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		idStr := r.PathValue("post_id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if p.ID != id || p.URLKey == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		p2, err := updatePost(r.Context(), p)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		j, _ := json.Marshal(redirectResponse{Location: p2.getURL()})
		w.Write(j)
	})

	http.HandleFunc("POST /delete/post/{post_id}", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(true, r, kvs); !ok {
			renderBadRequest(nil, w)
			return
		}

		idStr := r.PathValue("post_id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := deletePost(r.Context(), id); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	// === 誰でもアクセス可能 ===

	http.HandleFunc("GET /editor/demo", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(false, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}

		f, err := os.Open("static/demo.md")
		if err != nil {
			renderInternalServerError(s, w)
			return
		}
		defer f.Close()

		c, err := io.ReadAll(f)
		if err != nil {
			renderInternalServerError(s, w)
			return
		}

		renderTemplate(s, w, templateNameEditor, "エディタ", templateEditor{
			IsDemo: true,
			Post: post{
				Title:      "Demo",
				Text:       string(c),
				Visibility: postVisibilityPublic,
			},
		})
	})

	http.HandleFunc("GET /posts/unlisted/{url_key}", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(false, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}

		key := r.PathValue("url_key")

		p, err := getPost(r.Context(), key, postVisibilityUnlisted)
		if err != nil && errors.Is(err, errNotFound) {
			renderNotFound(s, w)
			return
		}
		if err != nil {
			renderInternalServerError(s, w)
			return
		}

		if p.Visibility != postVisibilityUnlisted {
			renderNotFound(s, w)
			return
		}

		renderTemplate(s, w, "post", p.Title+" | note.comame.xyz", templatePost{Post: *p, EditLink: fmt.Sprintf("/edit/post/%d", p.ID)})
	})

	http.HandleFunc("GET /posts/public/{url_key}", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(false, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}

		key := r.PathValue("url_key")

		p, err := getPost(r.Context(), key, postVisibilityPublic)
		if err != nil && errors.Is(err, errNotFound) {
			renderNotFound(s, w)
			return
		}
		if err != nil {
			renderInternalServerError(s, w)
			return
		}

		if p.Visibility != postVisibilityPublic {
			renderNotFound(s, w)
			return
		}

		renderTemplate(s, w, "post", p.Title+" | note.comame.xyz", templatePost{Post: *p, EditLink: fmt.Sprintf("/edit/post/%d", p.ID)})
	})

	http.HandleFunc("GET /static/", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	})

	http.HandleFunc("GET /out/dist/", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		http.StripPrefix("/out/dist/", http.FileServer(http.Dir("out/dist"))).ServeHTTP(w, r)
	})

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(false, r, kvs)
		if !ok {
			renderBadRequest(s, w)
			return
		}

		renderTemplate(s, w, templateNameTop, "note.comame.xyz", templateTop{})
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(false, r, kvs)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isPageRequest(r) {
			renderNotFound(s, w)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	})

	log.Println("start http://0.0.0.0:8080")
	http.ListenAndServe(":8080", http.DefaultServeMux)
}

type redirectResponse struct {
	Location string `json:"location"`
}
