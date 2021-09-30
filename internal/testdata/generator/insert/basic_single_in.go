package testdata

import (
	"time"
)

type User struct {
	Id        int       `sql:"id"`
	Email     string    `sql:"email"`
	FullName  string    `sql:"full_name"`
	CreatedAt time.Time `sql:"created_at"`
}

type InsertBasicSingleQuery struct {
	User *User `rel:"test_user:u"`
}
