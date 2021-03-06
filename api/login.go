package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/abhaybhu10/login/model"
	"github.com/abhaybhu10/login/persistence"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPServer struct {
	sessionStore persistence.Session
	userStore    persistence.User
}

func (server *HTTPServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var loginData model.Login
	if err = json.Unmarshal(body, &loginData); err != nil {
		fmt.Printf("error while parsing json %v\n", err.Error())
	}

	user, err := server.userStore.Get(loginData.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, fmt.Sprintf("user %s does not exist\n", loginData.Username))
		fmt.Printf("login failed with error: %s\n", err.Error())
		return
	}
	if user.Password != loginData.Password {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, fmt.Sprintf("wrong credential for %s user", loginData.Username))
		fmt.Printf("invalid password for user %s \n", loginData.Username)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	session := model.Session{
		UserId: loginData.Username,
		ID:     uuid.NewString(),
	}

	if err := server.sessionStore.Save(session); err != nil {
		fmt.Printf(err.Error())
	}
	fmt.Fprintf(w, fmt.Sprintf("user logined %s, session %s", loginData.Username, session.ID))
	fmt.Printf("login successfull for user %s\n", loginData.Username)
	w.Header().Add("sessionID", session.ID)
}

func (server *HTTPServer) signupHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user model.User

	json.Unmarshal(body, &user)
	server.userStore.Save(user)
	fmt.Printf("User created %v\n", user)
	w.WriteHeader(http.StatusCreated)
}

func (server *HTTPServer) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "hello world")
	return
}

func (server *HTTPServer) Setup() {
	r := mux.NewRouter()
	r.HandleFunc("/login", server.loginHandler)
	r.HandleFunc("/signup", server.signupHandler).Methods(http.MethodPost)
	r.HandleFunc("/index", server.withLogin(server.indexHandler)).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func NewHttpServer() *HTTPServer {
	return &HTTPServer{
		userStore:    persistence.GetUserStore(),
		sessionStore: persistence.GetSessionStore(),
	}
}

func (server *HTTPServer) withLogin(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("session")
		session, err := server.sessionStore.Get(sessionID)
		if sessionID == "" || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "invalid session")
			return
		}

		ctx := context.WithValue(r.Context(), "userid", session.UserId)

		fmt.Printf("logined with session %v\n", session)
		next(w, r.WithContext(ctx))
	})
}
