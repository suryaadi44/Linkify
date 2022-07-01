package entity

import (
	"time"
)

type User struct {
	Username string    `bson:"_id"`
	Password string    `bson:"pass"`
	Email    string    `bson:"email"`
	Rank     int       `bson:"rank"`
	Created  time.Time `bson:"created"`
}
