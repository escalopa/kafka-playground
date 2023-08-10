package kafka_playground

import (
	"encoding/json"
	"time"
)

type User struct {
	ID        string    `json:"iD"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"birthday"`
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
