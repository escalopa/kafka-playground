package kafka_playground

import (
	"encoding/json"
	"time"

	"github.com/brianvoe/gofakeit"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUser() User {
	return User{
		ID:        gofakeit.UUID(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		CreatedAt: time.Now(),
	}
}

func (u User) Encode() ([]byte, error) {
	return json.Marshal(&u)
}

func (u User) Length() int {
	b, _ := json.Marshal(&u)
	return len(b)
}

type Key string

func (k Key) Encode() ([]byte, error) {
	return []byte(k), nil
}

func (k Key) Length() int {
	return len(k)
}
