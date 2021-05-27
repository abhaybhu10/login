package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abhaybhu10/login/persistence"
)

type EnsureAuth struct {
	handler   http.Handler
	sessionDB persistence.Session
}

func (ea *EnsureAuth) validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("session")
		session, err := ea.sessionDB.Get(sessionID)
		if sessionID == "" || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid session")
			return
		}

		ctx := context.WithValue(r.Context(), "userid", session.UserId)

		fmt.Printf("logined with session %v", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
