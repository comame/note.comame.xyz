package server

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func readJSONFromBody(r *http.Request, v any) error {
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("content type is not application/json")
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

func setCommonHeaders(w http.ResponseWriter) {
	w.Header().Add("Referrer-Policy", "same-origin")
}

func validateRequest(requireLogin bool, r *http.Request, kvs *kvs) (s *session, ok bool) {
	// GET / OPTIONS 以外は CSRF チェックする
	if !(r.Method == http.MethodGet || r.Method == http.MethodOptions) {
		a := r.Header.Get("Origin")
		b := os.Getenv("ORIGIN")

		if b == "" || a == "" {
			return nil, false
		}

		if a != b {
			return nil, false
		}

		sfs := r.Header.Get("Sec-Fetch-Site")
		if sfs != "same-origin" {
			return nil, false
		}
	}

	// セッションのチェック
	s = resumeSession(r, kvs)
	if requireLogin {
		u, ok := s.getUserID()
		if !ok {
			return nil, false
		}
		if u != "comame" {
			return nil, false
		}
	}

	return s, true
}

func isPageRequest(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/html")
}

func renderBadRequest(s *session, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	renderTemplate(s, w, pageNotFound, "エラー", nil)
}

func renderInternalServerError(s *session, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	renderTemplate(s, w, pageNotFound, "エラー", nil)
}

func renderNotFound(s *session, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(s, w, pageNotFound, "Not Found", nil)
}

func staticHandler(d http.Dir) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := d.Open(r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer f.Close()

		// キャッシュ
		w.Header().Add("Cache-Control", "no-cache")
		stat, err := f.Stat()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cacheKey := stat.ModTime().Format(time.RFC3339)
		hash := base64.RawStdEncoding.EncodeToString(md5.New().Sum([]byte(cacheKey)))
		reqEtag := r.Header.Get("If-None-Match")
		if reqEtag == hash {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Add("ETag", hash)

		// gzip encoding
		w.Header().Add("Content-Encoding", "gzip")
		contentType := mime.TypeByExtension(path.Ext(r.URL.Path))
		w.Header().Add("Content-Type", contentType)

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		if _, err := io.Copy(gz, f); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
