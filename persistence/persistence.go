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
	cache *cache.RedisClient
	db    *sql.Mysql
}

type SessionStore struct {
	cache *cache.RedisClient
	db    *sql.Mysql
}

func GetSessionStore() Session {
	return &SessionStore{
		cache: cache.GetRedisClient(),
		db:    sql.GetMySql(),
	}
}

func GetUserStore() User {
	return &UserStore{
		cache: cache.GetRedisClient(),
		db:    sql.GetMySql(),
	}
}

func (u *UserStore) Save(user model.User) error {
	ctx := context.Background()
	if err := u.cache.Save(ctx, user.ID, model.User{}); err != nil {
		fmt.Printf("cache put failed %s\n", err.Error())
	}

	err := u.db.SaveUser(user)
	if err != nil {
		fmt.Printf("Error while saving to database %s\n", err.Error())
	}
	return err
}

func (s *SessionStore) Save(session model.Session) error {
	ctx := context.Background()
	if err := s.cache.Save(ctx, session.ID, model.Session{}); err != nil {
		fmt.Printf("cache put failed %s\n", err.Error())
	}

	err := s.db.SaveSession(session)
	if err != nil {
		fmt.Printf("Error while saving to database %s\n", err.Error())
	}
	return err
}

func (s *SessionStore) Get(key string) (*model.Session, error) {
	ctx := context.Background()
	value, err := s.cache.Get(ctx, key, model.Session{})

	if err == nil {
		fmt.Printf("Key %s found in cache\n", key)
		session := value.(model.Session)
		return &session, nil
	}
	fmt.Printf("Session %s not found in cache\n", key)

	session, err := s.db.GetSession(key)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Session %s not found", key))
	}
	s.cache.Save(ctx, key, model.Session{})
	return session, nil

}

func (u *UserStore) Get(userID string) (*model.User, error) {
	ctx := context.Background()
	value, err := u.cache.Get(ctx, userID, model.User{})

	if err == nil {
		fmt.Printf("User %s found in cache", userID)
		user := value.(model.User)
		return &user, nil
	}
	fmt.Printf("User %s not found in cache\n", userID)

	user, err := u.db.GetUser(userID)
	if err != nil {
		fmt.Printf("error %s while query database for user %s", err.Error(), userID)
		return nil, errors.New(fmt.Sprintf("user %s not found", userID))
	}
	fmt.Printf("user %s found in database %v \n", userID, user)
	return user, err
}
