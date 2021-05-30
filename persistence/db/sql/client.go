package sql

import (
	"errors"
	"fmt"
	"os"

	"github.com/abhaybhu10/login/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Password string
	Email    string
	Name     string
}

type Session struct {
	gorm.Model
	SessionID string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
}

type Mysql struct {
	db *gorm.DB
}

func (m *Mysql) SaveUser(user model.User) error {
	var userRes User
	if result := m.db.First(&userRes); result.RowsAffected != 0 {
		return errors.New(fmt.Sprintf("User %s already exist", user.ID))
	}

	userData := &User{
		ID:       user.ID,
		Password: user.Password,
		Email:    user.Email,
		Name:     user.Name,
	}
	result := m.db.Create(userData)
	return result.Error
}

func (m *Mysql) GetUser(userId string) (*model.User, error) {
	var user User
	if result := m.db.Find(&user); result.Error != nil {
		return nil, result.Error
	}

	return &model.User{
		ID:       user.ID,
		Password: user.Password,
		Email:    user.Email,
		Name:     user.Name,
	}, nil

}

func (m *Mysql) SaveSession(session model.Session) error {
	sessionData := &Session{
		SessionID: session.SessionID,
		UserID:    session.UserId,
	}
	result := m.db.Create(sessionData)

	return result.Error
}

func (m *Mysql) GetSession(Id string) (*model.Session, error) {
	var session Session

	if result := m.db.Find(&session); result.Error != nil {
		return nil, result.Error
	}
	return &model.Session{
		SessionID: session.SessionID,
		UserId:    session.UserID,
	}, nil
}

func GetMySql() *Mysql {
	userName, isPresent := os.LookupEnv("USERNAME")
	if !isPresent {
		panic("USERNAME not set")
	}
	password, isPresent := os.LookupEnv("PASSWORD")

	dbName, isPresent := os.LookupEnv("DB_NAME")
	dbUrl, isPresent := os.LookupEnv("DB_URL")

	conn := fmt.Sprintf("%s:%s@%s/%s", userName, password, dbUrl, dbName)
	fmt.Printf("Connecting to mysql %s\n", conn)

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
		panic("Error while starting up database")
	}
	return &Mysql{
		db: db,
	}
}
