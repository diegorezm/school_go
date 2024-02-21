package main

import (
	"encoding/json"
	"fmt"
	"github/diegorezm/school_go/api/db"
	studentshandler "github/diegorezm/school_go/api/handlers/studentsHandler"
	"github/diegorezm/school_go/api/handlers/templateHandler"
	"github/diegorezm/school_go/api/handlers/userHandler"
	studentModel "github/diegorezm/school_go/api/models/student"
	userModel "github/diegorezm/school_go/api/models/user"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type TemplateData struct {
	Students    []studentModel.StudentModel
	User        userModel.UserModel
	EditStudent studentModel.StudentModel
}

const BASE_TEMPLATE = "./templates/base.html"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
  port := os.Getenv("PORT")


	database := db.NewDatabase(password, dbname)
	connection := database.Connection

	uh := userhandler.NewUserHandler(connection)
	sh := studentshandler.NewStudentHandler(connection)

	studentsTable := templatehandler.NewTemplateHandler("./templates/students_table.html", BASE_TEMPLATE)
	loginForm := templatehandler.NewTemplateHandler("./templates/login.html", BASE_TEMPLATE)
	registerForm := templatehandler.NewTemplateHandler("./templates/register.html", BASE_TEMPLATE)

	allStudents, err := sh.GetAllStudents()

	if err != nil {
		log.Print(err.Error())
	}

	context := TemplateData{
		Students:    allStudents,
		User:        userModel.UserModel{},
		EditStudent: studentModel.StudentModel{},
	}
	// RENDER FUNC
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if context.User == (userModel.UserModel{}) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/students", http.StatusSeeOther)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginForm.Render(w, r, context)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		registerForm.Render(w, r, context)
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		if context.User == (userModel.UserModel{}) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		studentsTable.Render(w, r, context)
	})
	// -------------------------------------------------------------------------

	// HANDLER FUNC

	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		formEmail := r.PostFormValue("email")
		formPassword := r.PostFormValue("password")
		var tmpl *template.Template
		usr, err := uh.Login(formEmail, formPassword)
		if err != nil {
			errorString := fmt.Sprintf("<span class='text-danger'>%s</span>", err.Error())
			tmpl, _ = template.New("t").Parse(errorString)
			tmpl.Execute(w, nil)
			return
		}
		context.User = usr
		html := fmt.Sprintf(`
        <span>Logged in successfully! Welcome back %s!</span>
        <script>
          setTimeout(() => {
            window.location.href = "/students";
            }, 1000)
        </script>
    `, usr.Name)
		tmpl, _ = template.New("t").Parse(html)
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		formUsername := r.PostFormValue("username")
		formEmail := r.PostFormValue("email")
		formPassword := r.PostFormValue("password")
		var tmpl *template.Template

		newUser := userModel.NewUser(formUsername, formPassword, formEmail, -1)
		_, err := uh.Register(*newUser)
		if err != nil {
			errorString := fmt.Sprintf("<span class='text-danger'>%s</span>", err.Error())
			tmpl, _ = template.New("t").Parse(errorString)
			tmpl.Execute(w, nil)
		} else {
			html := `
        <span>User created!</span>
        <script>
          setTimeout(() => {
            window.location.href = "/login";
            }, 500)
        </script>
    `
			tmpl, _ = template.New("t").Parse(html)
			tmpl.Execute(w, nil)
		}
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		users, err := uh.GetAllUsers()
		if err != nil {
			log.Fatal(err.Error())
		}
		usersJSON, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
		w.Write(usersJSON)
	})

	http.HandleFunc("/students/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		formName := r.PostFormValue("name")
		formEmail := r.PostFormValue("email")
		formCourse := r.PostFormValue("course")
		formAgeStr := r.PostFormValue("age")
		formAge, err := strconv.Atoi(formAgeStr)

		var tmpl *template.Template
		if err != nil {
			errorString := fmt.Sprintf("<span class='text-danger'>Please provide a valid age!</span>")
			tmpl, _ = template.New("t").Parse(errorString)
			tmpl.Execute(w, nil)
			return
		}

		newStudent := studentModel.NewStudent(formAge, 0, formName, formEmail, formCourse)
		student, err := sh.CreateNewStudent(*newStudent)

		if err != nil {
			errorString := fmt.Sprintf("<span class='text-danger'>%s</span>", err.Error())
			tmpl, _ = template.New("t").Parse(errorString)
			tmpl.Execute(w, nil)
			return
		}
		context.Students = append(context.Students, student)
		html := `
        <span>Student created successfully!</span>
        <script>
          window.location.reload();
        </script>
    `
		tmpl, _ = template.New("t").Parse(html)
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/students/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var tmpl *template.Template
		studentID := r.PostFormValue("id")
    log.Print(studentID)
    err = sh.DeleteStudent(studentID)
    if err != nil {
			errorString := "<span class='text-danger'>Not able to delete student!</span>"
			tmpl, _ = template.New("t").Parse(errorString)
			tmpl.Execute(w, nil)
			return
    }
			html := `
      <span class='text-danger'>Student deleted!<span>
      <script>
        window.location.reload()
      </script> 
      `
			tmpl, _ = template.New("t").Parse(html)
			context.Students, _ = sh.GetAllStudents()
			tmpl.Execute(w, nil)
	})

	http.HandleFunc("/students/edit", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var tmpl *template.Template
			err := r.ParseForm()
			if err != nil {
				errorString := fmt.Sprintf("<span class='text-danger'>%s</span>", err.Error())
				tmpl, _ = template.New("t").Parse(errorString)
				tmpl.Execute(w, nil)
				return
			}
			studentID := r.Form.Get("student__id")
			getStudent, err := sh.GetStudentById(studentID)
			if err != nil {
				log.Print(err.Error())
				errorString := fmt.Sprintf("<span class='text-danger'>%s</span>", err.Error())
				tmpl, _ = template.New("t").Parse(errorString)
				tmpl.Execute(w, nil)
				return
			}
			context.EditStudent = getStudent

			html := fmt.Sprintf(`
      <div class="mb-3">
          <label for="exampleInputEmail1" class="form-label">Name</label>
          <input type="text" class="form-control" id="name" aria-describedby="name" name="name" value="%s">
      </div>

      <div class="mb-3">
          <label for="exampleInputEmail1" class="form-label">Email address</label>
          <input type="email" class="form-control" id="exampleInputEmail1" aria-describedby="emailHelp" name="email" value="%s">
      </div>

      <div class="mb-3">
          <label for="exampleInputEmail1" class="form-label">Course</label>
          <input type="text" class="form-control" id="course" aria-describedby="course" name="course" value="%s">
      </div>

      <div class="mb-3">
          <label for="exampleInputEmail1" class="form-label">Age</label>
          <input type="number" class="form-control" id="age" aria-describedby="age" name="age" value="%d">
      </div>

      <input type="text" class="form-control" id="st_id" aria-describedby="id" name="id" value="%s" hidden>
      <input type="text" class="form-control" id="st_id" aria-describedby="id" name="createdAt" value="%s" hidden>
      <input type="text" class="form-control" id="st_id" aria-describedby="id" name="updatedAt" value="%s" hidden>
      <button type="submit" class="btn btn-primary">
          Submit
      </button>
`,
				context.EditStudent.Name,
				context.EditStudent.Email, context.EditStudent.Course, context.EditStudent.Age,
				fmt.Sprint(context.EditStudent.Id), fmt.Sprint(context.EditStudent.CreatedAt), fmt.Sprint(context.EditStudent.UpdatedAt))

			tmpl, _ = template.New("t").Parse(html)
			tmpl.Execute(w, nil)
		case http.MethodPut:
			log.Print(context.EditStudent)
			formName := r.PostFormValue("name")
			formEmail := r.PostFormValue("email")
			formCourse := r.PostFormValue("course")
			formAgeStr := r.PostFormValue("age")
			formIDStr := r.PostFormValue("id")
			createdAtStr := r.PostFormValue("createdAt")
			updatedAtStr := r.PostFormValue("updatedAt")
			log.Printf(formIDStr)

			// Convert age string to int
			formAge, _ := strconv.Atoi(formAgeStr)
			formID, _ := strconv.Atoi(formIDStr)

			layout := "2006-01-02 15:04:05 -0700 MST"
			createAt, err := time.Parse(layout, createdAtStr)
			if err != nil {
				log.Printf("ERROR on createdAt parsing. %s", err.Error())
				return
			}
			updateAt, err := time.Parse(layout, updatedAtStr)
			if err != nil {
				log.Printf("ERROR on updatedAt parsing. %s", err.Error())
				return
			}

			updatedStudent := studentModel.StudentModel{
				Name:      formName,
				Email:     formEmail,
				Course:    formCourse,
				Age:       formAge,
				Id:        formID,
				CreatedAt: createAt,
				UpdatedAt: updateAt,
			}
			_, err = sh.UpdateStudent(updatedStudent)
			if err != nil {
				log.Printf(err.Error())
			}
			var tmpl *template.Template
			js := fmt.Sprintf(`
        Student %s updated!
        <script>
          window.location.reload()
        </script>
      `, context.EditStudent.Name)
			tmpl, _ = template.New("t").Parse(js)
			context.EditStudent = studentModel.StudentModel{}
			context.Students, _ = sh.GetAllStudents()
			tmpl.Execute(w, nil)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	})
	// -------------------------------------------------------------------------

	// creating a goroutine for the server
	go func() {
		fmt.Printf("Server running on port http://localhost%s\n", port)
		log.Fatal(http.ListenAndServe(port, nil))
	}()
	select {}
}
