package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/comame/note.comame.xyz/internal/oidc"
)

func handleLogin(kvs *kvs, oidcClientID, oidcRedirectURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func handleLogout(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			renderBadRequest(nil, w)
			return
		}
		destroySession(w, r, kvs)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func handleOIDCCallback(kvs *kvs, oidcClientID, oidcClientSecret, oidcRedirectURI string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		p, err := oidc.CallbackCode(c.Value, r.URL.Query(), oidcClientID, oidcClientSecret, oidcRedirectURI, kvs, oidcClientID)
		if err != nil {
			log.Println(err)
			renderInternalServerError(nil, w)
			return
		}
		startNewSession(w, p.Sub, kvs)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func handleNewPostPage(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(true, r, kvs)
		if !ok {
			renderBadRequest(s, w)
			return
		}
		renderTemplate(s, w, pageNewPost, "記事を作成", struct{}{})
	}
}

func handleCreatePost(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func handlePrivatePostPage(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(true, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}
		postPage(w, r, s)
	}
}

func handleEditPostPage(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		log.Println(p)
		renderTemplate(s, w, pageEditPost, "記事を編集", editPostPageData{Post: *p})
	}
}

func handleEditPost(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if p.ID != id {
			log.Println(p.ID, id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := updatePost(r.Context(), p); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		p2, err := getPostByID(r.Context(), id, p.Permission)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, _ := json.Marshal(redirectResponse{Location: p2.getURL()})
		w.Write(j)
	}
}

func handleDeletePost(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func handleAllPostsPage(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postListPage(w, r, kvs)
	}
}

func handleDemoEditorPage(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(false, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}
		renderTemplate(s, w, pageDemoEditor, "エディタ", nil)
	}
}

func handleUnlistedPostPage(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(false, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}
		postPage(w, r, s)
	}
}

func handlePublicPostPage(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		s, ok := validateRequest(false, r, kvs)
		if !ok {
			renderBadRequest(nil, w)
			return
		}
		postPage(w, r, s)
	}
}

func handleStatic(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))).ServeHTTP(w, r)
	}
}

func handleOutDist(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h := staticHandler(http.Dir("out/dist"))
		http.StripPrefix("/out/dist/", h).ServeHTTP(w, r)
	}
}

func handleAssets(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setCommonHeaders(w)
		if _, ok := validateRequest(false, r, kvs); !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		h := staticHandler(http.Dir("out/front/assets"))
		http.StripPrefix("/assets", h).ServeHTTP(w, r)
	}
}

func handleRootPage(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postListPage(w, r, kvs)
	}
}

func handleNotFound(kvs *kvs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

type redirectResponse struct {
	Location string `json:"location"`
}

// 記事ページの共通ハンドラ
func postPage(w http.ResponseWriter, r *http.Request, s *session) {
	var v permission
	switch strings.Split(r.URL.Path, "/")[2] {
	case "private":
		v = permissionPrivate
	case "unlisted":
		v = permissionURL
	case "public":
		v = permissionPublic
	}

	key := r.PathValue("url_key")

	p, err := getPostByURLKey(r.Context(), key, v)
	if err != nil && errors.Is(err, errNotFound) {
		renderNotFound(s, w)
		return
	}
	if err != nil {
		renderInternalServerError(s, w)
		return
	}

	if p.Permission != v {
		renderNotFound(s, w)
		return
	}

	renderTemplate(s, w, pagePost, p.Title+" | note.comame.xyz", postPageData{
		Post: *p,
	})
}

func postListPage(w http.ResponseWriter, r *http.Request, kvs *kvs) {
	setCommonHeaders(w)
	s, ok := validateRequest(false, r, kvs)
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

	var p []post

	if s.isLoggedIn() {
		p, err = con.getAllPostsForAdmin(r.Context())
		if err != nil {
			log.Println(err)
			renderInternalServerError(s, w)
			return
		}
	} else {
		p, err = con.getAllPostsForAnonymous(r.Context())
		if err != nil {
			log.Println(err)
			renderInternalServerError(s, w)
			return
		}
	}

	renderTemplate(s, w, pageAllPosts, "記事一覧", allPostsPageData{
		Posts: p,
	})
}
