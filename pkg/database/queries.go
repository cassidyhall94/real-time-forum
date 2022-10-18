package database

import "fmt"

type User struct {
	ID       string
	Email    string
	Username string
	Password string
	Nickname string
}

func GetUsers() ([]User, error) {
	users := []User{}
	rows, err := DB.Query(`SELECT * FROM users`)
	if err != nil {
		return users, fmt.Errorf("GetUsers DB Query error: %+v\n", err)
	}
	var id string
	var email string
	var username string
	var password string
	var nickname string
	for rows.Next() {
		err := rows.Scan(&id, &email, &username, &password, &nickname)
		if err != nil {
			return users, fmt.Errorf("GetUsers rows.Scan error: %+v\n", err)
		}
		users = append(users, User{
			ID:       id,
			Email:    email,
			Username: username,
			Password: password,
			Nickname: nickname,
		})
	}
	err = rows.Err()
	if err != nil {
		return users, err
	}
	return users, nil
}
