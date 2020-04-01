package users

import (
	"fmt"
	"strings"

	"github.com/amisini/bookstore_users-api/datasources/mysql/users_db"
	"github.com/amisini/bookstore_users-api/utils/date_utils"
	"github.com/amisini/bookstore_users-api/utils/errors"
)

const (
	indexUniqueEmail = "email_UNIQUE"
	errorNoRows      = "no rows in result set"
	queryInsertUser  = ("INSERT INTO users(first_name, last_name, email, date_created) VALUES (?, ?, ?, ?);")
	queryGetUser     = ("SELECT id, first_name, last_name, email, date_created FROM users WHERE id=?;")
	queryUpdateUser  = ("UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;")
	queryDeleteUser  = ("DELETE FROM users WHERE id=?;")
)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("User %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when getting user: %d :%s", user.Id, err.Error()))
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()
	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		if strings.Contains(err.Error(), indexUniqueEmail) {
			return errors.NewBadRequestError(fmt.Sprintf("Email %s already exists", user.Email))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when saving user: %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when saving user: %s", err.Error()))
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return errors.NewBadRequestError(fmt.Sprintf("Error updating user: %d", user.Id))
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	if _, err = stmt.Exec(user.Id); err != nil {
		return errors.NewBadRequestError(fmt.Sprintf("Error deleting user: %d", user.Id))
	}
	return nil
}
