package server

import (
	"net/http"
)

type session struct {
	ok     bool
	userID string
}

func startNewSession(w http.ResponseWriter, userID string, kvs *kvs) *session {
	if userID != "comame" {
		return nil
	}

	k, err := randomString(64)
	if err != nil {
		// 乱数を生成できないのは何かがおかしいので panic してしまう
		panic(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "s",
		Value:    k,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})

	s := &session{
		ok:     true,
		userID: userID,
	}

	kvs.SetSession(k, userID)

	return s
}

func destroySession(w http.ResponseWriter, r *http.Request, kvs *kvs) {
	if c, err := r.Cookie("s"); err == nil /* err IS nil */ {
		kvs.DeleteSession(c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "s",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func resumeSession(r *http.Request, kvs *kvs) *session {
	c, err := r.Cookie("s")
	if err != nil {
		return nil
	}

	u, ok := kvs.GetSession(c.Value)
	if !ok {
		return nil
	}

	return &session{
		ok:     true,
		userID: u,
	}
}

func (s *session) getUserID() (string, bool) {
	if s == nil {
		return "", false
	}
	if !s.ok {
		return "", false
	}

	return s.userID, true
}

func (s *session) isLoggedIn() bool {
	_, ok := s.getUserID()
	return ok
}
