package student

import "time"

type StudentModel struct {
	Id        int
	Name      string
	Email     string
	Course    string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewStudent(age, id int, name, email, course string) *StudentModel {
	return &StudentModel{
		Id:     id,
		Name:   name,
		Email:  email,
		Course: course,
		Age:    age,
	}
}
