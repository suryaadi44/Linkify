package entity

import (
	"time"
)

type User struct {
	UID      string    `bson:"_id"`
	Username string    `bson:"username"`
	Password string    `bson:"pass"`
	Email    string    `bson:"email"`
	Rank     int       `bson:"rank"`
	Created  time.Time `bson:"created"`
}
