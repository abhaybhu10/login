package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abhaybhu10/login/persistence"
)

type EnsureAuth struct {
	SessionDB persistence.Session
}

func (ea *EnsureAuth) Validate(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("session")
		session, err := ea.SessionDB.Get(sessionID)
		if sessionID == "" || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid session")
			return
		}

		ctx := context.WithValue(r.Context(), "userid", session.UserId)

		fmt.Printf("logined with session %v", session)
		next(w, r.WithContext(ctx))
	})
}
