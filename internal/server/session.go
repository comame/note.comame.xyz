package server

import "net/http"

type session struct {
	// TODO: まともな実装に置き換える
	ok     bool
	userID string
}

// TODO: まともな実装に置き換える
func startNewSession(w http.ResponseWriter, r *http.Request) *session {
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
	})

	s := &session{
		ok:     true,
		userID: k,
	}

	return s
}

func destroySession(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "s",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	})
}

func resumeSession(r *http.Request) *session {
	c, err := r.Cookie("s")
	if err != nil {
		return nil
	}

	return &session{
		ok:     true,
		userID: c.Value,
	}
}

func (s *session) getUserID() (string, bool) {
	if s == nil {
		return "", false
	}
	// TODO: 実装
	if !s.ok {
		return "", false
	}

	return "test", true
}

func (s *session) isLoggedIn() bool {
	_, ok := s.getUserID()
	return ok
}
