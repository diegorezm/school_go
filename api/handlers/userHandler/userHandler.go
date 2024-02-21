package userhandler

import (
	"database/sql"
	"fmt"
	userModel "github/diegorezm/school_go/api/models/user"
	"log"
)

const (
	SELECT_QUERY          = "SELECT id,name,email,password FROM school.users"
	SELECT_QUERY_BY_ID    = "SELECT id,name,email,password FROM school.users WHERE id= ?"
	SELECT_QUERY_BY_EMAIL = "SELECT id,name,email,password FROM school.users WHERE email= ?"
	INSERT_QUERY          = "insert into users (name, email, password) values (?, ?, ?)"
)

type UserHandler struct {
	connection *sql.DB
}

func NewUserHandler(conn *sql.DB) *UserHandler {
	return &UserHandler{connection: conn}
}

func (uh UserHandler) GetAllUsers() ([]userModel.UserModel, error) {
	rows, err := uh.connection.Query(SELECT_QUERY)
	if err != nil {
		return nil, fmt.Errorf("Erro while trying to get all users.")
	}
	defer rows.Close()
	var users []userModel.UserModel
	for rows.Next() {
		var user userModel.UserModel
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
			return nil, fmt.Errorf("Error while trying to scan users.")
		}
		users = append(users, user)
	}
	return users, nil
}

func (uh UserHandler) Login(email, password string) (userModel.UserModel, error) {
	row := uh.connection.QueryRow(SELECT_QUERY_BY_EMAIL, email)
	var user userModel.UserModel
	err := row.Scan(&user.Id, &user.Name,&user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return userModel.UserModel{}, fmt.Errorf("No user found with this email!")
		} else {
      log.Printf("Error scanning row: %v", err)
			return userModel.UserModel{}, fmt.Errorf("Error while trying to fetch user!")
		}
	}
	return user, nil
}

func (uh UserHandler) Register(newUser userModel.UserModel) (userModel.UserModel, error) {
	result, err := uh.connection.Exec(INSERT_QUERY, newUser.Name, newUser.Email, newUser.Password)
	if err != nil {
		return userModel.UserModel{}, fmt.Errorf("Error while inserting user!")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return userModel.UserModel{}, fmt.Errorf("Error while retrieving user from the databse!")
	}

	row := uh.connection.QueryRow(SELECT_QUERY_BY_ID, userID)
	var user userModel.UserModel

	if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password); err != nil {
		return userModel.UserModel{}, fmt.Errorf("Error while scanning row!")
	}

	return user, nil
}
