//Cahlil Tillett
//Tadeo Bennett

package postgresql

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"sysadmin.com/final/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) //iterates the hash 12  times
	if err != nil {
		return err
	}

	s := `INSERT INTO users(username, hashed_password) VALUES ($1, $2)`

	_, err = u.DB.Exec(s, username, hashedPassword)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicated key value violates unique constraint "users_email_key"`:
			return models.ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (u *UserModel) Authenticate(username string, password string) (bool, error) {
	var id int
	var hashedPassword []byte

	s := `SELECT id, hashed_password FROM users WHERE username = $1`

	err := u.DB.QueryRow(s, username).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, models.ErrInvalidCredentials
		} else {
			return false, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, models.ErrInvalidCredentials
		} else {
			return false, err
		}
	}

	return true, nil
}
