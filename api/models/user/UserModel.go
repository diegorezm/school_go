package user

type UserModel struct {
	Id       int
	Name     string
	Email    string
	Password string
}

func NewUser(name, password, email string, id int) *UserModel {
	return &UserModel{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
}
