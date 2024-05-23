package user

import (
	"database/sql"
	"fmt"
	"go-ecommerce/model"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*model.User, error) {
	rows, err := s.db.Query("SELECT * FROM Users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func scanRowIntoUser(rows *sql.Rows) (*model.User, error) {
	user := new(model.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *Store) GetUserById(id int) (*model.User, error) {
	rows, err := s.db.Query("SELECT * FROM Users WHERE email = ?", id)
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Store) CreateUser(user model.User) error {
	_, err := s.db.Exec("INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}
