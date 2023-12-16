package users

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (c UserService) GetUserByID(userID int) (*User, error) {
	rows, err := c.DB.Query("SELECT * FROM users WHERE id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	user := &User{}
	if rows.Next() {
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, errors.New("user not found")
}
