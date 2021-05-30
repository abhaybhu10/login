package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/abhaybhu10/login/model"
	"github.com/abhaybhu10/login/persistence/cache"
	"github.com/abhaybhu10/login/persistence/db/sql"
)

type Session interface {
	Save(model.Session) error
	Get(Id string) (*model.Session, error)
}

type User interface {
	Save(user model.User) error
	Get(userId string) (*model.User, error)
}

type UserStore struct {
	cache *cache.UserRedis
	db    *sql.Mysql
}

type SessionStore struct {
	cache *cache.SessionRedis
	db    *sql.Mysql
}

func GetSessionStore() Session {
	return &SessionStore{
		cache: cache.GetSessionRedis(),
		db:    sql.GetMySql(),
	}
}

func GetUserStore() User {
	return &UserStore{
		cache: cache.GetUserRedis(),
		db:    sql.GetMySql(),
	}
}

func (u *UserStore) Save(user model.User) error {
	ctx := context.Background()
	if err := u.cache.Save(ctx, user); err != nil {
		fmt.Printf("cache put failed")
	}

	err := u.db.SaveUser(user)
	if err != nil {
		fmt.Printf("Error while saving to database")
	}
	return err
}

func (s *SessionStore) Save(session model.Session) error {
	ctx := context.Background()
	if err := s.cache.Save(ctx, session); err != nil {
		fmt.Printf("cache put failed")
	}

	err := s.db.SaveSession(session)
	if err != nil {
		fmt.Printf("Error while saving to database")
	}
	return err
}

func (s *SessionStore) Get(key string) (*model.Session, error) {
	ctx := context.Background()
	session, err := s.cache.Get(ctx, key)

	if err == nil {
		fmt.Printf("Key %s found in cache", key)
		return session, nil
	}
	fmt.Printf("Session %s not found in cache", key)

	session, err = s.db.GetSession(key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Session %s not found", key))
	}
	return session, nil

}

func (u *UserStore) Get(userID string) (*model.User, error) {
	ctx := context.Background()
	user, err := u.cache.Get(ctx, userID)

	if err == nil {
		fmt.Printf("User %s found in cache", userID)
		return user, nil
	}
	fmt.Printf("User %s not found in cache", userID)

	user, err = u.db.GetUser(userID)
	return user, err

}
