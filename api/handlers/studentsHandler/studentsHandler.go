package studentshandler

import (
	"database/sql"
	"fmt"
	"github/diegorezm/school_go/api/models/student"
	"log"
	"time"
)

const (
	SELECT_QUERY          = "SELECT id,name,email,age,course,created_at,updated_at FROM school.students"
	SELECT_QUERY_BY_ID    = "SELECT id,name,email,age,course,created_at,updated_at FROM school.students WHERE id= ?"
	SELECT_QUERY_BY_EMAIL = "SELECT id,name,email,age,course,created_at,updated_at FROM school.students WHERE email= ?"
	INSERT_QUERY          = "insert into school.students (name, email, age, course) values (?, ?, ?, ?)"
	UPDATE_QUERY          = "UPDATE school.students SET name= ?, email= ?, course= ?, age= ? WHERE id= ?;"
	DELETE_QUERY          = "DELETE FROM school.students WHERE id = ?"
	LAYOUT                = "2006-01-02 15:04:05"
)

type StudentHandler struct {
	connection *sql.DB
}

func NewStudentHandler(conn *sql.DB) *StudentHandler {
	return &StudentHandler{connection: conn}
}

func (sh StudentHandler) scanRow(row *sql.Row) (student.StudentModel, error) {
	var newStudent student.StudentModel

	var createdAtStr, updatedAtStr string
	if err := row.Scan(&newStudent.Id, &newStudent.Name, &newStudent.Email, &newStudent.Age, &newStudent.Course, &createdAtStr, &updatedAtStr); err != nil {
		if err == sql.ErrNoRows {
			return student.StudentModel{}, fmt.Errorf("This student does not exist/was not found.")
		}
		return student.StudentModel{}, fmt.Errorf("Error while trying to scan student.")
	}

	createAt, err := time.Parse(LAYOUT, createdAtStr)
	if err != nil {
		return student.StudentModel{}, fmt.Errorf("Error parsing created_at value: %s", err.Error())
	}
	updateAt, err := time.Parse(LAYOUT, updatedAtStr)
	if err != nil {
		return student.StudentModel{}, fmt.Errorf("Error parsing updatedAt value: %s", err.Error())
	}

	newStudent.CreatedAt = createAt
	newStudent.UpdatedAt = updateAt

	return newStudent, nil
}

func (sh StudentHandler) GetAllStudents() ([]student.StudentModel, error) {
	rows, err := sh.connection.Query(SELECT_QUERY)
	if err != nil {
		return nil, fmt.Errorf("Error while trying to get all students data.")
	}
	defer rows.Close()
	var students []student.StudentModel
	for rows.Next() {
		var newStudent student.StudentModel
		var createdAtStr, updatedAtStr string
		if err := rows.Scan(&newStudent.Id, &newStudent.Name, &newStudent.Email, &newStudent.Age, &newStudent.Course, &createdAtStr, &updatedAtStr); err != nil {
			log.Printf(err.Error())
			return nil, fmt.Errorf("Error while trying to scan student.")
		}
		createAt, err := time.Parse(LAYOUT, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("Error parsing created_at value: %s", err.Error())
		}
		updateAt, err := time.Parse(LAYOUT, updatedAtStr)
		if err != nil {
			return nil, fmt.Errorf("Error parsing created_at value: %s", err.Error())
		}

		newStudent.CreatedAt = createAt
		newStudent.UpdatedAt = updateAt
		students = append(students, newStudent)
	}
	return students, nil
}

func (sh StudentHandler) CreateNewStudent(newStudent student.StudentModel) (student.StudentModel, error) {
	result, err := sh.connection.Exec(INSERT_QUERY, newStudent.Name, newStudent.Email, newStudent.Age, newStudent.Course)
	if err != nil {
		return student.StudentModel{}, fmt.Errorf("Error while inserting student!")
	}
	studentId, err := result.LastInsertId()
	if err != nil {
		return student.StudentModel{}, fmt.Errorf("Error while retrieving student from the databse!")
	}
	row := sh.connection.QueryRow(SELECT_QUERY_BY_ID, studentId)

	newStudent, err = sh.scanRow(row)
	if err != nil {
		return student.StudentModel{}, fmt.Errorf(err.Error())
	}

	return newStudent, nil
}

func (sh StudentHandler) GetStudentById(id string) (student.StudentModel, error) {
	row := sh.connection.QueryRow(SELECT_QUERY_BY_ID, id)
	var newStudent student.StudentModel
	newStudent, err := sh.scanRow(row)
	if err != nil {
		return student.StudentModel{}, fmt.Errorf(err.Error())
	}
	return newStudent, nil
}

func (sh StudentHandler) UpdateStudent(formStudent student.StudentModel) (student.StudentModel, error) {
	_, err := sh.connection.Exec(UPDATE_QUERY, formStudent.Name, formStudent.Email, formStudent.Course, formStudent.Age, formStudent.Id)
	if err != nil {

		return student.StudentModel{}, fmt.Errorf(err.Error())
	}
	return formStudent, nil
}

func (sh StudentHandler) DeleteStudent(id string) error {
	_, err := sh.connection.Exec(DELETE_QUERY, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("This student does not exist/was not found.")
		}
		return err
	}
	return nil
}
