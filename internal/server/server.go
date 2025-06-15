package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/comame/note.comame.xyz/internal/oidc"

	_ "github.com/go-sql-driver/mysql"
)

func Start() {
	// === ログイン ===

	oidcIssuer := os.Getenv("OIDC_ISSUER")
	oidcClientID := os.Getenv("OIDC_CLIENT_ID")
	oidcClientSecret := os.Getenv("OIDC_CLIENT_SECRET")
	oidcRedirectURI := fmt.Sprintf("%s/login/oidc-callback", os.Getenv("ORIGIN"))

	oidc.InitializeDiscovery(oidcIssuer)
	kvs := initKVS()

	http.HandleFunc("GET /login", handleLogin(kvs, oidcClientID, oidcRedirectURI))
	http.HandleFunc("GET /logout", handleLogout(kvs))
	http.HandleFunc("GET /login/oidc-callback", handleOIDCCallback(kvs, oidcClientID, oidcClientSecret, oidcRedirectURI))

	// === ログイン専用 ===
	http.HandleFunc("GET /new", handleNewPostPage(kvs))
	http.HandleFunc("POST /post/create", handleCreatePost(kvs))
	http.HandleFunc("GET /edit/post/{post_id}", handleEditPostPage(kvs))
	http.HandleFunc("POST /edit/post/{post_id}", handleEditPost(kvs))
	http.HandleFunc("POST /delete/post/{post_id}", handleDeletePost(kvs))

	// === 誰でもアクセス可能 ===
	http.HandleFunc("GET /all", handleAllPostsPage(kvs))
	http.HandleFunc("GET /posts/{url_key}", handlePostPage(kvs))
	http.HandleFunc("GET /editor/demo", handleDemoEditorPage(kvs))
	http.HandleFunc("GET /static/", handleStatic(kvs))
	http.HandleFunc("GET /out/dist/", handleOutDist(kvs))
	http.HandleFunc("GET /assets/", handleAssets(kvs))
	http.HandleFunc("GET /{$}", handleRootPage(kvs))
	http.HandleFunc("GET /", handleNotFound(kvs))

	log.Println("start http://0.0.0.0:8080")
	if err := http.ListenAndServe(":8080", http.DefaultServeMux); err != nil {
		panic(err)
	}
}
