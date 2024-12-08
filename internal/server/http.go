package server

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
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
	renderTemplate(s, w, templateNameError, "エラー", templateError{Title: "Bad Request", Message: ""})
}

func renderInternalServerError(s *session, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	renderTemplate(s, w, templateNameError, "エラー", templateError{Title: "Internal Server Error", Message: ""})
}

func renderNotFound(s *session, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(s, w, templateNameNotFound, "Not Found", nil)
}

func compressStaticHandler(d http.Dir) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Encoding", "gzip")
		ct := mime.TypeByExtension(path.Ext(r.URL.Path))
		w.Header().Add("Content-Type", ct)

		f, err := d.Open(r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer f.Close()

		gw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer gw.Close()

		if _, err := io.Copy(gw, f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})
}

func isCompressRequest(r *http.Request) bool {
	e := path.Ext(r.URL.Path)

	return e == ".wasm"
}
