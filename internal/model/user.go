package model

import "time"

type User struct {
	ID           string    `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    Username     string    `db:"username" json:"username"`
    Password    string      `db:"password" json:"password,omitempty"`
    Role         string    `db:"role" json:"role"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
