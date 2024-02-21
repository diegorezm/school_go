package userhandler

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"fmt"
	userModel "github/diegorezm/school_go/api/models/user"
	"log"
)

const (
	SELECT_QUERY          = "SELECT id,name,email,password FROM school.users"
	SELECT_QUERY_BY_ID    = "SELECT id,name,email,password FROM school.users WHERE id= ?"
	SELECT_QUERY_BY_EMAIL = "SELECT id,name,email,password FROM school.users WHERE email= ?"
	INSERT_QUERY          = "insert into users (name, email, password) values (?, ?, ?)"
	SECRET_KEY            = "abc&1*~#^2^#s0^=)^^7%b34"
)

// i have no idea of what is happening here
// stole everything from https://blog.logrocket.com/learn-golang-encryption-decryption/

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Encrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

func Decrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

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
		// probably shouldn't return the password here,
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
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return userModel.UserModel{}, fmt.Errorf("No user found with this email!")
		} else {
			log.Printf("Error scanning row: %v", err)
			return userModel.UserModel{}, fmt.Errorf("Error while trying to fetch user!")
		}
	}
	password_hash, err := Decrypt(user.Password, SECRET_KEY)
	if err != nil {
		log.Print(err.Error())
		return userModel.UserModel{}, fmt.Errorf("Error while decoding the password!")
	}
	if password != password_hash {
		return userModel.UserModel{}, fmt.Errorf("Wrong password!")
	}
	return user, nil
}

func (uh UserHandler) Register(newUser userModel.UserModel) (userModel.UserModel, error) {
	password_hash, err := Encrypt(newUser.Password, SECRET_KEY)
	if err != nil {
		log.Print(err.Error())
		return userModel.UserModel{}, fmt.Errorf("Error while encoding the password!")
	}
	result, err := uh.connection.Exec(INSERT_QUERY, newUser.Name, newUser.Email, &password_hash)
	if err != nil {
		log.Print(err.Error())
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
