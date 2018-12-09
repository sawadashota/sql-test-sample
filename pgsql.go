package sql_test_sample

import (
	"database/sql"
)

// DB is database connection
var DB *sql.DB

// User table columns
type User struct {
	Id   int
	Name string
	Sex  string
}

// GetUser by ID
func GetUser(id int) (*User, error) {
	var u User

	if err := DB.QueryRow("SELECT id, name, sex FROM users WHERE id = $1", id).Scan(&u.Id, &u.Name, &u.Sex); err != nil {
		return nil, err
	}

	return &u, nil
}

// InsertUser to database
func InsertUser(user *User) error {
	if _, err := DB.Query("INSERT INTO users (name, sex) VALUES ($1, $2)", user.Name, user.Sex); err != nil {
		return err
	}

	return nil
}
