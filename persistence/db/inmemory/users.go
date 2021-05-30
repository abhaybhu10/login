package db

import (
	"errors"

	"github.com/abhaybhu10/login/model"
)

type InMemoryDB struct {
	database map[string]model.User
}

func NewInMomoryDB() *InMemoryDB {
	return &InMemoryDB{
		database: map[string]model.User{},
	}
}
func (db *InMemoryDB) Save(user model.User) error {
	db.database[user.ID] = user
	return nil
}

func (db *InMemoryDB) Get(userId string) (*model.User, error) {
	user, ok := db.database[userId]
	if !ok {
		return nil, errors.New("User does not exist")
	}
	return &user, nil
}
